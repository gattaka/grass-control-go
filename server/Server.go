package server

import (
	vlcctrl "github.com/CedArctic/go-vlc-ctrl"
	"grass-control-go/indexer"
	"grass-control-go/ui"
	"grass-control-go/ui/common"
	"grass-control-go/ui/playlist"
	"grass-control-go/utils"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Resources struct {
	Styles  string
	Scripts string
	Favicon string
	Icons   string
}

func StartServer(port int, myVLC vlcctrl.VLC, indexer indexer.Indexer, resources Resources) {
	var errorPending *error

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Has(common.SearchParam) {
			// hledá se cokoliv dle názvu
			searchQuery := r.URL.Query().Get(common.SearchParam)
			if searchQuery != "" {
				io.WriteString(w, ui.ConstructPage(indexer.FindByString(searchQuery), true, searchQuery))
				return
			}
		} else if r.URL.Query().Has(common.DirParam) {
			// hledá se přímo dle adresáře
			dirQuery := r.URL.Query().Get(common.DirParam)
			if len(dirQuery) > 0 {
				io.WriteString(w, ui.ConstructPage(indexer.FindByPath(dirQuery), false, dirQuery))
				return
			}
		}

		// výchozí pohled
		io.WriteString(w, ui.ConstructPage(indexer.GetAllItems(), false, "/"))
	})

	http.HandleFunc("/resources/styles.css", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, resources.Styles) })
	http.HandleFunc("/resources/scripts.js", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, resources.Scripts) })
	http.HandleFunc("/resources/favicon.png", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, resources.Favicon) })
	http.HandleFunc("/resources/icons.png", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, resources.Icons) })

	ret := func(w http.ResponseWriter, err error) {
		if err != nil {
			*errorPending = err
		}
	}

	prepURL := func(path string) string {
		return "file:///" + utils.EncodeURL(indexer.GetPlayerRoot()+path)
	}

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		response, _ := myVLC.RequestMaker("/requests/status.json")
		operations := utils.UpdateStatus(errorPending, response)
		io.WriteString(w, operations)
	})

	http.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) { ret(w, myVLC.EmptyPlaylist()) })
	http.HandleFunc("/clearExceptPlaying", func(w http.ResponseWriter, r *http.Request) {
		statusJson, _ := myVLC.RequestMaker("/requests/status.json")
		status := utils.ParseStatus(statusJson)
		currentplid := status.Currentplid

		if currentplid == -1 {
			ret(w, myVLC.EmptyPlaylist())
			return
		}

		playlistJson, _ := myVLC.RequestMaker("/requests/playlist.json")
		playlist, _ := utils.ParsePlaylist(playlistJson)

		for _, item := range *playlist {
			id, _ := strconv.Atoi(item.Id)
			if id != currentplid {
				myVLC.Delete(id)
			}
		}
	})
	http.HandleFunc("/pause", func(w http.ResponseWriter, r *http.Request) { ret(w, myVLC.Pause()) })
	http.HandleFunc("/next", func(w http.ResponseWriter, r *http.Request) { ret(w, myVLC.Next()) })
	http.HandleFunc("/prev", func(w http.ResponseWriter, r *http.Request) { ret(w, myVLC.Previous()) })
	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) { ret(w, myVLC.Stop()) })
	// chyba v názvu funkce, ve skutečnosti volá Random (shuffle)
	http.HandleFunc("/shuffle", func(w http.ResponseWriter, r *http.Request) { ret(w, myVLC.ToggleLoop()) })
	http.HandleFunc("/loop", func(w http.ResponseWriter, r *http.Request) { ret(w, myVLC.ToggleRepeat()) })
	http.HandleFunc("/volume", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Has(common.ValueParam) {
			val, err := strconv.Atoi(query.Get(common.ValueParam))
			if err != nil {
				return
			}
			ret(w, myVLC.Volume(strconv.Itoa(val)))
		}
	})
	http.HandleFunc("/progress", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Has(common.ValueParam) {
			val, err := strconv.Atoi(query.Get(common.ValueParam))
			if err != nil {
				return
			}
			ret(w, myVLC.Seek(strconv.Itoa(val)))
		}
	})
	http.HandleFunc("/playlist", func(w http.ResponseWriter, r *http.Request) {
		resp, err := myVLC.RequestMaker("/requests/playlist.json")
		if err != nil {
			// TODO
			return
		}
		items, hash := utils.ParsePlaylist(resp)
		render := playlist.ConstructPlaylistTable(*items).Render()
		io.WriteString(w, utils.ConstructPlaylistJSON(render, hash))
	})

	fileOrSearch := func(w http.ResponseWriter, r *http.Request, idOperation func(w http.ResponseWriter, r *http.Request), searchOperation func(w http.ResponseWriter, r *http.Request)) {
		if r.URL.Query().Has(common.IdParam) {
			idOperation(w, r)
		} else if r.URL.Query().Has(common.SearchParam) {
			searchOperation(w, r)
		}
	}

	getParam := func(r *http.Request, paramName string) string {
		query := r.URL.Query()
		return query.Get(paramName)
	}
	getIdParam := func(r *http.Request) string { return getParam(r, common.IdParam) }

	vlcAdd := func(param string) { myVLC.Add(prepURL(param)) }
	vlcAddPlay := func(param string) { myVLC.AddStart(prepURL(param)) }

	// Když se do vlc přidá 1 soubor a 1 adresář, VLC mězi ně pro shuffle rozdělí pravděpodobnost 50:50 na spuštění -- dokud nepadne
	// adresář, "nerozbalí" ho na jednotlivé soubory a mezi ně znovu rozpočítá pravděpodobnost. První soubor se tedy spouští často opakovaně
	// než losování konečně padne na adresář
	directSeek := func(r *http.Request, vlcOperaton func(string)) {
		path := getIdParam(r)
		for _, item := range indexer.ExpandByPath(path) {
			vlcOperaton(item.GetPath())
		}
	}
	searchSeek := func(r *http.Request, vlcOperaton func(string)) {
		searchQuery := getParam(r, common.SearchParam)
		results := indexer.FindByString(searchQuery)
		for _, result := range results {
			for _, item := range indexer.ExpandByItem(result) {
				vlcOperaton(item.GetPath())
			}
		}
	}

	directAddToPlaylist := func(w http.ResponseWriter, r *http.Request) { directSeek(r, vlcAdd) }
	directAddToPlaylistPlay := func(w http.ResponseWriter, r *http.Request) { directSeek(r, vlcAddPlay) }

	searchAddToPlaylist := func(w http.ResponseWriter, r *http.Request) { searchSeek(r, vlcAdd) }
	searchAddToPlaylistPlay := func(w http.ResponseWriter, r *http.Request) { searchSeek(r, vlcAddPlay) }

	playFromPlaylist := func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(getIdParam(r))
		if err == nil {
			ret(w, myVLC.Play(id))
		}
	}
	removeFromPlaylist := func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(getIdParam(r))
		if err == nil {
			ret(w, myVLC.Delete(id))
		}
	}

	http.HandleFunc(common.AddEndpoint, func(w http.ResponseWriter, r *http.Request) {
		fileOrSearch(w, r, directAddToPlaylist, searchAddToPlaylist)
	})
	http.HandleFunc(common.AddAndPlayEndpoint, func(w http.ResponseWriter, r *http.Request) {
		fileOrSearch(w, r, directAddToPlaylistPlay, searchAddToPlaylistPlay)
	})
	http.HandleFunc(common.PlayEndpoint, func(w http.ResponseWriter, r *http.Request) { playFromPlaylist(w, r) })
	http.HandleFunc(common.RemoveEndpoint, func(w http.ResponseWriter, r *http.Request) { removeFromPlaylist(w, r) })

	http.HandleFunc("/reindex", func(w http.ResponseWriter, r *http.Request) {
		indexer.Reindex()
		io.WriteString(w, ui.ConstructPage(indexer.GetAllItems(), false, "/"))
	})

	http.HandleFunc("/quit", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Aplikace byla ukončena")
		os.Exit(0)
	})

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

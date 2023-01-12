package server

import (
	vlcctrl "github.com/CedArctic/go-vlc-ctrl"
	"grass-control-go/indexer"
	"grass-control-go/ui"
	"grass-control-go/utils"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
		if r.URL.Query().Has(ui.SearchParam) {
			// hledá se cokoliv dle názvu
			searchQuery := r.URL.Query().Get(ui.SearchParam)
			if searchQuery != "" {
				io.WriteString(w, ui.ConstructPage(indexer.FindByString(searchQuery), true, searchQuery))
				return
			}
		} else if r.URL.Query().Has(ui.DirParam) {
			// hledá se přímo dle adresáře
			dirQuery := r.URL.Query().Get(ui.DirParam)
			path := strings.Trim(dirQuery, "/")
			if len(path) > 0 {
				parts := strings.Split(path, "/")
				io.WriteString(w, ui.ConstructPage(indexer.FindByPath(parts), false, dirQuery))
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
	http.HandleFunc("/pause", func(w http.ResponseWriter, r *http.Request) { ret(w, myVLC.Pause()) })
	http.HandleFunc("/next", func(w http.ResponseWriter, r *http.Request) { ret(w, myVLC.Next()) })
	http.HandleFunc("/prev", func(w http.ResponseWriter, r *http.Request) { ret(w, myVLC.Previous()) })
	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) { ret(w, myVLC.Stop()) })
	// chyba v názvu funkce, ve skutečnosti volá Random (shuffle)
	http.HandleFunc("/shuffle", func(w http.ResponseWriter, r *http.Request) { ret(w, myVLC.ToggleLoop()) })
	http.HandleFunc("/loop", func(w http.ResponseWriter, r *http.Request) { ret(w, myVLC.ToggleRepeat()) })
	http.HandleFunc("/volume", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Has(ui.ValueParam) {
			val, err := strconv.Atoi(query.Get(ui.ValueParam))
			if err != nil {
				return
			}
			ret(w, myVLC.Volume(strconv.Itoa(val)))
		}
	})
	http.HandleFunc("/progress", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Has(ui.ValueParam) {
			val, err := strconv.Atoi(query.Get(ui.ValueParam))
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
		render := ui.ConstructPlaylist(items)
		io.WriteString(w, utils.ConstructPlaylistJSON(render, hash))
	})

	fileOrSearch := func(w http.ResponseWriter, r *http.Request, idOperation func(w http.ResponseWriter, r *http.Request), searchOperation func(w http.ResponseWriter, r *http.Request)) {
		if r.URL.Query().Has(ui.IdParam) {
			idOperation(w, r)
		} else if r.URL.Query().Has(ui.SearchParam) {
			searchOperation(w, r)
		}
	}

	getParam := func(r *http.Request, paramName string) string {
		query := r.URL.Query()
		return query.Get(paramName)
	}
	getIdParam := func(r *http.Request) string { return getParam(r, ui.IdParam) }

	addToPlaylist := func(w http.ResponseWriter, r *http.Request) { ret(w, myVLC.Add(prepURL(getIdParam(r)))) }
	addToPlaylistPlay := func(w http.ResponseWriter, r *http.Request) { ret(w, myVLC.AddStart(prepURL(getIdParam(r)))) }

	searchItems := func(w http.ResponseWriter, r *http.Request, consumer func(string) error) {
		searchQuery := getParam(r, ui.SearchParam)
		items := indexer.FindByString(searchQuery)
		for _, item := range items {
			ret(w, consumer(prepURL(item.GetPath())))
		}
	}
	searchAddToPlaylist := func(w http.ResponseWriter, r *http.Request) {
		searchItems(w, r, myVLC.Add)
	}
	searchAddToPlaylistPlay := func(w http.ResponseWriter, r *http.Request) {
		searchItems(w, r, func(t string) error { return myVLC.AddStart(t) })
	}

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

	http.HandleFunc(ui.AddEndpoint, func(w http.ResponseWriter, r *http.Request) { fileOrSearch(w, r, addToPlaylist, searchAddToPlaylist) })
	http.HandleFunc(ui.AddAndPlayEndpoint, func(w http.ResponseWriter, r *http.Request) {
		fileOrSearch(w, r, addToPlaylistPlay, searchAddToPlaylistPlay)
	})
	http.HandleFunc(ui.PlayEndpoint, func(w http.ResponseWriter, r *http.Request) { playFromPlaylist(w, r) })
	http.HandleFunc(ui.RemoveEndpoint, func(w http.ResponseWriter, r *http.Request) { removeFromPlaylist(w, r) })

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

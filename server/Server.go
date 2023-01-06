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

	addOrAddAndPlay := func(w http.ResponseWriter, r *http.Request, andPlay bool) {
		query := r.URL.Query()
		if query.Has(ui.IdParam) {
			target := prepURL(query.Get(ui.IdParam))
			if andPlay {
				ret(w, myVLC.AddStart(target))
			} else {
				ret(w, myVLC.Add(target))
			}
		} else if query.Has(ui.SearchParam) {
			searchQuery := query.Get(ui.SearchParam)
			items := indexer.FindByString(searchQuery)
			for _, item := range items {
				target := prepURL(item.GetPath())
				if andPlay {
					ret(w, myVLC.AddStart(target))
				} else {
					ret(w, myVLC.Add(target))
				}
			}
		}
	}

	http.HandleFunc(ui.AddEndpoint, func(w http.ResponseWriter, r *http.Request) { addOrAddAndPlay(w, r, false) })
	http.HandleFunc(ui.AddAndPlayEndpoint, func(w http.ResponseWriter, r *http.Request) { addOrAddAndPlay(w, r, true) })

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

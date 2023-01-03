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
	"strings"
)

type Resources struct {
	Styles  string
	Scripts string
	Favicon string
	Icons   string
}

func StartServer(myVLC vlcctrl.VLC, indexer indexer.Indexer, resources Resources) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Has(ui.SearchParam) {
			// hledá se cokoliv dle názvu
			searchQuery := r.URL.Query().Get(ui.SearchParam)
			if searchQuery != "" {
				ui.ConstructPage(indexer.FindByString(searchQuery), w, r, true, searchQuery)
				return
			}
		} else if r.URL.Query().Has(ui.DirParam) {
			// hledá se přímo dle adresáře
			dirQuery := r.URL.Query().Get(ui.DirParam)
			path := strings.Trim(dirQuery, "/")
			if len(path) > 0 {
				parts := strings.Split(path, "/")
				ui.ConstructPage(indexer.FindByPath(parts), w, r, false, dirQuery)
				return
			}
		}

		// výchozí pohled
		ui.ConstructPage(indexer.GetAllItems(), w, r, false, "/")
	})

	http.HandleFunc("/resources/styles.css", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, resources.Styles) })
	http.HandleFunc("/resources/scripts.js", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, resources.Scripts) })
	http.HandleFunc("/resources/favicon.png", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, resources.Favicon) })
	http.HandleFunc("/resources/icons.png", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, resources.Icons) })

	ret := func(w http.ResponseWriter, err error) {
		if err != nil {
			operations := "showError," + err.Error()
			operations += ";removeClass,info-div,info;"
			operations += ";removeClass,info-div,info-div-show;"
			operations += ";removeClass,info-div,info-div-hide;"
			operations += ";addClass,info-div,error;"
			operations += ";addClass,info-div,info-div-show;"
			io.WriteString(w, operations)
		}
	}

	prepURL := func(path string) string {
		return "file:///" + utils.EncodeURL(indexer.GetPlayerRoot()+path)
	}

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		response, _ := myVLC.RequestMaker("/requests/status.json")
		operations := utils.ProcessOperations(response)
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
		ui.ConstructPage(indexer.GetAllItems(), w, r, false, "/")
	})

	http.HandleFunc("/quit", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Aplikace byla ukončena")
		os.Exit(0)
	})

	log.Fatal(http.ListenAndServe(":8888", nil))
}

package main

import (
	vlcctrl "github.com/CedArctic/go-vlc-ctrl"
	"grass-control-go/indexer"
	"grass-control-go/ui"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

import _ "embed"

// https://github.com/CedArctic/go-vlc-ctrl
// https://pkg.go.dev/github.com/adrg/libvlc-go/v2
// https://tutorialedge.net/golang/creating-simple-web-server-with-golang/
// https://www.online-toolz.com/tools/unicode-html-entities-convertor.php
// https://en.wikipedia.org/wiki/Media_control_symbols

//go:embed resources/styles.css
var styles string

//go:embed resources/favicon.png
var favicon string

const vlcPort = 8080
const vlcPass = "vlcgatt"
const playerRoot = "D:/Hudba"

func initIndexer(indexer *indexer.Indexer) {
	indexer.Init(vlcPort, vlcPass, playerRoot)
}

func main() {

	indexer := indexer.Indexer{}
	initIndexer(&indexer)

	// Declare a local VLC instance on port 8080 with password "password"
	myVLC, err := vlcctrl.NewVLC("127.0.0.1", vlcPort, vlcPass)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Has("search") {
			// hledá se cokoliv dle názvu
			searchQuery := r.URL.Query().Get("search")
			if searchQuery != "" {
				ui.ConstructPage(indexer.FindByString(searchQuery), w, true, searchQuery)
				return
			}
		} else if r.URL.Query().Has("dir") {
			// hledá se přímo dle adresáře
			dirQuery := r.URL.Query().Get("dir")
			path := strings.Trim(dirQuery, "/")
			if len(path) > 0 {
				parts := strings.Split(path, "/")
				ui.ConstructPage(indexer.FindByPath(parts), w, false, dirQuery)
				return
			}
		}

		// výchozí pohled
		ui.ConstructPage(indexer.GetAllItems(), w, false, "/")
	})

	http.HandleFunc("/resources/styles.css", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, styles) })
	http.HandleFunc("/resources/favicon.png", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, favicon) })

	prepareURLForVLC := func(value string) string {
		// https://go.dev/play/p/pOfrn-Wsq5
		url := &url.URL{Path: "file:///D:/Hudba/" + value}
		// URL předsadí před sebe './'
		encoded := url.String()
		if encoded[:2] == "./" {
			encoded = encoded[2:]
		}
		// VLC má vadu a nebere URL mezery jako '+', zvládá jen '%20'
		//encoded := url.QueryEscape("file:///D:/Hudba/" + value)
		return encoded
	}

	http.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) { myVLC.EmptyPlaylist() })
	http.HandleFunc("/pause", func(w http.ResponseWriter, r *http.Request) { myVLC.Pause() })
	http.HandleFunc("/next", func(w http.ResponseWriter, r *http.Request) { myVLC.Next() })
	http.HandleFunc("/prev", func(w http.ResponseWriter, r *http.Request) { myVLC.Previous() })
	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) { myVLC.Stop() })
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		myVLC.Add(prepareURLForVLC(r.URL.Query().Get("id")))
	})
	http.HandleFunc("/addAndPlay", func(w http.ResponseWriter, r *http.Request) {
		myVLC.AddStart(prepareURLForVLC(r.URL.Query().Get("id")))
	})
	http.HandleFunc("/reindex", func(w http.ResponseWriter, r *http.Request) {
		initIndexer(&indexer)
		ui.ConstructPage(indexer.GetAllItems(), w, false, "/")
	})
	http.HandleFunc("/quit", func(w http.ResponseWriter, r *http.Request) { os.Exit(0) })

	log.Fatal(http.ListenAndServe(":8888", nil))
}

package main

import (
	"encoding/json"
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

import _ "embed"

// https://github.com/CedArctic/go-vlc-ctrl
// https://pkg.go.dev/github.com/adrg/libvlc-go/v2
// https://tutorialedge.net/golang/creating-simple-web-server-with-golang/
// https://www.online-toolz.com/tools/unicode-html-entities-convertor.php
// https://en.wikipedia.org/wiki/Media_control_symbols

//go:embed resources/styles.css
var styles string

//go:embed resources/scripts.js
var scripts string

//go:embed resources/favicon.png
var favicon string

//go:embed resources/icons.png
var icons string

const vlcPort = 8080
const vlcPass = "vlcgatt"
const playerRoot = "D:/Hudba"

type vlcJson struct {
	Fullscreen int `json:"fullscreen"`
	Stats      struct {
		Inputbitrate        float64 `json:"inputbitrate"`
		Sentbytes           int     `json:"sentbytes"`
		Lostabuffers        int     `json:"lostabuffers"`
		Averagedemuxbitrate int     `json:"averagedemuxbitrate"`
		Readpackets         int     `json:"readpackets"`
		Demuxreadpackets    int     `json:"demuxreadpackets"`
		Lostpictures        int     `json:"lostpictures"`
		Displayedpictures   int     `json:"displayedpictures"`
		Sentpackets         int     `json:"sentpackets"`
		Demuxreadbytes      int     `json:"demuxreadbytes"`
		Demuxbitrate        float64 `json:"demuxbitrate"`
		Playedabuffers      int     `json:"playedabuffers"`
		Demuxdiscontinuity  int     `json:"demuxdiscontinuity"`
		Decodedaudio        int     `json:"decodedaudio"`
		Sendbitrate         int     `json:"sendbitrate"`
		Readbytes           int     `json:"readbytes"`
		Averageinputbitrate int     `json:"averageinputbitrate"`
		Demuxcorrupted      int     `json:"demuxcorrupted"`
		Decodedvideo        int     `json:"decodedvideo"`
	} `json:"stats"`
	Audiodelay   int  `json:"audiodelay"`
	Apiversion   int  `json:"apiversion"`
	Currentplid  int  `json:"currentplid"`
	Time         int  `json:"time"`
	Volume       int  `json:"volume"`
	Length       int  `json:"length"`
	Random       bool `json:"random"`
	Audiofilters struct {
		Filter0 string `json:"filter_0"`
	} `json:"audiofilters"`
	Rate         int `json:"rate"`
	Videoeffects struct {
		Hue        int `json:"hue"`
		Saturation int `json:"saturation"`
		Contrast   int `json:"contrast"`
		Brightness int `json:"brightness"`
		Gamma      int `json:"gamma"`
	} `json:"videoeffects"`
	State       string  `json:"state"`
	Loop        bool    `json:"loop"`
	Version     string  `json:"version"`
	Position    float64 `json:"position"`
	Information struct {
		Chapter  int           `json:"chapter"`
		Chapters []interface{} `json:"chapters"`
		Title    int           `json:"title"`
		Category struct {
			Meta struct {
				DISCID      string `json:"DISCID"`
				Date        string `json:"date"`
				ArtworkUrl  string `json:"artwork_url"`
				Artist      string `json:"artist"`
				Album       string `json:"album"`
				TrackNumber string `json:"track_number"`
				Filename    string `json:"filename"`
				Title       string `json:"title"`
				Genre       string `json:"genre"`
			} `json:"meta"`
			Stream0 struct {
				Bitrate       string `json:"Bitrate"`
				Codec         string `json:"Codec"`
				Channels      string `json:"Channels"`
				BitsPerSample string `json:"Bits_per_sample"`
				Type          string `json:"Type"`
				SampleRate    string `json:"Sample_rate"`
			} `json:"Stream 0"`
		} `json:"category"`
		Titles []interface{} `json:"titles"`
	} `json:"information"`
	Repeat        bool          `json:"repeat"`
	Subtitledelay int           `json:"subtitledelay"`
	Equalizer     []interface{} `json:"equalizer"`
}

func initIndexer(indexer *indexer.Indexer) {
	indexer.Init(vlcPort, vlcPass, playerRoot)
}

func ternar(condition func() bool, trueResult string, falseResult string) string {
	if condition() {
		return trueResult
	} else {
		return falseResult
	}
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
	http.HandleFunc("/resources/scripts.js", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, scripts) })
	http.HandleFunc("/resources/favicon.png", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, favicon) })
	http.HandleFunc("/resources/icons.png", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, icons) })

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
		return "file:///" + utils.EncodeURL(playerRoot+path)
	}

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		var result vlcJson
		response, _ := myVLC.RequestMaker("/requests/status.json")
		json.Unmarshal([]byte(response), &result)
		operations := ""
		// Shuffle
		operations += ternar(func() bool { return result.Random }, "addClass", "removeClass") + ",shuffle-btn,checked;"
		// Loop
		operations += ternar(func() bool { return result.Loop }, "addClass", "removeClass") + ",loop-btn,checked;"
		// Paused/Playing
		isPlaying := result.State == "playing"
		pauseClassOper := ternar(func() bool { return isPlaying }, "addClass", "removeClass")
		playClassOper := ternar(func() bool { return !isPlaying }, "addClass", "removeClass")
		operations += pauseClassOper + ",play-pause-btn,pause-btn;"
		operations += playClassOper + ",play-pause-btn,play-btn;"
		// Current song
		artist := result.Information.Category.Meta.Artist
		title := result.Information.Category.Meta.Title
		album := result.Information.Category.Meta.Album
		filename := result.Information.Category.Meta.Filename
		if artist != "" {
			operations += "songInfo," + artist
			if title != "" {
				operations += ": " + title
			}
			if album != "" {
				operations += " (" + album + ")"
			}
		} else if filename != "" {
			operations += "songInfo," + filename
		}

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
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		ret(w, myVLC.Add(prepURL(r.URL.Query().Get("id"))))
	})
	http.HandleFunc("/addAndPlay", func(w http.ResponseWriter, r *http.Request) {
		ret(w, myVLC.AddStart(prepURL(r.URL.Query().Get("id"))))
	})
	http.HandleFunc("/reindex", func(w http.ResponseWriter, r *http.Request) {
		initIndexer(&indexer)
		ui.ConstructPage(indexer.GetAllItems(), w, false, "/")
	})
	http.HandleFunc("/quit", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Aplikace byla ukončena")
		os.Exit(0)
	})

	log.Fatal(http.ListenAndServe(":8888", nil))
}

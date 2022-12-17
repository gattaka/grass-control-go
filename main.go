package main

import (
	vlcctrl "github.com/CedArctic/go-vlc-ctrl"
	"grass-control-go/elements"
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

const vlcPort = 8080
const vlcPass = "vlcgatt"
const playerRoot = "D:/Hudba"

//go:embed resources/styles.css
var styles string

type item struct {
	isDir  bool
	name   string
	path   string
	items  []*item
	parent *item
}

func find(value string, root *item, items *[]*item) {
	value = strings.ToLower(value)
	for _, itm := range root.items {
		if strings.Contains(strings.ToLower(itm.name), value) {
			*items = append(*items, itm)
		}
		if itm.isDir {
			find(value, itm, items)
		}
	}
}

func indexDir(parent *item) {
	files, err := os.ReadDir(parent.path)
	if err != nil {
		log.Fatal(err)
	}

	if len(files) > 0 {
		parent.items = make([]*item, len(files))
		for i, file := range files {
			child := item{name: file.Name(), path: parent.path + "/" + file.Name(), isDir: file.IsDir(), parent: parent}
			parent.items[i] = &child
			if file.IsDir() {
				indexDir(&child)
			}
		}
	}
}

func main() {

	// Indexace
	root := item{name: "Hudba", path: playerRoot, isDir: true}
	indexDir(&root)

	// Declare a local VLC instance on port 8080 with password "password"
	myVLC, err := vlcctrl.NewVLC("127.0.0.1", vlcPort, vlcPass)
	if err != nil {
		log.Fatal(err)
	}

	ajax := func(url string) string {
		return `const xhttp = new XMLHttpRequest();
					xhttp.open('GET', '` + url + `', true);
					xhttp.send()`
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := elements.Html{}
		html.CssFiles = []string{"resources/styles.css"}

		header := elements.Header{Level: 1}
		header.Add(&elements.Text{Value: "GrassControl"})
		html.Add(&header)

		searchInput := elements.Input{JSfunc: "window.location.href = '?search=' + this.value;"}
		html.Add(&searchInput)

		controlsDiv := elements.Div{}
		controlsDiv.AddClass("controls-div")
		html.Add(&controlsDiv)

		controlsDiv.Add(&elements.Button{Value: "&#10006", JSfunc: ajax("clear")})
		controlsDiv.Add(&elements.Button{Value: "&#9205;", JSfunc: ajax("play")})
		controlsDiv.Add(&elements.Button{Value: "&#9208;", JSfunc: ajax("pause")})
		controlsDiv.Add(&elements.Button{Value: "&#9198;", JSfunc: ajax("prev")})
		controlsDiv.Add(&elements.Button{Value: "&#9209;", JSfunc: ajax("stop")})
		controlsDiv.Add(&elements.Button{Value: "&#9197;", JSfunc: ajax("next")})

		table := elements.Table[item]{}

		searchQuery := r.URL.Query().Get("search")
		if searchQuery != "" {
			find(searchQuery, &root, &table.Items)
		} else {
			table.Items = root.items
		}

		table.Columns = make([]elements.TableColumn[item], 2)
		table.Columns[0] = elements.TableColumn[item]{Name: "Název", Renderer: func(itm item) string {
			btn := elements.Button{Value: itm.name, JSfunc: ajax("add?id=" + itm.path)}
			return btn.Render()
		}}
		table.Columns[1] = elements.TableColumn[item]{Name: "Nadřazený adresář", Renderer: func(itm item) string {
			if itm.parent == nil {
				return ""
			}
			btn := elements.Button{Value: itm.parent.name, JSfunc: ajax("add?id=" + itm.parent.path)}
			return btn.Render()
		}}
		html.Add(&table)

		result := html.Render()
		io.WriteString(w, result)
	})

	http.HandleFunc("/resources/styles.css", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, styles)
	})

	http.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) {
		myVLC.EmptyPlaylist()
	})

	http.HandleFunc("/next", func(w http.ResponseWriter, r *http.Request) {
		myVLC.Next()
	})

	http.HandleFunc("/prev", func(w http.ResponseWriter, r *http.Request) {
		myVLC.Previous()
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		myVLC.Stop()
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		myVLC.Add(url.QueryEscape("file:///D:/Hudba/" + id))
	})

	http.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		err = myVLC.Play()
		if err != nil {
			log.Fatal(err)
		}
	})

	log.Fatal(http.ListenAndServe(":8888", nil))

}

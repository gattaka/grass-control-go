package main

import (
	"fmt"
	vlcctrl "github.com/CedArctic/go-vlc-ctrl"
	"grass-control-go/elements"
	"log"
	"net/http"
	"net/url"
)

// https://github.com/CedArctic/go-vlc-ctrl
// https://pkg.go.dev/github.com/adrg/libvlc-go/v2

const vlcPort = 8080
const vlcPass = "vlcgatt"

func main() {

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
		div := elements.Div{}

		header := elements.Header{Level: 1}
		header.Add(&elements.Text{Value: "GrassControl"})
		div.Add(&header)

		ids := []string{"JINE/Godfather.mp3", "JINE/Mafia - Main Theme.mp3"}
		for _, id := range ids {
			div.Add(&elements.Button{Value: id, JSfunc: ajax("add?id=" + id)})
			div.Add(&elements.Text{Value: "<br/>"})
		}

		result := div.Render()
		fmt.Fprintf(w, result)
	})

	http.HandleFunc("/next", func(w http.ResponseWriter, r *http.Request) {
		myVLC.Next()
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

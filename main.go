package main

import (
	vlcctrl "github.com/CedArctic/go-vlc-ctrl"
	"grass-control-go/indexer"
	"grass-control-go/server"
	"log"
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

func main() {

	indexer := indexer.Indexer{}
	indexer.Init(playerRoot)

	// Declare a local VLC instance on port 8080 with password "password"
	myVLC, err := vlcctrl.NewVLC("127.0.0.1", vlcPort, vlcPass)
	if err != nil {
		log.Fatal(err)
	}

	resources := server.Resources{
		Styles:  styles,
		Scripts: scripts,
		Favicon: favicon,
		Icons:   icons,
	}
	server.StartServer(myVLC, indexer, resources)
}

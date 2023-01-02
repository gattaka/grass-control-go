package main

import (
	vlcctrl "github.com/CedArctic/go-vlc-ctrl"
	"grass-control-go/indexer"
	"grass-control-go/server"
	"log"
	"os"
	"strconv"
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

func main() {

	var vlcPortConvErr error
	var vlcPort int
	var vlcPass string
	var playerRoot string

	if len(os.Args) != 4 {
		log.Fatal("Nebyly poskytnuty povinné argumenty -vlcPort=<port>, -vlcPass=<password> a -playerRoot=<playerroot>")
	}

	for i, val := range os.Args {
		if i == 0 {
			// první argument je název programu
			continue
		}
		keyVal := strings.Split(val, "=")
		if len(keyVal) == 2 {
			if strings.ToLower(keyVal[0]) == "-vlcport" {
				vlcPort, vlcPortConvErr = strconv.Atoi(keyVal[1])
			} else if strings.ToLower(keyVal[0]) == "-vlcpass" {
				vlcPass = keyVal[1]
			} else if strings.ToLower(keyVal[0]) == "-playerroot" {
				playerRoot = keyVal[1]
			}
		}
	}

	if vlcPortConvErr != nil || vlcPort < 1 {
		log.Fatal("VLC port musí být celé pozitivní číslo")
	}
	if len(vlcPass) == 0 {
		log.Fatal("VLC password nesmí být prázdný")
	}
	if len(playerRoot) == 0 {
		log.Fatal("Player root nesmí být prázdný")
	}

	indexer := indexer.Indexer{}
	indexer.Init(playerRoot)

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

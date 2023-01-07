package utils

import (
	"encoding/json"
	"hash"
	"hash/fnv"
	"log"
	"strconv"
)

type VlcPlaylistNode struct {
	Ro       string            `json:"ro"`
	Type     string            `json:"type"`
	Name     string            `json:"name"`
	Id       string            `json:"id"`
	Duration int               `json:"duration"`
	Uri      string            `json:"uri"`
	Current  string            `json:"current,omitempty"`
	Subnodes []VlcPlaylistNode `json:"children"`
}

type GrassControlPlaylistJson struct {
	Html string `json:"html"`
	Hash uint32 `json:"hash"`
}

func processPlaylist(node VlcPlaylistNode, items *[]*VlcPlaylistNode, hash *hash.Hash32) {
	if node.Type == "node" {
		for _, subnode := range node.Subnodes {
			processPlaylist(subnode, items, hash)
		}
	} else {
		(*hash).Write([]byte(node.Uri))
		// do hash započítávám i délku média, protože VLC z počátku neví a vypíše -1, až dodatečně čas dopočítá,
		// takže i já musím dodatečně aktualizovat seznam
		(*hash).Write([]byte(strconv.Itoa(node.Duration)))
		*items = append(*items, &node)
	}
}

func ParsePlaylist(vlcPlaylistJSON string) (*[]*VlcPlaylistNode, uint32) {
	var vlcPlaylist VlcPlaylistNode
	json.Unmarshal([]byte(vlcPlaylistJSON), &vlcPlaylist)
	var items []*VlcPlaylistNode

	// Hash se používá k indikaci změny playlistu -- aby se neaktualizoval i když není potřeba a neproblikával
	hash := fnv.New32a()
	processPlaylist(vlcPlaylist, &items, &hash)
	return &items, hash.Sum32()
}

func ConstructPlaylistJSON(payload string, hash uint32) string {
	result := GrassControlPlaylistJson{
		Html: payload,
		Hash: hash,
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

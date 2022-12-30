package indexer

import (
	"log"
	"os"
	"strings"
)

type Indexer struct {
	playerRoot string
	root       *Item
}

func (idx *Indexer) Init(playerRoot string) {
	idx.playerRoot = playerRoot
	idx.Reindex()
}

func (idx Indexer) GetPlayerRoot() string {
	return idx.playerRoot
}

func (idx *Indexer) Reindex() {
	idx.root = &Item{name: "Hudba", path: "", isDir: true}
	idx.indexDir(idx.root)
}

func (idx Indexer) GetAllItems() []*Item {
	return idx.root.items
}

func (idx *Indexer) indexDir(parent *Item) {
	files, err := os.ReadDir(idx.playerRoot + parent.path)
	if err != nil {
		log.Fatal(err)
	}

	if len(files) > 0 {
		parent.items = make([]*Item, len(files))
		for i, file := range files {
			dirPath := parent.path + "/" + file.Name()
			child := Item{name: file.Name(), path: dirPath, isDir: file.IsDir(), parent: parent}
			parent.items[i] = &child
			if file.IsDir() {
				idx.indexDir(&child)
			}
		}
	}
}

func (idx Indexer) FindByString(value string) []*Item {
	items := make([]*Item, 0)
	idx.findByString(value, idx.root, &items)
	return items
}

func (idx Indexer) findByString(value string, root *Item, items *[]*Item) {
	value = strings.ToLower(value)
	for _, itm := range root.items {
		if strings.Contains(strings.ToLower(itm.name), value) {
			*items = append(*items, itm)
		}
		if itm.isDir {
			idx.findByString(value, itm, items)
		}
	}
}

func (idx Indexer) FindByPath(path []string) []*Item {
	items := make([]*Item, 0)
	idx.findByPath(path, idx.root, &items)
	return items
}

func (idx Indexer) findByPath(path []string, root *Item, items *[]*Item) {
	value := path[0]
	for _, itm := range root.items {
		if itm.name == value && itm.isDir {
			if len(path) == 1 {
				*items = append(*items, itm.items...)
			} else {
				path = path[1:]
				idx.findByPath(path, itm, items)
			}
			return
		}
	}
}

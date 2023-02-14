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

func (idx Indexer) FindByItem(path string) (bool, *Item) {
	chunks := strings.Split(path, "/")
	item := idx.root
	for i, chunk := range chunks {
		// pokud cesta ukazuje na další potomky, ale to, co jsem částečně našel už další potomky nemá, vrať že jsi nenašel nic
		if i < len(chunks)-1 && !item.isDir {
			return false, idx.root
		}
		for _, child := range item.items {
			if child.name == chunk {
				item = child
				break
			}
		}
	}
	return true, item
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

func (idx Indexer) FindByPath(path string) []*Item {
	path = strings.Trim(path, "/")
	if path == "" {
		return idx.root.items
	}
	parts := strings.Split(path, "/")
	items := make([]*Item, 0)
	idx.findByPath(parts, idx.root, &items)
	return items
}

func (idx Indexer) findByPath(path []string, root *Item, items *[]*Item) {
	value := path[0]
	for _, itm := range root.items {
		// nalezena část cesty
		if itm.name == value {
			if itm.isDir {
				if len(path) == 1 {
					// adresář je buď sám cílem
					*items = append(*items, itm.items...)
				} else {
					// nebo je cíl někde uvnitř
					path = path[1:]
					idx.findByPath(path, itm, items)
				}
			} else {
				// ne-adresář je vždy cílem (nelze procházet)
				*items = append(*items, itm)
			}
			return
		}
	}
}

func (idx *Indexer) ExpandByPath(path string) []*Item {
	items := make([]*Item, 0)
	for _, item := range idx.FindByPath(path) {
		expandByPathRec(item, &items)
	}
	return items
}

func (idx *Indexer) ExpandByItem(item *Item) []*Item {
	items := make([]*Item, 0)
	expandByPathRec(item, &items)
	return items
}

func expandByPathRec(item *Item, items *[]*Item) {
	if item.IsDir() {
		for _, child := range item.items {
			expandByPathRec(child, items)
		}
	} else {
		*items = append(*items, item)
	}
}

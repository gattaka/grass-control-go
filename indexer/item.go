package indexer

type Item struct {
	isDir  bool
	name   string
	path   string
	items  []*Item
	parent *Item
}

func (i Item) GetName() string { return i.name }
func (i Item) GetPath() string { return i.path }
func (i Item) IsDir() bool     { return i.isDir }
func (i Item) HasParent() bool { return i.parent != nil }
func (i Item) GetParent() Item { return *i.parent }

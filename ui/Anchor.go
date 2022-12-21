package ui

type Anchor struct {
	Value   string
	Link    string
	classes []string
}

func (d *Anchor) Render() string        { return "<a href=\"" + d.Link + "\">" + d.Value + "</a>" }
func (d *Anchor) getClasses() *[]string { return &d.classes }
func (d *Anchor) AddClass(s string)     { addClass(d, s) }

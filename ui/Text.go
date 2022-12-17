package ui

type Text struct {
	Value   string
	classes []string
}

func (d *Text) Render() string        { return d.Value }
func (d *Text) getClasses() *[]string { return &d.classes }
func (d *Text) AddClass(s string)     { addClass(d, s) }

package ui

import "strconv"

type Header struct {
	Level    int
	elements []*Element
	classes  []string
}

func (d *Header) getElements() *[]*Element { return &d.elements }
func (d *Header) getClasses() *[]string    { return &d.classes }
func (d *Header) Add(e Element)            { add(d, e) }
func (d *Header) getTag() string           { return "h" + strconv.Itoa(d.Level) }
func (d *Header) Render() string           { return render(d) }
func (d *Header) AddClass(s string)        { addClass(d, s) }

package ui

type Div struct {
	elements []*Element
	classes  []string
}

func (d *Div) getElements() *[]*Element { return &d.elements }
func (d *Div) getClasses() *[]string    { return &d.classes }
func (d *Div) Add(e Element)            { add(d, e) }
func (d *Div) getTag() string           { return "div" }
func (d *Div) Render() string           { return render(d) }
func (d *Div) AddClass(s string)        { addClass(d, s) }

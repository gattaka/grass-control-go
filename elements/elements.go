package elements

import (
	"strconv"
)

type Element interface {
	Render() string
	Add(element Element)
	getElements() []Element
	setElements([]Element)
	getTag() string
}

func add(e Element, s Element) {
	if e.getElements() == nil {
		e.setElements(make([]Element, 0))
	}
	e.setElements(append(e.getElements(), s))
}

func render(e Element) string {
	result := ""
	result += "<" + e.getTag() + ">"
	if e.getElements() != nil {
		for _, s := range e.getElements() {
			result += s.Render()
		}
	}
	result += "</" + e.getTag() + ">"
	return result
}

// Text
type Button struct {
	Value  string
	JSfunc string
}

func (d *Button) getElements() []Element         { return nil }
func (d *Button) setElements(elements []Element) {}
func (d *Button) Add(e Element)                  {}
func (d *Button) getTag() string                 { return "" }
func (d *Button) Render() string {
	return "<input type=\"button\" value=\"" + d.Value + "\" onclick=\"" + d.JSfunc + "\">"
}

// Text
type Text struct {
	Value string
}

func (d *Text) getElements() []Element         { return nil }
func (d *Text) setElements(elements []Element) {}
func (d *Text) Add(e Element)                  {}
func (d *Text) getTag() string                 { return "" }
func (d *Text) Render() string                 { return d.Value }

// Div
type Div struct {
	elements []Element
}

func (d *Div) getElements() []Element      { return d.elements }
func (d *Div) setElements(items []Element) { d.elements = items }
func (d *Div) Add(e Element)               { add(d, e) }
func (d *Div) getTag() string              { return "div" }
func (d *Div) Render() string              { return render(d) }

// Headers
type Header struct {
	Level    int
	elements []Element
}

func (d *Header) getElements() []Element         { return d.elements }
func (d *Header) setElements(elements []Element) { d.elements = elements }
func (d *Header) Add(e Element)                  { add(d, e) }
func (d *Header) getTag() string                 { return "h" + strconv.Itoa(d.Level) }
func (d *Header) Render() string                 { return render(d) }

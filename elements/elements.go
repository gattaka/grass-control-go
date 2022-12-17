package elements

import (
	"strconv"
)

type Element interface {
	Render() string
}

type HasElements interface {
	Element
	getElements() []Element
	setElements([]Element)
	getTag() string
	Add(element Element)
}

func add(e HasElements, s Element) {
	if e.getElements() == nil {
		e.setElements(make([]Element, 0))
	}
	e.setElements(append(e.getElements(), s))
}

func render(e HasElements) string {
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

// Button
type Button struct {
	Value  string
	JSfunc string
}

func (d *Button) Render() string {
	return "<input type=\"button\" value=\"" + d.Value + "\" onclick=\"" + d.JSfunc + "\">"
}

// Text
type Text struct {
	Value string
}

func (d *Text) Render() string { return d.Value }

// Html
type Html struct {
	elements []Element
	Css      string
	CssFiles []string
}

func (d *Html) getElements() []Element      { return d.elements }
func (d *Html) setElements(items []Element) { d.elements = items }
func (d *Html) Add(e Element)               { add(d, e) }
func (d *Html) getTag() string              { return "html" }
func (d *Html) Render() string {
	result := "<html><head>"
	if len(d.CssFiles) > 0 {
		for _, css := range d.CssFiles {
			result += "<link rel=\"stylesheet\" href=\"" + css + "\">"
		}
	}
	result += "<style>"
	result += d.Css
	result += "</style>"
	result += "</head><body>"
	if d.elements != nil {
		for _, s := range d.getElements() {
			result += s.Render()
		}
	}
	result += "</body></html>"
	return result
}

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

// Table
type TableColumn[T any] struct {
	Name     string
	Renderer func(item T) string
}

type Table[T any] struct {
	Items   []*T
	Columns []TableColumn[T]
}

func (d *Table[T]) Render() string {
	result := "<div class='table-div'><div class='table-head-div'><div class='table-head-tr-div'>"
	for _, column := range d.Columns {
		result += "<div class='table-head-td-div'>" + column.Name + "</div>"
	}
	result += "</div></div><div class='table-body-div'>"
	for _, item := range d.Items {
		result += "<div class='table-body-tr-div'>"
		for _, column := range d.Columns {
			result += "<div class='table-body-td-div'>" + column.Renderer(*item) + "</div>"
		}
		result += "</div>"
	}
	result += "</div></div>"
	return result
}

// Input
type Input struct {
	JSfunc string
}

func (d *Input) Render() string {
	return "<input type=\"text\" onchange=\"" + d.JSfunc + "\">"
}

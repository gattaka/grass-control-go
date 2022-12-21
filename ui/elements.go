package ui

import "strings"

type Element interface {
	Render() string
	getClasses() *[]string
	AddClass(s string)
}

type HasElements interface {
	Element
	getTag() string
	getElements() *[]*Element
	Add(element Element)
}

func add(e HasElements, s Element) {
	children := e.getElements()
	if children == nil {
		*children = make([]*Element, 0)
	}
	*children = append(*children, &s)
}

func addClass(e Element, s string) {
	classes := e.getClasses()
	if classes == nil {
		*classes = make([]string, 0)
	}
	*classes = append(*classes, s)
}

func render(e HasElements) string {
	result := ""
	result += "<" + e.getTag()
	if e.getClasses() != nil && len(*e.getClasses()) > 0 {
		result += " class=\"" + strings.Join(*e.getClasses(), " ") + "\" "
	}
	result += ">"
	if e.getElements() != nil {
		for _, s := range *e.getElements() {
			result += (*s).Render()
		}
	}
	result += "</" + e.getTag() + ">"
	return result
}
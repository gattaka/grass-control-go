package ui

import "strings"

type Input struct {
	JSfunc  string
	classes []string
}

func (d *Input) Render() string {
	result := "<input type=\"text\" "
	if d.getClasses() != nil {
		result += "class=\"" + strings.Join(*d.getClasses(), " ") + "\" "
	}
	return result + "onchange=\"" + d.JSfunc + "\">"
}
func (d *Input) getClasses() *[]string { return &d.classes }
func (d *Input) AddClass(s string)     { addClass(d, s) }

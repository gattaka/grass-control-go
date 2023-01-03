package ui

import "strconv"

type RangeInput struct {
	cmn common
}

func (i *RangeInput) AddClass(s string) Element { i.cmn.addClass(s); return i }
func (i *RangeInput) Add(e Element)             { i.cmn.addElement(e) }
func (i *RangeInput) SetId(s string) Element    { i.cmn.setId(s); return i }
func (i *RangeInput) GetId() string             { return i.cmn.getId() }
func (i *RangeInput) SetAttribute(name string, value string) Element {
	i.cmn.setAttribute(name, value)
	return i
}
func (i *RangeInput) SetOnChange(value string) { i.cmn.setAttribute("onchange", value) }
func (i *RangeInput) SetValue(value string)    { i.cmn.setAttribute("value", value) }
func (i *RangeInput) SetMin(value int)         { i.cmn.setAttribute("min", strconv.Itoa(value)) }
func (i *RangeInput) SetMax(value int)         { i.cmn.setAttribute("max", strconv.Itoa(value)) }
func (i *RangeInput) SetName(value string)     { i.cmn.setAttribute("name", value) }
func (i *RangeInput) Render() string {
	i.cmn.setAttribute("type", "range")
	return i.cmn.render("input")
}

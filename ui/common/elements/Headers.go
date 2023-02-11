package elements

import "strconv"

type Header struct {
	cmn   common
	Level int
}

func (h *Header) AddClass(s string) Element { h.cmn.addClass(s); return h }
func (h *Header) Add(e Element)             { h.cmn.addElement(e) }
func (h *Header) SetId(s string) Element    { h.cmn.setId(s); return h }
func (h *Header) GetId() string             { return h.cmn.getId() }
func (h *Header) SetAttribute(name string, value string) Element {
	h.cmn.setAttribute(name, value)
	return h
}
func (h *Header) SetOnClick(value string) { h.cmn.setAttribute("onclick", value) }
func (h *Header) SetValue(value string)   { h.cmn.setAttribute("value", value) }
func (h *Header) Render() string          { return h.cmn.render("h" + strconv.Itoa(h.Level)) }

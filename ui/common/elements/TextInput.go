package elements

type TextInput struct {
	cmn common
}

func (i *TextInput) AddClass(s string) Element { i.cmn.addClass(s); return i }
func (i *TextInput) Add(e Element)             { i.cmn.addElement(e) }
func (i *TextInput) SetId(s string) Element    { i.cmn.setId(s); return i }
func (i *TextInput) GetId() string             { return i.cmn.getId() }
func (i *TextInput) SetAttribute(name string, value string) Element {
	i.cmn.setAttribute(name, value)
	return i
}
func (i *TextInput) SetOnChange(value string) { i.cmn.setAttribute("onchange", value) }
func (i *TextInput) SetValue(value string)    { i.cmn.setAttribute("value", value) }
func (i *TextInput) SetName(value string)     { i.cmn.setAttribute("name", value) }
func (i *TextInput) Render() string {
	i.cmn.setAttribute("type", "text")
	return i.cmn.render("input")
}

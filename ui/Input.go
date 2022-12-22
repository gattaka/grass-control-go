package ui

type Input struct {
	cmn common
}

func (i *Input) AddClass(s string) Element              { i.cmn.addClass(s); return i }
func (i *Input) Add(e Element)                          { i.cmn.addElement(e) }
func (i *Input) SetId(s string)                         { i.cmn.id = s }
func (i *Input) GetId() string                          { return i.cmn.id }
func (i *Input) SetAttribute(name string, value string) { i.cmn.setAttribute(name, value) }
func (i *Input) SetOnChange(value string)               { i.cmn.setAttribute("onchange", value) }
func (i *Input) SetValue(value string)                  { i.cmn.setAttribute("value", value) }
func (i *Input) Render() string {
	i.cmn.setAttribute("type", "text")
	return i.cmn.render("input")
}

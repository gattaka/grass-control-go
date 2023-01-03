package ui

type Form struct {
	cmn common
}

func (f *Form) AddClass(s string) Element { f.cmn.addClass(s); return f }
func (f *Form) Add(e Element)             { f.cmn.addElement(e) }
func (f *Form) SetId(s string) Element    { f.cmn.setId(s); return f }
func (f *Form) GetId() string             { return f.cmn.getId() }
func (f *Form) SetAttribute(name string, value string) Element {
	f.cmn.setAttribute(name, value)
	return f
}
func (f *Form) SetMethod(value string) { f.cmn.setAttribute("method", value) }
func (f *Form) SetAction(value string) { f.cmn.setAttribute("action", value) }
func (f *Form) Render() string {
	return f.cmn.render("form")
}

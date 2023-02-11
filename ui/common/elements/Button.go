package elements

type Button struct {
	cmn common
}

func NewButton(value string, onClick string) *Button {
	b := Button{}
	b.SetValue(value)
	b.SetOnClick(onClick)
	return &b
}

func (b *Button) AddClass(s string) Element { b.cmn.addClass(s); return b }
func (b *Button) SetId(s string) Element    { b.cmn.setId(s); return b }
func (b *Button) GetId() string             { return b.cmn.getId() }
func (b *Button) SetAttribute(name string, value string) Element {
	b.cmn.setAttribute(name, value)
	return b
}
func (b *Button) SetOnClick(value string) { b.cmn.setAttribute("onclick", value) }
func (b *Button) SetValue(value string)   { b.cmn.setAttribute("value", value) }
func (b *Button) Render() string {
	b.cmn.setAttribute("type", "button")
	return b.cmn.render("input")
}

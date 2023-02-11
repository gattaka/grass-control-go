package elements

type Anchor struct {
	cmn common
}

func NewAnchor(value string, link string) *Anchor {
	a := Anchor{}
	a.SetValue(value)
	a.SetLink(link)
	return &a
}

func NewAnchorJS(value string, onClickJS string) *Anchor {
	a := Anchor{}
	a.SetValue(value)
	a.SetOnClick(onClickJS)
	return &a
}

func (a *Anchor) AddClass(s string) Element { a.cmn.addClass(s); return a }
func (a *Anchor) SetId(s string) Element    { a.cmn.setId(s); return a }
func (a *Anchor) GetId() string             { return a.cmn.getId() }
func (a *Anchor) Add(e Element)             { a.cmn.addElement(e) }
func (a *Anchor) SetAttribute(name string, value string) Element {
	a.cmn.setAttribute(name, value)
	return a
}
func (a *Anchor) SetOnClick(value string) { a.cmn.setAttribute("onclick", value) }
func (a *Anchor) SetValue(value string)   { a.cmn.content = value }
func (a *Anchor) SetLink(value string)    { a.cmn.setAttribute("href", value) }
func (a *Anchor) Render() string          { return a.cmn.render("a") }

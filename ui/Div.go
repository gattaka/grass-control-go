package ui

type Div struct {
	cmn common
}

func (d *Div) AddClass(s string) Element { d.cmn.addClass(s); return d }
func (d *Div) Add(e Element)             { d.cmn.addElement(e) }
func (d *Div) SetId(s string) Element    { d.cmn.setId(s); return d }
func (d *Div) GetId() string             { return d.cmn.getId() }
func (d *Div) SetAttribute(name string, value string) Element {
	d.cmn.setAttribute(name, value)
	return d
}
func (d *Div) SetOnClick(value string) { d.cmn.setAttribute("onclick", value) }
func (d *Div) SetValue(value string)   { d.cmn.content = value }
func (d *Div) Render() string          { return d.cmn.render("div") }

package ui

type Span struct {
	cmn common
}

func NewSpan(value string) *Span {
	s := Span{}
	s.SetValue(value)
	return &s
}

func (s *Span) AddClass(c string) Element              { s.cmn.addClass(c); return s }
func (s *Span) Add(e Element)                          { s.cmn.addElement(e) }
func (s *Span) SetId(id string)                        { s.cmn.id = id }
func (s *Span) GetId() string                          { return s.cmn.id }
func (s *Span) SetAttribute(name string, value string) { s.cmn.setAttribute(name, value) }
func (s *Span) SetOnClick(value string)                { s.cmn.setAttribute("onclick", value) }
func (s *Span) SetValue(value string)                  { s.cmn.content = value }
func (s *Span) Render() string                         { return s.cmn.render("span") }

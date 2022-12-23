package ui

type Script struct {
	cmn common
}

func NewScript(value string) *Script {
	script := Script{}
	script.SetValue(value)
	return &script
}

func (s *Script) AddClass(st string) Element { s.cmn.addClass(st); return s }
func (s *Script) SetId(st string) Element    { s.cmn.setId(st); return s }
func (s *Script) GetId() string              { return s.cmn.getId() }
func (s *Script) SetAttribute(name string, value string) Element {
	s.cmn.setAttribute(name, value)
	return s
}
func (s *Script) SetValue(value string) { s.cmn.content = value }
func (s *Script) Render() string        { return s.cmn.render("script") }

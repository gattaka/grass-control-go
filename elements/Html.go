package elements

type Html struct {
	elements []*Element
	Css      string
	CssFiles []string
	classes  []string
}

func (d *Html) getElements() *[]*Element { return &d.elements }
func (d *Html) getClasses() *[]string    { return &d.classes }
func (d *Html) Add(e Element)            { add(d, e) }
func (d *Html) getTag() string           { return "html" }
func (d *Html) Render() string {
	result := "<html><head>"
	if len(d.CssFiles) > 0 {
		for _, css := range d.CssFiles {
			result += "<link rel=\"stylesheet\" href=\"" + css + "\">"
		}
	}
	result += "<style>"
	result += d.Css
	result += "</style>"
	result += "</head><body>"
	if d.elements != nil {
		for _, s := range *d.getElements() {
			result += (*s).Render()
		}
	}
	result += "</body></html>"
	return result
}

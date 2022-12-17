package ui

type Html struct {
	elements []*Element
	Css      string
	Headers  []string
	classes  []string
}

func (d *Html) getElements() *[]*Element { return &d.elements }
func (d *Html) getClasses() *[]string    { return &d.classes }
func (d *Html) Add(e Element)            { add(d, e) }
func (d *Html) getTag() string           { return "html" }
func (d *Html) Render() string {
	result := "<html><head>"
	if len(d.Headers) > 0 {
		for _, header := range d.Headers {
			result += header
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
func (d *Html) AddClass(s string) { addClass(d, s) }

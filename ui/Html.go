package ui

type Html struct {
	html    common
	body    common
	Css     string
	Headers []string
}

func (h *Html) Add(e Element)                          { h.html.addElement(e) }
func (h *Html) SetAttribute(name string, value string) { h.html.setAttribute(name, value) }
func (h *Html) Render() string {
	result := "<head>"
	if len(h.Headers) > 0 {
		for _, header := range h.Headers {
			result += header
		}
	}
	if h.Css != "" {
		result += "<style>"
		result += h.Css
		result += "</style>"
	}
	result += "</head>"
	h.html.content = result
	h.html.content += h.body.render("body")
	return h.html.render("html")
}

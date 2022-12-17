package elements

type Button struct {
	Value   string
	JSfunc  string
	classes []string
}

func (d *Button) Render() string {
	return "<input type=\"button\" value=\"" + d.Value + "\" onclick=\"" + d.JSfunc + "\">"
}
func (d *Button) getClasses() *[]string { return &d.classes }

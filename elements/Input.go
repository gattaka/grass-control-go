package elements

type Input struct {
	JSfunc  string
	classes []string
}

func (d *Input) Render() string {
	return "<input type=\"text\" onchange=\"" + d.JSfunc + "\">"
}

func (d *Input) getClasses() *[]string { return &d.classes }

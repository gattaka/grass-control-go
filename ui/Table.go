package ui

type TableColumn[T any] struct {
	Name     string
	Renderer func(item T) string
}

type Table[T any] struct {
	Items   []*T
	Columns []TableColumn[T]
	classes []string
}

func (d *Table[T]) Render() string {
	result := "<div class='table-div'><div class='table-head-div'><div class='table-head-tr-div'>"
	for _, column := range d.Columns {
		result += "<div class='table-head-td-div'>" + column.Name + "</div>"
	}
	result += "</div></div><div class='table-body-div'>"
	for _, item := range d.Items {
		result += "<div class='table-body-tr-div'>"
		for _, column := range d.Columns {
			result += "<div class='table-body-td-div'>" + column.Renderer(*item) + "</div>"
		}
		result += "</div>"
	}
	result += "</div></div>"
	return result
}
func (d *Table[T]) getClasses() *[]string { return &d.classes }
func (d *Table[T]) AddClass(s string)     { addClass(d, s) }

package ui

type TableColumn[T any] struct {
	Name     string
	Renderer func(item T) string
}

type Table[T any] struct {
	cmn     common
	Items   []*T
	Columns []TableColumn[T]
}

func (t *Table[T]) AddClass(s string) Element { t.cmn.addClass(s); return t }
func (t *Table[T]) SetId(s string) Element    { t.cmn.setId(s); return t }
func (t *Table[T]) GetId() string             { return t.cmn.getId() }
func (t *Table[T]) SetAttribute(name string, value string) Element {
	t.cmn.setAttribute(name, value)
	return t
}
func (t *Table[T]) SetValue(value string) { t.cmn.setAttribute("value", value) }
func (t *Table[T]) Render() string {
	result := "<div class='table-div'><div class='table-head-div'><div class='table-head-tr-div'>"
	for _, column := range t.Columns {
		result += "<div class='table-head-td-div'>" + column.Name + "</div>"
	}
	result += "</div></div><div class='table-body-div'>"
	for _, item := range t.Items {
		result += "<div class='table-body-tr-div'>"
		for _, column := range t.Columns {
			result += "<div class='table-body-td-div'>" + column.Renderer(*item) + "</div>"
		}
		result += "</div>"
	}
	result += "</div></div>"
	t.cmn.content = result
	return t.cmn.render("table")
}

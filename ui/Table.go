package ui

import "strconv"

type TableColumn[T any] struct {
	Name     string
	Renderer func(item T) string
	Width    int
}

type Table[T any] struct {
	cmn            common
	Items          []*T
	ItemIdProvider func(item T) string
	Columns        []TableColumn[T]
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
	defaultWidth := 100 / len(t.Columns)
	widths := make([]string, len(t.Columns))
	for i, column := range t.Columns {
		width := column.Width
		if width == 0 {
			width = defaultWidth
		}
		widths[i] = strconv.Itoa(width) + "%;"
	}

	result := "<div class='table-div'><div class='table-head-div'><div class='table-head-tr-div'>"
	for i, column := range t.Columns {
		result += "<div class='table-head-td-div' style='width:" + widths[i] + "'>" + column.Name + "</div>"
	}
	result += "</div></div><div class='table-body-div'>"
	for _, item := range t.Items {
		result += "<div class='table-body-tr-div' "
		if t.ItemIdProvider != nil {
			result += "id='" + t.ItemIdProvider(*item) + "' "
		}
		result += ">"
		for i, column := range t.Columns {
			result += "<div class='table-body-td-div' style='width:" + widths[i] + "'>" + column.Renderer(*item) + "</div>"
		}
		result += "</div>"
	}
	result += "</div></div>"
	t.cmn.content = result
	return t.cmn.render("table")
}

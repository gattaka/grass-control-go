package ui

import (
	"grass-control-go/indexer"
	"io"
	"net/http"
	"strings"
)

func ajaxCallback(url string, callback string) string {
	call := "ajaxCall('" + url + "'"
	if callback != "" {
		call += "," + callback
	}
	call += ")"
	return call
}

func ajax(url string) string {
	return ajaxCallback(url, "")
}

func ConstructPage(items []*indexer.Item, w http.ResponseWriter, fromSearch bool, query string) {
	html := Html{}
	html.Headers = []string{
		"<link rel=\"stylesheet\" href=\"resources/styles.css\"/>",
		"<script type=\"text/javascript\" src=\"resources/scripts.js\"></script>",
		"<link href=\"resources/favicon.png\" rel=\"icon\" sizes=\"16px\">",
	}

	menuDiv := Div{}
	menuDiv.AddClass("menu-div")
	html.Add(&menuDiv)

	menuDiv.Add(NewAnchor("GrassControl", "/"))
	menuDiv.Add(NewAnchor("Reindex", "/reindex"))
	menuDiv.Add(NewAnchor("Ukončit", "/quit"))

	mainDiv := Div{}
	mainDiv.AddClass("main-div")
	html.Add(&mainDiv)

	searchInput := Input{}
	searchInput.SetOnChange("window.location.href = '/?search=' + this.value;")
	searchInput.AddClass("search-div")
	mainDiv.Add(&searchInput)

	controlsDiv := Div{}
	controlsDiv.AddClass("controls-div")
	mainDiv.Add(&controlsDiv)

	controlsDiv.Add(NewButton("&#10006", ajax("clear")))
	controlsDiv.Add(NewButton("", ajax("pause")).SetId("pause-btn"))
	controlsDiv.Add(NewButton("", ajax("prev")).SetId("prev-btn"))
	controlsDiv.Add(NewButton("", ajax("stop")).SetId("stop-btn"))
	controlsDiv.Add(NewButton("", ajax("next")).SetId("next-btn"))
	controlsDiv.Add(NewButton("", ajax("loop")).SetId("loop-btn"))
	controlsDiv.Add(NewButton("", ajax("shuffle")).SetId("shuffle-btn"))

	// Výpis aktuálního umístění
	locationDiv := Div{}
	locationDiv.AddClass("location-div")
	mainDiv.Add(&locationDiv)
	if fromSearch {
		locationDiv.SetValue("Vypisuji výsledek vyhledávání \"" + query + "\"")
	} else {
		returnQuery := "/"
		lastIndex := strings.LastIndex(query, "/")
		if lastIndex > 0 {
			returnQuery = query[:lastIndex]
		}
		locationDiv.Add(NewButton("&#11181;", "window.location.href = '?dir="+returnQuery+"';"))
		locationDiv.Add(NewSpan("Vypisuji výsledek adresáře \"" + query + "\""))
	}

	table := Table[indexer.Item]{}
	table.Items = items

	table.Columns = make([]TableColumn[indexer.Item], 2)
	table.Columns[0] = TableColumn[indexer.Item]{Name: "Název", Renderer: func(itm indexer.Item) string {
		addAndPlayBtn := NewButton("&#9205;", ajax("addAndPlay?id="+itm.GetPath()))
		addBtn := NewButton("&#65291", ajax("add?id="+itm.GetPath()))
		render := addAndPlayBtn.Render() + addBtn.Render()
		if itm.IsDir() {
			dirBtn := NewButton(itm.GetName(), "window.location.href = '?dir="+itm.GetPath()+"';")
			render += dirBtn.Render()
		} else {
			render += itm.GetName()
		}
		return render
	}}
	table.Columns[1] = TableColumn[indexer.Item]{Name: "Nadřazený adresář", Renderer: func(itm indexer.Item) string {
		if !itm.HasParent() {
			return ""
		}
		addAndPlayBtn := NewButton("&#9205;", ajax("addAndPlay?id="+itm.GetParent().GetPath()))
		addBtn := NewButton("&#65291", ajax("add?id="+itm.GetParent().GetPath()))
		render := addAndPlayBtn.Render() + addBtn.Render()
		dirBtn := NewButton(itm.GetParent().GetName(), "window.location.href = '?dir="+itm.GetParent().GetPath()+"';")
		render += dirBtn.Render()
		return render
	}}
	mainDiv.Add(&table)

	result := html.Render()
	io.WriteString(w, result)
}

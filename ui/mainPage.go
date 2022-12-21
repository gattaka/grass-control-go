package ui

import (
	"grass-control-go/indexer"
	"io"
	"net/http"
	"strings"
)

func ajax(url string) string {
	return `const xhttp = new XMLHttpRequest();
					xhttp.open('GET', '` + url + `', true);
					xhttp.send()`
}

func ConstructPage(items []*indexer.Item, w http.ResponseWriter, fromSearch bool, query string) {
	html := Html{}
	html.Headers = []string{
		"<link rel=\"stylesheet\" href=\"resources/styles.css\"/>",
		"<link href=\"resources/favicon.png\" rel=\"icon\" sizes=\"16px\">",
	}

	menuDiv := Div{}
	menuDiv.AddClass("menu-div")
	html.Add(&menuDiv)

	logoAnchor := Anchor{Value: "GrassControl", Link: "/"}
	menuDiv.Add(&logoAnchor)

	reindexAnchor := Anchor{Value: "Reindex", Link: "/reindex"}
	menuDiv.Add(&reindexAnchor)

	quitAnchor := Anchor{Value: "Ukončit", Link: "/quit"}
	menuDiv.Add(&quitAnchor)

	mainDiv := Div{}
	mainDiv.AddClass("main-div")
	html.Add(&mainDiv)

	searchInput := Input{JSfunc: "window.location.href = '/?search=' + this.value;"}
	searchInput.AddClass("search-div")
	mainDiv.Add(&searchInput)

	controlsDiv := Div{}
	controlsDiv.AddClass("controls-div")
	mainDiv.Add(&controlsDiv)

	controlsDiv.Add(&Button{Value: "&#10006", JSfunc: ajax("clear")})
	controlsDiv.Add(&Button{Value: "&#9199;", JSfunc: ajax("pause")})
	controlsDiv.Add(&Button{Value: "&#9198;", JSfunc: ajax("prev")})
	controlsDiv.Add(&Button{Value: "&#9209;", JSfunc: ajax("stop")})
	controlsDiv.Add(&Button{Value: "&#9197;", JSfunc: ajax("next")})

	// Výpis aktuálního umístění
	locationDiv := Div{}
	locationDiv.AddClass("location-div")
	mainDiv.Add(&locationDiv)
	if fromSearch {
		locationDiv.Add(&Text{Value: "Vypisuji výsledek vyhledávání \"" + query + "\""})
	} else {
		returnQuery := "/"
		lastIndex := strings.LastIndex(query, "/")
		if lastIndex > 0 {
			returnQuery = query[:lastIndex]
		}
		locationDiv.Add(&Button{Value: "&#11181;", JSfunc: "window.location.href = '?dir=" + returnQuery + "';"})
		locationDiv.Add(&Text{Value: "Vypisuji výsledek adresáře \"" + query + "\""})
	}

	table := Table[indexer.Item]{}
	table.Items = items

	table.Columns = make([]TableColumn[indexer.Item], 2)
	table.Columns[0] = TableColumn[indexer.Item]{Name: "Název", Renderer: func(itm indexer.Item) string {
		addAndPlayBtn := Button{Value: "&#9205;", JSfunc: ajax("addAndPlay?id=" + itm.GetPath())}
		addBtn := Button{Value: "&#65291", JSfunc: ajax("add?id=" + itm.GetPath())}
		render := addAndPlayBtn.Render() + addBtn.Render()
		if itm.IsDir() {
			dirBtn := Button{Value: itm.GetName(), JSfunc: "window.location.href = '?dir=" + itm.GetPath() + "';"}
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
		addAndPlayBtn := Button{Value: "&#9205;", JSfunc: ajax("addAndPlay?id=" + itm.GetParent().GetPath())}
		addBtn := Button{Value: "&#65291", JSfunc: ajax("add?id=" + itm.GetParent().GetPath())}
		render := addAndPlayBtn.Render() + addBtn.Render()
		dirBtn := Button{Value: itm.GetParent().GetName(), JSfunc: "window.location.href = '?dir=" + itm.GetParent().GetPath() + "';"}
		render += dirBtn.Render()
		return render
	}}
	mainDiv.Add(&table)

	result := html.Render()
	io.WriteString(w, result)
}

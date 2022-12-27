package ui

import (
	"grass-control-go/indexer"
	"grass-control-go/utils"
	"io"
	"net/http"
	"strings"
)

func prepAjax(url string) string {
	return "ajaxCall('" + url + "')"
}

func prepAjaxWithParam(url string, paramToEncode string) string {
	return prepAjax(url + utils.EncodeURL(paramToEncode))
}

func prepDirNavigate(urlToEncode string) string {
	return "window.location.href='/?dir=" + utils.EncodeURL(urlToEncode) + "'"
}

func ConstructPage(items []*indexer.Item, w http.ResponseWriter, fromSearch bool, query string) {
	html := Html{}
	html.Headers = []string{
		"<link rel='stylesheet' href='resources/styles.css'/>",
		"<script type='text/javascript' src='resources/scripts.js'></script>",
		"<link href='resources/favicon.png' rel='icon' sizes='16px'>",
	}

	infoDiv := Div{}
	infoDiv.SetId("info-div")
	infoDiv.SetOnClick("this.classList.remove('info-div-show','error','info');")
	html.Add(&infoDiv)

	menuDiv := Div{}
	menuDiv.AddClass("menu-div")
	html.Add(&menuDiv)

	menuDiv.Add(NewAnchor("GrassControl", "/"))
	menuDiv.Add(NewAnchor("Reindex", "/reindex"))
	menuDiv.Add(NewAnchorJS("Ukončit", prepAjax("/quit")))

	mainDiv := Div{}
	mainDiv.AddClass("main-div")
	html.Add(&mainDiv)

	searchInput := Input{}
	searchInput.SetOnChange("window.location.href='/?search=' + encodeURIComponent(this.value);")
	searchInput.AddClass("search-div")
	mainDiv.Add(&searchInput)

	controlsDiv := Div{}
	controlsDiv.AddClass("controls-div")
	mainDiv.Add(&controlsDiv)

	controlsDiv.Add(NewButton("&#10006", prepAjax("clear")))
	controlsDiv.Add(NewButton("", prepAjax("pause")).SetId("pause-btn"))
	controlsDiv.Add(NewButton("", prepAjax("prev")).SetId("prev-btn"))
	controlsDiv.Add(NewButton("", prepAjax("stop")).SetId("stop-btn"))
	controlsDiv.Add(NewButton("", prepAjax("next")).SetId("next-btn"))
	controlsDiv.Add(NewButton("", prepAjax("loop")).SetId("loop-btn"))
	controlsDiv.Add(NewButton("", prepAjax("shuffle")).SetId("shuffle-btn"))

	currentSongDiv := Div{}
	currentSongDiv.SetId("current-song-div")
	controlsDiv.Add(&currentSongDiv)

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
		locationDiv.Add(NewButton("&#11181;", prepDirNavigate(returnQuery)))
		locationDiv.Add(NewSpan("Vypisuji výsledek adresáře \"" + query + "\""))
	}

	table := Table[indexer.Item]{}
	table.Items = items

	table.Columns = make([]TableColumn[indexer.Item], 2)
	table.Columns[0] = TableColumn[indexer.Item]{Name: "Název", Renderer: func(itm indexer.Item) string {
		addAndPlayBtn := NewButton("&#9205;", prepAjaxWithParam("addAndPlay?id=", itm.GetPath()))
		addBtn := NewButton("&#65291", prepAjaxWithParam("add?id=", itm.GetPath()))
		render := addAndPlayBtn.Render() + addBtn.Render()
		if itm.IsDir() {
			dirBtn := NewButton(itm.GetName(), prepDirNavigate(itm.GetPath()))
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
		addAndPlayBtn := NewButton("&#9205;", prepAjaxWithParam("addAndPlay?id=", itm.GetParent().GetPath()))
		addBtn := NewButton("&#65291", prepAjaxWithParam("add?id=", itm.GetParent().GetPath()))
		render := addAndPlayBtn.Render() + addBtn.Render()
		dirBtn := NewButton(itm.GetParent().GetName(), prepDirNavigate(itm.GetParent().GetPath()))
		render += dirBtn.Render()
		return render
	}}
	mainDiv.Add(&table)

	result := html.Render()
	io.WriteString(w, result)
}

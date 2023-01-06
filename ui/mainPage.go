package ui

import (
	"grass-control-go/indexer"
	"grass-control-go/utils"
	"strconv"
	"strings"
)

const AddEndpoint = "/add"
const AddAndPlayEndpoint = "/addAndPlay"
const PlayEndpoint = "/play"
const IdParam = "id"
const ValueParam = "value"
const DirParam = "dir"
const SearchParam = "grass-control-search"

func prepAjax(url string) string {
	return "ajaxCall('" + url + "')"
}

func prepAjaxWithParam(url string, paramToEncode string) string {
	return prepAjax(url + utils.EncodeURL(paramToEncode))
}

func prepDirNavigate(urlToEncode string) string {
	return "window.location.href='/?" + DirParam + "=" + utils.EncodeURL(urlToEncode) + "'"
}

func ConstructPlaylist(items *[]*utils.VlcPlaylistNode) string {
	table := Table[utils.VlcPlaylistNode]{}

	table.Items = *items

	table.ItemIdProvider = func(item utils.VlcPlaylistNode) string {
		return "playlist-item-" + item.Id
	}

	table.Columns = make([]TableColumn[utils.VlcPlaylistNode], 2)
	nameColumn := TableColumn[utils.VlcPlaylistNode]{Name: "Název", Renderer: func(itm utils.VlcPlaylistNode) string {
		playBtn := NewButton("&#9205;", prepAjaxWithParam(PlayEndpoint+"?"+IdParam+"=", itm.Id))
		render := playBtn.Render()
		render += itm.Name
		return render
	}}
	nameColumn.Width = 80

	table.Columns[0] = nameColumn
	durationColumn := TableColumn[utils.VlcPlaylistNode]{Name: "Délka", Renderer: func(itm utils.VlcPlaylistNode) string {
		sec := strconv.Itoa(itm.Duration % 60)
		if len(sec) == 1 {
			sec = "0" + sec
		}
		return strconv.Itoa(itm.Duration/60) + ":" + sec
	}}
	durationColumn.Width = 20
	table.Columns[1] = durationColumn

	return table.Render()
}

func ConstructPage(items []*indexer.Item, fromSearch bool, query string) string {
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

	menuLeftDiv := Div{}
	menuDiv.Add(&menuLeftDiv)
	menuLeftDiv.Add(NewAnchor("GrassControl", "/"))
	menuLeftDiv.Add(NewAnchor("Reindex", "/reindex"))

	menuRightDiv := Div{}
	menuDiv.Add(&menuRightDiv)
	menuRightDiv.Add(NewAnchorJS("Ukončit", prepAjax("/quit")))

	mainDiv := Div{}
	mainDiv.SetId("main-div")
	html.Add(&mainDiv)

	libraryDiv := Div{}
	libraryDiv.SetId("library-div")
	mainDiv.Add(&libraryDiv)

	playlistDiv := Div{}
	playlistDiv.SetId("playlist-div")
	mainDiv.Add(&playlistDiv)

	searchForm := Form{}
	searchForm.AddClass("search-form")
	searchForm.SetMethod("get")
	searchForm.SetAction("/")
	libraryDiv.Add(&searchForm)

	searchInput := TextInput{}
	searchInput.SetValue(query)
	searchInput.SetName(SearchParam)
	searchInput.SetAttribute("autocomplete", "do-not-autofill")
	searchForm.Add(&searchInput)

	currentSongDiv := Div{}
	currentSongDiv.SetId("current-song-div")
	libraryDiv.Add(&currentSongDiv)

	progressDiv := Div{}
	progressDiv.SetId("progress-div")
	libraryDiv.Add(&progressDiv)

	progressTimeSpan := Span{}
	progressTimeSpan.SetId("progress-time-span")
	progressDiv.Add(&progressTimeSpan)

	progressControl := RangeInput{}
	progressControl.SetMin(0)
	progressControl.SetMax(100)
	progressControl.SetId("progress-slider")
	progressControlChangeToggle := "elementsUnderChange['" + progressControl.GetId() + "']="
	progressControl.SetAttribute("onmousedown", progressControlChangeToggle+"true;")
	progressControl.SetAttribute("onmouseup", progressControlChangeToggle+"false;")
	progressControl.SetAttribute("onblur", progressControlChangeToggle+"false;")
	progressControl.SetOnChange("ajaxCall('progress?value='+this.value);")
	progressControl.SetAttribute("onwheel", "progressControlScroll(event, val => {ajaxCall('progress?value='+val)});")
	progressDiv.Add(&progressControl)

	progressLengthSpan := Span{}
	progressLengthSpan.SetId("progress-length-span")
	progressDiv.Add(&progressLengthSpan)

	controlsDiv := Div{}
	controlsDiv.AddClass("controls-div")
	libraryDiv.Add(&controlsDiv)

	controlsDiv.Add(NewButton("", prepAjax("pause")).SetId("play-pause-btn"))
	controlsDiv.Add(NewButton("", prepAjax("prev")).SetId("prev-btn"))
	controlsDiv.Add(NewButton("", prepAjax("stop")).SetId("stop-btn"))
	controlsDiv.Add(NewButton("", prepAjax("next")).SetId("next-btn"))
	controlsDiv.Add(NewButton("", prepAjax("loop")).SetId("loop-btn"))
	controlsDiv.Add(NewButton("", prepAjax("shuffle")).SetId("shuffle-btn"))

	volumeControlDiv := Div{}
	volumeControlDiv.SetId("volume-div")
	controlsDiv.Add(&volumeControlDiv)

	volumeControl := RangeInput{}
	volumeControl.SetMin(0)
	volumeControl.SetMax(320)
	volumeControl.SetId("volume-slider")
	volumeControlChangeToggle := "elementsUnderChange['" + volumeControl.GetId() + "']="
	volumeControl.SetAttribute("onmousedown", volumeControlChangeToggle+"true;")
	volumeControl.SetAttribute("onmouseup", volumeControlChangeToggle+"false;")
	volumeControl.SetAttribute("onblur", volumeControlChangeToggle+"false;")
	volumeControl.SetOnChange("ajaxCall('volume?value='+this.value);")
	volumeControl.SetAttribute("onwheel", "volumeControlScroll(event, val => {ajaxCall('volume?value='+val)});")
	volumeControlDiv.Add(&volumeControl)

	volumeSpan := Span{}
	volumeSpan.SetId("volume-span")
	volumeControlDiv.Add(&volumeSpan)

	// Výpis aktuálního umístění
	locationDiv := Div{}
	locationDiv.AddClass("location-div")
	libraryDiv.Add(&locationDiv)

	locationDiv.Add(NewButton("&#10006", prepAjax("clear")))

	tableBtnsParam := ""
	if fromSearch {
		tableBtnsParam = SearchParam
	} else {
		tableBtnsParam = IdParam
	}

	tableAddAndPlayBtn := NewButton("&#9205;", prepAjaxWithParam(AddAndPlayEndpoint+"?"+tableBtnsParam+"=", query))
	tableAddBtn := NewButton("&#65291", prepAjaxWithParam(AddEndpoint+"?"+tableBtnsParam+"=", query))
	locationDiv.Add(tableAddAndPlayBtn)
	locationDiv.Add(tableAddBtn)

	if fromSearch {
		locationDiv.Add(NewSpan("Vypisuji výsledek vyhledávání \"" + query + "\""))
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
		addAndPlayBtn := NewButton("&#9205;", prepAjaxWithParam(AddAndPlayEndpoint+"?"+IdParam+"=", itm.GetPath()))
		addBtn := NewButton("&#65291", prepAjaxWithParam(AddEndpoint+"?"+IdParam+"=", itm.GetPath()))
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
		addAndPlayBtn := NewButton("&#9205;", prepAjaxWithParam(AddAndPlayEndpoint+"?"+IdParam+"=", itm.GetParent().GetPath()))
		addBtn := NewButton("&#65291", prepAjaxWithParam(AddEndpoint+"?"+IdParam+"=", itm.GetParent().GetPath()))
		render := addAndPlayBtn.Render() + addBtn.Render()
		dirBtn := NewButton(itm.GetParent().GetName(), prepDirNavigate(itm.GetParent().GetPath()))
		render += dirBtn.Render()
		return render
	}}
	libraryDiv.Add(&table)

	result := html.Render()
	return result
}

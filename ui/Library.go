package ui

import (
	"grass-control-go/indexer"
	"grass-control-go/ui/common"
	"grass-control-go/ui/common/elements"
	"strings"
)

func constructLibrary(items []*indexer.Item, fromSearch bool, query string) elements.Element {
	libraryDiv := elements.Div{}
	libraryDiv.SetId("library-div")

	currentSongDiv := elements.Div{}
	currentSongDiv.SetId("current-song-div")
	libraryDiv.Add(&currentSongDiv)

	libraryDiv.Add(createProgressDiv())
	libraryDiv.Add(createControlsDiv())
	libraryDiv.Add(createSearchDiv(fromSearch, query))
	libraryDiv.Add(createLocationDiv(fromSearch, query))
	libraryDiv.Add(createTable(items))

	return &libraryDiv
}

func createProgressDiv() elements.Element {
	progressDiv := elements.Div{}
	progressDiv.SetId("progress-div")

	progressTimeSpan := elements.Span{}
	progressTimeSpan.SetId("progress-time-span")
	progressDiv.Add(&progressTimeSpan)

	progressControl := elements.RangeInput{}
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

	progressLengthSpan := elements.Span{}
	progressLengthSpan.SetId("progress-length-span")
	progressDiv.Add(&progressLengthSpan)

	return &progressDiv
}

func createControlsDiv() elements.Element {
	controlsDiv := elements.Div{}
	controlsDiv.AddClass("controls-div")

	controlsDiv.Add(elements.NewButton("", prepAjax("pause")).SetId("play-pause-btn"))
	controlsDiv.Add(elements.NewButton("", prepAjax("prev")).SetId("prev-btn"))
	controlsDiv.Add(elements.NewButton("", prepAjax("stop")).SetId("stop-btn"))
	controlsDiv.Add(elements.NewButton("", prepAjax("next")).SetId("next-btn"))
	controlsDiv.Add(elements.NewButton("", prepAjax("loop")).SetId("loop-btn"))
	controlsDiv.Add(elements.NewButton("", prepAjax("shuffle")).SetId("shuffle-btn"))

	controlsDiv.Add(createVolumeControl())

	return &controlsDiv
}

func createVolumeControl() elements.Element {
	volumeControlDiv := elements.Div{}
	volumeControlDiv.SetId("volume-div")

	volumeControl := elements.RangeInput{}
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

	volumeSpan := elements.Span{}
	volumeSpan.SetId("volume-span")
	volumeControlDiv.Add(&volumeSpan)

	return &volumeControlDiv
}

func createLocationDiv(fromSearch bool, query string) elements.Element {
	locationDiv := elements.Div{}
	locationDiv.AddClass("location-div")

	tableBtnsParam := ""
	if fromSearch {
		tableBtnsParam = common.SearchParam
	} else {
		tableBtnsParam = common.IdParam
	}

	tableAddAndPlayBtn := elements.NewButton(common.PlayUnicode, prepAjaxWithParam(common.AddAndPlayEndpoint+"?"+tableBtnsParam+"=", query))
	tableAddAndPlayBtn.AddClass(common.TableControlBtnClass)
	tableAddBtn := elements.NewButton(common.PlusUnicode, prepAjaxWithParam(common.AddEndpoint+"?"+tableBtnsParam+"=", query))
	tableAddBtn.AddClass(common.TableControlBtnClass)
	locationDiv.Add(tableAddAndPlayBtn)
	locationDiv.Add(tableAddBtn)

	if fromSearch {
		locationDiv.Add(elements.NewSpan("Vypisuji výsledek vyhledávání \"" + query + "\""))
	} else {
		returnQuery := "/"
		lastIndex := strings.LastIndex(query, "/")
		if lastIndex > 0 {
			returnQuery = query[:lastIndex]
		}
		locationDiv.Add(elements.NewButton("&#11181;", prepDirNavigate(returnQuery)))
		locationDiv.Add(elements.NewSpan("Vypisuji výsledek adresáře \"" + query + "\""))
	}

	return &locationDiv
}

func createTable(items []*indexer.Item) elements.Element {
	table := elements.Table[indexer.Item]{}
	table.SetId("library-table")
	table.Items = items

	table.Columns = make([]elements.TableColumn[indexer.Item], 2)
	table.Columns[0] = elements.TableColumn[indexer.Item]{Name: "Název", Renderer: func(itm indexer.Item) string {
		btnsDiv := elements.Div{}
		btnsDiv.AddClass(common.ControlBtnsDivClass)
		addAndPlayBtn := elements.NewButton(common.PlayUnicode, prepAjaxWithParam(common.AddAndPlayEndpoint+"?"+common.IdParam+"=", itm.GetPath()))
		addBtn := elements.NewButton(common.PlusUnicode, prepAjaxWithParam(common.AddEndpoint+"?"+common.IdParam+"=", itm.GetPath()))
		btnsDiv.Add(addAndPlayBtn)
		btnsDiv.Add(addBtn)
		render := btnsDiv.Render()
		if itm.IsDir() {
			dirBtn := elements.NewButton(itm.GetName(), prepDirNavigate(itm.GetPath()))
			render += dirBtn.Render()
		} else {
			render += elements.NewSpan(itm.GetName()).Render()
		}
		return render
	}}
	table.Columns[1] = elements.TableColumn[indexer.Item]{Name: "Nadřazený adresář", Renderer: func(itm indexer.Item) string {
		if !itm.HasParent() {
			return ""
		}
		btnsDiv := elements.Div{}
		btnsDiv.AddClass(common.ControlBtnsDivClass)
		addAndPlayBtn := elements.NewButton(common.PlayUnicode, prepAjaxWithParam(common.AddAndPlayEndpoint+"?"+common.IdParam+"=", itm.GetParent().GetPath()))
		addBtn := elements.NewButton(common.PlusUnicode, prepAjaxWithParam(common.AddEndpoint+"?"+common.IdParam+"=", itm.GetParent().GetPath()))
		btnsDiv.Add(addAndPlayBtn)
		btnsDiv.Add(addBtn)
		render := btnsDiv.Render()
		dirBtn := elements.NewButton(itm.GetParent().GetName(), prepDirNavigate(itm.GetParent().GetPath()))
		render += dirBtn.Render()
		return render
	}}

	return &table
}

func createSearchDiv(fromSearch bool, query string) elements.Element {
	searchForm := elements.Form{}
	searchForm.AddClass("search-form")
	searchForm.SetMethod("get")
	searchForm.SetAction("/")

	searchForm.Add(elements.NewSpan("Vyhledat"))

	searchInput := elements.TextInput{}
	if fromSearch {
		searchInput.SetValue(query)
	}
	searchInput.SetName(common.SearchParam)
	searchInput.SetId("search-input")
	searchInput.SetAttribute("autocomplete", "do-not-autofill")
	searchForm.Add(&searchInput)

	return &searchForm
}

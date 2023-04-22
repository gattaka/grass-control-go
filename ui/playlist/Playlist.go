package playlist

import (
	"grass-control-go/ui/common"
	"grass-control-go/ui/common/elements"
	uiUtils "grass-control-go/ui/utils"
	"grass-control-go/utils"
	"strconv"
)

func ConstructPlaylist() *elements.Div {
	playlistDiv := elements.Div{}
	playlistDiv.SetId("playlist-div")

	playlistDiv.Add(createSearchDiv())
	playlistDiv.Add(createControlsDiv())

	tableDiv := elements.Div{}
	tableDiv.SetId("playlist-table-div")
	playlistDiv.Add(&tableDiv)
	tableDiv.Add(ConstructPlaylistTable(nil))

	return &playlistDiv
}

func createControlsDiv() elements.Element {
	controlDiv := elements.Div{}
	controlDiv.AddClass("playlist-controls-div")

	controlDiv.Add(elements.NewButton("Vyčistit", uiUtils.PrepAjax("clear")))
	controlDiv.Add(elements.NewButton("Nechat jet hrající", uiUtils.PrepAjax("clearExceptPlaying")))

	return &controlDiv
}

func createSearchDiv() elements.Element {
	searchDiv := elements.Div{}
	searchDiv.AddClass("search-form")
	searchDiv.Add(elements.NewSpan("Vyhledat"))

	searchInput := elements.TextInput{}
	searchInput.SetName(common.SearchParam)
	searchInput.SetId("playlist-search-input")
	searchInput.SetAttribute("autocomplete", "do-not-autofill")
	searchInput.SetAttribute("onkeypress", "searchInPlaylist(event)")
	searchDiv.Add(&searchInput)

	return &searchDiv
}

func ConstructPlaylistTable(items []*utils.VlcPlaylistNode) *elements.Table[utils.VlcPlaylistNode] {
	table := elements.Table[utils.VlcPlaylistNode]{}
	table.SetId("playlist-table")

	table.Items = items

	table.ItemIdProvider = func(item utils.VlcPlaylistNode) string {
		return "playlist-item-" + item.Id
	}

	table.Columns = make([]elements.TableColumn[utils.VlcPlaylistNode], 2)
	nameColumn := elements.TableColumn[utils.VlcPlaylistNode]{Name: "Název", Renderer: func(itm utils.VlcPlaylistNode) string {
		btnsDiv := elements.Div{}
		btnsDiv.AddClass(common.ControlBtnsDivClass)
		removeBtn := elements.NewButton(common.CrossUnicode, uiUtils.PrepAjaxWithParam(common.RemoveEndpoint+"?"+common.IdParam+"=", itm.Id))
		removeBtn.AddClass(common.TableControlBtnClass)
		playBtn := elements.NewButton(common.PlayUnicode, uiUtils.PrepAjaxWithParam(common.PlayEndpoint+"?"+common.IdParam+"=", itm.Id))
		playBtn.AddClass(common.TableControlBtnClass)
		btnsDiv.Add(removeBtn)
		btnsDiv.Add(playBtn)
		render := btnsDiv.Render()
		render += elements.NewSpan(itm.Name).AddClass("playlist-item").Render()
		return render
	}}
	nameColumn.Width = 80

	table.Columns[0] = nameColumn
	durationColumn := elements.TableColumn[utils.VlcPlaylistNode]{Name: "Délka", Renderer: func(itm utils.VlcPlaylistNode) string {
		sec := strconv.Itoa(itm.Duration % 60)
		if len(sec) == 1 {
			sec = "0" + sec
		}
		return strconv.Itoa(itm.Duration/60) + ":" + sec
	}}
	durationColumn.Width = 20
	table.Columns[1] = durationColumn

	return &table
}

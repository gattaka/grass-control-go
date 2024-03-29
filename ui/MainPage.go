package ui

import (
	"grass-control-go/indexer"
	"grass-control-go/ui/common/elements"
	"grass-control-go/ui/library"
	"grass-control-go/ui/playlist"
	"grass-control-go/ui/utils"
)

func ConstructPage(items []*indexer.Item, fromSearch bool, query string) string {
	html := elements.Html{}
	html.Headers = []string{
		"<link rel='stylesheet' href='resources/styles.css'/>",
		"<script type='text/javascript' src='resources/scripts.js'></script>",
		"<link href='resources/favicon.png' rel='icon' sizes='16px'>",
	}

	infoDiv := elements.Div{}
	infoDiv.SetId("info-div")
	infoDiv.SetOnClick("this.classList.remove('info-div-show','error','info');")
	html.Add(&infoDiv)

	menuDiv := elements.Div{}
	menuDiv.AddClass("menu-div")
	html.Add(&menuDiv)

	menuLeftDiv := elements.Div{}
	menuDiv.Add(&menuLeftDiv)
	menuLeftDiv.Add(elements.NewAnchor("GrassControl", "/"))
	menuLeftDiv.Add(elements.NewAnchor("Reindex", "/reindex"))

	menuRightDiv := elements.Div{}
	menuDiv.Add(&menuRightDiv)
	menuRightDiv.Add(elements.NewAnchorJS("Ukončit", utils.PrepAjax("/quit")))

	mainDiv := elements.Div{}
	mainDiv.SetId("main-div")
	html.Add(&mainDiv)

	mainDiv.Add(library.ConstructLibrary(items, fromSearch, query))

	mainDiv.Add(playlist.ConstructPlaylist())

	result := html.Render()
	return result
}

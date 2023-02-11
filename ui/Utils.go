package ui

import (
	"grass-control-go/ui/common"
	"grass-control-go/utils"
)

func prepAjax(url string) string {
	return "ajaxCall('" + url + "')"
}

func prepAjaxWithParam(url string, paramToEncode string) string {
	return prepAjax(url + utils.EncodeURL(paramToEncode))
}

func prepDirNavigate(urlToEncode string) string {
	return "window.location.href='/?" + common.DirParam + "=" + utils.EncodeURL(urlToEncode) + "'"
}

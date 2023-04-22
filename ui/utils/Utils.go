package utils

import (
	"grass-control-go/ui/common"
	"grass-control-go/utils"
)

func PrepAjax(url string) string {
	return "ajaxCall('" + url + "')"
}

func PrepAjaxWithParam(url string, paramToEncode string) string {
	return PrepAjax(url + utils.EncodeURL(paramToEncode))
}

func PrepDirNavigate(urlToEncode string) string {
	return "window.location.href='/?" + common.DirParam + "=" + utils.EncodeURL(urlToEncode) + "'"
}

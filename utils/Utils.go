package utils

import "net/url"

func EncodeURL(value string) string {
	// https://go.dev/play/p/pOfrn-Wsq5
	url := &url.URL{Path: value}
	// URL předsadí před sebe './'
	encoded := url.String()
	if len(encoded) >= 2 && encoded[:2] == "./" {
		encoded = encoded[2:]
	}
	// VLC má vadu a nebere URL mezery jako '+', zvládá jen '%20'
	//encoded := url.QueryEscape("file:///D:/Hudba/" + value)
	return encoded
}

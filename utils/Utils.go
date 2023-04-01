package utils

import (
	"net/url"
	"strings"
)

func EncodeURL(value string) string {
	// https://go.dev/play/p/pOfrn-Wsq5
	encoded := url.QueryEscape(value)
	// VLC má vadu a nebere URL mezery jako '+', zvládá jen '%20'
	encoded = strings.ReplaceAll(encoded, "+", "%20")
	return encoded
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

package utils

import (
	"bytes"
	"strings"
)

// SplicePath splice
// eg: SplicePath("http://192.168.1.108:8080/", "chain1", "/transaction/", "6788795865a6d89e878ad9e999")
// get: http://192.168.1.108:8080/chain1/transaction/6788795865a6d89e878ad9e999/
func SplicePath(urls ...string) string {
	var bt bytes.Buffer
	for _, url := range urls {
		bt.WriteString(url)
		if !strings.HasSuffix(url, "/") {
			bt.WriteString("/")
		}
	}
	r := strings.ReplaceAll(bt.String(), "//", "/")
	r = strings.ReplaceAll(r, "http:/", "http://")
	r = strings.ReplaceAll(r, "https:/", "https://")
	return r
}

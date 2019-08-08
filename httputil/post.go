// Package httputil is simple http POST request wrapper
package httputil

import (
	"io"
	"log"
	"net/http"
	"strings"
)

// PostPage post to page
func PostPage(url, key string) (io.Reader, error) {
	res, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(key))
	if err != nil {
		log.Printf("PostPage error:%s\n", err.Error())
		return nil, err
	}
	return DecodeHTMLBody(res.Body)
}

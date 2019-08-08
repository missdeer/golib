// Package httputil is simple http decode request wrapper
package httputil

import (
	"bufio"
	"io"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/htmlindex"
)

func detectContentCharset(body io.Reader) (string, io.Reader) {
	r := bufio.NewReader(body)
	if data, err := r.Peek(1024); err == nil {
		_, name, _ := charset.DetermineEncoding(data, "")
		return name, r
	}
	return "utf-8", r
}

// DecodeHTMLBody returns an decoding reader of the html Body for the specified `charset`
// If `charset` is empty, DecodeHTMLBody tries to guess the encoding from the content
func DecodeHTMLBody(body io.Reader) (io.Reader, error) {
	// if charset == "" {
	charset, r := detectContentCharset(body)
	// }
	e, err := htmlindex.Get(charset)
	if err != nil {
		return nil, err
	}
	if name, _ := htmlindex.Name(e); name != "utf-8" {
		body = e.NewDecoder().Reader(r)
	} else {
		body = r
	}
	return body, nil
}

// Package httputil is simple http POST request wrapper
package httputil

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// PostPage post to page
func PostPage(url, key string) (io.Reader, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
	}
	client := &http.Client{Transport: tr, Timeout: 10 * time.Second}

	req, err := http.NewRequest("POST", url, strings.NewReader(key))
	if err != nil {
		log.Printf("PostPage error:%s\n", err.Error())
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		log.Printf("PostPage error:%s\n", err.Error())
		return nil, err
	}

	return DecodeHTMLBody(res.Body)
}

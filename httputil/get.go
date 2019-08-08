// Package httputil is simple http GET request wrapper
package httputil

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.108 Safari/537.36"

// GetPage get HTML page by url
func GetPage(url, ua string) (io.Reader, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("GetPage error:%s\n", err.Error())
		return nil, err

	}

	if ua == "" {
		ua = userAgent
	}
	req.Header.Set("User-Agent", ua)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("GetPage error:%s\n", err.Error())
		return nil, err
	}

	// defer resp.Body.Close()
	// return resp.Body, nil
	return DecodeHTMLBody(resp.Body)
}

// GetHostByURL get host by url
func GetHostByURL(url string) string {
	if url == "" {
		return ""
	}
	if strings.HasPrefix(url, "/") {
		return ""
	}
	urlStrArray := strings.Split(url, "//")
	if len(urlStrArray) != 2 {
		return ""
	}
	schema := urlStrArray[0]
	host := strings.Split(urlStrArray[1], "/")[0]
	return fmt.Sprintf("%s//%s", schema, host)
}

// GetBytes returns content as []byte
func GetBytes(u string, headers http.Header, timeout time.Duration, retryCount int) (c []byte, err error) {
	client := &http.Client{
		Timeout: timeout,
	}
	retry := 0
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Println("Could not parse novel page request:", err)
		return
	}

	req.Header = headers
doRequest:
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Could not send request:", err)
		retry++
		if retry < retryCount {
			time.Sleep(3 * time.Second)
			goto doRequest
		}
		return
	}

	if resp.StatusCode != 200 {
		resp.Body.Close()
		log.Println("response not 200:", resp.StatusCode, resp.Status)
		retry++
		if retry < retryCount {
			time.Sleep(3 * time.Second)
			goto doRequest
		}
		return
	}

	c, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Println("reading content failed")
		retry++
		if retry < retryCount {
			time.Sleep(3 * time.Second)
			goto doRequest
		}
		return
	}
	return
}

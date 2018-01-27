// Package httputil is simple http GET request wrapper
package httputil

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

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

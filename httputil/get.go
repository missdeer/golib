// Package httputil is simple http GET request wrapper
package httputil

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"
	"golang.org/x/net/publicsuffix"
)

var (
	errorNotIP    = errors.New("addr is not an IP")
	resolveResult = sync.Map{}
	once          = sync.Once{}
	globalClient  *http.Client
)

func patchAddress(addr string) (string, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return addr, err
	}
	ip := net.ParseIP(host)
	if ip.To4() != nil || ip.To16() != nil {
		return addr, errorNotIP
	}
	// query from cache
	if rr, ok := resolveResult.Load(host); ok {
		ips := rr.([]string)
		if len(ips) > 0 {
			return net.JoinHostPort(ips[rand.Intn(len(ips))], port), nil
		}
	}
	// resolve it via http://119.29.29.29/d?dn=api.baidu.com
	client := getGlobalHttpClient()
	req, err := http.NewRequest("GET", fmt.Sprintf("http://119.29.29.29/d?dn=%s", host), nil)
	if err != nil {
		log.Println(err)
		return addr, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return addr, err
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return addr, err
	}
	ips := string(content)
	ss := strings.Split(ips, ";")
	if len(ss) == 0 {
		return addr, err
	}
	resolveResult.Store(host, ss)
	return net.JoinHostPort(ss[0], port), nil
}

type dialer struct {
	addr   string
	socks5 proxy.Dialer
}

func (d *dialer) socks5DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	// TODO: golang.org/x/net/proxy need to add socks5DialContext
	return d.socks5Dial(network, addr)
}

func (d *dialer) socks5Dial(network, addr string) (net.Conn, error) {
	var err error
	if d.socks5 == nil {
		d.socks5, err = proxy.SOCKS5("tcp", d.addr, nil, proxy.Direct)
		if err != nil {
			return nil, err
		}
	}

	addr, _ = patchAddress(addr)
	return d.socks5.Dial(network, addr)
}

func socks5ProxyTransport(addr string) *http.Transport {
	d := &dialer{addr: addr}
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
		DialContext:     d.socks5DialContext,
		Dial:            d.socks5Dial,
	}
}

func createHttpClient(timeout time.Duration) *http.Client {
	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
	}
	client := &http.Client{
		Transport: tr,
		Jar:       jar,
		Timeout:   timeout,
	}

	httpProxy := os.Getenv("HTTP_PROXY")
	httpsProxy := os.Getenv("HTTPS_PROXY")
	socks5Proxy := os.Getenv("SOCKS5_PROXY")
	allProxy := os.Getenv("ALL_PROXY")
	if httpProxy != "" || httpsProxy != "" || allProxy != "" {
		if httpProxy == "" {
			httpProxy = httpsProxy
		}
		if httpProxy == "" {
			httpProxy = allProxy
		}
		if proxyURL, err := url.Parse(httpProxy); err == nil {
			transport := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
				Proxy:           http.ProxyURL(proxyURL),
			}
			dialer := &net.Dialer{}
			transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				addr, _ = patchAddress(addr)
				return dialer.DialContext(ctx, network, addr)
			}
			transport.Dial = func(network, addr string) (net.Conn, error) {
				addr, _ = patchAddress(addr)
				return dialer.Dial(network, addr)
			}
			client.Transport = transport
		}
	} else if socks5Proxy != "" {
		client.Transport = socks5ProxyTransport(socks5Proxy)
	}
	return client
}

func getGlobalHttpClient() *http.Client {
	once.Do(func() { globalClient = createHttpClient(30 * time.Second) })
	return globalClient
}

// GetPage get HTML page by url
func GetPage(url, ua string) (io.Reader, error) {
	client := createHttpClient(10 * time.Second)
	defer func() {
		client = nil
	}()
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
func GetHostByURL(u string) (host string) {
	theURL, e := url.Parse(u)
	if e != nil {
		log.Println(u, e)
		return
	}
	return fmt.Sprintf("%s://%s", theURL.Scheme, theURL.Host)
}

// GetBytes returns content as []byte
func GetBytes(u string, headers http.Header, timeout time.Duration, retryCount int) (c []byte, err error) {
	client := createHttpClient(timeout)
	defer func() {
		client = nil
	}()
	retry := 0
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Println("Could not parse novel page request:", err)
		return nil, err
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
		return nil, err
	}

	if resp.StatusCode != 200 {
		resp.Body.Close()
		log.Println("response not 200:", resp.StatusCode, resp.Status)
		retry++
		if retry < retryCount {
			time.Sleep(3 * time.Second)
			goto doRequest
		}
		return nil, fmt.Errorf("response code: %d, status: %s", resp.StatusCode, resp.Status)
	}

	c, err = io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Println("reading content failed")
		retry++
		if retry < retryCount {
			time.Sleep(3 * time.Second)
			goto doRequest
		}
		return nil, err
	}
	return c, nil
}

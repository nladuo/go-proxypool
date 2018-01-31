package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	simplejson "github.com/bitly/go-simplejson"
)

const (
	ConcurNum int = 20 //并发数
)

/**
 * 在HttpBin验证ip
 * proxy_url: http://[host]:[port]
 */
func validHTTPBin(proxyURL string) (bool, string) {
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(proxyURL)
	}
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(8 * time.Second),
	}
	resp, err := client.Get("http://httpbin.org/ip")

	if err != nil {
		return false, err.Error()
	}

	data, _ := ioutil.ReadAll(resp.Body)

	_, err = simplejson.NewJson(data) //判断是否是json

	if err != nil {
		return false, string(data)
	}

	return true, string(data)
}

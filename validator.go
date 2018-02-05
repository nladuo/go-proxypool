package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/**
 * 在脉脉网验证
 **/
func validMaiMai(proxyURL string) bool {
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(proxyURL)
	}
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(8 * time.Second),
	}
	resp, err := client.Get("https://maimai.cn/contact/comment_list/38253207?jsononly=1")

	if err != nil {
		return false
	}

	data, _ := ioutil.ReadAll(resp.Body)

	js, err := simplejson.NewJson(data) //判断是否是json

	if err != nil {
		return false
	}

	str, err := js.Get("result").String()

	if err != nil {
		return false
	}

	return str == "ok"
}

type Data struct {
	proxy Proxy
	key   string
	value bool
}

/**
 * 验证爬虫，每隔一段时间一次
 **/
func validCrawler(session *mgo.Session, success bool) {
	c := session.DB("go-proxytool").C("proxy")
	proxies := []Proxy{}
	err := c.Find(bson.M{"maimai": success}).All(&proxies)

	if err != nil {
		panic(err)
	}

	dataChan := make(chan Data, ConcurNum)
	occupyChan := make(chan bool, ConcurNum)
	exitChan := make(chan bool, 1)
	go func() { // 代理入库
	DONE:
		for {
			select {
			case data := <-dataChan:
				data.proxy.Update(session, data.key, data.value)
			case <-exitChan:
				break DONE
			}
		}
	}()

	for index, proxy := range proxies {
		occupyChan <- true //获取占用权
		go func(proxy Proxy, count int) {
			proxyURL := "http://" + proxy.IP + ":" + proxy.Port
			success, msg := validHTTPBin(proxyURL)
			fmt.Println(count, "valid http-bin", proxy, success, msg)

			dataChan <- Data{
				proxy: proxy,
				key:   "success",
				value: success,
			}

			if success {
				maimai := validMaiMai(proxyURL)
				dataChan <- Data{
					proxy: proxy,
					key:   "maimai",
					value: maimai,
				}
			}

			<-occupyChan //释放占用权
		}(proxy, index+1)
	}

}

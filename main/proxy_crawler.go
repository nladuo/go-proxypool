package main

import (
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/nladuo/go-proxypool"
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
		Timeout:   time.Duration(20 * time.Second),
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

/**
 * 爬IP，需要定制
 **/
func proxyCrawler(session *mgo.Session) {
	Iteration := 60  // 提取多少轮
	BatchCount := 32 // 一次提取多少个
	dataChan := make(chan proxypool.Proxy, ConcurNum)
	occupyChan := make(chan bool, ConcurNum)
	exitChan := make(chan bool, 1)
	go func() { // 代理入库
	DONE:
		for {
			select {
			case p := <-dataChan:
				p.Insert(session)
			case <-exitChan:
				break DONE
			}
		}
	}()
	count := 1 //记录第几个代理
	for i := 0; i < Iteration; i++ {
		resp, _ := http.Get(fmt.Sprintf("http://tvp.daxiangdaili.com/ip/?tid=557647932245581&num=%d&delay=1", BatchCount))

		data, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		proxys := strings.Split(string(data), "\r\n")

		for _, proxy := range proxys { //验证代理ip
			proxyURL := "http://" + proxy
			fmt.Println(proxyURL)
			occupyChan <- true //获取占用权
			go func(proxy string, count int) {
				success, msg := validHTTPBin(proxyURL)
				fmt.Println(count, msg)

				p := proxypool.Proxy{
					IP:         strings.Split(proxy, ":")[0],
					Port:       strings.Split(proxy, ":")[1],
					CreateTime: time.Now(),
					Success:    success,
					Msg:        msg,
				}
				<-occupyChan  //释放占用权
				dataChan <- p //单线程入库
			}(proxy, count)
			count++
		}

		// <-exit // 退出
	}
	exitChan <- true
	close(occupyChan) //关闭管道
	close(dataChan)
	close(exitChan)
}

/**
 * 验证爬虫，每隔一段时间一次
 **/
func validCrawler() {

}

func main() {

	session, err := mgo.Dial("")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	proxyCrawler(session)

}

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

	"github.com/nladuo/go-proxypool"
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

	return true, string(data)
}

/**
 * 爬IP，需要定制
 **/
func proxyCrawler(session *mgo.Session) {
	iteration := 20  // 提取多少轮
	batchCount := 25 // 一次提取多少个
	ch := make(chan proxypool.Proxy, batchCount)
	exit := make(chan int, 1)
	for i := 0; i < iteration; i++ {
		resp, _ := http.Get(fmt.Sprintf("http://tvp.daxiangdaili.com/ip/?tid=557647932245581&num=%d&delay=1", batchCount))

		data, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		proxys := strings.Split(string(data), "\r\n")

		go func() { // 代理入库
			count := 0
			for {
				select {
				case p := <-ch:
					p.Insert(session)
					count++
				}
				if count == batchCount {
					exit <- 1
				}
			}
		}()

		for _, proxy := range proxys { //验证代理ip
			proxyURL := "http://" + proxy
			fmt.Println(proxyURL)
			go func(proxy string) {
				success, msg := validHTTPBin(proxyURL)
				fmt.Println(i, msg)

				p := proxypool.Proxy{
					IP:         strings.Split(proxy, ":")[0],
					Port:       strings.Split(proxy, ":")[1],
					CreateTime: time.Now(),
					Success:    success,
					Msg:        msg,
				}

				ch <- p //一个任务完成
			}(proxy)
		}

		<-exit // 退出
	}
	close(ch) //关闭管道
	close(exit)
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

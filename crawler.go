package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

/**
 * 爬IP，需要定制
 * 这个用的是大象代理的,http://www.daxiangdaili.com/
 * 需要改成你买的代理ip
 **/
func proxyCrawler(session *mgo.Session) {
	Iteration := 3
	BatchCount := 500 // 一次提取多少个
	dataChan := make(chan Proxy, ConcurNum)
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
		resp, _ := http.Get(fmt.Sprintf("http://tvp.daxiangdaili.com/ip/?tid=[你的订单编号]&num=%d&delay=3", BatchCount))

		data, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		proxys := strings.Split(string(data), "\r\n")

		for _, proxy := range proxys { //验证代理ip
			occupyChan <- true //获取占用权
			go func(proxy string, count int) {
				proxyURL := "http://" + proxy
				success, msg := validHTTPBin(proxyURL)
				fmt.Println(count, "cralwed", proxy, success, msg)

				p := Proxy{
					IP:         strings.Split(proxy, ":")[0],
					Port:       strings.Split(proxy, ":")[1],
					CreateTime: time.Now(),
					Success:    success,
					Msg:        msg,
				}

				dataChan <- p //单线程入库
				<-occupyChan  //释放占用权
			}(proxy, count)
			count++
		}
		time.Sleep(5 * time.Minute)
	}
	exitChan <- true // 退出
}

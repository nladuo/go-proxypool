package main

import (
	"time"

	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"
)

func main() {

	session, err := mgo.Dial("")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	go func() {
		//代理爬虫，可以自己定制如何爬取，如果库里面有，可以选择不跑
		proxyCrawler(session)
		for {
			validCrawler(session)
			time.Sleep(10 * time.Minute) // 10分钟校验一次
		}
	}()

	router := gin.Default()
	router.Run(":4002")
}

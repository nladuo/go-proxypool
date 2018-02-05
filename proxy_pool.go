package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/**
 * 显示结果
 **/
type Result struct {
	Ip   string
	Port string
}

func main() {

	session, err := mgo.Dial("")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	go func() {
		//代理爬虫，可以自己定制如何爬取，如果库里面有，可以选择不跑
		for {
			proxyCrawler(session)
			time.Sleep(5 * time.Minute) // 5分钟爬取一次
		}
	}()

	go func() {
		for {
			validCrawler(session, true)  // 先校验当前可用的ip
			validCrawler(session, false) // 再校验其他ip
			time.Sleep(5 * time.Minute)  // 5分钟校验一次
		}
	}()

	router := gin.Default()

	router.GET("/proxy_pool", func(c *gin.Context) {
		count := c.DefaultQuery("count", "0")
		limit, err := strconv.ParseInt(count, 10, 16)
		if err != nil {
			limit = 100
		}

		collection := session.DB("go-proxytool").C("proxy")
		proxies := []Proxy{}
		err = collection.Find(bson.M{"maimai": true}).Limit(int(limit)).All(&proxies)
		results := []Result{}
		for _, proxy := range proxies {
			results = append(results, Result{
				Ip:   proxy.IP,
				Port: proxy.Port,
			})
		}
		c.JSON(200, gin.H{
			"success": true,
			"count":   len(results),
			"proxies": results,
		})

	})
	router.Run(":4002")
}

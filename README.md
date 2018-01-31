# go-proxypool
自用的代理ip池，作为一个基础框架。
## 依赖
golang + mongodb

## 安装
``` bash 
go get github.com/nladuo/go-proxypool
cd $GOPATH/src/github.com/nladuo/go-proxypool
make prepare
```
## 运行IP爬虫
从代理ip网站上提取ip，然后去httpbin校验可用性，存到mongodb数据库中。需要自己定制如何从代理IP网站爬取的爬虫。
``` bash
go run main/proxy_crawler.go
```

## 运行代理池
对指定的网站进行ip校验可用性，并提供http提取接口。需要自己定制校验的方法。
``` bash
go run main/proxy_pool.go
```

## 提取可用ip
http://127.0.0.1:4002/ip_pool?count=100
``` json
{
    "total_count": 523,
    "ips":[
        {
            "host": "12.33.81.3",
            "port": "80"
        },
        {
            "host": "54.64.7.99",
            "port": "8079"
        },
        ...............
    ]
}
```
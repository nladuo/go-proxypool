# go-proxypool
自用的代理ip池，作为一个基础框架。
## 依赖
golang + mongodb

## 安装
``` bash 
git clone https://github.com/nladuo/go-proxypool.git
cd go-proxypool && make prepare
```
## 运行代理池
- 首先从代理ip网站上提取ip，然后去httpbin校验可用性，存到mongodb数据库中。需要自己定制如何从代理IP网站爬取的爬虫。
- 对指定的网站进行ip校验可用性，并提供http提取接口。需要自己定制校验的方法。

``` bash
make
./proxy_pool
```

## 提取可用ip
http://127.0.0.1:4002/proxy_pool?count=100
``` json
{
  	"success": true,
	"count": 95,
	"proxies": [
        {
            "Ip": "166.111.80.162",
            "Port": "3128"
        }, 
        {
            "Ip": "140.143.96.216",
            "Port": "80"
        }, 
        .............
    ]
}
```
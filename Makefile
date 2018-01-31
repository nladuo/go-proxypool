default: build

prepare:
	go get -u gopkg.in/mgo.v2
	go get -u github.com/gin-gonic/gin
	go get -u github.com/bitly/go-simplejson

build:
	go build -o proxy_pool

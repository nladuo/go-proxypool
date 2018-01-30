default: prepare

prepare:
	go get -u gopkg.in/mgo.v2
	go get -u github.com/gin-gonic/gin
	go get -u github.com/bitly/go-simplejson
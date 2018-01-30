package proxypool

import (
	"time"

	"gopkg.in/mgo.v2"
)

type Proxy struct {
	IP         string
	Port       string
	CreateTime time.Time
	Msg        string
	Success    bool
}

func (this *Proxy) Insert(session *mgo.Session) {
	c := session.DB("go-proxytool").C("proxy")
	c.Insert(this)
}

func DeleteProxy(id int) {

}

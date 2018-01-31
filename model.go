package main

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Proxy struct {
	IP         string
	Port       string
	CreateTime time.Time
	Msg        string
	Success    bool
	MaiMai     bool
}

func (this *Proxy) Insert(session *mgo.Session) {
	c := session.DB("go-proxytool").C("proxy")
	proxy := Proxy{}
	err := c.Find(bson.M{"ip": this.IP}).One(&proxy)
	if err != nil { //不存在插入
		c.Insert(this)
	}
}

func (this *Proxy) Update(session *mgo.Session, key string, value bool) {
	c := session.DB("go-proxytool").C("proxy")
	data := bson.M{"$set": bson.M{key: value}}
	err := c.Update(bson.M{"ip": this.IP}, data)
	if err != nil {
		fmt.Println(err.Error())
	}
}

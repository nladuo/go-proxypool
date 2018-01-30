package proxypool

import (
	"time"
)

type IP struct {
	Id         int
	Address    string
	Port       string
	CreateTime time.Time
}

func (*IP) Insert()  {
	
}


func DeleteIP(id int)  {
	
}
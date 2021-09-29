package main

import (
	"labix.org/v2/mgo"
	"time"
)

func main() {
    dialInfo := &mgo.DialInfo{
        Addrs:     []string{"192.168.150.129:27017"},
        Direct:    true,
        Timeout:   time.Second * 1,
        Database:  "test", // database
        Source:    "admin", // table
        Username:  "root",
        Password:  "123456",
        //PoolLimit: 4096, // Session.SetPoolLimit
    }
    session, err := mgo.DialWithInfo(dialInfo)
    if nil != err {
        panic("dial error:" + err.Error())
    }
    defer session.Close()

    err = session.Ping()
    if nil != err {
        panic("ping error:" + err.Error())
    }
}

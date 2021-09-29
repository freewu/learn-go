package main

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type User struct
{
	Name string
	Phone string
}

/**
package main

import (
    "fmt"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "time"
)

func main() {
    dialInfo := &mgo.DialInfo{
        Addrs:     []string{"192.168.6.122"},
        Direct:    false,
        Timeout:   time.Second * 1,
        Database:  "yourdatabase",
        Source:    "admin",
        Username:  "username",
        Password:  "password",
        PoolLimit: 4096, // Session.SetPoolLimit
    }
    session, err := mgo.DialWithInfo(dialInfo)
    if nil != err {
        panic(err)
    }
    defer session.Close()
}
*/

func main() {
	session, err := mgo.Dial("mongodb://192.168.150.129:27017")
	if err != nil {
		panic("connect error: " + err.Error())
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	// 这样会报  BSON field 'authenticate.nonce' is an unknown field.
	//err = session.Login(&mgo.Credential{Username: "root", Password: "123456" })
	//if err != nil {
	//	panic("auth error: " + err.Error())
	//}

	session.SetMode(mgo.Monotonic, true)
	// 选择 DB
	db := session.DB("test")
	err = db.Login("root", "123456")
	if err != nil {
		panic("login error: " + err.Error())
	}
	// 选择表
	c := db.C("people")

	// insert
	err = c.Insert(&User{"Ale", "111111"}, &User{"Cla", "222222222"})
	if err != nil {
		panic("insert error: " + err.Error())
	}

	// find one
	result := User{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		panic("find error: " + err.Error())
	}
	fmt.Println("Phone:", result.Phone)

	// find empty
	err = c.Find(bson.M{"name": "Ale1"}).One(&result)
	if err != nil {
		fmt.Printf("find empty error: %v", err.Error())
	}
	fmt.Println("Phone:", result.Phone)
}
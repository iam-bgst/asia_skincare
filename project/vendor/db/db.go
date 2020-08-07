package db

import (
	"time"

	"gopkg.in/mgo.v2"
)

var Database struct {
	Session *mgo.Session
}
var Collection = make(map[string]*mgo.Collection)

func NewConnection() {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{"127.0.0.1:27017"},
		Timeout:  10 * time.Second,
		Database: "asia_sc",
	}
	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)

	Database.Session = session
}

func SetCollection(name string) {
	if Collection[name] == nil {
		Collection[name] = Database.Session.DB("asia_sc").C(name)
	}
}

func GetCollection(name string) *mgo.Collection {
	return Database.Session.DB("asia_sc").C(name)
}

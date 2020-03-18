package db

import (
	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/mgo.v2"
	"time"
)

var mgosess *mgo.Session

const (
	DB = "ethereum-block"
)

//初始化链接
func init() {
	dialinfo := mgo.DialInfo{
		Addrs:     []string{"127.0.0.1:27017"},
		Timeout:   500 * time.Millisecond,
		Username:  "admin",
		Source:    "admin",
		Password:  "123456",
		PoolLimit: 100,
	}
	session, e := mgo.DialWithInfo(&dialinfo)
	//path, _ := cnf.Cnf.GetValue("cnf", "DbPath")
	//session, e := mgo.Dial(cnf.CnfStr.Dbpath)
	if e != nil {
		log.Debug("connect mongo error:" + e.Error())
	}
	mgosess = session
}

//copy链接选择数据库和表
func Connect(db, collection string) (*mgo.Session, *mgo.Collection) {
	ms := mgosess.Copy()
	c := ms.DB(db).C(collection)
	ms.SetMode(mgo.Monotonic, true)
	return ms, c
}

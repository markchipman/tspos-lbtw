package db

import (
	"fmt"

	"gopkg.in/mgo.v2"

	"github.com/wormling/tspos-lbtw/config"
)

var (
	// Session Database session
	Session *mgo.Session
	// Mongo Connection options
	Mongo *mgo.DialInfo
)

// Connect connects to mongodb
func Connect() {
	config.BuildDefaultConf()
	c, _ := config.LoadConfYaml("./config.yaml")
	uri := c.Core.Database.Url

	mongo, err := mgo.ParseURL(uri)
	s, err := mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		panic(err.Error())
	}
	s.SetSafe(&mgo.Safe{})
	s.SetMode(mgo.Strong, true) // [Strong|Monotonic|Eventual]

	fmt.Println("Connected to", uri)
	Session = s
	Mongo = mongo
}

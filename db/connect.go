package db

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

var (
	// Session Database session
	Session *mgo.Session
	// Mongo Connection options
	Mongo *mgo.DialInfo
)

// Connect connects to mongodb
func Connect(uri string) {
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

func EnsureIndex() {
	session := Session.Clone()
	defer session.Close()

	c := session.DB("tspos_lbtw").C("tare_weights")

	index := mgo.Index{
		//Key:        []string{"brand", "category", "name", "bottle_size, empty_weight, full_weight, image_url, created_on, updated_on"},
		Key:        []string{"brand"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
		Collation: &mgo.Collation{
			Locale:   "en",
			Strength: 2,
		},
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

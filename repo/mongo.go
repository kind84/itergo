package repo

import (
	"log"

	"gopkg.in/mgo.v2"
)

// DB = database
var DB *mgo.Database

func init() {
	s, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}

	if err = s.Ping(); err != nil {
		panic(err)
	}

	DB = s.DB("iterpro")

	log.Println("Connected to mongo")
}

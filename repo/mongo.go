package repo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB = database
// var DB *mgo.Database
var DB *mongo.Database

func init() {
	// s, err := mgo.Dial("mongodb://localhost")
	c, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://mongodb:27017"))
	if err != nil {
		panic(err)
	}

	// if err = s.Ping(); err != nil {
	// 	panic(err)
	// }

	// DB = s.DB("iterpro")

	// Check the connection
	err = c.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	DB = c.Database("iterpro")

	log.Println("Connected to mongo")
}

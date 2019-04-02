package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Review object
type Review struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	AuthorID   bson.ObjectId `json:"authorID" bson:"author_id"`
	ReviewedID bson.ObjectId `json:"reviewedID" bson:"reviewed_id"`
	Body       string        `json:"body"`
	Score      int           `json:"score"`
}

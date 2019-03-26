package models

import (
	"gopkg.in/mgo.v2/bson"
)

// User object
type User struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	Username string        `json:"username" bson:"username"`
	Password string        `json:"password" bson:"password"`
	Role     string        `json:"role" bson:"role"`
}

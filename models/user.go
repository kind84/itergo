package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User object
type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password"`
	Role     string             `json:"role" bson:"role"`
}

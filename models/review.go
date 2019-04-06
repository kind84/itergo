package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Review object
type Review struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	AuthorID   primitive.ObjectID `json:"authorID" bson:"author_id"`
	ReviewedID primitive.ObjectID `json:"reviewedID" bson:"reviewed_id"`
	Body       string             `json:"body"`
	Score      int                `json:"score"`
}

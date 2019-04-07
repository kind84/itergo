package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Review object
type Review struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	AuthorID   string             `json:"authorID" bson:"author_id"`
	ReviewedID string             `json:"reviewedID" bson:"reviewed_id"`
	Body       string             `json:"body"`
	Score      int                `json:"score"`
}

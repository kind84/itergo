package repo

import (
	"context"
	"errors"

	"github.com/kind84/iterpro/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Reviews collection
var Reviews *mongo.Collection

func getReviewsCollection() {
	Reviews = DB.Collection("reviews")
}

// AddReview creates a new review into the reviews collection
func AddReview(r *models.Review) error {
	getReviewsCollection()
	r.ID = primitive.NewObjectID()

	_, err := Reviews.InsertOne(context.TODO(), *r)
	return err
}

func GetReviews(id string) ([]models.Review, error) {
	getReviewsCollection()
	rs := []models.Review{}
	var cur *mongo.Cursor

	if id != "" {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors.New("Invalid bson ObjectId")
		}

		cur, err = Reviews.Find(context.TODO(), bson.M{"reviewed_id": oid})
		if err != nil {
			return rs, err
		}
	} else {
		var err error
		cur, err = Reviews.Find(context.TODO(), bson.M{})
		if err != nil {
			return rs, err
		}
	}

	for cur.Next(context.TODO()) {
		var el models.Review
		err := cur.Decode(&el)
		if err != nil {
			return rs, err
		}
		rs = append(rs, el)
	}
	return rs, nil
}

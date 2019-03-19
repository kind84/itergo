package repo

import (
	"errors"

	"github.com/kind84/iterpro/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Reviews collection
var Reviews *mgo.Collection

func getReviewsCollection() {
	Reviews = DB.C("reviews")
}

// AddReview creates a new review into the reviews collection
func AddReview(r *models.Review) error {
	getReviewsCollection()
	r.ID = bson.NewObjectId()

	err := Reviews.Insert(*r)
	if err != nil {
		return err
	}
	return nil
}

func GetReviews(id string) ([]models.Review, error) {
	getReviewsCollection()
	rs := []models.Review{}

	if id != "" {
		if !bson.IsObjectIdHex(id) {
			return nil, errors.New("Invalid bson ObjectId")
		}

		oid := bson.ObjectIdHex(id)

		err := Reviews.Find(bson.M{"reviewed_id": oid}).All(&rs)
		if err != nil {
			return rs, err
		}
	} else {
		err := Reviews.Find(bson.M{}).All(&rs)
		if err != nil {
			return rs, err
		}
	}
	return rs, nil
}

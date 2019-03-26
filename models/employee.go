package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Employee object
type Employee struct {
	ID               bson.ObjectId `json:"id" bson:"_id"`
	FirstName        string        `json:"firstName" bson:"firstName"`
	LastName         string        `json:"lastName" bson:"lastName"`
	Employees2Review []Employee    `json:"employees2Review" bson:"employees2Review"`
	Email            string        `json:"email"`
}

func (e Employee) sendFeedback() {

}

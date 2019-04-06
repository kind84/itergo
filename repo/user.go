package repo

import (
	"context"

	"github.com/kind84/iterpro/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var Users *mongo.Collection

func getUsersCollection() {
	Users = DB.Collection("users")
}

func SignupUser(u *models.User) error {
	getUsersCollection()
	u.ID = primitive.NewObjectID()

	_, err := Users.InsertOne(context.TODO(), *u)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(e string) (models.User, error) {
	getUsersCollection()
	u := models.User{}

	err := Users.FindOne(context.TODO(), bson.M{"username": e}).Decode(&u)
	if err != nil {
		return u, err
	}
	return u, nil
}

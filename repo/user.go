package repo

import (
	"github.com/kind84/iterpro/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var Users *mgo.Collection

func getUsersCollection() {
	Users = DB.C("users")
}

func SignupUser(u *models.User) error {
	getUsersCollection()
	u.ID = bson.NewObjectId()

	err := Users.Insert(*u)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(e string) (models.User, error) {
	getUsersCollection()
	u := models.User{}

	err := Users.Find(bson.M{"username": e}).One(&u)
	if err != nil {
		return u, err
	}
	return u, nil
}

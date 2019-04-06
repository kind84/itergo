package repo

import (
	"context"
	"errors"

	"github.com/kind84/iterpro/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Employees collection
var Employees *mongo.Collection

func getEmployeesCollection() {
	Employees = DB.Collection("employees")
}

// GetEmployee returns an employee from the db, given its ID
func GetEmployee(id string) (models.Employee, error) {
	getEmployeesCollection()
	e := models.Employee{}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return e, errors.New("Invalid bson ObjectId")
	}

	err = Employees.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&e)
	if err != nil {
		return e, err
	}
	return e, nil
}

func GetEmployeeEmail(email string) (models.Employee, error) {
	getEmployeesCollection()
	e := models.Employee{}

	err := Employees.FindOne(context.TODO(), bson.M{"email": email}).Decode(&e)
	return e, err
}

// GetEmployees return a list of employees given their IDs
func GetEmployees(ids []primitive.ObjectID) ([]models.Employee, error) {
	getEmployeesCollection()
	es := []models.Employee{}
	var cur *mongo.Cursor

	if ids != nil {
		var err error
		cur, err = Employees.Find(context.TODO(), bson.D{{"_id", bson.D{{"$in", ids}}}})
		if err != nil {
			return es, err
		}
	} else {
		var err error
		cur, err = Employees.Find(context.TODO(), bson.M{})
		if err != nil {
			return es, err
		}
	}

	for cur.Next(context.TODO()) {
		var el models.Employee
		err := cur.Decode(&el)
		if err != nil {
			return es, err
		}
		es = append(es, el)
	}
	return es, nil
}

// AddEmployee creates a new employee and returns an error or nil
func AddEmployee(e *models.Employee) error {
	getEmployeesCollection()
	e.ID = primitive.NewObjectID()

	_, err := Employees.InsertOne(context.TODO(), *e)
	return err
}

// UpdateEmployee updates the given employee and returns the employee and an error or nil
func UpdateEmployee(e *models.Employee) error {
	getEmployeesCollection()
	_, err := Employees.UpdateOne(context.TODO(), bson.M{"_id": e.ID}, e)
	return err
}

// DeleteEmployee removes an employee from the collection and returns an error or nil
func DeleteEmployee(id string) error {
	getEmployeesCollection()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("Invalid bson ObjectId")
	}

	_, err = Employees.DeleteOne(context.TODO(), bson.M{"_id": oid})
	return err
}

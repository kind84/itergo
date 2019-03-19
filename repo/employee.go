package repo

import (
	"github.com/kind84/iterpro/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Employees collection
var Employees *mgo.Collection

func getEmployeesCollection() {
	Employees = DB.C("employees")
}

// GetEmployee returns an employee from the db, given its ID
func GetEmployee(id bson.ObjectId) (models.Employee, error) {
	getEmployeesCollection()
	e := models.Employee{}

	err := Employees.FindId(id).One(&e)
	if err != nil {
		return e, err
	}
	return e, nil
}

// GetEmployees return a list of employees given their IDs
func GetEmployees(ids []bson.ObjectId) ([]models.Employee, error) {
	getEmployeesCollection()
	es := []models.Employee{}

	if ids != nil {
		err := Employees.Find(bson.D{{"_id", bson.D{{"$in", ids}}}}).All(&es)
		if err != nil {
			return es, err
		}
	} else {
		err := Employees.Find(bson.M{}).All(&es)
		if err != nil {
			return es, err
		}
	}

	return es, nil
}

// AddEmployee creates a new employee and returns an error or nil
func AddEmployee(e *models.Employee) error {
	getEmployeesCollection()
	e.ID = bson.NewObjectId()

	err := Employees.Insert(*e)
	if err != nil {
		return err
	}
	return nil
}

// UpdateEmployee updates the given employee and returns the employee and an error or nil
func UpdateEmployee(e *models.Employee) error {
	getEmployeesCollection()
	err := Employees.UpdateId(e.ID, e)
	if err != nil {
		return err
	}
	return nil
}

// DeleteEmployee removes an employee from the collection and returns an error or nil
func DeleteEmployee(id bson.ObjectId) error {
	err := Employees.RemoveId(id)
	if err != nil {
		return err
	}
	return nil
}

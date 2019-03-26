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

func GetEmployeeEmail(email string) (models.Employee, error) {
	getEmployeesCollection()
	e := models.Employee{}

	err := Employees.Find(bson.M{"email": email}).One(&e)
	return e, err
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
	return err
}

// UpdateEmployee updates the given employee and returns the employee and an error or nil
func UpdateEmployee(e *models.Employee) error {
	getEmployeesCollection()
	err := Employees.UpdateId(e.ID, e)
	return err
}

// DeleteEmployee removes an employee from the collection and returns an error or nil
func DeleteEmployee(id bson.ObjectId) error {
	getEmployeesCollection()
	err := Employees.RemoveId(id)
	return err
}

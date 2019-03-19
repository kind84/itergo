package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/kind84/iterpro/repo"

	"github.com/julienschmidt/httprouter"
	"github.com/kind84/iterpro/models"
)

// Welcome message
func Welcome(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	io.WriteString(w, "Hello world")
}

// SendFeedback sends feedback to an employee
func SendFeedback(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var r models.Review

	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	e, err := repo.GetEmployee(r.AuthorID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if b, i := toBReviewed(&e, &r); b {
		err = repo.AddReview(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		e.Employees2Review[i] = e.Employees2Review[len(e.Employees2Review)-1]
		e.Employees2Review = e.Employees2Review[:len(e.Employees2Review)-1]

		err = repo.UpdateEmployee(&e)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// AddEmployee creates a new employee
func AddEmployee(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var e models.Employee

	err := json.NewDecoder(req.Body).Decode(&e)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = repo.AddEmployee(&e)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// UpdateEmployee updates an employee
func UpdateEmployee(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var e models.Employee

	err := json.NewDecoder(req.Body).Decode(&e)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	e, err = repo.GetEmployee(e.ID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = repo.UpdateEmployee(&e)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// DeleteEmployee removes an employee given its ID
func DeleteEmployee(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(id)

	err := repo.DeleteEmployee(oid)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetEmployee returns an employee given its ID
func GetEmployee(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(id)

	e, err := repo.GetEmployee(oid)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ej, _ := json.Marshal(e)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", ej)
}

// GetEmployees returns the list of all employees
func GetEmployees(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	es, err := repo.GetEmployees(nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	esj, _ := json.Marshal(es)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", esj)
}

// Get2BReviewed returns the list of employees to be reviewed by a given employee ID
func Get2BReviewed(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	oid := bson.ObjectIdHex(id)

	e, err := repo.GetEmployee(oid)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var ids []bson.ObjectId
	for _, ee := range e.Employees2Review {
		ids = append(ids, ee.ID)
	}
	es, err := repo.GetEmployees(ids)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	esj, _ := json.Marshal(es)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", esj)
}

func GetReviews(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	rs, err := repo.GetReviews(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	rsj, _ := json.Marshal(rs)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", rsj)
}

// Set2Review sets employees that will have to review the given employee
func Set2Review(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	r := struct {
		ToReview bson.ObjectId `json:"toReview"`
		Reviewer bson.ObjectId `json:"rviewer"`
	}{}

	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	t, err := repo.GetEmployee(r.ToReview)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	e, err := repo.GetEmployee(r.Reviewer)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	e.Employees2Review = append(e.Employees2Review, t)
	err = repo.UpdateEmployee(&e)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func toBReviewed(e *models.Employee, r *models.Review) (bool, int) {
	for i, e2r := range e.Employees2Review {
		if r.ReviewedID == e2r.ID {
			return true, i
		}
	}
	return false, -1
}
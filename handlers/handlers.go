package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/julienschmidt/httprouter"

	"github.com/kind84/iterpro/models"
	"github.com/kind84/iterpro/repo"
)

// Welcome message
func Welcome(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	io.WriteString(w, "Hello world")
}

// SendFeedback sends feedback to an employee
func SendFeedback(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var r models.Review

	ts := getToken(req)
	if ts == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, err := authorize("employee", ts, false)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, err.Error())
		return
	}

	err = json.NewDecoder(req.Body).Decode(&r)
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
		log.Println(errors.New("Feedback does not match"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// AddEmployee creates a new employee
func AddEmployee(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var e models.Employee

	ts := getToken(req)
	if ts == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, err := authorize("operator", ts, false)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, err.Error())
		return
	}

	err = json.NewDecoder(req.Body).Decode(&e)
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

	ej, _ := json.Marshal(e)

	w = setHeaders(w)
	fmt.Fprintf(w, "%s\n", ej)
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

	e, err = repo.GetEmployee(e.ID.String())
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

	err := repo.DeleteEmployee(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetEmployee returns an employee given its ID
func GetEmployee(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	e, err := repo.GetEmployee(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ej, _ := json.Marshal(e)

	w = setHeaders(w)
	fmt.Fprintf(w, "%s\n", ej)
}

// GetEmployeeEmail returns an employee given its email
func GetEmployeeEmail(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	email := p.ByName("email")

	e, err := repo.GetEmployeeEmail(email)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ej, _ := json.Marshal(e)

	w = setHeaders(w)
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

	w = setHeaders(w)
	fmt.Fprintf(w, "%s\n", esj)
}

// // Get2BReviewed returns the list of employees to be reviewed by a given employee ID
// func Get2BReviewed(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
// 	id := p.ByName("id")

// 	if !bson.IsObjectIdHex(id) {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	oid := bson.ObjectIdHex(id)

// 	e, err := repo.GetEmployee(oid)
// 	if err != nil {
// 		log.Println(err)
// 		w.WriteHeader(http.StatusNotFound)
// 		return
// 	}

// 	var ids []bson.ObjectId
// 	for _, ee := range e.Employees2Review {
// 		ids = append(ids, ee.ID)
// 	}
// 	es, err := repo.GetEmployees(ids)
// 	if err != nil {
// 		log.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	esj, _ := json.Marshal(es)
// 	w = setHeaders(w)
// 	fmt.Fprintf(w, "%s\n", esj)
// }

func GetReviews(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	rs, err := repo.GetReviews(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	rsj, _ := json.Marshal(rs)
	w = setHeaders(w)
	fmt.Fprintf(w, "%s\n", rsj)
}

func Username(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var un struct {
		Username string `json:"username"`
	}

	err := json.NewDecoder(req.Body).Decode(&un)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = repo.GetUser(un.Username)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err == nil {
		log.Println(err)
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, "Username already taken")
		return
	}
}

// Set2Review sets employees that will have to review the given employee
func Set2Review(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	r := struct {
		ToReview string `json:"toReview"`
		Reviewer string `json:"reviewer"`
	}{}

	ts := getToken(req)
	if ts == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, err := authorize("operator", ts, false)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, err.Error())
		return
	}

	err = json.NewDecoder(req.Body).Decode(&r)
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

	t.Employees2Review = nil

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
		if r.ReviewedID == e2r.ID.Hex() {
			return true, i
		}
	}
	return false, -1
}

func setHeaders(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	return w
}

func getToken(req *http.Request) string {
	auth := req.Header.Get("Authorization")
	if auth == "" {
		return ""
	}

	return strings.Split(auth, "Bearer ")[1]
}

package main

import (
	"net/http"

	"github.com/kind84/iterpro/handlers"

	"github.com/julienschmidt/httprouter"
)

func main() {
	mux := httprouter.New()

	mux.GET("/", handlers.Welcome)
	mux.GET("/employee/:id", handlers.GetEmployee)
	mux.GET("/employees", handlers.GetEmployees)
	mux.GET("/tobereviewed/:id", handlers.Get2BReviewed)
	mux.GET("/reviews/:id", handlers.GetReviews)
	mux.POST("/sendfeedback", handlers.SendFeedback)
	mux.POST("/addemployee", handlers.AddEmployee)
	mux.POST("/updateemployee", handlers.UpdateEmployee)
	mux.POST("/set2review", handlers.Set2Review)

	http.ListenAndServe(":8080", mux)
}

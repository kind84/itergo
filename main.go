package main

import (
	"net/http"

	"github.com/rs/cors"

	"github.com/kind84/iterpro/handlers"

	"github.com/julienschmidt/httprouter"
)

func main() {
	mux := httprouter.New()

	mux.GET("/", handlers.Welcome)
	mux.GET("/employee/:id", handlers.GetEmployee)
	mux.GET("/email/:email", handlers.GetEmployeeEmail)
	mux.GET("/employees", handlers.GetEmployees)
	// mux.GET("/tobereviewed/:id", handlers.Get2BReviewed)
	mux.GET("/reviews/:id", handlers.GetReviews)
	mux.POST("/sendfeedback", handlers.SendFeedback)
	mux.POST("/addemployee", handlers.AddEmployee)
	mux.POST("/updateemployee", handlers.UpdateEmployee)
	mux.POST("/set2review", handlers.Set2Review)
	mux.POST("/login", handlers.Login)
	mux.POST("/signup", handlers.Signup)
	mux.POST("/username", handlers.Username)

	// c := cors.New(cors.Options{
	// 	AllowCredentials: true,
	// 	AllowedHeaders:   []string{"*"},
	// 	AllowedMethods:   []string{"GET", "POST"},
	// 	AllowedOrigins:   []string{"http://localhost"},
	// })

	handler := cors.Default().Handler(mux)

	http.ListenAndServe(":8080", handler)
}

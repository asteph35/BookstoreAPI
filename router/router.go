package router

import (
	"go-postgres/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/book/{id}", middleware.GetBook).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/book", middleware.GetAllBooks).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newbook", middleware.CreateBook).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/book/{id}", middleware.UpdateBook).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deletebook/{id}", middleware.DeleteBook).Methods("DELETE", "OPTIONS")

	return router
}

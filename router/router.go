package router

import (
	"github.com/eceberker/gamecontextdb/middleware"
	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/users/{id}", middleware.GetUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/users", middleware.GetAllUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/users/create", middleware.CreateUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/score/submit", middleware.UpdateScore).Methods("POST", "OPTIONS")

	return router
}

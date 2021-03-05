package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/eceberker/gamecontextdb/helpers"
	"github.com/eceberker/gamecontextdb/models"
	"github.com/eceberker/gamecontextdb/services"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// CreateUser creates one user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// create an empty user of type models.User
	var user models.User

	// decode the json request to user
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	user = services.InsertUser(user)

	json.NewEncoder(w).Encode(user)
}

// GetAllUser will return all the users
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// get all the users in the db
	users := services.GetAllUsersCache()

	// send all the users as response
	json.NewEncoder(w).Encode(users)
}

// GetUser will return user by Id
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//gets parameter from request
	params := mux.Vars(r)

	//str to int for id
	id, err := helpers.StringToInt64(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	//get user with id
	user, err := services.GetUserCache(id)
	if err != nil {
		fmt.Printf("Unable to get user. %v", err)
		return
	} else if user.ID == 0 {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(user)
	return
}

// UpdateScore will return userId and updated score
func UpdateScore(w http.ResponseWriter, r *http.Request) {
	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	//create an empty scoreModel
	var scoreModel models.ScoreSubmit

	// decode the json request to scoreModel
	err := json.NewDecoder(r.Body).Decode(&scoreModel)
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call SubmitNewScore user function to update user score
	user := services.SubmitNewScore(scoreModel)

	// send the response
	json.NewEncoder(w).Encode(user)
}

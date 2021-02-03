package middleware

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/eceberker/gamecontextdb/cache"
	db "github.com/eceberker/gamecontextdb/db"
	"github.com/eceberker/gamecontextdb/helpers"
	"github.com/eceberker/gamecontextdb/models"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var ctx = context.Background()

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

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

	// call insert user function and pass the user
	insertID := InsertUser(user)

	// format a response object
	res := response{
		ID:      insertID,
		Message: "User created successfully",
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// GetAllUser will return all the users
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get all the users in the db
	users, err := getAllUsers()

	if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}

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

	//get user from cache with id
	user, err := getUserCache(id)
	if err != nil {
		fmt.Printf("Unable to get user. %v", err)
	}
	json.NewEncoder(w).Encode(user)
}

// UpdateScore will return userId and updated score
func UpdateScore(w http.ResponseWriter, r *http.Request) {
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
	// call update user function to update user score
	updatedScore := SubmitNewScore(user.Points, user.ID)

	msg := fmt.Sprintf("User updated succesfully. Total rows affected %v", updatedScore)

	res := response{
		ID:      int64(user.ID),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)

}
func SubmitNewScore(scoreWorth int64, id int64) int64 {
	// create the postgres db connection
	db := db.ConnectDb()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE users SET points = points + ($1) WHERE id = ($2)`

	res, err := db.Exec(sqlStatement, scoreWorth, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	//check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	fmt.Println("Total rows affected ", rowsAffected)

	return rowsAffected

}
func InsertUser(user models.User) int64 {

	// create redis connection
	rdb := cache.RedisConnection()

	// create the postgres db connection
	db := db.ConnectDb()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning userid will return the id of the inserted user
	sqlStatement := `INSERT INTO users (name, country, points) VALUES ($1, $2, $3) RETURNING id`

	// the inserted id will store in this id
	var id int64

	// execute the sql statement
	// Scan function will save the insert id in the id
	err := db.QueryRow(sqlStatement, user.Name, user.Country, user.Points).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rdb.HSet(ctx, helpers.Int64ToString(id), []string{"display_name", user.Name, "country", user.Country, "points", helpers.Int64ToString(user.Points)})

	fmt.Printf("Inserted a single record %v", id)

	// return the inserted id
	return id
}

func getUserCache(id int64) (models.User, error) {
	// create redis connection
	rdb := cache.RedisConnection()

	// Return model
	var user models.User

	// UserId to string
	userID := helpers.Int64ToString(id)

	// Check Redis to key exists or not
	_, err := rdb.Get(ctx, userID).Result()
	if err == redis.Nil {
		user, err = getUser(id)
		if err != nil {
			fmt.Printf("Unable to retrieve user from db %v", err)
		}
		rdb.HSet(ctx, userID, []string{"display_name", user.Name, "country", user.Country, "points", strconv.Itoa(int(user.Points))})
	}

	cachedUserName, err := rdb.HGet(ctx, userID, "display_name").Result()
	if err != nil {
		fmt.Printf("Unable to retrieve user struct into redis due to: %s \n", err)
	}
	cachedUserCountry, err := rdb.HGet(ctx, userID, "country").Result()
	if err != nil {
		fmt.Printf("Unable to retrieve user struct into redis due to: %s \n", err)
	}
	cachedUserPoints, err := rdb.HGet(ctx, userID, "points").Result()
	if err != nil {
		fmt.Printf("Unable to retrieve user struct into redis due to: %s \n", err)
	}
	points, err := helpers.StringToInt64(cachedUserPoints)
	if err != nil {
		fmt.Printf("Unable to convert string to int")
	}

	user.ID = id
	user.Name = cachedUserName
	user.Country = cachedUserCountry
	user.Points = points

	return user, nil
}

func getUser(id int64) (models.User, error) {

	db := db.ConnectDb()

	defer db.Close()

	var user models.User
	sqlStatement := `SELECT * FROM users WHERE id=$1`

	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&user.ID, &user.Name, &user.Country, &user.Points)

	switch err {
	case sql.ErrNoRows:
		fmt.Println(("No rows to return"))
	case nil:
		return user, nil
	default:
		log.Fatalf("Unable to scan %v", err)
	}
	return user, err
}

func getAllUsers() ([]models.User, error) {
	// create the postgres db connection
	db := db.ConnectDb()

	// close the db connection
	defer db.Close()

	var users []models.User

	// create the select sql query
	sqlStatement := `SELECT * FROM users`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var user models.User

		// unmarshal the row object to user
		err = rows.Scan(&user.ID, &user.Name, &user.Country, &user.Points)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the user in the users slice
		users = append(users, user)

	}
	// return empty user on error
	return users, err
}

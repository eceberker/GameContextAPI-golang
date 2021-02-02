package middleware

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/eceberker/gamecontextdb/models"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var ctx = context.Background()

func redisConnection() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Could not ping redis server due to err: %s \n", err)
	}

	return rdb
}
func createConnection() *sql.DB {

	// Open the connection
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/GameContextDb?sslmode=disable")
	if err != nil {
		panic(err)
	}
	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	// return the connection
	return db
}

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
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	//get user from cache with id
	user, err := getUserCache(int64(id))
	if err != nil {
		log.Fatalf("Unable to get user", err)
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

	fmt.Println(user.ID)
	fmt.Println(user.Points)
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
	db := createConnection()

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
	redis := redisConnection()

	// create the postgres db connection
	db := createConnection()

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

	redis.HSet(ctx, strconv.Itoa(int(id)), []string{"display_name", user.Name, "country", user.Country, "points", strconv.Itoa(int(user.Points))})

	fmt.Printf("Inserted a single record %v", id)

	// return the inserted id
	return id
}

func getUserCache(id int64) (models.User, error) {
	// create redis connection
	redis := redisConnection()

	if err := redis.HGet(ctx, strconv.Itoa(int(id)), "display_name").Err(); err != nil {
		fmt.Printf("Unable to store example struct into redis due to: %s \n", err)
	}

	var user models.User
	cachedUser, _ := redis.HGet(ctx, strconv.Itoa(int(id)), "display_name").Result()

	if err := user.UnmarshalBinary([]byte(cachedUser)); err != nil {
		fmt.Printf("Unable to unmarshal data into the new user struct due to: %s \n", err)
	}

	return user, nil
}

func getUser(id int64) (models.User, error) {

	db := createConnection()

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
	db := createConnection()

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

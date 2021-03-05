package services

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/eceberker/gamecontextdb/db"
	"github.com/eceberker/gamecontextdb/models"
)

// GetUserFromDb retrieves user from Db by Id
func GetUserFromDb(id int64) models.User {

	db := db.ConnectDb()

	defer db.Close()

	var user models.User
	sqlStatement := `SELECT * FROM users WHERE id=$1`

	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&user.ID, &user.Name, &user.Country, &user.Points)

	switch err {
	case sql.ErrNoRows:
		fmt.Println(("No rows to return"))
		user.ID = 0
		return user
	case nil:
		return user
	default:
		log.Fatalf("Unable to scan %v", err)
	}
	return user
}

// InsertUser inserts given user
func InsertUser(user models.User) models.User {

	fmt.Println(user.Name)
	// create the postgres db connection
	db := db.ConnectDb()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning userid will return inserted user
	sqlStatement := `INSERT INTO users (name, country, points) VALUES ($1, $2, $3) RETURNING *`

	// execute the sql statement
	// Scan function will save the insert id in the id
	row := db.QueryRow(sqlStatement, user.Name, user.Country, user.Points)
	err := row.Scan(&user.ID, &user.Name, &user.Country, &user.Points)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows found")
	case nil:
		{
			user = InsertUserToCache(user)
			fmt.Println(user.Name)
			// return the inserted user
			return user
		}
	default:
		log.Fatalf("Unable to scan %v", err)
	}
	return user
}

// SubmitNewScore submites score to user
func SubmitNewScore(score models.ScoreSubmit) models.User {
	// create the postgres db connection
	db := db.ConnectDb()

	// close the db connection
	defer db.Close()

	var user models.User

	// create the update sql query
	sqlStatement := `UPDATE users SET points = points + ($1) WHERE id = ($2) RETURNING *`

	err := db.QueryRow(sqlStatement, score.ScoreWorth, score.UserID).Scan(&user.ID, &user.Name, &user.Country, &user.Points)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows found")
	case nil:
		{
			UpdateUserScoreCache(score)
			fmt.Println("User updated succesfully.")
			return user
		}
	default:
		log.Fatalf("Unable to scan %v", err)
	}

	return user
}

// GetAllUsers retrieves all users from Db
func GetAllUsers() ([]models.User, error) {
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

		// CACHE TO BE IMPLEMENTED
	}
	// return empty user on error
	return users, err
}

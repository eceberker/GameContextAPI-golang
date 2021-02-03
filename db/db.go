package db

import (
	"database/sql"
	"fmt"
)

// ConnectDb creates database connection
func ConnectDb() *sql.DB {

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

package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

// ConnectDb creates database connection
func ConnectDb() *sql.DB {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	envs, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Unable to read .env file")
	}
	dbname := envs["DB_NAME"]
	dbsource := envs["DB_SOURCE"]

	// Open the connection
	db, err := sql.Open(dbname, dbsource)

	if err != nil {
		panic(err)
	}
	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to Db!")
	// return the connection
	return db
}

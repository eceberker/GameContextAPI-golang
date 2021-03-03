package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
)

// ConnectDb creates database connection
func ConnectDb() *sql.DB {

	fmt.Println("Initializing database connection . . .")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	envs, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Unable to read .env file")
	}
	dbDriver := envs["POSTGRES_DRIVER"]
	dbName := envs["POSTGRES_DB"]
	dbUser := envs["POSTGRES_USER"]
	dbPass := envs["POSTGRES_PASSWORD"]
	dbHost := envs["POSTGRES_HOST"]

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbName)

	var db *sql.DB
	var er error

	tries := 10
	for tries > 0 {
		// // Open the connection
		db, er = sql.Open(dbDriver, dsn)
		if err != nil {
			fmt.Printf("Unable to Open DB: %s... Retrying\n", er.Error())
			time.Sleep(time.Second * 2)
			tries--
		} else if er = db.Ping(); er != nil {
			// check the connection
			fmt.Printf("Unable to Ping DB: %s... Retrying\n", er.Error())
			time.Sleep(time.Second * 2)
		} else {
			er = nil
			break
		}
	}
	if er != nil {
		db = nil
		return db
	}
	fmt.Println("Connection to Database successful")
	return db
}

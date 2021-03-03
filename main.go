package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eceberker/gamecontextdb/router"
	"github.com/joho/godotenv"
)

func main() {

	r := router.Router()

	fmt.Println("Starting server on the port 8080...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	envs, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Unable to read .env file")
	}
	dbName := envs["POSTGRES_DB"]
	dbUser := envs["POSTGRES_USER"]
	dbPass := envs["POSTGRES_PASSWORD"]
	dbHost := envs["POSTGRES_HOST"]

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbName)
	fmt.Println(dsn)

	log.Fatal(http.ListenAndServe(":8080", r))

}

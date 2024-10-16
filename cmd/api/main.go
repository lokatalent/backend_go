package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/lokatalent/backend_go/cmd/api/util"
)

func main() {
	// load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %v\n", err)
	}

	// setup server configuration
	config := util.Config{}
	if err := config.Load(); err != nil {
		log.Fatalf("error parsing configuration variables: %v\n", err)
	}

	// open database connection
	db, err := openDB(&config)
	if err != nil {
		log.Fatalf("database connection error: %v\n", err)
	}
	defer db.Close()
	log.Println("database connection established")

	// start server
	if err := serveApp(&config, db); err != nil {
		log.Fatalf("error starting server: %v\n", err)
	}
}

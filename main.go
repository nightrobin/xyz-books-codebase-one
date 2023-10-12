package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"xyz-books/dbmigration"
	"xyz-books/router"
	
	"github.com/joho/godotenv"
)


func main() {
	// Load Environment Variables
	ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    exPath := filepath.Dir(ex)

	err = godotenv.Load(exPath + "/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Apply Migrations
	dbmigration.ApplyMigrations()

	router.GetRouter().Run(fmt.Sprintf("%s:%s", os.Getenv("API_HOST"), os.Getenv("API_PORT")))
}
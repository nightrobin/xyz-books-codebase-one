package main

import (
	// "fmt"
	"os"
	"log"
	"path/filepath"

	"github.com/joho/godotenv"

	"xyz-books/dbmigration"
)

func main() {
	ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    exPath := filepath.Dir(ex)

	err = godotenv.Load(exPath + "/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// fmt.Println(os.Getenv("DB_HOST"))
	// fmt.Println(os.Getenv("DB_PORT"))
	// fmt.Println(os.Getenv("DB_NAME"))
	// fmt.Println(os.Getenv("DB_USER"))
	// fmt.Println(os.Getenv("DB_PASS"))
	dbmigration.ApplyMigrations()
	
}
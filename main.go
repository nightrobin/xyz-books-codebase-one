package main

import (
	"fmt"
    // "bytes"
	// "html/template"
	"os"
	"log"
	// "net/http"
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

	// Index.html
	// parsedIndexTemplate, err := template.ParseFiles(exPath + "/templates/index.html")
    // if err != nil {
    //     log.Fatal(err)
    // }

	// tmpl := template.Must(parsedIndexTemplate, err)

	// type PageData struct {
	// 	Title string
	// 	Content string
	// }

    // http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    //     data := PageData{
    //         Title:   "Page Title",
    //         Content: "Hello, World!",
    //     }

    //     // Set the Content-Type header to indicate that we are sending HTML.
    //     w.Header().Set("Content-Type", "text/html; charset=utf-8")

    //     // Execute the template and send the result to the client.
    //     if err := tmpl.Execute(w, data); err != nil {
    //         http.Error(w, err.Error(), http.StatusInternalServerError)
    //     }
    // })


	// log.Print("Listening on :3000...")
	// err = http.ListenAndServe(":3000", nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	router.GetRouter().Run(fmt.Sprintf("%s:%s", os.Getenv("API_HOST"), os.Getenv("API_PORT")))

}
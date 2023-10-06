package method

import(
	// "fmt"
	// "encoding/json"
	"html/template"
	"net/http"
	"log"
	"sync"
	
	"xyz-books/model"


	"github.com/gin-gonic/gin"
)

type BookDisplay struct {
	ID				uint64
	Title			string
	Author			string
	Isbn13			string	`gorm:"column:isbn_13"`
	Isbn10			string	`gorm:"column:isbn_10"`
	PublicationYear	int16
	PublisherName	string
	Edition			string
	ListPrice		float32
	ImageURL		string
}

func BookIndex(c *gin.Context) {

	var wg sync.WaitGroup
	
	wg.Add(1)

	go func () {
		defer wg.Done()
		
		var books []BookDisplay

		Db.Table("books b").Select("b.id", "b.title", "GROUP_CONCAT(' ', CONCAT(a.first_name, ' ', IFNULL(a.middle_name, ''), ' ', a.last_name)) author", "b.isbn_13", "b.isbn_10", "b.publication_year", "p.name publisher_name", "b.edition", "b.list_price", "b.image_url").Joins("INNER JOIN book_authors ba ON b.id = ba.book_id").Joins("INNER JOIN authors a ON ba.author_id = a.id").Joins("INNER JOIN publishers p ON b.publisher_id = p.id").Group("b.id").Find(&books)

		type PageData struct {
			Books []BookDisplay
		}
		
		var data PageData
		data.Books = books

		w := c.Writer

		parsedIndexTemplate, err := template.ParseFiles(ExPath + "/templates/index.html")
		if err != nil {
			log.Fatal(err)
		}

		tmpl := template.Must(parsedIndexTemplate, err)
		
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}()
	
	wg.Wait()

	return
}

func AddBook(c *gin.Context) {
	data := model.Book{ID: 1, Title: "Book 1"}
	response := model.Response[model.Book]{
		Message: "Successfully added a book",
		Data:    data,
	}

	c.IndentedJSON(http.StatusOK, response)
	return
}

func UIViewBook(c *gin.Context) {
	isbn_13 := c.Param("isbn_13")

	var wg sync.WaitGroup
	
	wg.Add(1)

	go func () {
		defer wg.Done()
		
		var book BookDisplay

		Db.Table("books b").Select("b.id", "b.title", "GROUP_CONCAT(' ', CONCAT(a.first_name, ' ', IFNULL(a.middle_name, ''), ' ', a.last_name)) author", "b.isbn_13", "b.isbn_10", "b.publication_year", "p.name publisher_name", "b.edition", "b.list_price", "b.image_url").Joins("INNER JOIN book_authors ba ON b.id = ba.book_id").Joins("INNER JOIN authors a ON ba.author_id = a.id").Joins("INNER JOIN publishers p ON b.publisher_id = p.id").Where("b.isbn_13 = ?", isbn_13).Group("b.id").First(&book)

		w := c.Writer

		parsedIndexTemplate, err := template.ParseFiles(ExPath + "/templates/books/display.html")
		if err != nil {
			log.Fatal(err)
		}

		tmpl := template.Must(parsedIndexTemplate, err)
		
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := tmpl.Execute(w, book); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}()
	
	wg.Wait()

	return
}

func ReadBook(c *gin.Context) {
}

func UpdateBook() {
}

func DeleteBook() {
}
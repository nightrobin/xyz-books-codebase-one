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

func BookIndex(c *gin.Context) {

	var wg sync.WaitGroup
	
	wg.Add(1)

	go func () {
		defer wg.Done()
		
		var books []model.Book
		Db.Table("books").Find(&books)
		
		// booksJson, err := json.Marshal(books)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// 	return
		// }

		// booksJsonStr := string(booksJson)

		type PageData struct {
			Books []model.Book
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

func ReadBook() {
}

func UpdateBook() {
}

func DeleteBook() {
}
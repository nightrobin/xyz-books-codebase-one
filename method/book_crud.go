package method

import(
	// "fmt"
	"net/http"
	// "encoding/json"
	
	"xyz-books/model"

	"github.com/gin-gonic/gin"
)

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
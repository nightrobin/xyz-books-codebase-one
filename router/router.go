package router

import (
	// "fmt"
	// "os"

	// "path/filepath"
	// "time"

	// "github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	
	"xyz-books/method"
)

func GetRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", method.BookIndex)
	router.POST("/books", method.AddBook)
	router.GET("/books-ui/view/:isbn_13", method.UIViewBook)
	// router.GET("/books-ui/update/:isbn_13", method.UIUpdateBook)
	// router.GET("/books-ui/delete/:isbn_13", method.UIDeleteBook)

	// router.GET("/books/display/{isbn_13:[0-9a-zA-Z]+}", method.DisplayBook)
	// router.GET("/books-ui/{isbn_13:[0-9]+}", method.DisplayBook)

	return router
}

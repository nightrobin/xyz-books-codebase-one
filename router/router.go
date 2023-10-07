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

	router.GET("/", method.UIBookIndex)
	

	authorsUI := router.Group("/authors-ui")
	{
		authorsUI.GET("/add-form", method.UIAddAuthorForm)
		authorsUI.POST("/submit-add-form", method.UISubmitAddAuthorForm)
	}

	booksUI := router.Group("/books-ui")
	{
		booksUI.GET("/add-form", method.UIAddBookForm)
		booksUI.POST("/submit-add-form", method.UISubmitAddBookForm)
		booksUI.GET("/view/:isbn_13", method.UIViewBook)
	}

	publishersUI := router.Group("/publishers-ui")
	{
		publishersUI.GET("/add-form", method.UIAddPublisherForm)
		publishersUI.POST("/submit-form", method.UISubmitAddPublisherForm)
	}
	// router.GET("/books-ui/update/:isbn_13", method.UIUpdateBook)
	// router.GET("/books-ui/delete/:isbn_13", method.UIDeleteBook)

	// router.GET("/books/display/{isbn_13:[0-9a-zA-Z]+}", method.DisplayBook)
	// router.GET("/books-ui/{isbn_13:[0-9]+}", method.DisplayBook)

	return router
}

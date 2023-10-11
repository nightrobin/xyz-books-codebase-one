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
	router.GET("/insert-records", method.UIAddRecords) // Use to add huge number of records. For testing uses.
	
	authorsUI := router.Group("/ui/authors")
	{
		authorsUI.GET("/", method.UIAuthorIndex)
		authorsUI.GET("/add-form", method.UIAddAuthorForm)
		authorsUI.POST("/submit-add-form", method.UISubmitAddAuthorForm)
		authorsUI.GET("/update-form/:id", method.UIUpdateAuthorForm)
		authorsUI.POST("/submit-update-form/:id", method.UISubmitUpdateAuthorForm)
		authorsUI.GET("/view/:id", method.UIViewAuthor)		
		authorsUI.GET("/delete/:id", method.UIDeleteAuthor)
	}

	authorsAPI := router.Group("/api/authors")
	{
		authorsAPI.GET("/", method.GetAuthors)
		authorsAPI.GET("/:id", method.GetAuthor)
		// authorsUI.POST("/", method.AddAuthor)
		// authorsUI.PATCH("/", method.UpdateAuthor)
		// authorsUI.DELETE("/:id", method.DeleteAuthor)		
	}
	
	booksUI := router.Group("/books/ui")
	{
		booksUI.GET("/add-form", method.UIAddBookForm)
		booksUI.POST("/submit-add-form", method.UISubmitAddBookForm)
		booksUI.GET("/update-form/:isbn_13", method.UIUpdateBookForm)
		booksUI.POST("/submit-update-form/:isbn_13", method.UISubmitUpdateBookForm)
		booksUI.GET("/view/:isbn_13", method.UIViewBook)
		booksUI.GET("/delete/:isbn_13", method.UIDeleteBook)
	}

	publishersUI := router.Group("/publishers/ui")
	{
		publishersUI.GET("/", method.UIPublisherIndex)
		publishersUI.GET("/add-form", method.UIAddPublisherForm)
		publishersUI.POST("/submit-add-form", method.UISubmitAddPublisherForm)
		publishersUI.GET("/update-form/:id", method.UIUpdatePublisherForm)
		publishersUI.POST("/submit-update-form/:id", method.UISubmitUpdatePublisherForm)
		publishersUI.GET("/view/:id", method.UIViewPublisher)
		publishersUI.GET("/delete/:id", method.UIDeletePublisher)
	}
	// router.GET("/books/ui/update/:isbn_13", method.UIUpdateBook)
	// router.GET("/books/ui/delete/:isbn_13", method.UIDeleteBook)

	// router.GET("/books/display/{isbn_13:[0-9a-zA-Z]+}", method.DisplayBook)
	// router.GET("/books/ui/{isbn_13:[0-9]+}", method.DisplayBook)

	return router
}

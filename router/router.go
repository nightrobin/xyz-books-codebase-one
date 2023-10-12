package router

import (
	"xyz-books-codebase-one/method"

	"github.com/gin-gonic/gin"
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

	booksUI := router.Group("/ui/books")
	{
		booksUI.GET("/", method.UIBookIndex)
		booksUI.GET("/add-form", method.UIAddBookForm)
		booksUI.POST("/submit-add-form", method.UISubmitAddBookForm)
		booksUI.GET("/update-form/:isbn_13", method.UIUpdateBookForm)
		booksUI.POST("/submit-update-form/:isbn_13", method.UISubmitUpdateBookForm)
		booksUI.GET("/view/:isbn_13", method.UIViewBook)
		booksUI.GET("/delete/:isbn_13", method.UIDeleteBook)
	}

	publishersUI := router.Group("/ui/publishers")
	{
		publishersUI.GET("/", method.UIPublisherIndex)
		publishersUI.GET("/add-form", method.UIAddPublisherForm)
		publishersUI.POST("/submit-add-form", method.UISubmitAddPublisherForm)
		publishersUI.GET("/update-form/:id", method.UIUpdatePublisherForm)
		publishersUI.POST("/submit-update-form/:id", method.UISubmitUpdatePublisherForm)
		publishersUI.GET("/view/:id", method.UIViewPublisher)
		publishersUI.GET("/delete/:id", method.UIDeletePublisher)
	}

	authorsAPI := router.Group("/api/authors")
	{
		authorsAPI.GET("/", method.GetAuthors)
		authorsAPI.GET("/:id", method.GetAuthor)
		authorsAPI.POST("/", method.AddAuthor)
		authorsAPI.PATCH("/:id", method.UpdateAuthor)
		authorsAPI.DELETE("/:id", method.DeleteAuthor)		
	}

	booksAPI := router.Group("/api/books")
	{
		booksAPI.GET("/", method.GetBooks)
		booksAPI.GET("/:isbn_13", method.GetBook)
		booksAPI.POST("/", method.AddBook)
		booksAPI.PATCH("/:isbn_13", method.UpdateBook)
		booksAPI.DELETE("/:isbn_13", method.DeleteBook)		
	}

	publishersAPI := router.Group("/api/publishers")
	{
		publishersAPI.GET("/", method.GetPublishers)
		publishersAPI.GET("/:id", method.GetPublisher)
		publishersAPI.POST("/", method.AddPublisher)
		publishersAPI.PATCH("/:id", method.UpdatePublisher)
		publishersAPI.DELETE("/:id", method.DeletePublisher)		
	}

	return router
}

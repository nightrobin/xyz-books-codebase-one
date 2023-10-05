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

	return router
}

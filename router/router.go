package router

import (
	// "fmt"
	// "os"

	// "path/filepath"
	// "time"

	// "github.com/joho/godotenv"

	"xyz-books/method"
	"github.com/gin-gonic/gin"

)

func GetRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/books", method.AddBook)

	return router
}

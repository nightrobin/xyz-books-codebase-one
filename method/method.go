package method

import (
	// "fmt"
	"log"
	"os"
	"path/filepath"

	"xyz-books/orm"
	
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var Db *gorm.DB
var ExPath string

func init() {
	Db = orm.ConnectToDB()

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	
	ExPath = filepath.Dir(ex)

	err = godotenv.Load(ExPath + "/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

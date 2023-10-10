package method

import (
	// "fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"xyz-books/orm"
	
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var Db *gorm.DB
var ExPath string

var recordLimitPerPage int
var countShownPageNumber int

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

	recordLimitPerPage, _ = strconv.Atoi(os.Getenv("RECORD_LIMIT_PER_PAGE"))
	countShownPageNumber, _ = strconv.Atoi(os.Getenv("COUNT_SHOWN_PAGE_NUMBER"))
}

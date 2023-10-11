package method

import (
	// "fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"xyz-books/orm"
	
	"github.com/joho/godotenv"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var Db *gorm.DB
var ExPath string
var Validate *validator.Validate

var recordLimitPerPage int
var countShownPageNumber int


func init() {
	Db = orm.ConnectToDB()
	Validate = validator.New()

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

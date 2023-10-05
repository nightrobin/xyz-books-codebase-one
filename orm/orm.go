package orm

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init(){
	ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    exPath := filepath.Dir(ex)

	err = godotenv.Load(exPath + "/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func ConnectToDB() *gorm.DB {
	var config string = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=True", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: config,
		DefaultStringSize: 256,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func GetDB() *gorm.DB {
	return db
}

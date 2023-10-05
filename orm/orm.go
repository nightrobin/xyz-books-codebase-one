package orm

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectToDB() *gorm.DB {
	var config string = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=True", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               config,
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

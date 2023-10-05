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
	var config string = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=True", os.Getenv("G12PH_DB_USER"), os.Getenv("G12PH_DB_PASS"), os.Getenv("G12PH_DB_HOST"), os.Getenv("G12PH_DB_PORT"), os.Getenv("G12PH_DB_NAME"))
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

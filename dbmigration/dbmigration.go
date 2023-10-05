package dbmigration

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func ApplyMigrations(){
	ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    exPath := filepath.Dir(ex)
	migrationsDirectory := "file://" + exPath + "/migrations"

	var config string = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=True&multiStatements=true", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	db, _ := sql.Open("mysql", config)
    driver, _ := mysql.WithInstance(db, &mysql.Config{})

	m, _ := migrate.NewWithDatabaseInstance(
		migrationsDirectory,
        "mysql", 
        driver,
    )

	m.Up()
}
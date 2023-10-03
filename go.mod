module xyz-books/main

go 1.21.1

require (
	github.com/joho/godotenv v1.5.1
	xyz-books/dbmigration v0.0.0-00010101000000-000000000000
)

require (
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/golang-migrate/migrate/v4 v4.16.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
)

replace xyz-books/dbmigration => ./dbmigration

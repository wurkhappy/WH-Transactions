package DB

import (
	"database/sql"
	_ "github.com/bmizerany/pq"
	// "log"
)

var DB *sql.DB
var Name string = "wurkhappy"

func init() {
	Setup()
	CreateStatements()
}

func Setup() {
	var err error
	DB, err = sql.Open("postgres", "user=postgres dbname="+Name+" sslmode=disable")
	if err != nil {
		panic(err)
	}
}

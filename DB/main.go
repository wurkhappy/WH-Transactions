package DB

import (
	"database/sql"
	_ "github.com/bmizerany/pq"
	// "log"
)

var DB *sql.DB
var Name string = "wurkhappy"

func Setup(production bool) {
	Connect(production)
	CreateStatements()
}

func Connect(production bool) {
	var err error
	if production {
		DB, err = sql.Open("postgres", "user=wurkhappy password=whcollab dbname="+Name+" sslmode=disable")
	} else {
		DB, err = sql.Open("postgres", "user=postgres dbname="+Name+" sslmode=disable")
	}
	if err != nil {
		panic(err)
	}
}

func Close() {
	UpsertTransaction.Close()
	FindTransactionByID.Close()
	FindTransactionByPaymentID.Close()

	DB.Close()
}

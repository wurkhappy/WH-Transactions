package DB

import (
	"database/sql"
	_ "github.com/bmizerany/pq"
	// "log"
)

var UpsertTransaction *sql.Stmt
var FindTransactionByID *sql.Stmt
var FindTransactionByPaymentID *sql.Stmt

func CreateStatements() {
	var err error
	UpsertTransaction, err = DB.Prepare("SELECT upsert_transaction($1, $2)")
	if err != nil {
		panic(err)
	}

	FindTransactionByID, err = DB.Prepare("SELECT data FROM transaction WHERE id = $1")
	if err != nil {
		panic(err)
	}

	FindTransactionByPaymentID, err = DB.Prepare("SELECT data FROM transaction WHERE data->>'paymentID' = $1")
	if err != nil {
		panic(err)
	}
}

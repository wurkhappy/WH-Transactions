package models

import (
	"testing"
)

func test_ConvertToCredit(t *testing.T){
	transaction := new(Transaction)
	transaction.ID = "123"
	transaction.Amount = 100
	credit := transaction.ConvertToCredit()
	if credit.Amount != int(transaction.Amount) * 100 {
		t.Error("wrong amount set on credit")
	}
	if credit.AppearsOnStatementAs != "Wurk Happy" {
		t.Error("wrong appears as on credit")
	}
	if credit.Meta["id"] != transaction.ID {
		t.Error("metadata not set correctly")
	}
}
package models

import (
	"testing"
)

func test_ConvertToCredit(t *testing.T) {
	transaction := new(Transaction)
	transaction.ID = "123"
	transaction.Amount = 100
	credit := transaction.ConvertToCredit()
	if credit.Amount >= int(transaction.Amount)*100 {
		t.Error("fee not added to credit")
	}
	if credit.AppearsOnStatementAs != "Wurk Happy" {
		t.Error("wrong appears as on credit")
	}
	if credit.Meta["id"] != transaction.ID {
		t.Error("metadata not set correctly")
	}
}

func test_CalculateFee(t *testing.T) {
	trans := new(Transaction)
	trans.Amount = 100
	trans.PaymentType = "CardBalanced"
	fee1 := trans.CalculateFee()
	if fee1 != 13.23 {
		t.Error("wrong wurkhappy fee", fee1)
	}

	trans.PaymentType = "BankBalanced"
	fee2 := trans.CalculateFee()
	if fee2 != 11.3 {
		t.Error("wrong bank fee returned", fee2)
	}
}

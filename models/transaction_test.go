package models

import (
	"testing"
)

func test_ConvertToCredit(t *testing.T){
	transaction := new(Transaction)
	transaction.ID = "123"
	transaction.Amount = 100
	credit := transaction.ConvertToCredit()
	if credit.Amount >= int(transaction.Amount) * 100 {
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
	fee1 := trans.CalculateFee()
	if fee1 > 50 || fee1 <= 0 {
		t.Error("wrong wurkhappy fee")
	}

	trans.PaymentType = "BankBalanced"
	fee2 := trans.CalculateFee()
	if fee2 != 10 {
		t.Error("wrong bank fee returned")
	}

	trans.PaymentType = "CardBalanced"
	fee3 := trans.CalculateFee()
	if fee3 != (100*0.029 + 0.33 + 5) {
		t.Error("wrong bank fee returned")
	}
}
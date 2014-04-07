package models

import (
	"encoding/json"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/WH-Config"
	"github.com/wurkhappy/WH-Transactions/DB"
	"github.com/wurkhappy/balanced-go"
	"log"
	"time"
)

type Transaction struct {
	ID                     string    `json:"id"`
	DebitSourceBalancedID  string    `json:"debitSourceURI"`
	DebitSourceID          string    `json:"debitSourceID"`
	ClientID               string    `json:"clientID"`
	FreelancerID           string    `json:"freelancerID"`
	AgreementID            string    `json:"agreementID"`
	PaymentID              string    `json:"paymentID"`
	Amount                 float64   `json:"amount"`
	CreditSourceBalancedID string    `json:"creditSourceURI"`
	CreditSourceID         string    `json:"creditSourceID"`
	PaymentType            string    `json:"paymentType"`
	DebitURI               string    `json:"debitURI"`
	CreditURI              string    `json:"creditURI"`
	DebitDate              time.Time `json:"debitDate"`
	CreditDate             time.Time `json:"creditDate"`
}

var BalancedCardType string = "CardBalanced"
var BalancedBankType string = "BankBalanced"

func NewTransaction() *Transaction {
	id, _ := uuid.NewV4()

	return &Transaction{
		ID: id.String(),
	}

}

func NewTransactionFromRequest(m map[string]interface{}) *Transaction {
	id, _ := uuid.NewV4()

	return &Transaction{
		ID:             id.String(),
		AgreementID:    m["agreementID"].(string),
		Amount:         m["amount"].(float64),
		ClientID:       m["clientID"].(string),
		CreditSourceID: m["creditSourceID"].(string),
		FreelancerID:   m["freelancerID"].(string),
		PaymentID:      m["paymentID"].(string),
	}

}

func (t *Transaction) Save() error {
	jsonByte, _ := json.Marshal(t)
	r, err := DB.UpsertTransaction.Query(t.ID, string(jsonByte))
	r.Close()
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func FindTransactionByID(id string) (t *Transaction, err error) {
	var s string
	err = DB.FindTransactionByID.QueryRow(id).Scan(&s)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(s), &t)
	return t, nil
}

func FindTransactionByPaymentID(id string) (t *Transaction, err error) {
	var s string
	err = DB.FindTransactionByPaymentID.QueryRow(id).Scan(&s)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(s), &t)
	return t, nil
}

func (t *Transaction) CreateBankAccount() *balanced.BankAccount {
	bank_account := new(balanced.BankAccount)
	bank_account.ID = t.CreditSourceBalancedID
	return bank_account
}

func (t *Transaction) ConvertToDebit() *balanced.Debit {
	debit := new(balanced.Debit)
	debit.Amount = int(t.Amount * 100)
	debit.AppearsOnStatementAs = "Wurk Happy"
	bankAccount := new(balanced.BankAccount)
	creditCard := new(balanced.Card)
	if t.PaymentType == BalancedCardType {
		creditCard.ID = t.DebitSourceBalancedID
		debit.Owner = creditCard
	} else {
		bankAccount.ID = t.DebitSourceBalancedID
		debit.Owner = bankAccount
	}
	debit.Meta = map[string]string{
		"id":              t.ID,
		"creditSourceURI": t.CreditSourceBalancedID,
	}
	return debit
}

func (t *Transaction) ConvertToCredit() *balanced.Credit {
	credit := new(balanced.Credit)
	fee := t.CalculateFee()
	credit.Amount = int((t.Amount - fee) * 100)
	credit.AppearsOnStatementAs = "Wurk Happy"
	credit.Meta = map[string]string{
		"id": t.ID,
	}
	return credit
}

func (t *Transaction) CalculateFee() float64 {
	var fee float64
	if t.PaymentType == "CardBalanced" {
		fee = (t.Amount * 0.029) + 0.30 + 0.25
	} else if t.PaymentType == "BankBalanced" {
		fee = (t.Amount * 0.01) + 0.30 + 0.25
		if fee > 5 {
			fee = 5
		}
	}
	whFee := calculateWurkHappyFee(t.Amount)
	fee = fee + whFee
	return fee
}

func (t *Transaction) GetCreditSourceInfo() error {
	resp, statusCode := sendServiceRequest("GET", config.PaymentInfoService, "/user/"+t.FreelancerID+"/bank_account/"+t.CreditSourceID+"/uri", nil)
	if statusCode > 400 {
		return fmt.Errorf("Could not find source")
	}
	var m map[string]interface{}
	json.Unmarshal(resp, &m)
	t.CreditSourceBalancedID = m["balanced_id"].(string)
	return nil
}

func (t *Transaction) GetDebitSourceInfo() error {
	var path string
	if t.PaymentType == BalancedBankType {
		path = "/user/" + t.ClientID + "/bank_account/" + t.DebitSourceID + "/uri"
	} else {
		path = "/user/" + t.ClientID + "/cards/" + t.DebitSourceID + "/uri"
	}
	resp, statusCode := sendServiceRequest("GET", config.PaymentInfoService, path, nil)
	if statusCode > 400 {
		return fmt.Errorf("Could not find source")
	}
	var m map[string]interface{}
	json.Unmarshal(resp, &m)
	t.DebitSourceBalancedID = m["balanced_id"].(string)
	return nil
}

func (t *Transaction) CreateDebit() (string, error) {
	debit := t.ConvertToDebit()
	bErrors := balanced.Create(debit)
	return debit.ID, formatBalancedErrors(bErrors)
}

func calculateWurkHappyFee(amount float64) float64 {
	var whFee float64 = amount * 0.01
	if amount > 10 {
		whFee = 10
	}
	return whFee
}

package models

import (
	"encoding/json"
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/Balanced-go"
	"github.com/wurkhappy/WH-Transactions/DB"
	"log"
)

type Transaction struct {
	ID              string  `json:"id"`
	DebitSourceURI  string  `json:"debitSourceURI"`
	ClientID        string  `json:"clientID"`
	FreelancerID    string  `json:"freelancerID"`
	AgreementID     string  `json:"agreementID"`
	PaymentID       string  `json:"paymentID"`
	Amount          float64 `json:"amount"`
	CreditSourceURI string  `json:"creditSourceURI"`
	PaymentType     string  `json:"paymentType"`
	DebitURI        string  `json:"debitURI"`
	CreditURI       string  `json:"creditURI"`
}

func NewTransactionFromRequest(m map[string]interface{}) *Transaction {
	id, _ := uuid.NewV4()

	return &Transaction{
		ID:              id.String(),
		AgreementID:     m["agreementID"].(string),
		Amount:          m["amount"].(float64),
		ClientID:        m["clientID"].(string),
		CreditSourceURI: m["creditSourceURI"].(string),
		FreelancerID:    m["freelancerID"].(string),
		PaymentID:       m["paymentID"].(string),
	}

}

func (t *Transaction) Save() error {
	jsonByte, _ := json.Marshal(t)
	_, err := DB.UpsertTransaction.Query(t.ID, string(jsonByte))
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
	bank_account.CreditsURI = t.CreditSourceURI
	return bank_account
}

func (t *Transaction) ConvertToDebit() *balanced.Debit {
	debit := new(balanced.Debit)
	debit.Amount = int(t.Amount) * 100
	debit.AppearsOnStatementAs = "Wurk Happy"
	debit.SourceUri = t.DebitSourceURI
	debit.Meta = map[string]string{
		"id":              t.ID,
		"creditSourceURI": t.CreditSourceURI,
	}
	return debit
}

func (t *Transaction) ConvertToCredit() *balanced.Credit {
	credit := new(balanced.Credit)
	fee := t.CalculateFee()
	credit.Amount = int(t.Amount-fee) * 100
	credit.AppearsOnStatementAs = "Wurk Happy"
	credit.Meta = map[string]string{
		"id": t.ID,
	}
	return credit
}

func (t *Transaction) CalculateFee() float64 {
	var fee float64
	if t.PaymentType == "CardBalanced" {
		fee = (t.Amount * 0.029) + 0.33
	} else if t.PaymentType == "BankBalanced" {
		fee = (t.Amount * 0.01) + 0.30
		if fee > 5 {
			fee = 5
		}
	}
	whFee := calculateWurkHappyFee(t.Amount)
	fee = fee + whFee
	return fee
}

func calculateWurkHappyFee(amount float64) float64 {
	var whFee float64
	if amount >= 1000 {
		whFee = 50
	} else if amount >= 500 {
		whFee = 25
	} else if amount >= 100 {
		whFee = 10
	} else if amount >= 10 {
		whFee = 5
	}
	return whFee
}

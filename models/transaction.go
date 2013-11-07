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

func FindTransactionByPaymentID(id string) (t *Transaction, err error) {
	var s string
	err = DB.FindTransactionByID.QueryRow(id).Scan(&s)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(s), &t)
	return t, nil
}

func (t *Transaction) ConvertToDebit() *balanced.Debit {
	debit := new(balanced.Debit)
	debit.Amount = int(t.Amount)
	debit.AppearsOnStatementAs = "Wurk Happy"
	debit.SourceUri = t.DebitSourceURI
	debit.Meta = map[string]string{
		"id":              t.ID,
		"creditSourceURI": t.CreditSourceURI,
	}
	return debit
}

package models

import (
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/Balanced-go"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

type Transaction struct {
	ID              string  `json:"id" bson:"_id"`
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

func (t *Transaction) SaveToDB(db *mgo.Database) error {
	coll := db.C("transactions")
	if _, err := coll.UpsertId(t.ID, &t); err != nil {
		log.Printf("db error is: %v", err)
		return err
	}
	return nil
}

func FindTransactionByPaymentID(id string, db *mgo.Database) (t *Transaction, err error) {
	err = db.C("transactions").Find(bson.M{"paymentid": id}).One(&t)
	if err != nil {
		return nil, err
	}

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

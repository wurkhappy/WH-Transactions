package handlers

import (
	"github.com/wurkhappy/WH-Transactions/models"
	// "github.com/streadway/amqp"
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/Balanced-go"
	"labix.org/v2/mgo"
	"log"
	"net/http"
)

func CreateTransaction(params map[string]string, body map[string]interface{}, db *mgo.Database) error {
	log.Print(body)
	transaction := models.NewTransactionFromRequest(body)
	err := transaction.SaveToDB(db)
	if err != nil {
		return err
	}

	return nil
}

func SendPayment(params map[string]string, body map[string]interface{}, db *mgo.Database) error {
	paymentID := params["id"]
	transaction, _ := models.FindTransactionByPaymentID(paymentID, db)
	transaction.DebitSourceURI = body["debitSourceURI"].(string)
	debit := transaction.ConvertToDebit()
	customer := getClient(transaction.ClientID)
	bError := customer.Debit(debit)
	if bError != nil {
		return fmt.Errorf(bError.Description+" %s", bError.StatusCode)
	}
	transaction.DebitURI = debit.URI
	err := transaction.SaveToDB(db)
	if err != nil {
		return err
	}

	return nil
}

func getClient(clientID string) *balanced.Customer {
	client := &http.Client{}
	r, _ := http.NewRequest("GET", "http://localhost:3120/user/"+clientID, nil)
	resp, err := client.Do(r)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}

	m := make(map[string]interface{})
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&m)
	log.Print(m)
	customer := new(balanced.Customer)
	customer.URI = m["uri"].(string)
	customer.DebitsURI = m["debitsURI"].(string)
	return customer
}

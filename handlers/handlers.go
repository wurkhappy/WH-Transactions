package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/Balanced-go"
	"github.com/wurkhappy/WH-Transactions/models"
	"github.com/wurkhappy/WH-Config"
	"github.com/wurkhappy/mdp"
	"log"
)

func CreateTransaction(params map[string]string, body map[string]interface{}) error {
	log.Print(body)
	transaction := models.NewTransactionFromRequest(body)
	err := transaction.Save()
	if err != nil {
		return err
	}

	return nil
}

func SendPayment(params map[string]string, body map[string]interface{}) error {
	log.Print("send payment")
	paymentID := params["id"]
	transaction, _ := models.FindTransactionByPaymentID(paymentID)
	log.Print(paymentID)
	transaction.DebitSourceURI = body["debitSourceURI"].(string)
	log.Print(transaction.DebitSourceURI)
	debit := transaction.ConvertToDebit()
	log.Print(debit)
	customer := getClient(transaction.ClientID)
	debit.Amount = debit.Amount * 100
	bError := customer.Debit(debit)
	if bError != nil {
		log.Printf("berror is %s", bError)
		return fmt.Errorf(bError.Description+" %s", bError.StatusCode)
	}
	transaction.DebitURI = debit.URI
	err := transaction.Save()
	if err != nil {
		log.Printf("db error is %s", err)
		return err
	}

	return nil
}

func getClient(clientID string) *balanced.Customer {
	resp, statusCode := sendServiceRequest("GET", config.PaymentInfoService, "/user/"+clientID, nil)
	if statusCode >= 400 {
		return nil
	}

	var m map[string]interface{}
	json.Unmarshal(resp, &m)
	customer := new(balanced.Customer)
	customer.URI = m["uri"].(string)
	customer.DebitsURI = m["debitsURI"].(string)
	return customer
}

type ServiceResp struct {
	StatusCode float64 `json:"status_code"`
	Body       []byte  `json:"body"`
}

func sendServiceRequest(method, service, path string, body []byte) (response []byte, statusCode int) {
	client := mdp.NewClient(config.MDPBroker, false)
	defer client.Close()
	m := map[string]interface{}{
		"Method": method,
		"Path":   path,
		"Body":   body,
	}
	req, _ := json.Marshal(m)
	request := [][]byte{req}
	reply := client.Send([]byte(service), request)
	if len(reply) == 0 {
		return nil, 404
	}
	resp := new(ServiceResp)
	json.Unmarshal(reply[0], &resp)
	return resp.Body, int(resp.StatusCode)
}

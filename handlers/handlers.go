package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/Balanced-go"
	"github.com/wurkhappy/WH-Config"
	"github.com/wurkhappy/WH-Transactions/models"
	"github.com/wurkhappy/mdp"
	"log"
	"time"
)

func CreateTransaction(params map[string]string, body []byte) error {
	log.Print(string(body))
	var m map[string]interface{}
	json.Unmarshal(body, &m)
	log.Println(m)
	transaction := models.NewTransactionFromRequest(m)
	err := transaction.GetCreditSourceURI()
	if err != nil {
		return err
	}
	err = transaction.Save()
	if err != nil {
		return err
	}

	return nil
}

func SendPayment(params map[string]string, body []byte) error {
	log.Print("send payment")
	var m map[string]interface{}
	json.Unmarshal(body, &m)

	paymentID := params["id"]
	transaction, _ := models.FindTransactionByPaymentID(paymentID)
	if transaction.Amount == 0 {
		return nil
	}
	transaction.DebitSourceID = m["debitSourceID"].(string)
	err := transaction.GetDebitSourceURI()
	if err != nil {
		return err
	}
	transaction.PaymentType = m["paymentType"].(string)
	debit := transaction.ConvertToDebit()
	customer := getClient(transaction.ClientID)
	bError := customer.Debit(debit)
	if bError != nil {
		log.Printf("berror is %s", bError)
		return fmt.Errorf(bError.Description+" %s", bError.StatusCode)
	}
	transaction.DebitDate = time.Now()
	transaction.DebitURI = debit.URI
	err = transaction.Save()
	if err != nil {
		log.Printf("db error is %s", err)
		return err
	}

	return nil
}

type DebitCallback struct {
	Debit *balanced.Debit `json:"entity"`
	Type  string          `json:"type"`
}

func ProcessCredit(params map[string]string, body []byte) error {
	fmt.Println(string(body))
	var callback *DebitCallback
	json.Unmarshal(body, &callback)
	id := callback.Debit.Meta["id"]
	transaction, _ := models.FindTransactionByID(id)
	credit := transaction.ConvertToCredit()
	bank_account := transaction.CreateBankAccount()
	bError := bank_account.Credit(credit)
	if bError != nil {
		log.Printf("berror is %s", bError)
		return fmt.Errorf(bError.Description+" %s", bError.StatusCode)
	}
	transaction.CreditDate = time.Now()
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

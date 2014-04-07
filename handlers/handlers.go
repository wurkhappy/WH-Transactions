package handlers

import (
	"encoding/json"
	"github.com/wurkhappy/WH-Transactions/models"
	"log"
	"net/http"
	"time"
)

func ReceivedCreditInfo(params map[string]interface{}, body []byte) ([]byte, error, int) {
	var message struct {
		PaymentID            string  `json:"paymentID"`
		Amount               float64 `json:"amount"`
		UserID               string  `json:"userID"`
		CreditSourceBalanced string  `json:"creditSourceBalanced,omitempty"`
	}
	json.Unmarshal(body, &message)

	transaction := models.NewTransaction()
	transaction.PaymentID = message.PaymentID
	transaction.Amount = message.Amount
	transaction.FreelancerID = message.UserID
	transaction.CreditSourceBalancedID = message.CreditSourceBalanced

	err := transaction.Save()
	if err != nil {
		return nil, err, http.StatusBadRequest
	}

	return nil, nil, http.StatusOK

}

func CreateTransaction(params map[string]string, body []byte) error {
	var m map[string]interface{}
	json.Unmarshal(body, &m)
	transaction := models.NewTransactionFromRequest(m)
	err := transaction.GetCreditSourceInfo()
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
	var m map[string]interface{}
	json.Unmarshal(body, &m)

	paymentID := params["id"]
	transaction, _ := models.FindTransactionByPaymentID(paymentID)
	if transaction.Amount == 0 {
		return nil
	}
	transaction.DebitSourceID = m["debitSourceID"].(string)
	err := transaction.GetDebitSourceInfo()
	if err != nil {
		return err
	}
	transaction.PaymentType = m["paymentType"].(string)
	debitID, err := transaction.CreateDebit()
	if err != nil {
		return err
	}
	transaction.DebitDate = time.Now()
	transaction.DebitSourceBalancedID = debitID
	err = transaction.Save()
	if err != nil {
		log.Printf("db error is %s", err)
		return err
	}

	return nil
}

// type DebitCall// 	Type  string          `json:"type"`
// }

// func ProcessCredit(params map[string]string, body []byte) error {
// 	fmt.Println(string(body))
// 	var callback *DebitCallback
// 	json.Unmarshal(body, &callback)
// 	id := callback.Debit.Meta["id"]
// 	transaction, _ := models.FindTransactionByID(id)
// 	credit := transaction.ConvertToCredit()
// 	bank_account := transaction.CreateBankAccount()
// 	bError := bank_account.Credit(credit)
// 	if bError != nil {
// 		log.Printf("berror is %s", bError)
// 		return fmt.Errorf(bError.Description+" %s", bError.StatusCode)
// 	}
// 	transaction.CreditDate = time.Now()
// 	err := transaction.Save()
// 	if err != nil {
// 		log.Printf("db error is %s", err)
// 		return err
// 	}
// 	return nil
// }back struct {
// 	Debit *balanced.Debit `json:"entity"`

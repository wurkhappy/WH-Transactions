package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/WH-Transactions/models"
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
	if transaction.CreditSourceBalancedID == "" {
		events := Events{&Event{"transaction.missing_credit", []byte(`{"userID":"` + transaction.FreelancerID + `"}`)}}
		go events.Publish()
	}

	err := transaction.Save()
	if err != nil {
		return nil, err, http.StatusBadRequest
	}

	return nil, nil, http.StatusOK

}

func ReceivedDebitInfo(params map[string]interface{}, body []byte) ([]byte, error, int) {
	var message struct {
		PaymentID           string  `json:"paymentID"`
		PaymentType         string  `json:"paymentType"`
		Amount              float64 `json:"amount"`
		UserID              string  `json:"userID"`
		DebitSourceBalanced string  `json:"debitSourceBalanced,omitempty"`
	}
	json.Unmarshal(body, &message)

	transaction, err := models.FindTransactionByPaymentID(message.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("Error finding transaction %s", err.Error()), http.StatusBadRequest
	}
	transaction.DebitSourceBalancedID = message.DebitSourceBalanced
	transaction.PaymentType = message.PaymentType

	debitID, err := transaction.CreateDebit()
	if err != nil {
		return nil, fmt.Errorf("Error creating debit %s", err.Error()), http.StatusBadRequest
	}

	transaction.DebitDate = time.Now()
	transaction.DebitSourceBalancedID = debitID

	err = transaction.Save()
	if err != nil {
		return nil, fmt.Errorf("Error saving transaction %s", err.Error()), http.StatusBadRequest
	}

	return nil, nil, http.StatusOK
}

func CancelPayment(params map[string]interface{}, body []byte) ([]byte, error, int) {
	var message struct {
		PaymentID string `json:"paymentID"`
		UserID    string `json:"userID"`
	}
	json.Unmarshal(body, &message)

	err := models.DeleteTransactionWithPaymentID(message.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("Error deleting transaction %s", err.Error()), http.StatusBadRequest
	}

	return nil, nil, http.StatusOK
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

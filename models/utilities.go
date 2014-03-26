package models

import (
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/WH-Config"
	"github.com/wurkhappy/balanced-go"
	"github.com/wurkhappy/mdp"
)

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

func formatBalancedErrors(bErrors []*balanced.BalancedError) error {
	var errorMessages string
	for _, bError := range bErrors {
		errorMessages += bError.Description + ", "
	}
	if errorMessages == "" {
		return nil
	}
	return fmt.Errorf(errorMessages)
}

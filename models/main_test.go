package models

import (
	"github.com/wurkhappy/WH-Agreements/DB"
	"testing"
)

func init() {

	DB.Name = "testdb"
	DB.Setup(false)
	DB.CreateStatements()
}

func TestIntegrationTests(t *testing.T) {
	if !testing.Short() {

		DB.DB.Exec("DELETE from transaction")
	}
}

func TestUnitTests(t *testing.T) {
	test_ConvertToCredit(t)
	test_CalculateFee(t)
}

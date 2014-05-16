package main

import (
	"github.com/ant0ine/go-urlrouter"
	"github.com/wurkhappy/WH-Transactions/handlers"
)

var router urlrouter.Router = urlrouter.Router{
	Routes: []urlrouter.Route{
		urlrouter.Route{
			PathExp: "paymentinfo.credit",
			Dest:    handlers.ReceivedCreditInfo,
		},
		urlrouter.Route{
			PathExp: "paymentinfo.debit",
			Dest:    handlers.ReceivedDebitInfo,
		},
		urlrouter.Route{
			PathExp: "payment.cancelled",
			Dest:    handlers.CancelPayment,
		},
	},
}

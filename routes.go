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
	},
}

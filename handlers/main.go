package handlers

import (
	"github.com/streadway/amqp"
	"github.com/wurkhappy/WH-Config"
)

var Connection *amqp.Connection

func Setup() {
	dialRMQ()
}

func dialRMQ() {
	var err error
	Connection, err = amqp.Dial(config.RMQBroker)
	if err != nil {
		panic(err)
	}
}

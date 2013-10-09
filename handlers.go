package main

import (
	"github.com/streadway/amqp"
	"log"
)

func testKey(d amqp.Delivery) {
	log.Print("testKey")
	log.Print(string(d.Body))
	d.Ack(true)
}

func other(d amqp.Delivery) {
	log.Print("other")
	log.Print(string(d.Body))
	d.Ack(true)
}

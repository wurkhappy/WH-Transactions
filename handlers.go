package main

import (
	// "github.com/streadway/amqp"
	"log"
)

func testKey(params map[string]string, body map[string]interface{}) error{
	log.Print("testKey")
	log.Print(body)
	// d.Ack(true)
	return nil
}

func other(params map[string]string, body map[string]interface{}) error{
	log.Print("other")
	log.Print(body)
	// d.Ack(true)
	return nil
}

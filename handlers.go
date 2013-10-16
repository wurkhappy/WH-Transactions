package main

import (
	"github.com/nu7hatch/gouuid"
	// "github.com/streadway/amqp"
	"labix.org/v2/mgo"
	"log"
)

func testKey(params map[string]string, body map[string]interface{}, db *mgo.Database) error {
	log.Print("testKey")
	log.Print(body)
	// d.Ack(false)
	return nil
}

func CreateTransaction(params map[string]string, body map[string]interface{}, db *mgo.Database) error {
	log.Print("other")
	log.Print(body)
	coll := db.C("transactions")
	id, _ := uuid.NewV4()
	if _, err := coll.UpsertId(id, &body); err != nil {
		log.Printf("db error is: %v", err)
		return err
	}
	// d.Ack(false)
	return nil
}

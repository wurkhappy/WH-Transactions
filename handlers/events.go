package handlers

import (
	"github.com/streadway/amqp"
	"github.com/wurkhappy/WH-Config"
	"log"
)

type Event struct {
	Name string
	Body []byte
}

type Events []*Event

func (e Events) Publish() {
	ch := getChannel()
	defer ch.Close()
	for _, event := range e {
		event.PublishOnChannel(ch)
	}
}

func (e *Event) PublishOnChannel(ch *amqp.Channel) {
	if ch == nil {
		ch = getChannel()
		defer ch.Close()
	}

	log.Println("Publish event:", e.Name, config.MainExchange)

	err := ch.ExchangeDeclare(
		config.MainExchange, // name
		"topic",             // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // noWait
		nil,                 // arguments
	)
	if err != nil {
		log.Println(err)
	}

	err = ch.Publish(
		config.MainExchange, // exchange
		e.Name,              // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        e.Body,
		})
	if err != nil {
		log.Println(err)
	}
}

func getChannel() *amqp.Channel {
	ch, err := Connection.Channel()
	if err != nil {
		dialRMQ()
		ch, err = Connection.Channel()
		if err != nil {
			log.Print(err.Error())
		}
	}

	return ch
}

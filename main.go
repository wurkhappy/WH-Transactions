package main

import (
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	// "time"
)

var (
	uri          = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	exchange     = flag.String("exchange", "test-exchange", "Durable, non-auto-deleted AMQP exchange name")
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	queue        = flag.String("queue", "test-queue", "Ephemeral AMQP queue name")
	consumerTag  = flag.String("consumer-tag", "simple-consumer", "AMQP consumer tag (should not be blank)")
)

var routeMapper map[string]func(amqp.Delivery) = map[string]func(amqp.Delivery){
	"test-key": testKey,
	"other": other,
}

func init() {
	flag.Parse()
}

func main() {

	log.Printf("dialing %q", *uri)
	conn, err := amqp.Dial(*uri)
	if err != nil {
		fmt.Errorf("Dial: %s", err)
	}
	c, err := NewConsumer(conn, *exchange, *exchangeType, *queue, *consumerTag)
	if err != nil {
		log.Fatalf("%s", err)
	}
	for key, _ := range routeMapper{
		c.bindToQueue(*exchange, *queue, key)
	}

	deliveries, err := c.channel.Consume(
		*queue, // name
		c.tag,  // consumerTag,
		false,  // noAck
		false,  // exclusive
		false,  // noLocal
		false,  // noWait
		nil,    // arguments
	)
	if err != nil {
		fmt.Errorf("Deliveries: %s", err)
	}

	go router(deliveries, c.done)

	select {}
}

func router(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		route := routeMapper[d.RoutingKey]
		go route(d)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}

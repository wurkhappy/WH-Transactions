package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ant0ine/go-urlrouter"
	"github.com/streadway/amqp"
	"github.com/wurkhappy/Balanced-go"
	rbtmq "github.com/wurkhappy/Rabbitmq-go-wrapper"
	"github.com/wurkhappy/WH-Transactions/handlers"
	"log"
)

var (
	uri          = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	exchange     = flag.String("exchange", "transactions", "Durable, non-auto-deleted AMQP exchange name")
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	queue        = flag.String("queue", "transactions", "Ephemeral AMQP queue name")
	consumerTag  = flag.String("consumer-tag", "simple-consumer", "AMQP consumer tag (should not be blank)")
)

var Config = map[string]string{
	"DBName": "UserDB",
	"DBURL":  "192.168.139.152:27017",
}

//order matters so most general should go towards the bottom
var router urlrouter.Router = urlrouter.Router{
	Routes: []urlrouter.Route{
		urlrouter.Route{
			PathExp: "/transactions",
			Dest: map[string]interface{}{
				"POST": handlers.CreateTransaction,
			},
		},
		urlrouter.Route{
			PathExp: "/payment/:id/transaction",
			Dest: map[string]interface{}{
				"PUT": handlers.SendPayment,
			},
		},
	},
}

func init() {
	flag.Parse()
}

func main() {
	balanced.Username = "ak-test-x9PqPQUtpvUtnXsZqBL4rXGAE8WvvqoJ"
	var err error
	log.Printf("dialing %q", *uri)
	conn, err := amqp.Dial(*uri)
	if err != nil {
		fmt.Errorf("Dial: %s", err)
	}
	c, err := rbtmq.NewConsumer(conn, *exchange, *exchangeType, *queue, *consumerTag)
	if err != nil {
		log.Fatalf("%s", err)
	}

	deliveries := c.Consume(*queue)

	err = router.Start()
	if err != nil {
		panic(err)
	}
	log.Print("route")
	for d := range deliveries {
		go routeMapper(d)
	}
}

func routeMapper(d amqp.Delivery) {
	route, params, err := router.FindRoute(d.RoutingKey)
	if err != nil || route == nil {
		log.Printf("first error is: %v", err)
		return
	}

	var m map[string]interface{}
	json.Unmarshal(d.Body, &m)
	body := m["Body"].(map[string]interface{})
	routedMap := route.Dest.(map[string]interface{})
	handler := routedMap[m["Method"].(string)].(func(map[string]string, map[string]interface{}) error)
	err = handler(params, body)
	if err != nil {
		log.Printf("second error is: %v", err)
		d.Nack(false, false)
	}
	d.Ack(false)
}

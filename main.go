package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ant0ine/go-urlrouter"
	"github.com/streadway/amqp"
	"github.com/wurkhappy/Balanced-go"
	rbtmq "github.com/wurkhappy/Rabbitmq-go-wrapper"
	"github.com/wurkhappy/WH-Config"
	"github.com/wurkhappy/WH-Transactions/handlers"
	"log"
)

var (
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	consumerTag  = flag.String("consumer-tag", "simple-consumer", "AMQP consumer tag (should not be blank)")
	production   = flag.Bool("production", false, "Production settings")
)

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

func main() {
	flag.Parse()
	if *production {
		config.Prod()
	} else {
		config.Test()
	}
	balanced.Username = config.BalancedUsername
	var err error
	conn, err := amqp.Dial(config.TransactionsBroker)
	if err != nil {
		fmt.Errorf("Dial: %s", err)
	}
	c, err := rbtmq.NewConsumer(conn, config.TransactionsExchange, *exchangeType, config.TransactionsQueue, *consumerTag)
	if err != nil {
		log.Fatalf("%s", err)
	}

	deliveries := c.Consume(config.TransactionsQueue)

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

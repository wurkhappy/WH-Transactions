package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ant0ine/go-urlrouter"
	"github.com/streadway/amqp"
	rbtmq "github.com/wurkhappy/Rabbitmq-go-wrapper"
	"github.com/wurkhappy/WH-Config"
	"github.com/wurkhappy/WH-Transactions/DB"
	"github.com/wurkhappy/WH-Transactions/handlers"
	"github.com/wurkhappy/balanced-go"
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
		// urlrouter.Route{
		// 	PathExp: "/debit/process",
		// 	Dest: map[string]interface{}{
		// 		"POST": handlers.ProcessCredit,
		// 	},
		// },
	},
}

func main() {
	flag.Parse()
	if *production {
		config.Prod()
	} else {
		config.Test()
	}
	DB.Setup(*production)
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
	for d := range deliveries {
		go routeMapper(d)
	}
}

type ServiceReq struct {
	Method string
	Path   string
	Body   []byte
}

func routeMapper(d amqp.Delivery) {
	route, params, err := router.FindRoute(d.RoutingKey)
	if err != nil || route == nil {
		log.Printf("ERROR is: %v", err)
		return
	}

	var req *ServiceReq
	err = json.Unmarshal(d.Body, &req)
	if err != nil {
		fmt.Println(err)
		d.Nack(false, false)
	}

	fmt.Println(req.Path, req.Method, req.Body)

	routedMap := route.Dest.(map[string]interface{})
	handler := routedMap[req.Method].(func(map[string]string, []byte) error)
	err = handler(params, req.Body)
	if err != nil {
		log.Printf("ERROR is: %v", err)
		d.Nack(false, false)
	}
	d.Ack(false)
}

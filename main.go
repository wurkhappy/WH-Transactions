package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ant0ine/go-urlrouter"
	"github.com/streadway/amqp"
	rbtmq "github.com/wurkhappy/Rabbitmq-go-wrapper"
	"log"
	// "time"
)

var (
	uri          = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	exchange     = flag.String("exchange", "transactions", "Durable, non-auto-deleted AMQP exchange name")
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	queue        = flag.String("queue", "transactions", "Ephemeral AMQP queue name")
	consumerTag  = flag.String("consumer-tag", "simple-consumer", "AMQP consumer tag (should not be blank)")
)

//order matters so most general should go towards the bottom
var router urlrouter.Router = urlrouter.Router{
	Routes: []urlrouter.Route{
		urlrouter.Route{
			PathExp: "/other",
			Dest: map[string]interface{}{
				"POST": testKey,
			},
		},
		urlrouter.Route{
			PathExp: "/transaction/hi",
			Dest: map[string]interface{}{
				"POST": testKey,
			},
		},
		urlrouter.Route{
			PathExp: "/transaction/:id",
			Dest: map[string]interface{}{
				"POST": other,
			},
		},
	},
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
	c, err := rbtmq.NewConsumer(conn, *exchange, *exchangeType, *queue, *consumerTag)
	if err != nil {
		log.Fatalf("%s", err)
	}

	deliveries := c.Consume(*queue)

	err = router.Start()
	if err != nil {
		panic(err)
	}

	go routeMapper(deliveries)

	select {}
}

func routeMapper(deliveries <-chan amqp.Delivery) {
	for d := range deliveries {
		log.Print(d.RoutingKey)
		route, params, err := router.FindRoute(d.RoutingKey)
		if err != nil || route == nil {
			log.Print(err)
			continue
		}

		var m map[string]interface{}
		json.Unmarshal(d.Body, &m)
		body := m["Body"].(map[string]interface{})
		routedMap := route.Dest.(map[string]interface{})
		handler := routedMap[m["Method"].(string)].(func(map[string]string, map[string]interface{})(error))

		err = handler(params, body)
		if err != nil {
			d.Nack(true, true)
		}
		d.Ack(true)

	}
	log.Printf("handle: deliveries channel closed")
}

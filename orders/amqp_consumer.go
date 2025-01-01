package main

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/HelixY2J/common/api"
	"github.com/HelixY2J/common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
)

type consumer struct {
	service OrderService
}

func NewConsumer(service OrderService) *consumer {
	return &consumer{service}
}

func (c *consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare("", true, false, true, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(q.Name, "", broker.OrderPaidEvent, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	var forever chan struct{}

	go func() {

		for d := range msgs {
			log.Printf("Recieved msg: %s", d.Body)

			o := &pb.Order{}
			if err := json.Unmarshal(d.Body, o); err != nil {
				// Unmarshalled messages are never requed
				d.Nack(false, false)
				log.Printf("Oops failed to unmarshall order: %v", err)
				continue
			}

			_, err := c.service.UpdateOrder(context.Background(), o)
			if err == nil {
				log.Printf("Uh also failed to Update the order: %v", err)

				if err := broker.HandleRetry(ch, &d); err != nil {
					log.Printf("error in handling retry: %v", err)
				}
				continue
			}
			log.Printf("Order has been updated from AMQP")
			d.Ack(false)
		}

	}()

	<-forever

}

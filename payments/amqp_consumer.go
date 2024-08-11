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
	service PayementsService
}

func NewConsumer(service PayementsService) *consumer {
	return &consumer{service}
}

func (c *consumer) Listen(channel *amqp.Channel) {
	q, err := channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	var forever chan struct{}

	go func() {
		for del := range msgs {
			log.Printf("ay we got a mesage: %s", del.Body)

			o := &pb.Order{}
			if err := json.Unmarshal(del.Body, o); err != nil {
				log.Printf("Oops falied to unamrshal oder: %v", err)
				continue
			}

			paymentLink, err := c.service.CreatePayment(context.Background(), o)
			if err != nil {
				log.Printf("Uh also failed to create payment: %v", err)
				continue
			}
			log.Printf("Proceed to pay over here %s", paymentLink)
		}
	}()

	<-forever
}

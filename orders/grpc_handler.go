package main

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/HelixY2J/common/api"
	"github.com/HelixY2J/common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer

	service OrderService
	channel *amqp.Channel
}

func NewGRPCHandler(grpcServer *grpc.Server, service OrderService, channel *amqp.Channel) {
	handler := &grpcHandler{
		service: service,
		channel: channel,
	}
	pb.RegisterOrderServiceServer(grpcServer, handler)

}

func (h *grpcHandler) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {

	log.Printf("Yohoo another One! Order is %v", p)
	o := &pb.Order{
		ID: "42",
	}
	marshalledOrder, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	q, err := h.channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	h.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         marshalledOrder,
		DeliveryMode: amqp.Persistent,
	})

	return o, nil
}

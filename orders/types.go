package main

import (
	"context"

	pb "github.com/HelixY2J/common/api"
)

type OrderService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)

	ValidateOrder(context.Context, *pb.CreateOrderRequest) ([]*pb.Item, error)
}
type OrderStore interface {
	// methods for interacting with the database
	Create(context.Context) error
}

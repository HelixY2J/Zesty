package main

import (
	"context"

	pb "github.com/HelixY2J/common/api"
)

type OrderService interface {
	CreateOrder(context.Context) error

	ValidateOrder(context.Context, *pb.CreateOrderRequest) error
}
type OrderStore interface {
	// methods for interacting with the database
	Create(context.Context) error
}

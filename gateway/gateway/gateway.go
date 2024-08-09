package gateway

import (
	"context"

	pb "github.com/HelixY2J/common/api"
)

type OrderGateway interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
}

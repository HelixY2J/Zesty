package main

import (
	"context"

	pb "github.com/HelixY2J/common/api"
)

type PayementsService interface {
	CreatePayment(context.Context, *pb.Order) (string, error) // return the link user will follow to pay
}

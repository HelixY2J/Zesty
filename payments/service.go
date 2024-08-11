package main

import (
	"context"

	pb "github.com/HelixY2J/common/api"
)

type service struct {
}

func NewService() *service {

	return &service{}
}

func (s *service) CreatePayment(context.Context, *pb.Order) (string, error) {
	// connect to the payment processor

	return "", nil
}

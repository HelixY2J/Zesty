package main

import (
	"context"

	pb "github.com/HelixY2J/common/api"
	"github.com/HelixY2J/zesty-payments/processor"
)

type service struct {
	processor processor.PaymentProcessor
}

func NewService(processor processor.PaymentProcessor) *service {

	return &service{processor}
}

func (s *service) CreatePayment(ctx context.Context, o *pb.Order) (string, error) {
	// connect to the payment processor
	link, err := s.processor.CreatePaymentLink(o)

	if err != nil {
		return "", err
	}

	// update order with link
	return link, nil
}

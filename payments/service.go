package main

import (
	"context"

	pb "github.com/HelixY2J/common/api"
	"github.com/HelixY2J/zesty-payments/gateway"
	"github.com/HelixY2J/zesty-payments/processor"
)

type service struct {
	processor processor.PaymentProcessor
	gateway   gateway.OrdersGateway
}

func NewService(processor processor.PaymentProcessor, gateway gateway.OrdersGateway) *service {

	return &service{processor, gateway}
}

func (s *service) CreatePayment(ctx context.Context, o *pb.Order) (string, error) {
	// connect to the payment processor
	link, err := s.processor.CreatePaymentLink(o)

	if err != nil {
		return "", err
	}
	// update order with link

	err = s.gateway.UpdateOrderAfterPaymentLink(ctx, o.ID, link)
	if err != nil {
		return "", err
	}

	return link, nil
}

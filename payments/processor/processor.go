package processor

import pb "github.com/HelixY2J/common/api"

type PaymentProcessor interface {
	CreatePaymentLink(*pb.Order) (string, error)
}

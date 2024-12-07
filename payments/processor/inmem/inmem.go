package inmem

import pb "github.com/HelixY2J/common/api"

type Inmem struct{}

func NewInmem() *Inmem {
	return &Inmem{}
}

func (i *Inmem) CreatePaymentLink(*pb.Order) (string, error) {
	return "this is a dummy-link", nil
}

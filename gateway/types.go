package main

import pb "github.com/HelixY2J/common/api"

type CreateOrderRequest struct {
	Order         *pb.Order `json:"order"`
	RedirectToURL string    `json:"redirectToURL"`
}

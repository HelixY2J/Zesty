package main

import (
	"errors"
	"net/http"

	"github.com/HelixY2J/common"
	pb "github.com/HelixY2J/common/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	// gateway - layer of service discovery
	// handler struct for managing HTTP routes
	client pb.OrderServiceClient
}

func NewHandler(client pb.OrderServiceClient) *handler {
	return &handler{client}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.HandleCreateOrder)
}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")

	var items []*pb.ItemsWithQuantity
	if err := common.ReadJSON(r, &items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// payload validation
	if err := validateItems(items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	o, err := h.client.CreateOrder(r.Context(),
		&pb.CreateOrderRequest{
			CustomerID: customerID,
			Items:      items, // sent thru POST req
		})

	// Handling gRPC error status codes
	rStatus := status.Convert(err)

	if rStatus != nil {
		if rStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, rStatus.Message())
			return
		}

		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	common.WriteJSON(w, http.StatusOK, o)

}

func validateItems(items []*pb.ItemsWithQuantity) error {
	if len(items) == 0 {
		return common.ErrNoItems
	}

	for _, i := range items {
		if i.ID == "" {
			return errors.New("item ID is required")
		}

		if i.Quantity <= 0 {
			return errors.New("items  should have a valid quantity")
		}
	}

	return nil
}

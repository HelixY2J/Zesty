package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	pb "github.com/HelixY2J/common/api"
	"github.com/HelixY2J/common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
)

type PaymentHTTPHandler struct {
	channel *amqp.Channel
}

func NewPaymentHTTPHandler(channel *amqp.Channel) *PaymentHTTPHandler {
	return &PaymentHTTPHandler{channel}
}

func (h *PaymentHTTPHandler) registerRoutes(router *http.ServeMux) {
	router.HandleFunc("/webhook", h.handleCheckoutWebhook)

}

func (h *PaymentHTTPHandler) handleCheckoutWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)

	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	body, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in reading the request body:%v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	//fmt.Fprint(os.Stdout, "Got the body: %s\n", body)
	log.Printf("Received body: %s", body)

	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), endpointStripeSecret)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // returnns a 400 on  a bad signature
		return
	}

	// if event.Type == "checkout.session.completed" {
	// 	var session stripe.CheckoutSession
	// 	err := json.Unmarshal(event.Data.Raw, &session)
	// 	if err != nil {
	// 		fmt.Fprint(os.Stderr, "Error in parsing webhook JSON: %v\n", err)
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		return
	// 	}

	// 	if session.PaymentStatus == "paid" {
	// 		log.Printf("Payment for checkout session %v completed", session.ID)
	// 		// publish the msg
	// 	}
	// }
	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		var session stripe.CheckoutSession
		// Unmarshal the event raw data into the CheckoutSession object
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
			return
		}

		// Handle successful payment
		if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
			log.Printf("Payment for Checkout session %v completed", session.ID)
			//  process the successful payment

			orderID := session.Metadata["orderID"]
			customerID := session.Metadata["customerID"]
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			o := &pb.Order{
				ID:          orderID,
				CustomerID:  customerID,
				Status:      "paid",
				PaymentLink: "",
			}
			marshalledOrder, err := json.Marshal(o)
			if err != nil {
				log.Fatal("error in marshalling the order obj")
			}
			//  publish a message
			h.channel.PublishWithContext(ctx, broker.OrderPaidEvent, "", false, false, amqp.Publishing{
				ContentType:  "application/json",
				Body:         marshalledOrder,
				DeliveryMode: amqp.Persistent,
			})
			log.Printf("message published order is Paid")
		}

	}

	w.WriteHeader(http.StatusOK)

}

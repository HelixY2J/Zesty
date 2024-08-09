package main

import (
	"context"
	"log"
	"net/http"
	"time"

	common "github.com/HelixY2J/common"
	"github.com/HelixY2J/common/discovery"
	"github.com/HelixY2J/common/discovery/consul"
	"github.com/HelixY2J/zesty-gateway/gateway"
	_ "github.com/joho/godotenv/autoload"
)

var (
	serviceName = "gateway"
	httpAddr    = common.EnvString("HTTP_ADDR", ":8082")
	consulAddr  = common.EnvString("CONSUL_ADRR", "localhost:8500")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	instanceID := discovery.GenerateInstanceID(serviceName)

	if err := registry.Register(ctx, instanceID, serviceName, httpAddr); err != nil {
		panic(err)
	}

	go func() {
		/// checking health status with go routines
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("Failed to health check")
			}
			time.Sleep(time.Second * 1)
		}
	}()
	defer registry.Unregister(ctx, instanceID, serviceName)

	mux := http.NewServeMux()
	orderGateway := gateway.NewGRPCGateway(registry)
	handler := NewHandler(orderGateway)
	handler.registerRoutes(mux)

	log.Printf("starting HTTP server at %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("WE cant start server :<")
	}
}

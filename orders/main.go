package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/HelixY2J/common"
	"github.com/HelixY2J/common/discovery"
	"github.com/HelixY2J/common/discovery/consul"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
)

var (
	serviceName = "orders"
	grpcAddr    = common.EnvString("GRPC_ADDR", "localhost:2000") // exposing gRPC ports
	consulAddr  = common.EnvString("CONSUL_ADRR", "localhost:8500")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	instanceID := discovery.GenerateInstanceID(serviceName)

	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
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

	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

	store := NewStore()
	svc := NewService(store)
	NewGRPCHandler(grpcServer, svc)
	svc.CreateOrder(context.Background())

	log.Println("GRPC server at", grpcAddr)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf(err.Error())
	}
}

//

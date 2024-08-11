package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/HelixY2J/common"
	"github.com/HelixY2J/common/broker"
	"github.com/HelixY2J/common/discovery"
	"github.com/HelixY2J/common/discovery/consul"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
)

var (
	serviceName = "payments"
	grpcAddr    = common.EnvString("GRPC_ADDR", "localhost:2003") // exposing gRPC ports
	consulAddr  = common.EnvString("CONSUL_ADRR", "localhost:8500")
	amqpUser    = common.EnvString("RABBITMQ_USER", "guest")
	amqpPass    = common.EnvString("RABBITMQ_PASS", "guest")
	amqpHost    = common.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort    = common.EnvString("RABBITMQ_PORT", "5672")
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
				log.Fatalf("Failed to health check")
			}
			time.Sleep(time.Second * 1)
		}
	}()
	defer registry.Unregister(ctx, instanceID, serviceName)

	// broker connection
	channel, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		channel.Close()

	}()

	svc := NewService()

	amqoConsumer := NewConsumer(svc)
	go amqoConsumer.Listen(channel)

	// gRPC server
	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

	log.Println("GRPC server at", grpcAddr)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf(err.Error())
	}
}

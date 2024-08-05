package main

import (
	"log"
	"net/http"

	common "github.com/HelixY2J/common"
	pb "github.com/HelixY2J/common/api"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	httpAddr         = common.EnvString("HTTP_ADDR", ":8082")
	orderServiceAddr = "localhost:2000"
)

func main() {
	conn, err := grpc.NewClient(orderServiceAddr, grpc.WithTransportCredentials((insecure.
		NewCredentials())))
	if err != nil {
		log.Fatalf("fialed to start dial server:%v", err)
	}
	defer conn.Close()

	log.Println("Dialing orders service at", orderServiceAddr)

	c := pb.NewOrderServiceClient(conn)

	mux := http.NewServeMux()
	handler := NewHandler(c)
	handler.registerRoutes(mux)

	log.Printf("starting HTTP server at %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("WE cant start server :<")
	}
}

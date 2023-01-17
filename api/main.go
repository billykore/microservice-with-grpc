package main

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"microservice-with-grpc/api/handler"
	"microservice-with-grpc/api/router"
	pb "microservice-with-grpc/gen/customer/v1"
)

func main() {
	// conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials())) // local
	conn, err := grpc.Dial("172.22.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials())) //docker
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewCustomerClient(conn)
	customer := handler.NewCustomerHandler(client)

	h := handler.Handlers{Customer: customer}
	r := router.New(h)
	log.Printf("server listening at :8080")
	if err = r.Run(":8080"); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

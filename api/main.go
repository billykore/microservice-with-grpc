package main

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"microservice-with-grpc/api/handler"
	"microservice-with-grpc/api/router"
	authpb "microservice-with-grpc/gen/auth/v1"
	customerpb "microservice-with-grpc/gen/customer/v1"
)

func main() {
	//authConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials())) // local
	authConn, err := grpc.Dial("172.22.0.1:50052", grpc.WithTransportCredentials(insecure.NewCredentials())) // docker
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer authConn.Close()
	authClient := authpb.NewAuthClient(authConn)
	auth := handler.NewAuthHandler(authClient)

	//clientConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials())) // local
	clientConn, err := grpc.Dial("172.22.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials())) //docker
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer clientConn.Close()
	customerClient := customerpb.NewCustomerClient(clientConn)
	customer := handler.NewCustomerHandler(customerClient)

	h := handler.Handlers{Auth: auth, Customer: customer}
	r := router.New(h)
	log.Printf("server listening at :8080")
	if err = r.Run(":8080"); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

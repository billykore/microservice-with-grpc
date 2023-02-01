package main

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	customerpb "microservice-with-grpc/gen/customer/v1"
)

func CustomerClient() (customerpb.CustomerClient, func()) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials())) //docker
	//conn, err := grpc.Dial("172.22.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials())) //docker
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	closer := func() {
		err = conn.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
	}
	client := customerpb.NewCustomerClient(conn)
	return client, closer
}

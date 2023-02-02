package main

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	transactionpb "microservice-with-grpc/gen/transaction/v1"
)

func TransactionClient() (transactionpb.TransactionClient, func()) {
	conn, err := grpc.Dial("localhost:50055", grpc.WithTransportCredentials(insecure.NewCredentials())) //local
	//notificationConn, err := grpc.Dial("172.22.0.1:50054", grpc.WithTransportCredentials(insecure.NewCredentials())) //docker
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	closer := func() {
		err = conn.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
	}
	client := transactionpb.NewTransactionClient(conn)
	return client, closer
}

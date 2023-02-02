package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"microservice-with-grpc/database"
	pb "microservice-with-grpc/gen/transaction/v1"
)

type transactionServer struct {
	pb.UnimplementedTransactionServer
	Service TransactionService
}

func NewTransactionServer(service TransactionService) pb.TransactionServer {
	return &transactionServer{Service: service}
}

func (t *transactionServer) Transfer(ctx context.Context, in *pb.TransferRequest) (*pb.TransferResponse, error) {
	err := t.Service.Transfer(ctx, &Request{
		TrxId:              in.GetTrxId(),
		SourceAccount:      in.GetSourceAccount(),
		DestinationAccount: in.GetDestinationAccount(),
		Amount:             in.GetAmount(),
	})
	if err != nil {
		return &pb.TransferResponse{
			Success: false,
			Message: "Transfer failed",
		}, err
	}
	return &pb.TransferResponse{
		Success: true,
		Message: "Transfer successful",
	}, nil
}

func main() {
	db := database.New(database.MySQL, &database.Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "localhost",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
	repo := NewTransactionRepo(db)
	service := NewTransactionService(repo)
	transaction := NewTransactionServer(service)

	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTransactionServer(s, transaction)
	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

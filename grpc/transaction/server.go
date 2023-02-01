package main

import (
	"context"
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
		SourceAccount:      in.GetAccountSource(),
		DestinationAccount: in.GetAccountDestination(),
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

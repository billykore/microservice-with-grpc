package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"

	"microservice-with-grpc/database"
	pb "microservice-with-grpc/gen/payment/v1"
)

type paymentServer struct {
	pb.UnimplementedPaymentServer
	Service PaymentService
}

func NewPaymentServer(service PaymentService) pb.PaymentServer {
	return &paymentServer{Service: service}
}

func (p *paymentServer) Qris(ctx context.Context, in *pb.QrisRequest) (*pb.QrisResponse, error) {
	req := &Request{
		MerchantId:         in.GetMerchantId(),
		TrxNumber:          in.GetTrxNumber(),
		AccountSource:      in.GetAccountSource(),
		AccountDestination: in.GetAccountDestination(),
		Amount:             in.GetAmount(),
	}
	succeed, err := p.Service.Qris(ctx, req)
	if err != nil {
		log.Printf("[server error] qris payment error")
		return nil, err
	}
	if !succeed {
		log.Printf("[server error] qris payment failed")
		return &pb.QrisResponse{Success: false}, nil
	}
	return &pb.QrisResponse{Success: true}, nil
}

func main() {
	db := database.New(database.Postgres, &database.Config{
		DatabaseUser:     "postgres",
		DatabasePassword: "postgres",
		DatabaseHost:     "172.22.0.1",
		DatabasePort:     "5432",
		DatabaseName:     "grpc_payment_service",
	})
	repo := NewPaymentRepo(db)
	service := NewPaymentService(repo)
	payment := NewPaymentServer(service)

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPaymentServer(s, payment)
	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

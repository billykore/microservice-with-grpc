package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"microservice-with-grpc/database"
	pb "microservice-with-grpc/gen/customer/v1"
)

type customerServer struct {
	pb.UnimplementedCustomerServer
	service customerService
}

func newCustomerServer(service customerService) pb.CustomerServer {
	return &customerServer{service: service}
}

func (c *customerServer) accountCreation(ctx context.Context, in *pb.AccountCreationRequest) (*pb.AccountCreationResponse, error) {
	err := c.service.accountCreation(ctx, in)
	if err != nil {
		log.Printf("[server error] account creation error: %v", err)
		return &pb.AccountCreationResponse{
			Success: false,
			Message: "account creation failed",
		}, err
	}
	return &pb.AccountCreationResponse{
		Success: true,
		Message: "account creation succeed",
	}, nil
}

func (c *customerServer) accountInquiry(ctx context.Context, in *pb.InquiryRequest) (*pb.InquiryResponse, error) {
	acc, err := c.service.accountInquiry(ctx, in.GetAccountNumber())
	if err != nil {
		log.Printf("[server error] account inquiry error: %v", err)
		return &pb.InquiryResponse{}, err
	}
	return &pb.InquiryResponse{
		Cif:            acc.Cif,
		AccountNumber:  acc.AccountNumber,
		AccountType:    acc.Type,
		Name:           acc.Customer.Name,
		Currency:       acc.Currency,
		Status:         acc.Status,
		Blocked:        acc.Blocked,
		Balance:        acc.Balance,
		MinimumBalance: acc.MinimumBalance,
		ProductType:    acc.ProductType,
	}, nil
}

func main() {
	db := database.New(database.MySQL, &database.Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "172.22.0.1",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
	repo := newCustomerRepo(db)
	service := newCustomerService(repo)
	customer := newCustomerServer(service)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCustomerServer(s, customer)
	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

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
	Service CustomerService
}

func NewCustomerServer(service CustomerService) pb.CustomerServer {
	return &customerServer{Service: service}
}

func (c *customerServer) AccountCreation(ctx context.Context, in *pb.AccountCreationRequest) (*pb.AccountCreationResponse, error) {
	err := c.Service.AccountCreation(ctx, in)
	if err != nil {
		log.Printf("[server error] account creation error: %v", err)
		return &pb.AccountCreationResponse{
			Success: false,
			Message: "Account creation failed",
		}, err
	}
	return &pb.AccountCreationResponse{
		Success: true,
		Message: "Account creation succeed",
	}, nil
}

func (c *customerServer) AccountInquiry(ctx context.Context, in *pb.InquiryRequest) (*pb.InquiryResponse, error) {
	account, err := c.Service.AccountInquiry(ctx, in.GetAccountNumber())
	if err != nil {
		log.Printf("[server error] account inquiry error: %v", err)
		return &pb.InquiryResponse{}, err
	}
	return &pb.InquiryResponse{
		Cif:            account.Cif,
		AccountNumber:  account.AccountNumber,
		AccountType:    account.Type,
		Name:           account.Customer.Name,
		Currency:       account.Currency,
		Status:         account.Status,
		Blocked:        account.Blocked,
		Balance:        account.Balance,
		MinimumBalance: account.MinimumBalance,
		ProductType:    account.ProductType,
	}, nil
}

func main() {
	db, _ := database.New(&database.Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "172.22.0.1",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
	repo := NewCustomerRepo(db)
	service := NewCustomerService(repo)
	customer := NewCustomerServer(service)

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

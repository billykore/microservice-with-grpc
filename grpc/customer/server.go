package main

import (
	"context"
	"log"
	"microservice-with-grpc/database"
	"net"

	"google.golang.org/grpc"

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
		log.Printf("[server] account creation error: %v", err)
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

var (
	db, _ = database.New(&database.Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "localhost",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
	repo     = NewCustomerRepo(db)
	service  = NewCustomerService(repo)
	customer = NewCustomerServer(service)
)

func main() {
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

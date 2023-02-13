package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"microservice-with-grpc/database"
	pb "microservice-with-grpc/gen/auth/v1"
)

type authServer struct {
	pb.UnimplementedAuthServer
	service authService
}

func newAuthServer(service authService) pb.AuthServer {
	return &authServer{service: service}
}

func (a *authServer) GetToken(ctx context.Context, in *pb.TokenRequest) (*pb.TokenResponse, error) {
	req := &request{
		username:  in.GetUsername(),
		password:  in.GetPassword(),
		grantType: in.GetGrantType(),
	}
	token, err := a.service.getToken(ctx, req)
	if err != nil {
		log.Printf("[server error] get token error: %v", err)
		return nil, err
	}
	return &pb.TokenResponse{
		Token:     token.Token,
		TokenType: token.Type,
		ExpiredAt: token.ExpiresIn,
	}, nil
}

func main() {
	db := database.New(database.Postgres, &database.Config{
		DatabaseUser:     "postgres",
		DatabasePassword: "postgres",
		DatabaseHost:     "172.22.0.1",
		DatabasePort:     "5432",
		DatabaseName:     "grpc_auth_service",
	})
	repo := newAuthRepo(db)
	service := newAuthService(repo)
	auth := newAuthServer(service)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAuthServer(s, auth)
	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package main

import (
	"context"
	"log"
	pb "microservice-with-grpc/gen/auth/v1"
)

type authServer struct {
	pb.UnimplementedAuthServer
	Service AuthService
}

func NewAuthServer(service AuthService) pb.AuthServer {
	return &authServer{Service: service}
}

func (a *authServer) GetToken(ctx context.Context, in *pb.TokenRequest) (*pb.TokenResponse, error) {
	req := &Request{
		Username:  in.GetUsername(),
		Password:  in.GetPassword(),
		GrantType: in.GetGrantType(),
	}
	token, err := a.Service.GetToken(ctx, req)
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

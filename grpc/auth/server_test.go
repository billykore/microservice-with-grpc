package main

import (
	"context"
	"errors"
	"log"
	"microservice-with-grpc/database"
	"microservice-with-grpc/entity"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	pb "microservice-with-grpc/gen/auth/v1"
)

type authServiceMock struct {
	mock.Mock
}

func (m *authServiceMock) GetToken(ctx context.Context, req *Request) (*entity.Token, error) {
	args := m.Mock.Called(ctx, req)
	if args.Get(0) == nil && args.Get(1) != nil {
		return nil, errors.New("failed to create token")
	}
	return args.Get(0).(*entity.Token), nil
}

func server(ctx context.Context, service AuthService) (pb.AuthClient, func()) {
	buffer := 1024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()
	pb.RegisterAuthServer(baseServer, NewAuthServer(service))
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err = lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	client := pb.NewAuthClient(conn)
	return client, closer
}

func TestAuthServer_GetToken(t *testing.T) {
	in := &pb.TokenRequest{
		Username:  "user",
		Password:  "password",
		GrantType: "password",
	}

	expectation := &pb.TokenResponse{
		Token:     "example-token-777",
		TokenType: "Bearer token",
		ExpiredAt: 900,
	}

	service := &authServiceMock{Mock: mock.Mock{}}
	service.On("GetToken", mock.Anything, mock.Anything).Return(&entity.Token{
		Token:     "example-token-777",
		Type:      "Bearer token",
		ExpiresIn: 900,
	}, nil)

	ctx := context.Background()
	auth, closer := server(ctx, service)
	defer closer()

	out, err := auth.GetToken(ctx, in)
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, expectation.Token, out.Token)
	assert.Equal(t, expectation.TokenType, out.TokenType)
	assert.Equal(t, expectation.ExpiredAt, out.ExpiredAt)
}

func TestAuthServerIntegrationTest_GetToken(t *testing.T) {
	in := &pb.TokenRequest{
		Username:  "user",
		Password:  "password",
		GrantType: "password",
	}

	db := database.New(database.Postgres, &database.Config{
		DatabaseUser:     "postgres",
		DatabasePassword: "postgres",
		DatabaseHost:     "localhost",
		DatabasePort:     "5432",
		DatabaseName:     "grpc_auth_service",
	})
	assert.NotNil(t, db)
	repo := NewAuthRepo(db)
	assert.NotNil(t, repo)
	service := NewAuthService(repo)
	assert.NotNil(t, service)

	ctx := context.Background()
	auth, closer := server(ctx, service)
	defer closer()

	out, err := auth.GetToken(ctx, in)
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.NotEmpty(t, out.Token)
	assert.Equal(t, "Bearer token", out.TokenType)
	assert.Equal(t, 900.0, out.ExpiredAt)
}

func TestAuthServerIntegrationTest_GetTokenUserNotExist(t *testing.T) {
	in := &pb.TokenRequest{
		Username:  "userNotExist",
		Password:  "duh",
		GrantType: "password",
	}

	db := database.New(database.Postgres, &database.Config{
		DatabaseUser:     "postgres",
		DatabasePassword: "postgres",
		DatabaseHost:     "localhost",
		DatabasePort:     "5432",
		DatabaseName:     "grpc_auth_service",
	})
	assert.NotNil(t, db)
	repo := NewAuthRepo(db)
	assert.NotNil(t, repo)
	service := NewAuthService(repo)
	assert.NotNil(t, service)

	ctx := context.Background()
	auth, closer := server(ctx, service)
	defer closer()

	out, err := auth.GetToken(ctx, in)
	assert.Error(t, err)
	assert.Nil(t, out)
}

func TestAuthServerIntegrationTest_GetTokenWrongPassword(t *testing.T) {
	in := &pb.TokenRequest{
		Username:  "user",
		Password:  "wrongPassword",
		GrantType: "password",
	}

	db := database.New(database.Postgres, &database.Config{
		DatabaseUser:     "postgres",
		DatabasePassword: "postgres",
		DatabaseHost:     "localhost",
		DatabasePort:     "5432",
		DatabaseName:     "grpc_auth_service",
	})
	assert.NotNil(t, db)
	repo := NewAuthRepo(db)
	assert.NotNil(t, repo)
	service := NewAuthService(repo)
	assert.NotNil(t, service)

	ctx := context.Background()
	auth, closer := server(ctx, service)
	defer closer()

	out, err := auth.GetToken(ctx, in)
	assert.Error(t, err)
	assert.Nil(t, out)
}

func TestAuthServerIntegrationTest_GetTokenInvalidGrantType(t *testing.T) {
	in := &pb.TokenRequest{
		Username:  "user",
		Password:  "password",
		GrantType: "grantType",
	}

	db := database.New(database.Postgres, &database.Config{
		DatabaseUser:     "postgres",
		DatabasePassword: "postgres",
		DatabaseHost:     "localhost",
		DatabasePort:     "5432",
		DatabaseName:     "grpc_auth_service",
	})
	assert.NotNil(t, db)
	repo := NewAuthRepo(db)
	assert.NotNil(t, repo)
	service := NewAuthService(repo)
	assert.NotNil(t, service)

	ctx := context.Background()
	auth, closer := server(ctx, service)
	defer closer()

	out, err := auth.GetToken(ctx, in)
	assert.Error(t, err)
	assert.Nil(t, out)
}

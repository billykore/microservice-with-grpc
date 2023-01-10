package main

import (
	"context"
	"errors"
	"log"
	"microservice-with-grpc/database"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	pb "microservice-with-grpc/gen/customer/v1"
)

type customerServiceMock struct {
	mock.Mock
}

func (m *customerServiceMock) AccountCreation(ctx context.Context, data *pb.AccountCreationRequest) error {
	args := m.Mock.Called(ctx, data)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func server(ctx context.Context, service CustomerService) (pb.CustomerClient, func()) {
	buffer := 1024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()
	pb.RegisterCustomerServer(baseServer, NewCustomerServer(service))
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
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	client := pb.NewCustomerClient(conn)
	return client, closer
}

func TestCustomerServer_AccountCreation(t *testing.T) {
	in := &pb.AccountCreationRequest{
		Nik:            "5310121711980001",
		Name:           "Evanbill Antonio Kore",
		Pob:            "Bajawa",
		Dob:            "17/11/1998",
		Address:        "Jakarta",
		Profession:     "Officer",
		Gender:         pb.Gender_MALE,
		Religion:       pb.Religion_PROTESTANT,
		MarriageStatus: pb.MarriageStatus_NOT_MARRIED,
		Citizen:        pb.Citizen_WNI,
	}

	expectation := &pb.AccountCreationResponse{
		Success: true,
		Message: "Account creation succeed",
	}

	service := &customerServiceMock{Mock: mock.Mock{}}
	service.Mock.On("AccountCreation", mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()
	customer, closer := server(ctx, service)
	defer closer()

	out, err := customer.AccountCreation(ctx, in)
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, expectation.Success, out.Success)
	assert.Equal(t, expectation.Message, out.Message)
}

func TestCustomerServer_AccountCreationFailed(t *testing.T) {
	in := &pb.AccountCreationRequest{}
	expectation := &pb.AccountCreationResponse{
		Success: false,
		Message: "Account creation failed",
	}

	service := &customerServiceMock{Mock: mock.Mock{}}
	service.Mock.On("AccountCreation", mock.Anything, mock.Anything).Return(errors.New("account creation failed"))

	ctx := context.Background()
	customer, closer := server(ctx, service)
	defer closer()

	out, err := customer.AccountCreation(ctx, in)
	assert.Error(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, expectation.Success, out.Success)
	assert.Equal(t, expectation.Message, out.Message)
}

func TestCustomerServerIntegrationTest_AccountCreation(t *testing.T) {
	in := &pb.AccountCreationRequest{
		Nik:            "0101011509990001",
		Name:           "Melita Wandriani Baru",
		Pob:            "Ruteng",
		Dob:            "15/19/1999",
		Address:        "Kupang",
		Profession:     "Veterinarian",
		Gender:         pb.Gender_FEMALE,
		Religion:       pb.Religion_CATHOLIC,
		MarriageStatus: pb.MarriageStatus_NOT_MARRIED,
		Citizen:        pb.Citizen_WNI,
	}

	expectation := &pb.AccountCreationResponse{
		Success: true,
		Message: "Account creation succeed",
	}

	db, err := database.New(&database.Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "localhost",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
	assert.NoError(t, err)
	err = database.Migrate(db, &Customer{})
	assert.NoError(t, err)
	repo := NewCustomerRepo(db)
	assert.NotNil(t, repo)
	service := NewCustomerService(repo)
	assert.NotNil(t, service)

	ctx := context.Background()
	customer, closer := server(ctx, service)
	defer closer()

	out, err := customer.AccountCreation(ctx, in)
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, expectation.Success, out.Success)
	assert.Equal(t, expectation.Message, out.Message)
}

func TestCustomerServerIntegrationTest_AccountCreationFailed(t *testing.T) {
	// customer nik already exist, account creation must fail.
	in := &pb.AccountCreationRequest{
		Nik:            "0101011509990001",
		Name:           "Melita Wandriani Baru",
		Pob:            "Ruteng",
		Dob:            "15/19/1999",
		Address:        "Kupang",
		Profession:     "Veterinarian",
		Gender:         pb.Gender_FEMALE,
		Religion:       pb.Religion_CATHOLIC,
		MarriageStatus: pb.MarriageStatus_NOT_MARRIED,
		Citizen:        pb.Citizen_WNI,
	}

	expectation := &pb.AccountCreationResponse{
		Success: false,
		Message: "Account creation failed",
	}

	db, err := database.New(&database.Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "localhost",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
	assert.NoError(t, err)
	err = database.Migrate(db, &Customer{})
	assert.NoError(t, err)
	repo := NewCustomerRepo(db)
	assert.NotNil(t, repo)
	service := NewCustomerService(repo)
	assert.NotNil(t, service)

	ctx := context.Background()
	customer, closer := server(ctx, service)
	defer closer()

	out, err := customer.AccountCreation(ctx, in)
	assert.Error(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, expectation.Success, out.Success)
	assert.Equal(t, expectation.Message, out.Message)
}

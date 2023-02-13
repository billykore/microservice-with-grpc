package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"microservice-with-grpc/database"
	pb "microservice-with-grpc/gen/customer/v1"
)

type customerServiceMock struct {
	mock.Mock
}

func (m *customerServiceMock) accountCreation(ctx context.Context, data *pb.AccountCreationRequest) error {
	args := m.Mock.Called(ctx, data)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func (m *customerServiceMock) accountInquiry(ctx context.Context, accountNumber string) (*account, error) {
	args := m.Mock.Called(ctx, accountNumber)
	if args.Get(0) == nil && args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*account), nil
}

func server(ctx context.Context, service customerService) (pb.CustomerClient, func()) {
	buffer := 1024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()
	pb.RegisterCustomerServer(baseServer, newCustomerServer(service))
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
		Message: "account creation succeed",
	}

	service := &customerServiceMock{Mock: mock.Mock{}}
	service.Mock.On("accountCreation", mock.Anything, mock.Anything).Return(nil)

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
		Message: "account creation failed",
	}

	service := &customerServiceMock{Mock: mock.Mock{}}
	service.Mock.On("accountCreation", mock.Anything, mock.Anything).Return(errors.New("account creation failed"))

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
		Dob:            "15/09/1999",
		Address:        "Kupang",
		Profession:     "Veterinarian",
		Gender:         pb.Gender_FEMALE,
		Religion:       pb.Religion_CATHOLIC,
		MarriageStatus: pb.MarriageStatus_NOT_MARRIED,
		Citizen:        pb.Citizen_WNI,
	}

	expectation := &pb.AccountCreationResponse{
		Success: true,
		Message: "account creation succeed",
	}

	db := database.New(database.MySQL, &database.Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "localhost",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
	assert.NotNil(t, db)
	err := database.Migrate(db, &customer{})
	assert.NoError(t, err)
	repo := newCustomerRepo(db)
	assert.NotNil(t, repo)
	service := newCustomerService(repo)
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
	// account nik already exist, account creation must fail.
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
		Message: "account creation failed",
	}

	db := database.New(database.MySQL, &database.Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "localhost",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
	assert.NotNil(t, db)
	err := database.Migrate(db, &customer{})
	assert.NoError(t, err)
	repo := newCustomerRepo(db)
	assert.NotNil(t, repo)
	service := newCustomerService(repo)
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

func TestCustomerServer_AccountInquiry(t *testing.T) {
	in := &pb.InquiryRequest{
		AccountNumber: "001001000001300",
	}

	expectation := &pb.InquiryResponse{
		Cif:            "0000000001",
		AccountNumber:  "001001000001300",
		AccountType:    "S",
		Name:           "John Doe",
		Currency:       "IDR",
		Status:         "1",
		Blocked:        "0",
		Balance:        "100000000.00",
		MinimumBalance: "0.00",
		ProductType:    "000005",
	}

	service := &customerServiceMock{Mock: mock.Mock{}}
	service.On("accountInquiry", mock.Anything, in.AccountNumber).Return(expectation, nil)

	ctx := context.Background()
	customer, closer := server(ctx, service)
	defer closer()

	out, err := customer.AccountInquiry(ctx, in)
	fmt.Println(out)
	assert.NoError(t, err)
	assert.NotNil(t, out)
}

func TestCustomerServer_AccountInquiryFailed(t *testing.T) {
	in := &pb.InquiryRequest{
		AccountNumber: "001001000001300",
	}

	service := &customerServiceMock{Mock: mock.Mock{}}
	service.On("accountInquiry", mock.Anything, in.AccountNumber).Return(nil, errors.New("account not found"))

	ctx := context.Background()
	customer, closer := server(ctx, service)
	defer closer()

	out, err := customer.AccountInquiry(ctx, in)
	assert.Error(t, err)
	assert.Nil(t, out)
}

func TestCustomerServerIntegrationTest_AccountInquiry(t *testing.T) {
	in := &pb.InquiryRequest{
		AccountNumber: "001001000002300",
	}

	expectation := &pb.InquiryResponse{
		Cif:            "0000000002",
		AccountNumber:  "001001000002300",
		AccountType:    "S",
		Name:           "FLORENCE FEDORA AGUSTINA",
		Currency:       "IDR",
		Status:         "1",
		Blocked:        "0",
		Balance:        "0.00",
		MinimumBalance: "0.00",
		ProductType:    "000005",
	}

	db := database.New(database.MySQL, &database.Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "localhost",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
	assert.NotNil(t, db)
	repo := newCustomerRepo(db)
	assert.NotNil(t, repo)
	service := newCustomerService(repo)
	assert.NotNil(t, service)

	ctx := context.Background()
	customer, closer := server(ctx, service)
	defer closer()

	out, err := customer.AccountInquiry(ctx, in)
	fmt.Println(out)
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, expectation.Cif, out.Cif)
	assert.Equal(t, expectation.AccountNumber, out.AccountNumber)
	assert.Equal(t, expectation.Name, out.Name)
}

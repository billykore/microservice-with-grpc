package main

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"microservice-with-grpc/database"
	pb "microservice-with-grpc/gen/transaction/v1"
)

type transactionServiceMock struct {
	mock.Mock
}

func (m *transactionServiceMock) Transfer(ctx context.Context, req *Request) error {
	args := m.Mock.Called(ctx, req)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func server(ctx context.Context, service TransactionService) (pb.TransactionClient, func()) {
	buffer := 1024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()
	pb.RegisterTransactionServer(baseServer, NewTransactionServer(service))
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

	client := pb.NewTransactionClient(conn)
	return client, closer
}

func TestTransactionServer_Transfer(t *testing.T) {
	type mockData struct {
		err error
	}

	type args struct {
		ctx context.Context
		in  *pb.TransferRequest
	}

	type expectation struct {
		out *pb.TransferResponse
		err error
	}

	ctx := context.Background()

	// test cases
	tests := map[string]struct {
		mock     mockData
		args     args
		expected expectation
	}{
		"transfer_successful": {
			mock: mockData{
				err: nil,
			},
			args: args{
				ctx: ctx,
				in: &pb.TransferRequest{
					TrxId:              "example-id-123",
					SourceAccount:      "123456789",
					DestinationAccount: "987654321",
					Amount:             "500000",
				},
			},
			expected: expectation{
				out: &pb.TransferResponse{
					Success: true,
					Message: "Transfer successful",
				},
				err: nil,
			},
		},
		"transfer_failed": {
			mock: mockData{
				err: errors.New("error from service"),
			},
			args: args{
				ctx: ctx,
				in:  &pb.TransferRequest{},
			},
			expected: expectation{
				out: &pb.TransferResponse{
					Success: false,
					Message: "Transfer failed",
				},
				err: errors.New("error from service"),
			},
		},
	}

	// run the test
	for scenario, test := range tests {
		service := &transactionServiceMock{Mock: mock.Mock{}}
		transaction := NewTransactionServer(service)
		t.Run(scenario, func(t *testing.T) {
			// mock program
			service.Mock.On("Transfer", mock.Anything, mock.Anything).Return(test.mock.err)
			// call transfer function
			got, err := transaction.Transfer(test.args.ctx, test.args.in)
			// check the output
			if got != nil {
				assert.Equal(t, test.expected.out.Success, got.Success)
				assert.Equal(t, test.expected.out.Message, got.Message)
			}
			assert.Equal(t, test.expected.err, err)
		})
	}
}

func TestTransactionServerIntegrationTest_Transfer(t *testing.T) {
	type args struct {
		ctx context.Context
		in  *pb.TransferRequest
	}

	type expectation struct {
		out *pb.TransferResponse
		err error
	}

	ctx := context.Background()

	// test cases
	tests := map[string]struct {
		args     args
		expected expectation
	}{
		"transfer_successful": {
			args: args{
				ctx: ctx,
				in: &pb.TransferRequest{
					TrxId:              "example-id-123",
					SourceAccount:      "001001000001300",
					DestinationAccount: "001001000002300",
					Amount:             "50000",
				},
			},
			expected: expectation{
				out: &pb.TransferResponse{
					Success: true,
					Message: "Transfer successful",
				},
				err: nil,
			},
		},
		"transfer_failed": {
			args: args{
				ctx: ctx,
				in:  &pb.TransferRequest{},
			},
			expected: expectation{
				out: nil,
				err: errors.New("rpc error: code = Unknown desc = transfer error"),
			},
		},
	}

	db := database.New(database.MySQL, &database.Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "localhost",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
	assert.NotNil(t, db)
	repo := NewTransactionRepo(db)
	assert.NotNil(t, repo)
	service := NewTransactionService(repo)
	assert.NotNil(t, service)

	transaction, closer := server(ctx, service)
	defer closer()

	// run the test
	for scenario, test := range tests {
		t.Run(scenario, func(t *testing.T) {
			// call transfer function
			got, err := transaction.Transfer(test.args.ctx, test.args.in)
			// check the output
			if got != nil {
				assert.Equal(t, test.expected.out.Success, got.Success)
				assert.Equal(t, test.expected.out.Message, got.Message)
			}
			if err != nil {
				assert.Equal(t, test.expected.err.Error(), err.Error())
			}
		})
	}
}

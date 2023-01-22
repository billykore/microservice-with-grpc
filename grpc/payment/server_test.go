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
	pb "microservice-with-grpc/gen/payment/v1"
)

type paymentServiceMock struct {
	mock.Mock
}

func (m *paymentServiceMock) Qris(ctx context.Context, req *Request) (bool, error) {
	args := m.Mock.Called(ctx, req)
	if args.Get(0) == false && args.Get(1) != nil {
		return false, args.Get(1).(error)
	}
	if args.Get(0) == false && args.Get(1) == nil {
		return false, nil
	}
	return true, nil
}

func server(ctx context.Context, service PaymentService) (pb.PaymentClient, func()) {
	buffer := 1024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()
	pb.RegisterPaymentServer(baseServer, NewPaymentServer(service))
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

	client := pb.NewPaymentClient(conn)
	return client, closer
}

func TestPaymentServer_Qris(t *testing.T) {
	type mockData struct {
		ctx context.Context
		req *Request
		res struct {
			success bool
			err     error
		}
	}

	type args struct {
		ctx context.Context
		in  *pb.QrisRequest
	}

	type expectation struct {
		out *pb.QrisResponse
		err error
	}

	ctx := context.Background()

	// test cases
	tests := map[string]struct {
		mock     mockData
		args     args
		expected expectation
	}{
		"success": {
			mock: mockData{
				ctx: ctx,
				req: &Request{},
				res: struct {
					success bool
					err     error
				}{
					success: true,
					err:     nil,
				},
			},
			args: args{
				ctx: ctx,
				in:  &pb.QrisRequest{},
			},
			expected: expectation{
				out: &pb.QrisResponse{Success: true},
				err: nil,
			},
		},
		"qris_error": {
			mock: mockData{
				ctx: ctx,
				req: &Request{},
				res: struct {
					success bool
					err     error
				}{
					success: false,
					err:     errors.New("error"),
				},
			},
			args: args{
				ctx: ctx,
				in:  &pb.QrisRequest{},
			},
			expected: expectation{
				out: nil,
				err: errors.New("error"),
			},
		},
		"qris_failed": {
			mock: mockData{
				ctx: ctx,
				req: &Request{},
				res: struct {
					success bool
					err     error
				}{
					success: false,
					err:     nil,
				},
			},
			args: args{
				ctx: ctx,
				in:  &pb.QrisRequest{},
			},
			expected: expectation{
				out: &pb.QrisResponse{Success: false},
				err: nil,
			},
		},
	}

	// run the test
	for scenario, test := range tests {
		service := &paymentServiceMock{Mock: mock.Mock{}}
		payment := NewPaymentServer(service)
		t.Run(scenario, func(t *testing.T) {
			// mock program
			service.On("Qris", test.mock.ctx, test.mock.req).Return(test.mock.res.success, test.mock.res.err)
			// calling Qris function
			got, err := payment.Qris(test.args.ctx, test.args.in)
			// test the output
			if got != nil {
				assert.Equal(t, test.expected.out.Success, got.Success)
			}
			assert.Equal(t, test.expected.err, err)
		})
	}
}

func TestPaymentServerIntegrationTest_Qris(t *testing.T) {
	type args struct {
		ctx context.Context
		in  *pb.QrisRequest
	}

	type expectation struct {
		out *pb.QrisResponse
		err error
	}

	ctx := context.Background()

	// test cases
	tests := map[string]struct {
		args     args
		expected expectation
	}{
		"success": {
			args: args{
				ctx: ctx,
				in: &pb.QrisRequest{
					MerchantId:         "M-001",
					TrxNumber:          "000007",
					AccountSource:      "001001000001300",
					AccountDestination: "001001000002300",
					Amount:             "50000",
				},
			},
			expected: expectation{
				out: &pb.QrisResponse{Success: true},
				err: nil,
			},
		},
	}

	db := database.New(database.Postgres, &database.Config{
		DatabaseUser:     "postgres",
		DatabasePassword: "postgres",
		DatabaseHost:     "localhost",
		DatabasePort:     "5432",
		DatabaseName:     "grpc_payment_service",
	})
	assert.NotNil(t, db)
	repo := NewPaymentRepo(db)
	assert.NotNil(t, repo)
	service := NewPaymentService(repo)
	assert.NotNil(t, service)

	payment, closer := server(ctx, service)
	defer closer()

	// run the test
	for scenario, test := range tests {
		t.Run(scenario, func(t *testing.T) {
			// calling Qris function
			got, err := payment.Qris(test.args.ctx, test.args.in)
			// test the output
			if got != nil {
				assert.Equal(t, test.expected.out.Success, got.Success)
			}
			if err != nil {
				assert.Equal(t, test.expected.err.Error(), err.Error())
			}
		})
	}
}

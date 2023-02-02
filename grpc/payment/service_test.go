package main

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"microservice-with-grpc/database"
	"microservice-with-grpc/entity"
)

type paymentRepoMock struct {
	mock.Mock
}

func (m *paymentRepoMock) InsertQrisLog(ctx context.Context, qrisLog *entity.QrisLog) error {
	args := m.Mock.Called(ctx, qrisLog)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func TestPaymentService_Qris(t *testing.T) {
	type mockData struct {
		ctx     context.Context
		qrisLog *entity.QrisLog
		res     struct {
			err error
		}
	}

	type args struct {
		ctx context.Context
		req *Request
	}

	type expectation struct {
		out bool
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
				ctx:     ctx,
				qrisLog: &entity.QrisLog{},
				res: struct {
					err error
				}{
					err: nil,
				},
			},
			args: args{
				ctx: ctx,
				req: &Request{},
			},
			expected: expectation{
				out: true,
				err: nil,
			},
		},
		"failed": {
			mock: mockData{
				ctx:     ctx,
				qrisLog: &entity.QrisLog{},
				res: struct {
					err error
				}{
					err: errors.New("error"),
				},
			},
			args: args{
				ctx: ctx,
				req: &Request{},
			},
			expected: expectation{
				out: false,
				err: errors.New("qris payment error"),
			},
		},
	}

	// run the test
	for scenario, test := range tests {
		repo := &paymentRepoMock{Mock: mock.Mock{}}
		service := NewPaymentService(repo)
		t.Run(scenario, func(t *testing.T) {
			// mock program
			repo.Mock.On("InsertQrisLog", test.mock.ctx, test.mock.qrisLog).Return(test.mock.res.err)
			// call Qris function
			got, err := service.Qris(test.args.ctx, test.args.req)
			// test the output
			assert.Equal(t, test.expected.out, got)
			assert.Equal(t, test.expected.err, err)
		})
	}
}

func TestPaymentServiceIntegrationTest_Qris(t *testing.T) {
	type args struct {
		ctx context.Context
		req *Request
	}

	type expectation struct {
		out bool
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
				req: &Request{
					MerchantId:         "M-001",
					TrxNumber:          "000077",
					SourceAccount:      "001001000001300",
					DestinationAccount: "001001000002300",
					Amount:             "50000",
				},
			},
			expected: expectation{
				out: true,
				err: nil,
			},
		},
		"failed": {
			args: args{
				ctx: ctx,
				req: nil,
			},
			expected: expectation{
				out: false,
				err: errors.New("qris payment error"),
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

	// run the test
	for scenario, test := range tests {
		t.Run(scenario, func(t *testing.T) {
			// call Qris function
			got, err := service.Qris(test.args.ctx, test.args.req)
			// test the output
			assert.Equal(t, test.expected.out, got)
			assert.Equal(t, test.expected.err, err)
		})
	}
}

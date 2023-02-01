package main

import (
	"context"
	"errors"
	"microservice-with-grpc/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type transactionRepoMock struct {
	mock.Mock
}

func (m *transactionRepoMock) UpdateBalanceByAccountNumber(ctx context.Context, accountNumber, balance string) error {
	args := m.Mock.Called(ctx, accountNumber, balance)
	if args.Get(0) != nil {
		return args.Error(1)
	}
	return nil
}

func TestTransactionService_Transfer(t *testing.T) {
	type mockData struct {
		err error
	}

	type args struct {
		ctx context.Context
		req *Request
	}

	type expectation struct {
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
				err: nil,
			},
			args: args{
				ctx: ctx,
				req: &Request{
					TrxId:              "example-id-123",
					SourceAccount:      "001001000001300",
					DestinationAccount: "001001000002300",
					Amount:             "50000",
				},
			},
			expected: expectation{
				err: nil,
			},
		},
		"transfer_failed": {
			mock: mockData{
				err: errors.New("error from repo"),
			},
			args: args{
				ctx: ctx,
				req: &Request{
					TrxId:              "example-id-123",
					SourceAccount:      "123456789",
					DestinationAccount: "987654321",
					Amount:             "500000",
				},
			},
			expected: expectation{
				err: errors.New("transfer error"),
			},
		},
	}

	// run the test
	for scenario, test := range tests {
		repo := &transactionRepoMock{Mock: mock.Mock{}}
		service := NewTransactionService(repo)
		t.Run(scenario, func(t *testing.T) {
			// mock program
			repo.Mock.On("UpdateBalanceByAccountNumber", mock.Anything, mock.Anything, mock.Anything).Return(test.mock.err)
			// call transfer function
			err := service.Transfer(test.args.ctx, test.args.req)
			// check the output
			assert.Equal(t, test.expected.err, err)
		})
	}
}

func TestTransactionServiceIntegrationTest_Transfer(t *testing.T) {
	type args struct {
		ctx context.Context
		req *Request
	}

	type expectation struct {
		err error
	}

	ctx := context.Background()

	// test cases
	tests := map[string]struct {
		args     args
		expected expectation
	}{
		"transfer_success": {
			args: args{
				ctx: ctx,
				req: &Request{
					TrxId:              "example-id-123",
					SourceAccount:      "001001000001300",
					DestinationAccount: "001001000002300",
					Amount:             "50000",
				},
			},
			expected: expectation{
				err: nil,
			},
		},
		"transfer_failed": {
			args: args{
				ctx: ctx,
				req: &Request{
					TrxId:              "example-id-123",
					SourceAccount:      "123456789",
					DestinationAccount: "987654321",
					Amount:             "500000",
				},
			},
			expected: expectation{
				err: errors.New("transfer failed"),
			},
		},
	}

	// run the test
	for scenario, test := range tests {
		t.Run(scenario, func(t *testing.T) {
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

			// call transfer function
			err := service.Transfer(test.args.ctx, test.args.req)
			// check the output
			assert.Equal(t, test.expected.err, err)
		})
	}
}

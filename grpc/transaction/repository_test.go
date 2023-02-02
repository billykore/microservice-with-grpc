package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"microservice-with-grpc/database"
	"testing"
)

func TestTransactionRepo_UpdateBalanceByAccountNumber(t *testing.T) {
	type args struct {
		ctx           context.Context
		accountNumber string
		balance       string
	}

	type expectation struct {
		err error
	}

	ctx := context.Background()

	tests := map[string]struct {
		args     args
		expected expectation
	}{
		"update_balance_success": {
			args: args{
				ctx:           ctx,
				accountNumber: "001001000002300",
				balance:       "1000000",
			},
			expected: expectation{
				err: nil,
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

	// run the test
	for scenario, test := range tests {
		t.Run(scenario, func(t *testing.T) {
			// call update balance function
			err := repo.UpdateBalanceByAccountNumber(test.args.ctx, test.args.accountNumber, test.args.balance)
			// check the output
			assert.Equal(t, test.expected.err, err)
		})
	}
}

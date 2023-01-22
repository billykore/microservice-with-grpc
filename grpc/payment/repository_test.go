package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"microservice-with-grpc/database"
	"microservice-with-grpc/entity"
)

func TestDBConnection(t *testing.T) {
	db := database.New(database.Postgres, &database.Config{
		DatabaseUser:     "postgres",
		DatabasePassword: "postgres",
		DatabaseHost:     "localhost",
		DatabasePort:     "5432",
		DatabaseName:     "grpc_payment_service",
	})
	assert.NotNil(t, db)
}

func TestPaymentRepo_InsertQrisLog(t *testing.T) {
	db := database.New(database.Postgres, &database.Config{
		DatabaseUser:     "postgres",
		DatabasePassword: "postgres",
		DatabaseHost:     "localhost",
		DatabasePort:     "5432",
		DatabaseName:     "grpc_payment_service",
	})
	assert.NotNil(t, db)
	err := db.AutoMigrate(&entity.QrisLog{})
	assert.NoError(t, err)

	repo := NewPaymentRepo(db)
	assert.NotNil(t, repo)

	err = repo.InsertQrisLog(context.Background(), &entity.QrisLog{
		MerchantId:         "M-001",
		TrxNumber:          "000001",
		AccountSource:      "001001000001300",
		AccountDestination: "001001000002300",
		Amount:             "50000",
	})
	assert.NoError(t, err)
}

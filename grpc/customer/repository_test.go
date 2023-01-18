package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"microservice-with-grpc/database"
)

func TestCustomerRepo_CreateCustomer(t *testing.T) {
	db := database.New(database.MySQL, &database.Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "localhost",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
	assert.NotNil(t, db)
	err := database.Migrate(db, &Customer{})

	repo := NewCustomerRepo(db)
	err = repo.CreateCustomer(context.Background(), &Customer{
		Nik:            "0101010505970001",
		Name:           "NI LUH PUTU GIRI GITA SARASWATI",
		Pob:            "BAJAWA",
		Dob:            "05/05/1997",
		Address:        "DENPASAR",
		Profession:     "SOFTWARE DEVELOPER",
		Gender:         "FEMALE",
		Religion:       "HINDU",
		MarriageStatus: "NOT_MARRIED",
		Citizen:        "WNI",
		Cif:            "0000000001",
	})
	assert.NoError(t, err)
}

func TestCustomerRepo_GetLastCif(t *testing.T) {
	db := database.New(database.MySQL, &database.Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "localhost",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
	assert.NotNil(t, db)
	err := database.Migrate(db, &Customer{})

	repo := NewCustomerRepo(db)
	lastCif, err := repo.GetLastCif(context.Background())
	assert.NotEmpty(t, lastCif)
	assert.NoError(t, err)
}

func TestCustomerRepo_CreateAccount(t *testing.T) {
	db := database.New(database.MySQL, &database.Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "localhost",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
	assert.NotNil(t, db)
	err := database.Migrate(db, &Account{})

	repo := NewCustomerRepo(db)
	err = repo.CreateAccount(context.Background(), &Account{
		Cif:              "0000000001",
		AccountNumber:    "001001000001300",
		Type:             "S",
		Balance:          "100000",
		MinimumBalance:   "0",
		AvailableBalance: "100000",
		Status:           "1",
		Currency:         "IDR",
		ProductType:      "000005",
		Blocked:          "0",
	})
	assert.NoError(t, err)
}

func TestCustomerRepo_InquiryByAccountNumber(t *testing.T) {
	db := database.New(database.MySQL, &database.Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "localhost",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
	assert.NotNil(t, db)

	repo := NewCustomerRepo(db)
	account, err := repo.InquiryByAccountNumber(context.Background(), "001001000002300")
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.IsType(t, &Account{}, account)
}

func TestCustomerRepo_GetCustomerByAccountNumber(t *testing.T) {
	db := database.New(database.MySQL, &database.Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "localhost",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
	assert.NotNil(t, db)

	repo := NewCustomerRepo(db)
	customer, err := repo.GetCustomerByAccountNumber(context.Background(), "001001000002300")
	assert.NoError(t, err)
	assert.NotNil(t, customer)
	assert.IsType(t, &Customer{}, customer)
	assert.Equal(t, "FLORENCE FEDORA AGUSTINA", customer.Name)
}

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
	err := database.Migrate(db, &customer{})

	repo := newCustomerRepo(db)
	err = repo.createCustomer(context.Background(), &customer{
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
	err := database.Migrate(db, &customer{})

	repo := newCustomerRepo(db)
	lastCif, err := repo.getLastCif(context.Background())
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
	err := database.Migrate(db, &customer{})

	repo := newCustomerRepo(db)
	err = repo.createAccount(context.Background(), &account{
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

	repo := newCustomerRepo(db)
	acc, err := repo.inquiryByAccountNumber(context.Background(), "001001000002300")
	assert.NoError(t, err)
	assert.NotNil(t, acc)
	assert.IsType(t, &account{}, acc)
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

	repo := newCustomerRepo(db)
	cust, err := repo.getCustomerByAccountNumber(context.Background(), "001001000002300")
	assert.NoError(t, err)
	assert.NotNil(t, cust)
	assert.IsType(t, &customer{}, cust)
	assert.Equal(t, "FLORENCE FEDORA AGUSTINA", cust.Name)
}

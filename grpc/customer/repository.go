package main

import (
	"context"
	"log"

	"gorm.io/gorm"
)

type CustomerRepo interface {
	CreateCustomer(ctx context.Context, customer *Customer) error
	GetLastCif(ctx context.Context) (string, error)
	GetLastAccount(ctx context.Context) (string, error)
	CreateAccount(ctx context.Context, account *Account) error
}

type customerRepo struct {
	DB *gorm.DB
}

func NewCustomerRepo(DB *gorm.DB) CustomerRepo {
	return &customerRepo{DB: DB}
}

func (r *customerRepo) CreateCustomer(ctx context.Context, customer *Customer) error {
	tx := r.DB.WithContext(ctx).Create(customer)
	if err := tx.Error; err != nil {
		log.Printf("[repository error] error create account. %v", err)
		return err
	}
	return nil
}

func (r *customerRepo) GetLastCif(ctx context.Context) (string, error) {
	customer := &Customer{}
	tx := r.DB.WithContext(ctx).Last(customer)
	if err := tx.Error; err != nil {
		log.Printf("[repository error] error get last cif. %v", err)
		return "", err
	}
	return customer.Cif, nil
}

func (r *customerRepo) GetLastAccount(ctx context.Context) (string, error) {
	account := &Account{}
	tx := r.DB.WithContext(ctx).Last(account)
	if err := tx.Error; err != nil {
		log.Printf("[repository error] error get last account. %v", err)
		return "", err
	}
	return account.AccountNumber, nil
}

func (r *customerRepo) CreateAccount(ctx context.Context, account *Account) error {
	tx := r.DB.WithContext(ctx).Create(account)
	if err := tx.Error; err != nil {
		log.Printf("[repository error] error create new account. %v", err)
		return err
	}
	return nil
}

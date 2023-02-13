package main

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

type customerRepo interface {
	createCustomer(ctx context.Context, customer *customer) error
	getLastCif(ctx context.Context) (string, error)
	getLastAccount(ctx context.Context) (string, error)
	createAccount(ctx context.Context, account *account) error
	inquiryByAccountNumber(ctx context.Context, accountNumber string) (*account, error)
	getCustomerByAccountNumber(ctx context.Context, accountNumber string) (*customer, error)
}

type customerRepoImpl struct {
	db *gorm.DB
}

func newCustomerRepo(db *gorm.DB) customerRepo {
	return &customerRepoImpl{db: db}
}

func (r *customerRepoImpl) createCustomer(ctx context.Context, customer *customer) error {
	tx := r.db.WithContext(ctx).Create(customer)
	if err := tx.Error; err != nil {
		log.Printf("[repository error] error create customer: %v", err)
		return errors.New("error create customer")
	}
	return nil
}

func (r *customerRepoImpl) getLastCif(ctx context.Context) (string, error) {
	customer := new(customer)
	tx := r.db.WithContext(ctx).Last(customer)
	if err := tx.Error; err != nil {
		log.Printf("[repository error] error get last cif: %v", err)
		return "", errors.New("error get last cif")
	}
	return customer.Cif, nil
}

func (r *customerRepoImpl) getLastAccount(ctx context.Context) (string, error) {
	account := new(account)
	tx := r.db.WithContext(ctx).Last(account)
	if err := tx.Error; err != nil {
		log.Printf("[repository error] error get last account: %v", err)
		return "", errors.New("error get last account")
	}
	return account.AccountNumber, nil
}

func (r *customerRepoImpl) createAccount(ctx context.Context, account *account) error {
	tx := r.db.WithContext(ctx).Create(account)
	if err := tx.Error; err != nil {
		log.Printf("[repository error] error create new account: %v", err)
		return errors.New("create new account")
	}
	return nil
}

func (r *customerRepoImpl) inquiryByAccountNumber(ctx context.Context, accountNumber string) (*account, error) {
	account := new(account)
	tx := r.db.WithContext(ctx).First(account, "account_number = ?", accountNumber)
	if err := tx.Error; err != nil {
		log.Printf("[repository error] error inquiry by account number: %v", err)
		return nil, errors.New("inquiry by account number")
	}
	return account, nil
}

func (r *customerRepoImpl) getCustomerByAccountNumber(ctx context.Context, accountNumber string) (*customer, error) {
	customer := new(customer)
	tx := r.db.Table("customers").
		Select("*").
		Joins("LEFT JOIN accounts ON customers.cif = accounts.cif").
		Where("accounts.account_number = ?", accountNumber).
		WithContext(ctx).
		Scan(customer)
	if err := tx.Error; err != nil {
		log.Printf("[repository error] error get customer by account number. %v", err)
		return nil, errors.New("error get customer by account number")
	}
	return customer, nil
}

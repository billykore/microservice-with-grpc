package main

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

type CustomerRepo interface {
	CreateCustomer(ctx context.Context, customer *Customer) error
	GetLastCif(ctx context.Context) (string, error)
	GetLastAccount(ctx context.Context) (string, error)
	CreateAccount(ctx context.Context, account *Account) error
	InquiryByAccountNumber(ctx context.Context, accountNumber string) (*Account, error)
	GetCustomerByAccountNumber(ctx context.Context, accountNumber string) (*Customer, error)
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
		log.Printf("[repository error] error create customer: %v", err)
		return errors.New("error create customer")
	}
	return nil
}

func (r *customerRepo) GetLastCif(ctx context.Context) (string, error) {
	customer := new(Customer)
	tx := r.DB.WithContext(ctx).Last(customer)
	if err := tx.Error; err != nil {
		log.Printf("[repository error] error get last cif: %v", err)
		return "", errors.New("error get last cif")
	}
	return customer.Cif, nil
}

func (r *customerRepo) GetLastAccount(ctx context.Context) (string, error) {
	account := new(Account)
	tx := r.DB.WithContext(ctx).Last(account)
	if err := tx.Error; err != nil {
		log.Printf("[repository error] error get last account: %v", err)
		return "", errors.New("error get last account")
	}
	return account.AccountNumber, nil
}

func (r *customerRepo) CreateAccount(ctx context.Context, account *Account) error {
	tx := r.DB.WithContext(ctx).Create(account)
	if err := tx.Error; err != nil {
		log.Printf("[repository error] error create new account: %v", err)
		return errors.New("create new account")
	}
	return nil
}

func (r *customerRepo) InquiryByAccountNumber(ctx context.Context, accountNumber string) (*Account, error) {
	account := new(Account)
	tx := r.DB.WithContext(ctx).First(account, "account_number = ?", accountNumber)
	if err := tx.Error; err != nil {
		log.Printf("[repository error] error inquiry by account number: %v", err)
		return nil, errors.New("inquiry by account number")
	}
	return account, nil
}

func (r *customerRepo) GetCustomerByAccountNumber(ctx context.Context, accountNumber string) (*Customer, error) {
	customer := new(Customer)
	tx := r.DB.Table("customers").
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

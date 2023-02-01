package main

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
	"microservice-with-grpc/entity"
)

type TransactionRepo interface {
	UpdateBalanceByAccountNumber(ctx context.Context, accountNumber, balance string) error
}

type transactionRepo struct {
	DB *gorm.DB
}

func NewTransactionRepo(DB *gorm.DB) TransactionRepo {
	return &transactionRepo{DB: DB}
}

func (r *transactionRepo) UpdateBalanceByAccountNumber(ctx context.Context, accountNumber, balance string) error {
	account := new(entity.Account)
	tx := r.DB.Model(account).Where("account_number = ?", accountNumber).WithContext(ctx).Update("balance", balance)
	if err := tx.Error; err != nil {
		log.Printf("[transaction repo error] error update balance by account: %v", err)
		return errors.New("error update balance by account")
	}
	return nil
}

package main

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"

	"microservice-with-grpc/entity"
)

type PaymentRepo interface {
	InsertQrisLog(ctx context.Context, qrisLog *entity.QrisLog) error
}

type paymentRepo struct {
	DB *gorm.DB
}

func NewPaymentRepo(DB *gorm.DB) PaymentRepo {
	return &paymentRepo{DB: DB}
}

func (r *paymentRepo) InsertQrisLog(ctx context.Context, qrisLog *entity.QrisLog) error {
	tx := r.DB.WithContext(ctx).Create(qrisLog)
	if err := tx.Error; err != nil {
		log.Printf("[repository error] error insert qris log: %v", err)
		return errors.New("error insert qris log")
	}
	return nil
}

package main

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"

	"microservice-with-grpc/entity"
)

type AuthRepo interface {
	GetUser(ctx context.Context, username string) (*entity.User, error)
	InsertTokenLog(ctx context.Context, log *entity.TokenLog) error
}

type authRepo struct {
	DB *gorm.DB
}

func NewAuthRepo(DB *gorm.DB) AuthRepo {
	return &authRepo{DB: DB}
}

func (r *authRepo) GetUser(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	tx := r.DB.WithContext(ctx).First(&user, "username = ?", username)
	if err := tx.Error; err != nil {
		log.Printf("[repository error] error get user: %v", err)
		return nil, errors.New("error get user")
	}
	return &user, nil
}

func (r *authRepo) InsertTokenLog(ctx context.Context, tokenLog *entity.TokenLog) error {
	tx := r.DB.WithContext(ctx).Create(tokenLog)
	if err := tx.Error; err != nil {
		log.Printf("[repository error] error insert token log: %v,", err)
		return errors.New("error insert token log")
	}
	return nil
}

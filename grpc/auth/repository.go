package main

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"

	"microservice-with-grpc/entity"
)

type authRepo interface {
	getUser(ctx context.Context, username string) (*entity.User, error)
	insertTokenLog(ctx context.Context, log *entity.TokenLog) error
}

type authRepoImpl struct {
	db *gorm.DB
}

func newAuthRepo(DB *gorm.DB) authRepo {
	return &authRepoImpl{db: DB}
}

func (r *authRepoImpl) getUser(ctx context.Context, username string) (*entity.User, error) {
	user := new(entity.User)
	tx := r.db.WithContext(ctx).First(user, "username = ?", username)
	if err := tx.Error; err != nil {
		log.Printf("[repository error] error get user: %v", err)
		return nil, errors.New("error get user")
	}
	return user, nil
}

func (r *authRepoImpl) insertTokenLog(ctx context.Context, tokenLog *entity.TokenLog) error {
	tx := r.db.WithContext(ctx).Create(tokenLog)
	if err := tx.Error; err != nil {
		log.Printf("[repository error] error insert token log: %v,", err)
		return errors.New("error insert token log")
	}
	return nil
}

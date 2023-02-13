package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"microservice-with-grpc/database"
	"microservice-with-grpc/entity"
)

func TestAuthRepo_GetUser(t *testing.T) {
	db := database.New(database.Postgres, &database.Config{
		DatabaseUser:     "postgres",
		DatabasePassword: "postgres",
		DatabaseHost:     "localhost",
		DatabasePort:     "5432",
		DatabaseName:     "grpc_auth_service",
	})
	assert.NotNil(t, db)
	repo := newAuthRepo(db)

	user, err := repo.getUser(context.Background(), "user")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user", user.Username)
}

func TestAuthRepo_InsertTokenLog(t *testing.T) {
	db := database.New(database.Postgres, &database.Config{
		DatabaseUser:     "postgres",
		DatabasePassword: "postgres",
		DatabaseHost:     "localhost",
		DatabasePort:     "5432",
		DatabaseName:     "grpc_auth_service",
	})
	assert.NotNil(t, db)
	err := db.Table("token_logs").AutoMigrate(&entity.TokenLog{})
	assert.NoError(t, err)
	repo := newAuthRepo(db)
	assert.NotNil(t, repo)

	token, err := generateToken("user")
	assert.NoError(t, err)
	err = repo.insertTokenLog(context.Background(), &entity.TokenLog{
		Token:          token,
		User:           "user",
		TokenExpiresIn: tokenExpiresTime.Seconds(),
	})
	assert.NoError(t, err)
}

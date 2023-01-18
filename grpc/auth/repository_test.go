package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"microservice-with-grpc/database"
	"testing"
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
	repo := NewAuthRepo(db)

	user, err := repo.GetUser(context.Background(), "user")
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
	err := db.Table("token_logs").AutoMigrate(&Log{})
	assert.NoError(t, err)
	repo := NewAuthRepo(db)
	assert.NotNil(t, repo)

	token, err := GenerateToken("user")
	assert.NoError(t, err)
	err = repo.InsertTokenLog(context.Background(), &Log{
		Token:          token,
		User:           "user",
		TokenExpiresIn: TokenExpiresTime.Seconds(),
	})
	assert.NoError(t, err)
}

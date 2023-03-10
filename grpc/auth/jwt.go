package main

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(username string) (string, error) {
	token, err := createToken(username)
	return token, err
}

const TokenExpiresTime = 15 * time.Minute

func createToken(subject string) (string, error) {
	jwtKey := "this-is-secret-key"
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiresTime)), // 10 minutes
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   subject,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		log.Printf("[jwt error] create token error. %v", err)
		return "", errors.New("cannot create token")
	}
	return tokenString, nil
}

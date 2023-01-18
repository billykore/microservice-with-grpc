package main

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type AuthService interface {
	GetToken(ctx context.Context, req *Request) (*Token, error)
}

type authService struct {
	Repo AuthRepo
}

func NewAuthService(repo AuthRepo) AuthService {
	return &authService{Repo: repo}
}

func (s *authService) GetToken(ctx context.Context, req *Request) (*Token, error) {
	if req.GrantType != "password" {
		log.Printf("[service error] invalid grant type: %v", req.GrantType)
		return nil, errors.New("invalid grant type")
	}
	user, err := s.Repo.GetUser(ctx, req.Username)
	if err != nil {
		log.Printf("[service error] error get token: %v", err)
		return nil, errors.New("error get token")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		log.Printf("[service error] error get token: %v", err)
		return nil, errors.New("error get token")
	}
	tokenString, err := GenerateToken(user.Username)
	if err != nil {
		log.Printf("[service error] error get token: %v", err)
		return nil, errors.New("error get token")
	}
	token := &Token{
		Token:     tokenString,
		Type:      "Bearer token",
		ExpiresIn: TokenExpiresTime.Seconds(),
	}
	err = s.Repo.InsertTokenLog(ctx, &Log{
		Token:          token.Token,
		User:           user.Username,
		TokenExpiresIn: token.ExpiresIn,
	})
	if err != nil {
		log.Printf("[service error] failed insert log")
	}
	return token, nil
}

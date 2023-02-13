package main

import (
	"context"
	"golang.org/x/crypto/bcrypt"

	"microservice-with-grpc/entity"
)

type authService interface {
	getToken(ctx context.Context, req *request) (*entity.Token, error)
}

type authServiceImpl struct {
	repo authRepo
}

func newAuthService(repo authRepo) authService {
	return &authServiceImpl{repo: repo}
}

func (s *authServiceImpl) getToken(ctx context.Context, req *request) (*entity.Token, error) {
	token, err := s.generateToken(ctx, req)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *authServiceImpl) generateToken(ctx context.Context, req *request) (*entity.Token, error) {
	token, err := s.generateTokenAndInsertTokenLog(ctx, req)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *authServiceImpl) generateTokenAndInsertTokenLog(ctx context.Context, req *request) (*entity.Token, error) {
	err := s.checkGrantType(req)
	if err != nil {
		return nil, err
	}
	user, err := s.getUserAndCheckPassword(ctx, req)
	if err != nil {
		return nil, err
	}
	tokenStr, err := s.generateTokenString(user)
	if err != nil {
		return nil, err
	}
	token := entity.NewToken(tokenStr)
	// insert token log
	err = s.insertTokenLog(ctx, token, user)
	if err != nil {
		return nil, err
	}
	return token, nil
}

const GrantTypePassword = "password"

func (s *authServiceImpl) checkGrantType(req *request) error {
	if req.grantType != GrantTypePassword {
		return invalidGrantTypeException(req)
	}
	return nil
}

func (s *authServiceImpl) getUserAndCheckPassword(ctx context.Context, req *request) (*entity.User, error) {
	user, err := s.getUserFromRepo(ctx, req)
	if err != nil {
		return nil, err
	}
	err = s.checkUserPassword(user, req)
	if err != nil {
		return nil, err
	}
	return nil, err
}

func (s *authServiceImpl) getUserFromRepo(ctx context.Context, req *request) (*entity.User, error) {
	user, err := s.repo.getUser(ctx, req.username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authServiceImpl) checkUserPassword(user *entity.User, req *request) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.password))
	if err != nil {
		return err
	}
	return nil
}

func (s *authServiceImpl) generateTokenString(user *entity.User) (string, error) {
	tokenString, err := generateToken(user.Username)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *authServiceImpl) insertTokenLog(ctx context.Context, token *entity.Token, user *entity.User) error {
	err := s.repo.insertTokenLog(ctx, &entity.TokenLog{
		Token:          token.Token,
		User:           user.Username,
		TokenExpiresIn: token.ExpiresIn,
	})
	if err != nil {
		return err
	}
	return nil
}

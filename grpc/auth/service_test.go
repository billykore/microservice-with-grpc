package main

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"microservice-with-grpc/database"
	"microservice-with-grpc/entity"
)

type authRepoMock struct {
	mock.Mock
}

func (m *authRepoMock) GetUser(ctx context.Context, username string) (*entity.User, error) {
	args := m.Mock.Called(ctx, username)
	if args.Get(0) == false && args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*entity.User), nil
}

func (m *authRepoMock) InsertTokenLog(ctx context.Context, tokenLog *entity.TokenLog) error {
	args := m.Mock.Called(ctx, tokenLog)
	if args.Get(0) == nil {
		return args.Get(0).(error)
	}
	return nil
}

func TestAuthService_GetToken(t *testing.T) {
	type repoOut struct {
		user *entity.User
		err  error
	}

	type args struct {
		ctx context.Context
		req *Request
	}

	type expectation struct {
		token *entity.Token
		err   error
	}

	tests := map[string]struct {
		repoOut  repoOut
		args     args
		expected expectation
	}{
		"success": {
			repoOut: repoOut{
				user: &entity.User{
					Username: "user",
					Password: "$2a$10$EkaaT.WwU4y5kRnOdXpoKuBg7IBTwVr2ixlcPS7DjVKnPHytG5X4K",
				},
				err: nil,
			},
			args: args{
				ctx: context.Background(),
				req: &Request{
					Username:  "user",
					Password:  "password",
					GrantType: "password",
				},
			},
			expected: expectation{
				token: &entity.Token{
					Token:     "",
					Type:      "Bearer token",
					ExpiresIn: 900,
				},
				err: nil,
			},
		},
		"user_not_found": {
			repoOut: repoOut{
				user: nil,
				err:  errors.New("error get user"),
			},
			args: args{
				ctx: context.Background(),
				req: &Request{
					Username:  "userNotExist",
					Password:  "duh",
					GrantType: "password",
				},
			},
			expected: expectation{
				token: nil,
				err:   errors.New("error get token"),
			},
		},
	}

	for scenario, test := range tests {
		t.Run(scenario, func(t *testing.T) {
			repo := &authRepoMock{Mock: mock.Mock{}}
			repo.On("GetUser", test.args.ctx, test.args.req.Username).Return(test.repoOut.user, test.repoOut.err)
			service := NewAuthService(repo)
			out, err := service.GetToken(test.args.ctx, test.args.req)
			assert.Equal(t, test.expected.err, err)
			if out != nil {
				assert.NotEmpty(t, out.Token)
				assert.Equal(t, test.expected.token.Type, out.Type)
				assert.Equal(t, test.expected.token.ExpiresIn, out.ExpiresIn)
			}
		})
	}
}

func TestAuthServiceIntegrationTest_GetToken(t *testing.T) {
	type args struct {
		ctx context.Context
		req *Request
	}

	type expectation struct {
		token *entity.Token
		err   error
	}

	tests := map[string]struct {
		args     args
		expected expectation
	}{
		"success": {
			args: args{
				ctx: context.Background(),
				req: &Request{
					Username:  "user",
					Password:  "password",
					GrantType: "password",
				},
			},
			expected: expectation{
				token: &entity.Token{
					Token:     "",
					Type:      "Bearer token",
					ExpiresIn: 900,
				},
				err: nil,
			},
		},
		"user_not_found": {
			args: args{
				ctx: context.Background(),
				req: &Request{
					Username:  "userNotExist",
					Password:  "duh",
					GrantType: "password",
				},
			},
			expected: expectation{
				token: nil,
				err:   errors.New("error get token"),
			},
		},
		"wrong_password": {
			args: args{
				ctx: context.Background(),
				req: &Request{
					Username:  "user",
					Password:  "wrongPassword",
					GrantType: "password",
				},
			},
			expected: expectation{
				token: nil,
				err:   errors.New("error get token"),
			},
		},
		"invalid_grant_type": {
			args: args{
				ctx: context.Background(),
				req: &Request{
					Username:  "user",
					Password:  "password",
					GrantType: "client_credentials",
				},
			},
			expected: expectation{
				token: nil,
				err:   errors.New("invalid grant type"),
			},
		},
	}

	db := database.New(database.Postgres, &database.Config{
		DatabaseUser:     "postgres",
		DatabasePassword: "postgres",
		DatabaseHost:     "localhost",
		DatabasePort:     "5432",
		DatabaseName:     "grpc_auth_service",
	})
	assert.NotNil(t, db)
	repo := NewAuthRepo(db)
	assert.NotNil(t, repo)
	service := NewAuthService(repo)
	assert.NotNil(t, service)

	for scenario, test := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := service.GetToken(test.args.ctx, test.args.req)
			assert.Equal(t, test.expected.err, err)
			if out != nil {
				assert.NotEmpty(t, out.Token)
				assert.Equal(t, test.expected.token.Type, out.Type)
				assert.Equal(t, test.expected.token.ExpiresIn, out.ExpiresIn)
			}
		})
	}
}

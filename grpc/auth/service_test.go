package main

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"microservice-with-grpc/database"
)

type authRepoMock struct {
	mock.Mock
}

func (m *authRepoMock) GetUser(ctx context.Context, username string) (*User, error) {
	args := m.Mock.Called(ctx, username)
	if args.Get(0) == false && args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*User), nil
}

func (m *authRepoMock) InsertTokenLog(ctx context.Context, log *Log) error {
	//TODO implement me
	panic("implement me")
}

func TestAuthService_GetToken(t *testing.T) {
	type repoOut struct {
		user *User
		err  error
	}

	type args struct {
		ctx context.Context
		req *Request
	}

	type expectation struct {
		token *Token
		err   error
	}

	tests := map[string]struct {
		repoOut  repoOut
		args     args
		expected expectation
	}{
		"success": {
			repoOut: repoOut{
				user: &User{
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
				token: &Token{
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

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			repo := &authRepoMock{Mock: mock.Mock{}}
			repo.On("GetUser", tt.args.ctx, tt.args.req.Username).Return(tt.repoOut.user, tt.repoOut.err)
			service := NewAuthService(repo)
			out, err := service.GetToken(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.expected.err, err)
			if out != nil {
				assert.NotEmpty(t, out.Token)
				assert.Equal(t, tt.expected.token.Type, out.Type)
				assert.Equal(t, tt.expected.token.ExpiresIn, out.ExpiresIn)
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
		token *Token
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
				token: &Token{
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

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := service.GetToken(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.expected.err, err)
			if out != nil {
				assert.NotEmpty(t, out.Token)
				assert.Equal(t, tt.expected.token.Type, out.Type)
				assert.Equal(t, tt.expected.token.ExpiresIn, out.ExpiresIn)
			}
		})
	}
}

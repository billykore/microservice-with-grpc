package entity

import (
	"gorm.io/gorm"
	"time"
)

const (
	TokenExpiresTime = 15 * time.Minute
	TokenType        = "Bearer token"
)

type Token struct {
	Token     string
	Type      string
	ExpiresIn float64
}

func NewToken(token string) *Token {
	return &Token{
		Token:     token,
		Type:      TokenType,
		ExpiresIn: TokenExpiresTime.Seconds(),
	}
}

type TokenLog struct {
	gorm.Model
	Token          string  `gorm:"unique;not null"`
	User           string  `gorm:"not null"`
	TokenExpiresIn float64 `gorm:"not null"`
}

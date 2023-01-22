package entity

import "gorm.io/gorm"

type Token struct {
	Token     string
	Type      string
	ExpiresIn float64
}

type TokenLog struct {
	gorm.Model
	Token          string  `gorm:"unique;not null"`
	User           string  `gorm:"not null"`
	TokenExpiresIn float64 `gorm:"not null"`
}

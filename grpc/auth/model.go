package main

import "gorm.io/gorm"

type Request struct {
	Username  string
	Password  string
	GrantType string
}

type Token struct {
	Token     string
	Type      string
	ExpiresIn float64
}

type User struct {
	Username string
	Password string
}

type Log struct {
	gorm.Model
	Token          string  `gorm:"unique;not null"`
	User           string  `gorm:"not null"`
	TokenExpiresIn float64 `gorm:"not null"`
}

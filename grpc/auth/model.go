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
	Token        *Token `gorm:"foreignKey:Token;unique;not null"`
	Date         string `gorm:"not null"`
	User         string `gorm:"not null"`
	TokenExpired bool   `gorm:"not null"`
}

package entity

import "gorm.io/gorm"

type QrisLog struct {
	gorm.Model
	MerchantId         string `gorm:"not null"`
	TrxNumber          string `gorm:"unique;not null"`
	AccountSource      string `gorm:"not null"`
	AccountDestination string `gorm:"not null"`
	Amount             string `gorm:"not null"`
}

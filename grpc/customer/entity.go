package main

import (
	"strings"

	"gorm.io/gorm"

	pb "microservice-with-grpc/gen/customer/v1"
	"microservice-with-grpc/internal"
)

type customer struct {
	gorm.Model
	Nik            string `gorm:"unique;not null"`
	Name           string `gorm:"not null"`
	Pob            string `gorm:"not null"` // place of birth
	Dob            string `gorm:"not null"` // date of birth
	Address        string `gorm:"not null"`
	Profession     string `gorm:"not null"`
	Gender         string `gorm:"not null"`
	Religion       string `gorm:"not null"`
	MarriageStatus string `gorm:"not null"`
	Citizen        string `gorm:"not null"`
	Cif            string `gorm:"unique;not null"`
}

func buildCustomer(data *pb.AccountCreationRequest) *customer {
	return &customer{
		Model:          gorm.Model{},
		Nik:            data.Nik,
		Name:           strings.ToUpper(data.Name),
		Pob:            strings.ToUpper(data.Pob),
		Dob:            data.Dob,
		Address:        strings.ToUpper(data.Address),
		Profession:     strings.ToUpper(data.Profession),
		Gender:         data.Gender.String(),
		Religion:       data.Religion.String(),
		MarriageStatus: data.MarriageStatus.String(),
		Citizen:        data.Citizen.String(),
	}
}

type account struct {
	gorm.Model
	Customer         *customer `gorm:"foreignKey:Cif"`
	Cif              string    `gorm:"unique;not null"`
	AccountNumber    string    `gorm:"unique;not null"`
	Type             string    `gorm:"not null"`
	Balance          string    `gorm:"not null"`
	MinimumBalance   string    `gorm:"not null"`
	AvailableBalance string    `gorm:"not null"`
	Status           string    `gorm:"not null"`
	Currency         string    `gorm:"not null"`
	ProductType      string    `gorm:"not null"`
	Blocked          string    `gorm:"not null;default:0"`
}

func buildAccount(cif, accNumber string, accType internal.AccountType) *account {
	return &account{
		Model:            gorm.Model{},
		Cif:              cif,
		AccountNumber:    accNumber,
		Type:             accType.String(),
		Balance:          "0.00",
		MinimumBalance:   "0.00",
		AvailableBalance: "0.00",
		Status:           internal.StatusOpen.String(), // default status is open
		Currency:         "IDR",
		ProductType:      "000005",
		Blocked:          internal.NotBlocked.String(),
	}
}

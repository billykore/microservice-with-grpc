package main

import (
	"context"
	"errors"
	"log"

	pb "microservice-with-grpc/gen/customer/v1"
)

type CustomerService interface {
	AccountCreation(ctx context.Context, req *pb.AccountCreationRequest) error
	AccountInquiry(ctx context.Context, accountNumber string) (*Account, error)
}

type customerService struct {
	Repo CustomerRepo
}

func NewCustomerService(repo CustomerRepo) CustomerService {
	return &customerService{Repo: repo}
}

func (s *customerService) AccountCreation(ctx context.Context, req *pb.AccountCreationRequest) error {
	// create new cif and register account.
	customer := BuildCustomer(req)
	newCif, err := s.createCif(ctx)
	if err != nil {
		log.Printf("[service error] account creation error. %v", err)
		return err
	}
	customer.Cif = newCif
	err = s.Repo.CreateCustomer(ctx, customer)
	if err != nil {
		log.Printf("[service error] account creation error. %v", err)
		return errors.New("customerService.AccountCreation returns error. please check the logs")
	}
	// create new account number for new account cif.
	newAccountNumber, err := s.createAccountNumber(ctx, SavingAccount)
	if err != nil {
		log.Printf("[service error] account creation error. %v", err)
		return errors.New("customerService.AccountCreation returns error. please check the logs")
	}
	account := BuildAccount(customer.Cif, newAccountNumber, SavingAccount)
	err = s.Repo.CreateAccount(ctx, account)
	if err != nil {
		log.Printf("[service error] account creation error. %v", err)
		return errors.New("customerService.AccountCreation returns error. please check the logs")
	}
	// should return no error if account creation is successful.
	return nil
}

func (s *customerService) createCif(ctx context.Context) (string, error) {
	lastCif, err := s.Repo.GetLastCif(ctx)
	if err != nil {
		log.Printf("[service error] create new cif error. %v", err)
		return "", err
	}
	newCif := BuildNewCif(lastCif)
	return newCif, nil
}

func (s *customerService) createAccountNumber(ctx context.Context, accType AccountType) (string, error) {
	lastAccountNumber, err := s.Repo.GetLastAccount(ctx)
	if err != nil {
		log.Printf("[service error] create new account number error. %v", err)
		return "", err
	}
	var newAccount string
	switch accType {
	case SavingAccount:
		newAccount = BuildNewAccountNumber(lastAccountNumber)
		break
	case GiroAccount:
		newAccount = BuildNewAccountNumber(lastAccountNumber)
		break
	}
	return newAccount, nil
}

func (s *customerService) AccountInquiry(ctx context.Context, accountNumber string) (*Account, error) {
	account, err := s.Repo.InquiryByAccountNumber(ctx, accountNumber)
	if err != nil {
		log.Printf("[service error] inquiry account error. %v", err)
		return nil, errors.New("customerService.AccountCreation returns error. please check the logs")
	}
	customer, err := s.Repo.GetCustomerByAccountNumber(ctx, accountNumber)
	if err != nil {
		log.Printf("[service error] inquiry account error. %v", err)
		return nil, errors.New("customerService.AccountCreation returns error. please check the logs")
	}
	account.Customer = customer
	return account, nil
}

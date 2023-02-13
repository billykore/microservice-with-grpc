package main

import (
	"context"
	pb "microservice-with-grpc/gen/customer/v1"
	"microservice-with-grpc/internal"
)

type customerService interface {
	accountCreation(ctx context.Context, req *pb.AccountCreationRequest) error
	accountInquiry(ctx context.Context, accountNumber string) (*account, error)
}

type customerServiceImpl struct {
	repo customerRepo
}

func newCustomerService(repo customerRepo) customerService {
	return &customerServiceImpl{repo: repo}
}

func (s *customerServiceImpl) accountCreation(ctx context.Context, req *pb.AccountCreationRequest) error {
	err := s.createCustomerWithAccount(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *customerServiceImpl) createCustomerWithAccount(ctx context.Context, req *pb.AccountCreationRequest) error {
	customer, err := s.createCustomerWithCif(ctx, req)
	if err != nil {
		return err
	}
	err = s.createCustomerCifAccount(ctx, customer, internal.SavingAccount)
	if err != nil {
		return err
	}
	return nil
}

func (s *customerServiceImpl) createCustomerWithCif(ctx context.Context, req *pb.AccountCreationRequest) (*customer, error) {
	customer := buildCustomer(req)
	newCif, err := s.createCif(ctx)
	if err != nil {
		return nil, err
	}
	customer.Cif = newCif
	return customer, nil
}

func (s *customerServiceImpl) createCif(ctx context.Context) (string, error) {
	lastCif, err := s.repo.getLastCif(ctx)
	if err != nil {
		return "", err
	}
	newCif := internal.BuildNewCif(lastCif)
	return newCif, nil
}

func (s *customerServiceImpl) createCustomerCifAccount(ctx context.Context, customer *customer, accountType internal.AccountType) error {
	newAccountNumber, err := s.createAccountNumber(ctx, accountType)
	if err != nil {
		return err
	}
	account := buildAccount(customer.Cif, newAccountNumber, accountType)
	err = s.repo.createAccount(ctx, account)
	if err != nil {
		return err
	}
	return nil
}

func (s *customerServiceImpl) createAccountNumber(ctx context.Context, accType internal.AccountType) (string, error) {
	lastAccountNumber, err := s.repo.getLastAccount(ctx)
	if err != nil {
		return "", err
	}
	var newAccount string
	switch accType {
	case internal.SavingAccount:
		newAccount = internal.BuildNewAccountNumber(lastAccountNumber)
		break
	case internal.GiroAccount:
		newAccount = internal.BuildNewAccountNumber(lastAccountNumber)
		break
	}
	return newAccount, nil
}

func (s *customerServiceImpl) accountInquiry(ctx context.Context, accountNumber string) (*account, error) {
	account, err := s.inquiryCustomerWithAccount(ctx, accountNumber)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *customerServiceImpl) inquiryCustomerWithAccount(ctx context.Context, accountNumber string) (*account, error) {
	account, err := s.repo.inquiryByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, err
	}
	customer, err := s.repo.getCustomerByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, err
	}
	account.Customer = customer
	return account, nil
}

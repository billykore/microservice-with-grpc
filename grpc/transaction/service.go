package main

import (
	"context"
	"errors"
	"log"
	"sync"

	"microservice-with-grpc/entity"
	customerpb "microservice-with-grpc/gen/customer/v1"
)

type TransactionService interface {
	Transfer(ctx context.Context, req *Request) error
}

type transactionService struct {
	Repo TransactionRepo
}

func NewTransactionService(repo TransactionRepo) TransactionService {
	return &transactionService{
		Repo: repo,
	}
}

func (s *transactionService) Transfer(ctx context.Context, req *Request) error {
	sourceAccount, err := s.getSourceAccount(ctx, req.SourceAccount)
	if err != nil {
		return errors.New("transfer error")
	}
	destinationAccount, err := s.getDestinationAccount(ctx, req.DestinationAccount)
	if err != nil {
		return errors.New("transfer error")
	}

	errCh := make(chan error)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		TransferBalance(sourceAccount, destinationAccount, req.Amount)
		err := s.Repo.UpdateBalanceByAccountNumber(ctx, sourceAccount.AccountNumber, sourceAccount.Balance)
		if err != nil {
			log.Printf("[transaction service error] update balance error: %v", err)
			errCh <- err
		}
		err = s.Repo.UpdateBalanceByAccountNumber(ctx, destinationAccount.AccountNumber, destinationAccount.Balance)
		if err != nil {
			log.Printf("[transaction service error] update balance error: %v", err)
			errCh <- err
		}
	}()

	wg.Wait()
	select {
	case err = <-errCh:
		log.Printf("[transaction service error] transfer error: %v", err)
		return errors.New("transfer error")
	default:
		return nil
	}
}

func (s *transactionService) getSourceAccount(ctx context.Context, sourceAccount string) (*entity.Account, error) {
	customerClient, closer := CustomerClient()
	defer closer()
	account, err := customerClient.AccountInquiry(ctx, &customerpb.InquiryRequest{AccountNumber: sourceAccount})
	if err != nil {
		log.Printf("[transaction service error] inquiry account error: %v", err)
		return nil, errors.New("inquiry account error")
	}
	return &entity.Account{
		Cif:            account.GetCif(),
		AccountNumber:  account.GetAccountNumber(),
		AccountType:    account.GetAccountType(),
		Name:           account.GetName(),
		Currency:       account.GetCurrency(),
		Status:         account.GetStatus(),
		Blocked:        account.GetBlocked(),
		Balance:        account.GetBalance(),
		MinimumBalance: account.GetMinimumBalance(),
		ProductType:    account.GetProductType(),
	}, nil
}

func (s *transactionService) getDestinationAccount(ctx context.Context, destinationAccount string) (*entity.Account, error) {
	customerClient, closer := CustomerClient()
	defer closer()
	account, err := customerClient.AccountInquiry(ctx, &customerpb.InquiryRequest{AccountNumber: destinationAccount})
	if err != nil {
		log.Printf("[transaction service error] inquiry account error: %v", err)
		return nil, errors.New("inquiry account error")
	}
	return &entity.Account{
		Cif:            account.GetCif(),
		AccountNumber:  account.GetAccountNumber(),
		AccountType:    account.GetAccountType(),
		Name:           account.GetName(),
		Currency:       account.GetCurrency(),
		Status:         account.GetStatus(),
		Blocked:        account.GetBlocked(),
		Balance:        account.GetBalance(),
		MinimumBalance: account.GetMinimumBalance(),
		ProductType:    account.GetProductType(),
	}, nil
}

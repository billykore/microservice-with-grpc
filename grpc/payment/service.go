package main

import (
	"context"
	"errors"
	"log"

	"microservice-with-grpc/entity"
	transactionpb "microservice-with-grpc/gen/transaction/v1"
)

type PaymentService interface {
	Qris(ctx context.Context, req *Request) (bool, error)
}

type paymentService struct {
	Repo PaymentRepo
}

func NewPaymentService(repo PaymentRepo) PaymentService {
	return &paymentService{Repo: repo}
}

func (s *paymentService) Qris(ctx context.Context, req *Request) (bool, error) {
	if req == nil {
		return false, errors.New("qris payment error. request shouldn't be empty")
	}
	err := s.transfer(ctx, req)
	if err != nil {
		return false, errors.New("qris payment error")
	}
	err = s.insertLog(ctx, req)
	if err != nil {
		return false, errors.New("qris payment error")
	}
	return true, nil
}

func (s *paymentService) transfer(ctx context.Context, req *Request) error {
	transferClient, closer := TransactionClient()
	defer closer()
	response, err := transferClient.Transfer(ctx, &transactionpb.TransferRequest{
		TrxId:              req.TrxNumber,
		SourceAccount:      req.SourceAccount,
		DestinationAccount: req.DestinationAccount,
		Amount:             req.Amount,
	})
	if err != nil {
		log.Printf("[payment service error] transfer error: %v", err)
		return errors.New("transfer error. " + response.Message)
	}
	log.Printf("[payment service] transfer success. %v", response.Message)
	return nil
}

func (s *paymentService) insertLog(ctx context.Context, req *Request) error {
	qrisLog := &entity.QrisLog{
		MerchantId:         req.MerchantId,
		TrxNumber:          req.TrxNumber,
		AccountSource:      req.SourceAccount,
		AccountDestination: req.DestinationAccount,
		Amount:             req.Amount,
	}
	err := s.Repo.InsertQrisLog(ctx, qrisLog)
	if err != nil {
		log.Printf("[service error] insert qris log error: %v", err)
		return errors.New("error insert qris log")
	}
	return nil
}

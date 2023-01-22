package main

import (
	"context"
	"errors"
	"log"
	"microservice-with-grpc/entity"
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
	qrisLog := &entity.QrisLog{
		MerchantId:         req.MerchantId,
		TrxNumber:          req.TrxNumber,
		AccountSource:      req.AccountSource,
		AccountDestination: req.AccountDestination,
		Amount:             req.Amount,
	}
	err := s.Repo.InsertQrisLog(ctx, qrisLog)
	if err != nil {
		log.Printf("[service error] qris payment error, %v", err)
		return false, errors.New("qris payment error")
	}
	return true, nil
}

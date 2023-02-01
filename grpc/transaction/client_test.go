package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	pb "microservice-with-grpc/gen/customer/v1"
)

func TestCustomerClient(t *testing.T) {
	client, closer := CustomerClient()
	assert.NotNil(t, client)
	assert.NotNil(t, closer)
}

func TestCustomerClient_InquiryAccount(t *testing.T) {
	client, closer := CustomerClient()
	defer closer()
	account, err := client.AccountInquiry(context.Background(), &pb.InquiryRequest{AccountNumber: "001001000002300"})
	assert.NotNil(t, account)
	assert.Nil(t, err)
}

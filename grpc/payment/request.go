package main

type Request struct {
	MerchantId         string
	TrxNumber          string
	AccountSource      string
	AccountDestination string
	Amount             string
}

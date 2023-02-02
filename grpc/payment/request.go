package main

type Request struct {
	MerchantId         string
	TrxNumber          string
	SourceAccount      string
	DestinationAccount string
	Amount             string
}

package request

type QrisPayment struct {
	MerchantId         string
	TrxNumber          string
	AccountSource      string
	AccountDestination string
	Amount             string
}

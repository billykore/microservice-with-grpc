package request

type QrisPayment struct {
	MerchantId         string `validate:"required"`
	TrxNumber          string `validate:"required"`
	AccountSource      string `validate:"required"`
	AccountDestination string `validate:"required"`
	Amount             string `validate:"required"`
}

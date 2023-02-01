package handler

// Handlers are struct to put all the api handler that will be use by the router.
//
// Register api handlers here.
type Handlers struct {
	Auth         *AuthHandler
	Customer     *CustomerHandler
	Payment      *PaymentHandler
	Notification *NotificationHandler
}

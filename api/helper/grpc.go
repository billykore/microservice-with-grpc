package helper

import (
	"microservice-with-grpc/api/request"

	customerpb "microservice-with-grpc/gen/customer/v1"
	paymentpb "microservice-with-grpc/gen/payment/v1"
)

var (
	genderValue = map[string]customerpb.Gender{
		"Male":   customerpb.Gender_MALE,
		"Female": customerpb.Gender_FEMALE,
	}
	religionValue = map[string]customerpb.Religion{
		"Islam":      customerpb.Religion_ISLAM,
		"Protestant": customerpb.Religion_PROTESTANT,
		"Catholic":   customerpb.Religion_CATHOLIC,
		"Hindu":      customerpb.Religion_HINDU,
		"Buddha":     customerpb.Religion_BUDDHA,
		"Konghucu":   customerpb.Religion_KONGHUCU,
	}
	marriageStatusValue = map[string]customerpb.MarriageStatus{
		"NotMarried": customerpb.MarriageStatus_NOT_MARRIED,
		"Married":    customerpb.MarriageStatus_MARRIED,
	}
	citizenValue = map[string]customerpb.Citizen{
		"WNI": customerpb.Citizen_WNI,
		"WNA": customerpb.Citizen_WNA,
	}
)

func BuildCustomerGrpcRequest(req *request.CreateAccount) *customerpb.AccountCreationRequest {
	return &customerpb.AccountCreationRequest{
		Nik:            req.Nik,
		Name:           req.Name,
		Pob:            req.Pob,
		Dob:            req.Dob,
		Address:        req.Address,
		Profession:     req.Profession,
		Gender:         genderValue[req.Gender],
		Religion:       religionValue[req.Religion],
		MarriageStatus: marriageStatusValue[req.MarriageStatus],
		Citizen:        citizenValue[req.Citizen],
	}
}

func BuildQrisPaymentGrpcRequest(req *request.QrisPayment) *paymentpb.QrisRequest {
	return &paymentpb.QrisRequest{
		MerchantId:         req.MerchantId,
		TrxNumber:          req.TrxNumber,
		AccountSource:      req.AccountSource,
		AccountDestination: req.AccountDestination,
		Amount:             req.Amount,
	}
}

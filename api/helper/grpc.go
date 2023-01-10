package helper

import (
	"microservice-with-grpc/api/request"
	pb "microservice-with-grpc/gen/customer/v1"
)

var (
	genderValue = map[string]pb.Gender{
		"Male":   pb.Gender_MALE,
		"Female": pb.Gender_FEMALE,
	}
	religionValue = map[string]pb.Religion{
		"Islam":      pb.Religion_ISLAM,
		"Protestant": pb.Religion_PROTESTANT,
		"Catholic":   pb.Religion_CATHOLIC,
		"Hindu":      pb.Religion_HINDU,
		"Buddha":     pb.Religion_BUDDHA,
		"Konghucu":   pb.Religion_KONGHUCU,
	}
	marriageStatusValue = map[string]pb.MarriageStatus{
		"NotMarried": pb.MarriageStatus_NOT_MARRIED,
		"Married":    pb.MarriageStatus_MARRIED,
	}
	citizenValue = map[string]pb.Citizen{
		"WNI": pb.Citizen_WNI,
		"WNA": pb.Citizen_WNA,
	}
)

func BuildGrpcRequest(req *request.CreateAccount) *pb.AccountCreationRequest {
	return &pb.AccountCreationRequest{
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

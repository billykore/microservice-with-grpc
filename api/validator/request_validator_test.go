package validator

import (
	"github.com/stretchr/testify/assert"
	"microservice-with-grpc/api/request"
	"testing"
)

func TestValidateRequestBody(t *testing.T) {
	type args struct {
		requestBody any
	}

	type expectation struct {
		valid    bool
		errorMsg string
	}

	// test cases
	tests := map[string]struct {
		args     args
		expected expectation
	}{
		"valid_qris_payment_request": {
			args: args{
				requestBody: &request.QrisPayment{
					MerchantId:         "example-id-123",
					TrxNumber:          "12345",
					AccountSource:      "123456789",
					AccountDestination: "987654321",
					Amount:             "50000",
				},
			},
			expected: expectation{
				valid:    true,
				errorMsg: "",
			},
		},
		"empty_qris_payment_request": {
			args: args{
				requestBody: &request.QrisPayment{
					MerchantId:         "",
					TrxNumber:          "",
					AccountSource:      "",
					AccountDestination: "",
					Amount:             "",
				},
			},
			expected: expectation{
				valid:    false,
				errorMsg: "MerchantId cannot be empty",
			},
		},
		"valid_create_account_request": {
			args: args{
				requestBody: &request.CreateAccount{
					Nik:            "123456789",
					Name:           "Oyen",
					Pob:            "Surabaya",
					Dob:            "31/08/2000",
					Address:        "Jakarta",
					Profession:     "Banker",
					Gender:         "Female",
					Religion:       "Protestant",
					MarriageStatus: "Not Married",
					Citizen:        "WNI",
				},
			},
			expected: expectation{
				valid:    true,
				errorMsg: "",
			},
		},
		"empty_create_account_request": {
			args: args{
				requestBody: &request.CreateAccount{
					Nik:            "",
					Name:           "",
					Pob:            "",
					Dob:            "",
					Address:        "",
					Profession:     "",
					Gender:         "",
					Religion:       "",
					MarriageStatus: "",
					Citizen:        "",
				},
			},
			expected: expectation{
				valid:    false,
				errorMsg: "Nik cannot be empty",
			},
		},
	}

	// run the test
	for scenario, test := range tests {
		t.Run(scenario, func(t *testing.T) {
			// call the validator function
			valid, errMsg := ValidateRequestBody(test.args.requestBody)
			// check the output
			assert.Equal(t, test.expected.valid, valid)
			assert.Equal(t, test.expected.errorMsg, errMsg)
		})
	}
}

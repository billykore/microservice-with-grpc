package main

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"microservice-with-grpc/database"
	pb "microservice-with-grpc/gen/customer/v1"
	"microservice-with-grpc/internal"
)

type customerRepoMock struct {
	mock.Mock
}

func (m *customerRepoMock) CreateCustomer(ctx context.Context, model *Customer) error {
	args := m.Mock.Called(ctx, model)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func (m *customerRepoMock) GetLastCif(ctx context.Context) (string, error) {
	args := m.Mock.Called(ctx)
	if args.Get(0) == "" {
		return args.Get(0).(string), args.Get(1).(error)
	}
	if args.Get(1) != nil {
		return args.Get(0).(string), args.Get(1).(error)
	}
	return args.Get(0).(string), nil
}

func (m *customerRepoMock) GetLastAccount(ctx context.Context) (string, error) {
	args := m.Mock.Called(ctx)
	if args.Get(0) == "" && args.Get(1) != nil {
		return "", args.Get(1).(error)
	}
	return args.Get(0).(string), nil
}

func (m *customerRepoMock) CreateAccount(ctx context.Context, account *Account) error {
	args := m.Mock.Called(ctx, account)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func (m *customerRepoMock) InquiryByAccountNumber(ctx context.Context, accountNumber string) (*Account, error) {
	args := m.Mock.Called(ctx, accountNumber)
	if args.Get(0) == nil && args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*Account), nil
}

func (m *customerRepoMock) GetCustomerByAccountNumber(ctx context.Context, accountNumber string) (*Customer, error) {
	//TODO implement me
	panic("implement me")
}

func TestCustomerService_AccountCreation(t *testing.T) {
	type args struct {
		ctx context.Context
		in  *pb.AccountCreationRequest
	}

	type expectation struct {
		err error
	}

	tests := map[string]struct {
		args     args
		expected expectation
	}{
		"success": {
			args: args{
				ctx: context.Background(),
				in: &pb.AccountCreationRequest{
					Nik:            "5310121711980001",
					Name:           "Evanbill Antonio Kore",
					Pob:            "Bajawa",
					Dob:            "17/11/1998",
					Address:        "Jakarta",
					Profession:     "Officer",
					Gender:         pb.Gender_MALE,
					Religion:       pb.Religion_PROTESTANT,
					MarriageStatus: pb.MarriageStatus_NOT_MARRIED,
					Citizen:        pb.Citizen_WNI,
				},
			},
			expected: expectation{
				err: nil,
			},
		},
		"failed": {
			args: args{
				ctx: context.Background(),
				in:  &pb.AccountCreationRequest{},
			},
			expected: expectation{
				err: errors.New("failed to create new account"),
			},
		},
	}

	for scenario, test := range tests {
		t.Run(scenario, func(t *testing.T) {
			customer := BuildCustomer(test.args.in)
			customer.Cif = internal.BuildNewCif("0000011111")
			repo := &customerRepoMock{Mock: mock.Mock{}}
			repo.Mock.On("CreateCustomer", test.args.ctx, customer).Return(test.expected.err)
			repo.Mock.On("GetLastCif", test.args.ctx).Return("0000011111", nil)
			repo.Mock.On("GetLastAccount", test.args.ctx).Return("001001001111300", nil)
			repo.Mock.On("CreateAccount", test.args.ctx, mock.Anything).Return(nil)
			service := NewCustomerService(repo)
			out := service.AccountCreation(test.args.ctx, test.args.in)
			assert.NotNil(t, service)
			assert.Equal(t, test.expected.err, out)
		})
	}
}

func TestCustomerServiceIntegrationTest_AccountCreation(t *testing.T) {
	type args struct {
		ctx context.Context
		in  *pb.AccountCreationRequest
	}

	type expectation struct {
		err error
	}

	// test cases
	tests := map[string]struct {
		args     args
		expected expectation
	}{
		"success": {
			args: args{
				ctx: context.Background(),
				in: &pb.AccountCreationRequest{
					Nik:            "0101013108000001",
					Name:           "Florence Fedora Agustina",
					Pob:            "Surabaya",
					Dob:            "31/08/2000",
					Address:        "Jakarta",
					Profession:     "Banker",
					Gender:         pb.Gender_FEMALE,
					Religion:       pb.Religion_PROTESTANT,
					MarriageStatus: pb.MarriageStatus_NOT_MARRIED,
					Citizen:        pb.Citizen_WNI,
				},
			},
			expected: expectation{
				err: nil,
			},
		},
		"failed": {
			args: args{
				ctx: context.Background(),
				// account nik already exist, account creation must fail.
				in: &pb.AccountCreationRequest{
					Nik:            "0101013108000001",
					Name:           "Florence Fedora Agustina",
					Pob:            "Surabaya",
					Dob:            "31/08/2000",
					Address:        "Jakarta",
					Profession:     "Banker",
					Gender:         pb.Gender_FEMALE,
					Religion:       pb.Religion_PROTESTANT,
					MarriageStatus: pb.MarriageStatus_NOT_MARRIED,
					Citizen:        pb.Citizen_WNI,
				},
			},
			expected: expectation{
				err: errors.New("customerService.AccountCreation returns error. please check the logs"),
			},
		},
	}

	db := database.New(database.MySQL, &database.Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "localhost",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
	assert.NotNil(t, db)
	assert.NotNil(t, db)
	err := database.Migrate(db, &Customer{})
	assert.NoError(t, err)
	repo := NewCustomerRepo(db)
	assert.NotNil(t, repo)
	service := NewCustomerService(repo)
	assert.NotNil(t, service)

	for scenario, test := range tests {
		t.Run(scenario, func(t *testing.T) {
			out := service.AccountCreation(test.args.ctx, test.args.in)
			assert.Equal(t, test.expected.err, out)
		})
	}
}

func TestCustomerService_AccountInquiry(t *testing.T) {
	type args struct {
		ctx           context.Context
		accountNumber string
	}

	type expectation struct {
		account *Account
		err     error
	}

	tests := map[string]struct {
		args     args
		expected expectation
	}{
		"success": {
			args: args{
				ctx:           context.Background(),
				accountNumber: "001001000001300",
			},
			expected: expectation{
				account: &Account{
					Cif:              "0000000003",
					AccountNumber:    "001001000002300",
					Type:             "S",
					Balance:          "0.00",
					MinimumBalance:   "0.00",
					AvailableBalance: "0.00",
					Status:           "1",
					Currency:         "IDR",
					ProductType:      "000005",
					Blocked:          "0",
				},
				err: nil,
			},
		},
		"notfound": {
			args: args{
				ctx:           context.Background(),
				accountNumber: "001001000001399",
			},
			expected: expectation{
				account: nil,
				err:     errors.New("account account not found"),
			},
		},
	}

	for scenario, test := range tests {
		t.Run(scenario, func(t *testing.T) {
			repo := &customerRepoMock{Mock: mock.Mock{}}
			repo.Mock.On("InquiryByAccountNumber", test.args.ctx, test.args.accountNumber).Return(test.expected.account, test.expected.err)
			service := NewCustomerService(repo)
			out, err := service.AccountInquiry(test.args.ctx, test.args.accountNumber)
			assert.Equal(t, test.expected.account, out)
			assert.Equal(t, test.expected.err, err)
		})
	}
}

func TestCustomerServiceIntegrationTest_AccountInquiry(t *testing.T) {
	type args struct {
		ctx           context.Context
		accountNumber string
	}

	type expectation struct {
		account *Account
		err     error
	}

	tests := map[string]struct {
		args     args
		expected expectation
	}{
		"success": {
			args: args{
				ctx:           context.Background(),
				accountNumber: "001001000002300",
			},
			expected: expectation{
				account: &Account{
					Cif:              "0000000003",
					AccountNumber:    "001001000002300",
					Type:             "S",
					Balance:          "0.00",
					MinimumBalance:   "0.00",
					AvailableBalance: "0.00",
					Status:           "1",
					Currency:         "IDR",
					ProductType:      "000005",
					Blocked:          "0",
				},
				err: nil,
			},
		},
		"notfound": {
			args: args{
				ctx:           context.Background(),
				accountNumber: "001001000001399",
			},
			expected: expectation{
				account: nil,
				err:     errors.New("customerService.AccountCreation returns error. please check the logs"),
			},
		},
	}

	db := database.New(database.MySQL, &database.Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "localhost",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
	assert.NotNil(t, db)
	assert.NotNil(t, db)
	repo := NewCustomerRepo(db)
	assert.NotNil(t, repo)
	service := NewCustomerService(repo)
	assert.NotNil(t, service)

	for scenario, test := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := service.AccountInquiry(test.args.ctx, test.args.accountNumber)
			assert.IsType(t, test.expected.account, out)
			assert.Equal(t, test.expected.err, err)
			if out != nil {
				assert.Equal(t, test.expected.account.AccountNumber, out.AccountNumber)
				assert.Equal(t, test.expected.account.Cif, out.Cif)
				assert.Equal(t, "FLORENCE FEDORA AGUSTINA", out.Customer.Name)
			}
		})
	}
}

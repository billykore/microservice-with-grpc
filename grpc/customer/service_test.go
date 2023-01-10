package main

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"microservice-with-grpc/database"
	pb "microservice-with-grpc/gen/customer/v1"
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
		return "", errors.New("error get last account")
	}
	return args.Get(0).(string), nil
}

func (m *customerRepoMock) CreateAccount(ctx context.Context, account *Account) error {
	args := m.Mock.Called(ctx, account)
	if args.Get(0) != nil {
		return errors.New("error get last account")
	}
	return nil
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

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			customer := BuildCustomer(tt.args.in)
			customer.Cif = BuildNewCif("0000011111")
			repo := &customerRepoMock{Mock: mock.Mock{}}
			repo.Mock.On("CreateCustomer", tt.args.ctx, customer).Return(tt.expected.err)
			repo.Mock.On("GetLastCif", tt.args.ctx).Return("0000011111", nil)
			repo.Mock.On("GetLastAccount", tt.args.ctx).Return("001001001111300", nil)
			repo.Mock.On("CreateAccount", tt.args.ctx, mock.Anything).Return(nil)
			service := NewCustomerService(repo)
			out := service.AccountCreation(tt.args.ctx, tt.args.in)
			assert.NotNil(t, service)
			assert.Equal(t, tt.expected.err, out)
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
				// customer nik already exist, account creation must fail.
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

	db, err := database.New(&database.Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "localhost",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
	assert.NoError(t, err)
	assert.NotNil(t, db)
	err = database.Migrate(db, &Customer{})
	assert.NoError(t, err)
	repo := NewCustomerRepo(db)
	assert.NotNil(t, repo)
	service := NewCustomerService(repo)
	assert.NotNil(t, service)

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out := service.AccountCreation(tt.args.ctx, tt.args.in)
			assert.Equal(t, tt.expected.err, out)
		})
	}
}

func TestName(t *testing.T) {
	lastAccount := "001001003935300"
	fmt.Println(lastAccount[8:12])
}

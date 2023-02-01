package request

type CreateAccount struct {
	Nik            string `validate:"required"`
	Name           string `validate:"required"`
	Pob            string `validate:"required"`
	Dob            string `validate:"required"`
	Address        string `validate:"required"`
	Profession     string `validate:"required"`
	Gender         string `validate:"required"`
	Religion       string `validate:"required"`
	MarriageStatus string `validate:"required"`
	Citizen        string `validate:"required"`
}

package validator

import "github.com/go-playground/validator/v10"

func ValidateRequestBody(body any) (bool, string) {
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		errors := err.(validator.ValidationErrors)
		for _, e := range errors {
			switch e.Tag() {
			case "required":
				return false, e.Field() + " cannot be empty"
			default:
			}
		}
	}
	return true, ""
}

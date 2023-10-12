package method

import (
	"xyz-books-codebase-one/model"

	"github.com/go-playground/validator/v10"
)

func FieldValidator(s interface{}) []model.ApiError {
	var apiErrors []model.ApiError
	err := Validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var fieldError model.ApiError
			fieldError.Param = err.StructField()
			fieldError.Message = FieldValidatorMessage(err)
			apiErrors = append(apiErrors, fieldError)
		}
	}
	return apiErrors
}

func FieldValidatorMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "This field is required"
	case "max":
		return "This field exceeds the maximum limit"
	}
	return fieldError.Error()
}



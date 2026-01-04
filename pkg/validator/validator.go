package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

var validate *validator.Validate

func Init() {
	validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func Get() *validator.Validate {
	if validate == nil {
		Init()
	}
	return validate
}

func Validate(data interface{}) []ErrorResponse {
	var errors []ErrorResponse

	if err := Get().Struct(data); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, ErrorResponse{
				Field:   err.Field(),
				Tag:     err.Tag(),
				Message: generateMessage(err),
			})
		}
	}

	return errors
}

func generateMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "email":
		return err.Field() + " must be a valid email"
	case "min":
		return err.Field() + " must be at least " + err.Param() + " characters"
	case "max":
		return err.Field() + " must be at most " + err.Param() + " characters"
	case "eqfield":
		return err.Field() + " must match " + err.Param()
	default:
		return err.Field() + " is invalid"
	}
}
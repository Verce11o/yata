package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

func NewValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return validate
}

func ValidationErrors(errs validator.ValidationErrors) string {
	var errMessages []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMessages = append(errMessages, fmt.Sprintf("field %s is a required field", err.Field()))
		case "min":
			errMessages = append(errMessages, fmt.Sprintf("field %s must be longer than %s", err.Field(), err.Param()))
		case "max":
			errMessages = append(errMessages, fmt.Sprintf("field %s must be shorter than %s", err.Field(), err.Param()))
		default:
			errMessages = append(errMessages, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}
	return strings.Join(errMessages, ",")
}

type ValidationError struct {
	errs validator.ValidationErrors
}

func (e ValidationError) Error() string {
	return ValidationErrors(e.errs)
}

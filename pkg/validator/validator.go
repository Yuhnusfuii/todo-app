package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidationError holds field-level errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validate validates a struct and returns a slice of ValidationErrors or nil
func Validate(s interface{}) []ValidationError {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}
	var errs []ValidationError
	for _, e := range err.(validator.ValidationErrors) {
		errs = append(errs, ValidationError{
			Field:   e.Field(),
			Message: msgForTag(e.Tag(), e.Param()),
		})
	}
	return errs
}

func msgForTag(tag, param string) string {
	switch tag {
	case "required":
		return "this field is required"
	case "email":
		return "invalid email address"
	case "min":
		return "value is too short (min " + param + ")"
	case "max":
		return "value is too long (max " + param + ")"
	case "oneof":
		return "must be one of: " + param
	default:
		return "invalid value"
	}
}

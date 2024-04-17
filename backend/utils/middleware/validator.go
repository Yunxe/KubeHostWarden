package middleware

import (
	"github.com/go-playground/validator/v10"
)

func ValidateReq(req interface{}) error {
	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		return err
	}
	return nil
}

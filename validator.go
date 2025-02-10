package main

import "github.com/go-playground/validator/v10"

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(a any) error {
	return cv.Validator.Struct(a)
}

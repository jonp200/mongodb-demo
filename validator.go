package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/labstack/echo/v4"
)

func NewValidator() echo.Validator {
	v := validator.New()
	_ = v.RegisterValidation("not_blank", validators.NotBlank)

	return &customValidator{validator: v}
}

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(a any) error {
	return cv.validator.Struct(a)
}

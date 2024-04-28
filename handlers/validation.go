package handlers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator"
	ct "github.com/kanawat2566/assessment-tax/constants"
	"github.com/labstack/echo/v4"
)

func validateInput(input interface{}) error {

	validate := validator.New()
	var msgs []string

	err := validate.Struct(input)
	if err != nil {

		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErrs {

				fieldName := e.Field()
				tagValue := e.Tag()
				message := fmt.Sprintf("Field '%s' failed validation for '%s'", fieldName, tagValue)
				msgs = append(msgs, message)
			}
			return errors.New(strings.Join(msgs, ",\n"))
		}
	}
	return nil
}

func BindWithValidate(c echo.Context, rq interface{}) error {

	if err := c.Bind(rq); err != nil {
		return errors.New(ct.ErrInvalidFormatReq)
	}
	if err := validateInput(rq); err != nil {
		return err
	}
	return nil
}

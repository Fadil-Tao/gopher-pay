package customvalidator

import (
	"regexp"
	"github.com/go-playground/validator/v10"
)

var amountRegex = regexp.MustCompile(`^[0-9]+(\.[0-9]{1,2})?$`) // e.g., 100, 100.50

func AmountValidator(fl validator.FieldLevel) bool {
	return amountRegex.MatchString(fl.Field().String())
}

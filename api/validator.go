package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/minhnghia2k3/simple_bank/util"
)

var validCurrency validator.Func = func(fieldlevel validator.FieldLevel) bool {
	if currency, ok := fieldlevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}

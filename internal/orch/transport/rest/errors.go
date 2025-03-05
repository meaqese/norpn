package rest

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrExpressionIsInvalid = errors.New("expression is not valid")
	ErrJsonValidation      = errors.New("json validation error")
)

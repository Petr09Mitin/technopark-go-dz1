package mathparser

import "errors"

var (
	ErrInvalidExpression = errors.New("invalid input expression")
	ErrInvalidParenths   = errors.New("invalid parentheses sequence in input expression")
	ErrZeroDivision      = errors.New("zero division error")
)

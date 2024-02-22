package mathparser_test

import (
	"testing"

	"github.com/Petr09Mitin/technopark-go-dz1/calculator/mathparser"
	"github.com/stretchr/testify/assert"
)

type Expected struct {
	result float64
	err    error
}

var testCases = map[string]struct {
	input    string
	expected Expected
}{
	"addition": {
		input: "2+2",
		expected: Expected{
			result: 4,
			err:    nil,
		},
	},
	"substraction": {
		input: "2-1",
		expected: Expected{
			result: 1,
			err:    nil,
		},
	},
	"multiplication": {
		input: "2*2-1",
		expected: Expected{
			result: 3,
			err:    nil,
		},
	},
	"division": {
		input: "2/2+1",
		expected: Expected{
			result: 2,
			err:    nil,
		},
	},
	"negative number addition": {
		input: "-5 + 3",
		expected: Expected{
			result: -2,
			err:    nil,
		},
	},
	"negative number addition 2": {
		input: "3 + -5",
		expected: Expected{
			result: -2,
			err:    nil,
		},
	},
	"fractions addition": {
		input: "1 + 1.5",
		expected: Expected{
			result: 2.5,
			err:    nil,
		},
	},
	"fractions division": {
		input: "2.4 / 4.8",
		expected: Expected{
			result: 0.5,
			err:    nil,
		},
	},
	"negative fractions division": {
		input: "2.4 / -4.8",
		expected: Expected{
			result: -0.5,
			err:    nil,
		},
	},
	"parentheses": {
		input: "(2+2)*2",
		expected: Expected{
			result: 8,
			err:    nil,
		},
	},
	"complex": {
		input: "((5 - 3) / (3 + 5)) * 3 * 3",
		expected: Expected{
			result: 2.25,
			err:    nil,
		},
	},
	"complex 2": {
		input: "(-1.4) / 0.7 + 2.5 * 2 - 5",
		expected: Expected{
			result: -2,
			err:    nil,
		},
	},
	"complex 3": {
		input: "0.1 * 10 + 4 / -2 + 2 * 1.0",
		expected: Expected{
			result: 1,
			err:    nil,
		},
	},
	"complex 4": {
		input: "(2 / -2) * -2 * .1 / .1",
		expected: Expected{
			result: 2,
			err:    nil,
		},
	},
	"zero division": {
		input: "1 / 0",
		expected: Expected{
			result: 0,
			err:    mathparser.ErrZeroDivision,
		},
	},
	"invalid operator": {
		input: "1 % 0",
		expected: Expected{
			result: 0,
			err:    mathparser.ErrInvalidExpression,
		},
	},
	"invalid operand": {
		input: "1 + 12,12",
		expected: Expected{
			result: 0,
			err:    mathparser.ErrInvalidExpression,
		},
	},
	"invalid left parenthese": {
		input: "1 + (2",
		expected: Expected{
			result: 0,
			err:    mathparser.ErrInvalidParenths,
		},
	},
	"invalid left parenthese 2": {
		input: "1(2",
		expected: Expected{
			result: 0,
			err:    mathparser.ErrInvalidExpression,
		},
	},
	"invalid right parenthese": {
		input: "1 + 2)",
		expected: Expected{
			result: 0,
			err:    mathparser.ErrInvalidParenths,
		},
	},
	"invalid expression": {
		input: "1 + 2 +",
		expected: Expected{
			result: 0,
			err:    mathparser.ErrInvalidExpression,
		},
	},
	"empty string": {
		input: "",
		expected: Expected{
			result: 0,
			err:    mathparser.ErrInvalidExpression,
		},
	},
}

func TestCalculateExpression(t *testing.T) {
	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			result, err := mathparser.CalculateExpression(test.input)
			assert.Equal(t, test.expected.result, result)
			if test.expected.err != nil {
				assert.ErrorIs(t, err, test.expected.err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

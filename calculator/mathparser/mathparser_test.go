package mathparser_test

import (
	"errors"
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
	"zero division": {
		input: "1 / 0",
		expected: Expected{
			result: 0,
			err:    errors.New("zero division error"),
		},
	},
	"invalid operator": {
		input: "1 % 0",
		expected: Expected{
			result: 0,
			err:    errors.New("invalid operator in input expression"),
		},
	},
	"invalid operand": {
		input: "1 + 12,12",
		expected: Expected{
			result: 0,
			err:    errors.New("invalid operands in input expression"),
		},
	},
	"invalid left parenthese": {
		input: "1 + (2",
		expected: Expected{
			result: 0,
			err:    errors.New("invalid parentheses sequence in input expression"),
		},
	},
	"invalid left parenthese 2": {
		input: "1(2",
		expected: Expected{
			result: 0,
			err:    errors.New("invalid input expression"),
		},
	},
	"invalid right parenthese": {
		input: "1 + 2)",
		expected: Expected{
			result: 0,
			err:    errors.New("invalid parentheses sequence in input expression"),
		},
	},
	"invalid expression": {
		input: "1 + 2 +",
		expected: Expected{
			result: 0,
			err:    errors.New("invalid input expression"),
		},
	},
}

func TestCalculateExpression(t *testing.T) {
	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			result, err := mathparser.CalculateExpression(test.input)
			assert.Equal(t, test.expected.result, result)
			assert.Equal(t, test.expected.err, err)
		})
	}
}

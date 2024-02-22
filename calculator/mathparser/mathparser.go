package mathparser

import (
	"regexp"
	"strconv"

	"github.com/Petr09Mitin/technopark-go-dz1/calculator/datastruct"
)

const (
	plus         = "+"
	minus        = "-"
	multiply     = "*"
	divide       = "/"
	leftParenth  = "("
	rightParenth = ")"
)

var anyNumberRegexp, _ = regexp.Compile(`[-]?[\d]+\.?[\d]*$`)

var OperatorsData = map[string]struct {
	priority          uint
	isLeftAssociative bool
}{
	"+": {
		priority:          1,
		isLeftAssociative: true,
	},
	"-": {
		priority:          1,
		isLeftAssociative: true,
	},
	"*": {
		priority:          2,
		isLeftAssociative: true,
	},
	"/": {
		priority:          2,
		isLeftAssociative: false,
	},
}

// parseTokensFromExpression parses mathmatical tokens from expression
func parseTokensFromExpression(expression string) (tokens []string, err error) {
	number := ""
	operandsCount, operatorsCount := 0, 0
	for _, char := range expression {
		char := string(char)
		switch char {
		case plus:
			if number != "" {
				tokens = append(tokens, number)
				operandsCount++
				number = ""
			}
			operatorsCount++
			tokens = append(tokens, char)
		case minus:
			if number != "" {
				tokens = append(tokens, number)
				operandsCount++
				number = ""
			}

			var lastToken string
			if len(tokens) > 0 {
				lastToken = tokens[len(tokens)-1]
			}
			if lastToken == rightParenth || anyNumberRegexp.MatchString(lastToken) {
				operatorsCount++
				tokens = append(tokens, char)
			} else {
				number = minus
			}
		case multiply:
			if number != "" {
				tokens = append(tokens, number)
				operandsCount++
				number = ""
			}
			operatorsCount++
			tokens = append(tokens, char)
		case divide:
			if number != "" {
				tokens = append(tokens, number)
				operandsCount++
				number = ""
			}
			operatorsCount++
			tokens = append(tokens, char)
		case leftParenth:
			if number != "" {
				tokens = append(tokens, number)
				operandsCount++
				number = ""
			}
			tokens = append(tokens, char)
		case rightParenth:
			if number != "" {
				tokens = append(tokens, number)
				operandsCount++
				number = ""
			}
			tokens = append(tokens, char)
		case " ":
			if number != "" {
				tokens = append(tokens, number)
				operandsCount++
				number = ""
			}
		default:
			number += char
		}
	}

	if number != "" {
		tokens = append(tokens, number)
		operandsCount++
	}

	if operandsCount != operatorsCount+1 {
		err = ErrInvalidExpression
	}
	return
}

func shouldParseStack(stack *datastruct.Stack[string], parsedToken string) (shouldParse bool) {
	return !stack.IsEmpty() && (OperatorsData[stack.Peek()].priority > OperatorsData[parsedToken].priority || OperatorsData[stack.Peek()].priority == OperatorsData[parsedToken].priority && OperatorsData[parsedToken].isLeftAssociative)
}

func parseRPNFromExpression(expression string) (rpn *datastruct.Queue[string], err error) {
	tokens, err := parseTokensFromExpression(expression)
	if err != nil {
		return
	}

	stack := datastruct.NewStack[string]()
	queue := datastruct.NewQueue[string]()

	for _, token := range tokens {
		if anyNumberRegexp.MatchString(token) {
			queue.Enqueue(token)
			continue
		}

		if _, ok := OperatorsData[token]; ok {
			for shouldParseStack(stack, token) {
				queue.Enqueue(stack.Pop())
			}
			stack.Push(token)
			continue
		}

		if token == leftParenth {
			stack.Push(token)
			continue
		}

		if token == rightParenth {
			stackToken := stack.Pop()
			for ; !stack.IsEmpty() && stackToken != leftParenth; stackToken = stack.Pop() {
				queue.Enqueue(stackToken)
			}

			if stackToken != leftParenth {
				err = ErrInvalidParenths
				return
			}
			continue
		}

		err = ErrInvalidExpression
		return
	}

	for !stack.IsEmpty() {
		stackToken := stack.Pop()
		if stackToken == leftParenth {
			err = ErrInvalidParenths
			return
		}
		queue.Enqueue(stackToken)
	}
	return queue, nil
}

func evalOperator(operand1, operand2, operator string) (result string, err error) {
	floatOperand1, err1 := strconv.ParseFloat(operand1, 64)
	floatOperand2, err2 := strconv.ParseFloat(operand2, 64)
	if err1 != nil || err2 != nil {
		err = ErrInvalidExpression
	}
	switch operator {
	case plus:
		result = strconv.FormatFloat(floatOperand1+floatOperand2, 'f', -1, 64)
	case minus:
		result = strconv.FormatFloat(floatOperand1-floatOperand2, 'f', -1, 64)
	case multiply:
		result = strconv.FormatFloat(floatOperand1*floatOperand2, 'f', -1, 64)
	case divide:
		if floatOperand2 == 0 {
			err = ErrZeroDivision
			return
		}
		result = strconv.FormatFloat(floatOperand1/floatOperand2, 'f', -1, 64)
	default:
		err = ErrInvalidExpression
		return
	}
	return
}

func CalculateExpression(expression string) (result float64, err error) {
	rpn, err := parseRPNFromExpression(expression)
	if err != nil {
		return
	}

	stack := datastruct.NewStack[string]()

	for token := rpn.Dequeue(); token != ""; token = rpn.Dequeue() {
		if anyNumberRegexp.MatchString(token) {
			stack.Push(token)
			continue
		}

		if _, ok := OperatorsData[token]; ok {
			operand2, operand1 := stack.Pop(), stack.Pop()
			evalResult, evalErr := evalOperator(operand1, operand2, token)

			if evalErr != nil {
				err = evalErr
				return
			}

			stack.Push(evalResult)
			continue
		}

		err = ErrInvalidExpression
		return
	}

	result, err = strconv.ParseFloat(stack.Pop(), 64)
	return
}

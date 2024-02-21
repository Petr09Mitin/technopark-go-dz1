package main

import (
	"fmt"
	"os"

	"github.com/Petr09Mitin/technopark-go-dz1/calculator/mathparser"
)

func main() {
	expression := os.Args[1]
	if expression == "" {
		fmt.Println("no expression provided")
		return
	}

	result, err := mathparser.CalculateExpression(expression)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}

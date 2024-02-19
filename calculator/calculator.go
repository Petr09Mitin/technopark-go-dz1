package main

import (
	"flag"
	"fmt"

	"github.com/Petr09Mitin/technopark-go-dz1/calculator/mathparser"
)

func main() {
	flag.Parse()
	expression := flag.Arg(0)
	if expression == "" {
		fmt.Print("no expression provided")
		return
	}

	result, err := mathparser.CalculateExpression(expression)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Print(result)
}

package main

import (
	"fmt"
	"mathInterpreter/parser"
)

func main() {
	test, err := parser.New("22.44 +	50*(30+47)")

	if err != nil {
		println(err)
	}

	fmt.Print(test)
}

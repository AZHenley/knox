package main

import (
	"fmt"
	"io/ioutil"
	"knox/lexer"
	"knox/parser"
)

func main() {
	code, err := ioutil.ReadFile("examples/simple.knox")
	if err != nil {
		panic(err)
	}
	l := lexer.New(string(code))
	p := parser.New(l)
	ast := p.Program()

	fmt.Println(ast.Children)
}

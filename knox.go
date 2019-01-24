package main

import (
	"io/ioutil"
	"knox/ast"
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
	a := p.Program()

	ast.Print(a)
}

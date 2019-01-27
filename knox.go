package main

import (
	"io/ioutil"
	"knox/ast"
	"knox/lexer"
	"knox/parser"
	"knox/typechecker"
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

	typechecker.Analyze(&a)
	//cfa.Analyze(&a)

}

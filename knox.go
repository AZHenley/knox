package main

import (
	"fmt"
	"io/ioutil"
	"knox/ast"
	"knox/emitter"
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
	fmt.Print(emitter.Generate(&a))

}

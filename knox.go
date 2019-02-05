package main

import (
	"fmt"
	"io/ioutil"
	"knox/ast"
	"knox/emitter"
	"knox/lexer"
	"knox/parser"
	"knox/typechecker"
	"time"
)

func main() {
	code, err := ioutil.ReadFile("examples/simple.knox")
	if err != nil {
		panic(err)
	}

	start := time.Now()
	l := lexer.New(string(code))
	p := parser.New(l)
	a := p.Program()
	elapsedParsing := time.Since(start)

	ast.Print(a)

	start = time.Now()
	typechecker.Analyze(&a)
	elapsedTypeChecking := time.Since(start)

	//cfa.Analyze(&a)

	start = time.Now()
	output := emitter.Generate(&a)
	elapsedEmitting := time.Since(start)

	fmt.Println(output)

	fmt.Println(elapsedParsing)
	fmt.Println(elapsedTypeChecking)
	fmt.Println(elapsedEmitting)
}

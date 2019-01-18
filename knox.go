package main

import (
	"fmt"
	"io/ioutil"
	"knox/lexer"
	"knox/parser"
)

func main() {
	code, err := ioutil.ReadFile("examples/basic.knox")
	if err != nil {
		panic(err)
	}
	l := lexer.New(string(code))
	p := parser.New(l)
	p.Program()
	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
	}
}

func printParserErrors(errors []string) {
	for _, msg := range errors {
		fmt.Println("\t" + msg + "\n")
	}
}

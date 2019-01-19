package main

import (
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
}

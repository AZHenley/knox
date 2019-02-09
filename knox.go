package main

import (
	"flag"
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
	// Flags
	timeFlag := flag.Bool("time", false, "Print the time taken by each compiler phase.")
	astFlag := flag.Bool("ast", false, "Print the AST.")
	goFlag := flag.Bool("go", false, "Print the Go code.")
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		panic("Specify file to be compiled.")
	}
	code, err := ioutil.ReadFile(args[0]) // TODO: Support multiple files.
	if err != nil {
		panic(err)
	}

	start := time.Now()
	l := lexer.New(string(code))
	p := parser.New(l)
	a := p.Program()
	elapsedParsing := time.Since(start)

	if *astFlag {
		ast.Print(a)
	}

	start = time.Now()
	typechecker.Analyze(&a)
	elapsedTypeChecking := time.Since(start)

	//cfa.Analyze(&a)

	start = time.Now()
	output := emitter.Generate(&a)
	elapsedEmitting := time.Since(start)

	if *goFlag {
		fmt.Println(output)
	}

	if *timeFlag {
		fmt.Println(elapsedParsing)
		fmt.Println(elapsedTypeChecking)
		fmt.Println(elapsedEmitting)
	}
}

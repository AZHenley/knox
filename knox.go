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
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"
)

func main() {
	// Flags
	timeFlag := flag.Bool("time", false, "Print the time taken by each compiler phase.")
	astFlag := flag.Bool("ast", false, "Print the AST.")
	goFlag := flag.Bool("go", false, "Print the Go code.")
	outFlag := flag.String("out", "", "Path for output files.")
	nameFlag := flag.String("name", "", "Name for output executable.")
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		panic("Specify file to be compiled.")
	}
	code, err := ioutil.ReadFile(args[0]) // TODO: Support multiple files.
	if err != nil {
		panic(err)
	}

	// Lex, parse, and generate the AST.
	start := time.Now()
	l := lexer.New(string(code))
	p := parser.New(l)
	a := p.Program()
	elapsedParsing := time.Since(start)

	if *astFlag {
		ast.Print(a)
	}

	// Type check.
	start = time.Now()
	typechecker.Analyze(&a)
	elapsedTypeChecking := time.Since(start)

	// Control flow analysis.
	//cfa.Analyze(&a)

	// Generate code.
	start = time.Now()
	output := emitter.Generate(&a)
	elapsedEmitting := time.Since(start)

	if *goFlag {
		fmt.Println(output)
	}

	// Output code.
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	local := filepath.Dir(ex) // Get current path.
	outputDir := path.Join(local, *outFlag)
	outputFile := path.Join(outputDir, "out.go") // TODO: Go files should use Knox file names.
	binName := *nameFlag
	if binName == "" {
		binName = "out"
	}
	outputBin := path.Join(outputDir, binName)
	werr := ioutil.WriteFile(outputFile, []byte(output), 0644)
	if werr != nil {
		panic(werr)
	}

	// Invoke Go compiler.
	cmd := exec.Command("go", "build", "-o", outputBin, outputFile)
	_, sterr := cmd.Output()
	if sterr != nil {
		panic(sterr)
	}

	if *timeFlag {
		fmt.Println(elapsedParsing)
		fmt.Println(elapsedTypeChecking)
		fmt.Println(elapsedEmitting)
	}
}

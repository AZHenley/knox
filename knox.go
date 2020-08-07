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
	codeFlag := flag.Bool("code", false, "Print the C code.")
	outFlag := flag.String("out", "", "Path for output files.")
	nameFlag := flag.String("name", "", "Name for output executable.")
	binaryFlag := flag.Bool("binary", true, "Generates executable.")
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		panic("Specify file to be compiled.")
	}
	code, err := ioutil.ReadFile(args[0]) // TODO: Support multiple files.
	//code, err := ioutil.ReadFile("examples/chain.knox")
	if err != nil {
		panic(err)
	}

	// Lex, parse, and generate the AST.
	start := time.Now()
	l := lexer.New(string(code) + "\n")
	p := parser.New(l)
	a := p.Program()
	elapsedParsing := time.Since(start)

	if *astFlag {
		ast.Print(a)
	}

	// Builtin functions.
	// TODO: Build builtin first in case user's program conflicts.
	//a = *builtin.Init(&a) // TODO: Uncomment this to get stdlib back.

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

	if *codeFlag {
		fmt.Println(output)
	}

	// Output code.
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	local := filepath.Dir(ex) // Get current path.
	outputDir := path.Join(local, *outFlag)
	codeFile := path.Join(outputDir, "out.c") // TODO: C files should use Knox file names.
	binName := *nameFlag
	if binName == "" {
		binName = "a.out"
	}
	outputBin := path.Join(outputDir, binName) // TODO: Make the flag specify a file, not just a path.
	werr := ioutil.WriteFile(codeFile, []byte(output), 0644)
	if werr != nil {
		panic(werr)
	}

	// Invoke compiler.
	if *binaryFlag {
		cmd := exec.Command("clang", codeFile, "-o", outputBin)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cerr := cmd.Run()
		if cerr != nil {
			panic(cerr)
		}
	}

	if *timeFlag {
		fmt.Printf("Parsing took: %v\n", elapsedParsing)
		fmt.Printf("Parsing took: %v\n", elapsedTypeChecking)
		fmt.Printf("Parsing took: %v\n", elapsedEmitting)
	}
}

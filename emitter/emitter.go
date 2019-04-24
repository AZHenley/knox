package emitter

import (
	"fmt"
	"knox/ast"
	"strings"
)

var level = 0
var prototypes []string
var currentMethods []string

func indent() string {
	return strings.Repeat("\t", level)
}

// Generate outputs code given an AST.
func Generate(node *ast.Node) string {
	return program(node)
}

func header() string {
	code := ""
	code += "#include <stdlib.h>\n#include <stdio.h>\n#include <string.h>\n#include <stdint.h>\n#include <stdbool.h>\n#include \"knoxutil.h\"\n\n" // TODO: #130 Only include what is needed.
	return code
}

// TODO: Declare all functions first, before definitions, to avoid ordering issues.
// TODO: Main needs to return int.

func program(node *ast.Node) string {
	head := header()

	var code string
	for _, child := range node.Children {
		if child.Type == ast.FUNCDECL {
			code += funcDecl(&child)
		} else if child.Type == ast.CLASS {
			code += classDecl(&child)
			for _, method := range currentMethods {
				code += method
			}
		}
	}

	// Generate function prototypes since C requires functions to be declared before use.
	for _, prototype := range prototypes {
		head += prototype + "\n"
	}
	head += "\n"

	return head + code
}

func classDecl(node *ast.Node) string {
	currentMethods = nil
	code := "struct " + node.Children[0].TokenStart.Literal + " " + classBlock(&node.Children[1])
	return code
}

func classBlock(node *ast.Node) string {
	var code string
	//level++
	code += "{\n"
	for _, child := range node.Children {
		if child.Type == ast.VARDECL {
			code += "\t" + varDecl(&child)
		} else if child.Type == ast.FUNCDECL {
			method := funcDecl(&child) // TODO: Add self reference as first param.
			currentMethods = append(currentMethods, method)
		} else {
			// Should not happen.
			fmt.Println(child.Type)
		}
	}
	//level--
	code += "};\n\n"
	return code
}

func funcDecl(node *ast.Node) string {
	code := ""

	// Return types.
	// TODO: Create a struct for multiple return values.
	if len(node.Children[2].Children) > 1 {
		code += "("
		for index, item := range node.Children[2].Children {
			code += item.Children[0].TokenStart.Literal
			if index < len(node.Children[2].Children)-1 {
				code += ", "
			}
		}
		code += ") "
	} else if len(node.Children[2].Children) == 1 {
		returnType := node.Children[2].Children[0].Children[0].TokenStart.Literal
		// Knox allows main to be void or int but C requires int.
		if node.Children[0].TokenStart.Literal == "main" {
			returnType = "int"
		}
		code += returnType + " "
	}

	// Function name.
	code += node.Children[0].TokenStart.Literal + "("

	// Parameters.
	for i := 0; i < len(node.Children[1].Children); i++ {
		paramName := node.Children[1].Children[i].Children[0].TokenStart.Literal
		paramType := node.Children[1].Children[i].Children[1].Children[0].TokenStart.Literal
		if paramType == "string" {
			paramType = "char *"
		}
		code += paramType + " " + paramName

		if i+1 < len(node.Children[1].Children) {
			code += ", "
		}
	}
	code += ")"

	// Save this as the prototype.
	prototypes = append(prototypes, code+";")

	// Body.
	code += " " + block(&node.Children[3])

	return code
}

func block(node *ast.Node) string {
	var code string
	level++
	code += "{\n"
	for _, s := range node.Children {
		code += indent() + statement(&s)
	}
	level--
	code += indent() + "}\n\n"
	return code
}

// Block without newline.
func blockIf(node *ast.Node) string {
	var code string
	level++
	code += "{\n"
	for _, s := range node.Children {
		code += indent() + statement(&s)
	}
	level--
	code += indent() + "}"
	return code
}

func statement(node *ast.Node) string {
	var code string
	switch node.Type {
	case ast.VARDECL:
		code = varDecl(node)
	case ast.VARASSIGN:
		code = varAssign(node)
	case ast.IFSTATEMENT:
		code = ifStatement(node)
	case ast.WHILESTATEMENT:
		code = whileStatement(node)
	case ast.JUMPSTATEMENT:
		code = jumpStatement(node)
	case ast.LEFTEXPR:
		code = expr(&node.Children[0]) + "\n"
	case ast.FUNCCALL:
		code = funcCall(node) + "\n"
	}
	return code
}

func funcCall(node *ast.Node) string {
	funcName := node.Children[0].TokenStart.Literal

	// Normal function calls.
	if node.Children[0].Type != ast.DOTOP {
		var argList string
		for index, child := range node.Children {
			if index == 0 {
				continue
			}
			argList += expr(&child.Children[0])
			if index < len(node.Children)-1 {
				argList += ", "
			}
		}
		return funcName + "(" + argList + ")"
	}
	// Either a package or a method.

	// TODO: Handle builtin functions in a better way.
	if node.Children[0].Children[0].TokenStart.Literal == "stl" {
		switch node.Children[0].Children[1].TokenStart.Literal {
		case "print":
			var funcName string
			var argList string
			funcName = "printf"
			argList = expr(&node.Children[1])
			return funcName + "(\"%s\", " + argList + ");"
		}
	}
	return ""
}

func ifStatement(node *ast.Node) string {
	code := "if("
	code += expr(&node.Children[0]) + ") " + blockIf(&node.Children[1]) // Condition and block

	for i := 2; i < len(node.Children); i += 2 {
		if i+1 != len(node.Children) { // Else if
			code += " else if(" + expr(&node.Children[i]) + ") " + blockIf(&node.Children[i+1])

		} else { // Else
			code += " else " + blockIf(&node.Children[i])
		}
	}

	return code + "\n"
}

func whileStatement(node *ast.Node) string {
	return "while(" + expr(&node.Children[0]) + ") " + blockIf(&node.Children[1]) + "\n"
}

func jumpStatement(node *ast.Node) string {
	code := node.TokenStart.Literal + " "
	for index, child := range node.Children {
		code += expr(&child.Children[0])
		if index < len(node.Children)-1 {
			code += ", "
		}
	}
	return code + ";\n"
}

func varAssign(node *ast.Node) string {
	code := ""
	// TODO: Fix multiple assignment.
	for i := 0; i < len(node.Children)-1; i++ {
		code += node.Children[i].Children[0].TokenStart.Literal
		if i+1 < len(node.Children)-1 {
			code += ", "
		}
	}
	return code + " = " + expr(&node.Children[len(node.Children)-1].Children[0]) + ";\n"
}

func varDecl(node *ast.Node) string {
	code := ""
	// TODO: Fix multiple declarations.
	for i := 0; i < len(node.Children)-1; i += 2 {
		varName := node.Children[i].TokenStart.Literal
		varType := node.Children[i+1].Children[0].TokenStart.Literal
		if varType == "string" {
			varType = "const char *"
		}
		code += varType + " " + varName
		if i+2 < len(node.Children)-1 {
			code += ", "
		}
	}

	varExpr := expr(&node.Children[len(node.Children)-1].Children[0])
	//return varName + " " + varType + " := " + varExpr + "\n"
	return code + " = " + varExpr + ";\n"
}

func expr(node *ast.Node) string {
	if node.Type == ast.BINARYOP {
		// If concatenating strings.
		if node.TokenStart.Literal == "concat" { // Type checker will convert + for strings to concat.
			return "concat(" + expr(&node.Children[0]) + ", " + expr(&node.Children[1]) + ")"
		}
		// Else any other binary op.
		return "(" + expr(&node.Children[0]) + node.TokenStart.Literal + expr(&node.Children[1]) + ")"
	} else if node.Type == ast.UNARYOP {
		return "(" + node.TokenStart.Literal + expr(&node.Children[0]) + ")"
	} else if node.Type == ast.FUNCCALL {
		return funcCall(node)
	} else if node.Type == ast.EXPRESSION {
		return expr(&node.Children[0])
	} else { // Primary.
		if node.Type == ast.STRING {
			return "\"" + node.TokenStart.Literal + "\""
		}
		if node.Type == ast.NIL {
			return "null"
		}
		return node.TokenStart.Literal
	}
}

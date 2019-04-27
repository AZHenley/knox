package emitter

import (
	"fmt"
	"knox/ast"
	"strings"
)

var level = 0
var currentName string      // Current var declaration name.
var nextLine string         // Placeholder to add next line once current line completes.
var prototypes []string     // Keep track of function prototypes so that order doesn't matter.
var currentMethods []string // Keep track of methods in current class
var currentMembers []string // Keep track of members in current class
var currentClass string     // Name of current class

func indent() string {
	return strings.Repeat("\t", level)
}

// Generate outputs code given an AST.
func Generate(node *ast.Node) string {
	return program(node)
}

func header() string {
	code := ""
	code += "#include <stdlib.h>\n#include <stdio.h>\n#include <string.h>\n#include <stdint.h>\n#include <stdbool.h>\n#include <stddef.h>\n#include \"knoxutil.h\"\n\n" // TODO: #130 Only include what is needed.
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

			// Fake constructor.
			constructor := "void _" + currentClass + "(struct " + currentClass + "* self)"
			prototypes = append(prototypes, constructor+";")
			code += constructor + " {\n"
			for _, member := range currentMembers {
				code += "\tself->" + member
			}
			code += "}\n\n"

			// All other methods.
			for _, method := range currentMethods {
				code += method
			}

			currentClass = ""
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
	currentMembers = nil
	currentClass = node.Children[0].TokenStart.Literal
	code := "struct " + node.Children[0].TokenStart.Literal + " " + classBlock(&node.Children[1])
	return code
}

func classBlock(node *ast.Node) string {
	var code string
	//level++
	code += "{\n"
	for _, child := range node.Children {
		if child.Type == ast.VARDECL {
			//code += "\t" + varDecl(&child)
			rawMethod := "\t" + varDecl(&child)
			code += strings.Split(rawMethod, " =")[0] + ";\n"
		} else if child.Type == ast.FUNCDECL {
			method := funcDecl(&child)
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
	if currentClass != "" {
		code += "struct " + currentClass + "* self, "
	}
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
	currentName = ""

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
		code = expr(&node.Children[0]) + ";\n"
	case ast.FUNCCALL:
		code = funcCall(node) + "\n"
	}

	if nextLine != "" {
		code += nextLine + "\n"
		nextLine = ""
	}
	return code
}

// TODO: This needs to be rewritten to recursively handle dotop.
func funcCall(node *ast.Node) string {

	// Normal function calls.
	if node.Children[0].Type != ast.DOTOP {
		funcName := node.Children[0].TokenStart.Literal
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

	// If a method.
	// myobj.foo(a, b)  ->  foo(myobj, a, b)
	if node.Children[0].Type == ast.DOTOP {
		funcName := node.Children[0].Children[1].TokenStart.Literal
		var argList string
		argList += node.Children[0].Children[0].TokenStart.Literal + ", "
		for index, child := range node.Children[1].Children {
			argList += expr(&child)
			if index < len(node.Children[1].Children)-1 {
				argList += ", "
			}
		}
		return funcName + "(" + argList + ")"
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
		currentName = node.Children[i].Children[0].TokenStart.Literal
		if i+1 < len(node.Children)-1 {
			code += ", "
		}
	}
	if node.Children[0].Children[0].Type == ast.DOTOP {
		code = expr(&node.Children[0].Children[0])
	}
	return code + " = " + expr(&node.Children[len(node.Children)-1].Children[0]) + ";\n"
}

func varDecl(node *ast.Node) string {
	code := ""
	member := ""
	// TODO: Fix multiple declarations.
	for i := 0; i < len(node.Children)-1; i += 2 {
		varName := node.Children[i].TokenStart.Literal
		currentName = varName

		varType := node.Children[i+1].Children[0].TokenStart.Literal
		if varType == "string" {
			varType = "const char *"
		} else if varType != "int" && varType != "bool" && varType != "float" {
			// TODO: This needs to be changed.
			// If a reference type
			varType = "struct " + varType + " *"
		}
		code += varType + " " + varName
		member += varName
		if i+2 < len(node.Children)-1 {
			code += ", "
		}
	}

	varExpr := expr(&node.Children[len(node.Children)-1].Children[0])
	member += " = " + varExpr + ";\n"
	if currentClass != "" {
		currentMembers = append(currentMembers, member)
	}
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
	} else if node.Type == ast.DOTOP {
		return "(" + expr(&node.Children[0]) + "->" + expr(&node.Children[1]) + ")"
	} else if node.Type == ast.EXPRESSION {
		return expr(&node.Children[0])
	} else if node.Type == ast.NEW {
		nextLine = "\t_" + node.Children[0].Children[0].TokenStart.Literal + "(" + currentName + ");"
		return "malloc(sizeof(struct " + node.Children[0].Children[0].TokenStart.Literal + "))"
	} else { // Primary.
		if node.Type == ast.STRING {
			return "\"" + node.TokenStart.Literal + "\""
		}
		if node.Type == ast.NIL {
			return "NULL"
		}
		return node.TokenStart.Literal
	}
}

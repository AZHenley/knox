package emitter

import (
	"knox/ast"
	"strings"
)

// Emitter object.
type Emitter struct {
	output string
	level  int
}

func (e *Emitter) emit(code string) {
	e.output += strings.Repeat("\t", e.level) + code
}

// Generate outputs code given an AST.
func Generate(node *ast.Node) string {
	return program(node)
}

func header() string {
	code := ""
	code += "package main\n\n"
	code += "import (\n\t\"fmt\"\n)\n\n"
	return code
}

func program(node *ast.Node) string {
	var code string
	code += header()

	for _, child := range node.Children {
		if child.Type == ast.FUNCDECL {
			code += funcDecl(&child)
		} else if child.Type == ast.CLASS {

		}
	}

	return code
}

func funcDecl(node *ast.Node) string {
	code := "func " + node.Children[0].TokenStart.Literal + "("

	// Parameters.
	for i := 0; i < len(node.Children[1].Children); i++ {
		paramName := node.Children[1].Children[i].Children[0].TokenStart.Literal
		paramType := node.Children[1].Children[i].Children[1].Children[0].TokenStart.Literal
		code += paramName + " " + paramType

		if i+1 < len(node.Children[1].Children) {
			code += ", "
		}
	}
	code += ") "

	// Return types.
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
		if returnType != "void" {
			code += node.Children[2].Children[0].Children[0].TokenStart.Literal + " "
		}
	}

	// Body.
	code += block(&node.Children[3])

	return code
}

func block(node *ast.Node) string {
	var code string
	code += "{\n"
	for _, s := range node.Children {
		code += statement(&s)
	}
	code += "}\n\n"
	return code
}

// Special block without newline.
func blockIf(node *ast.Node) string {
	var code string
	code += "{\n"
	for _, s := range node.Children {
		code += statement(&s)
	}
	code += "}"
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
	case ast.FUNCCALL:
		code = funcCall(node) + "\n"
	}
	return code
}

func funcCall(node *ast.Node) string {
	funcName := node.Children[0].TokenStart.Literal

	// TODO: Remove this once Go modules can be imported.
	if funcName == "print" {
		funcName = "fmt.Println"
	}

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

func ifStatement(node *ast.Node) string {
	// TODO: Block will not put else if and else on correct line.

	code := "if "
	code += expr(&node.Children[0]) + " " + blockIf(&node.Children[1]) // Condition and block

	for i := 2; i < len(node.Children); i += 2 {
		if i+1 != len(node.Children) { // Else if
			code += " else if " + expr(&node.Children[i]) + " " + blockIf(&node.Children[i+1])

		} else { // Else
			code += " else " + blockIf(&node.Children[i])
		}
	}

	return code + "\n"
}

func whileStatement(node *ast.Node) string {
	return "for " + expr(&node.Children[0]) + " " + blockIf(&node.Children[1]) + "\n"
}

func jumpStatement(node *ast.Node) string {
	code := node.TokenStart.Literal + " "
	for index, child := range node.Children {
		code += expr(&child.Children[0])
		if index < len(node.Children)-1 {
			code += ", "
		}
	}
	return code + "\n"
}

func varAssign(node *ast.Node) string {
	code := ""
	for i := 0; i < len(node.Children)-1; i++ {
		code += node.Children[i].Children[0].TokenStart.Literal
		if i+1 < len(node.Children)-1 {
			code += ", "
		}
	}
	return code + " = " + expr(&node.Children[len(node.Children)-1].Children[0]) + "\n"
}

func varDecl(node *ast.Node) string {
	code := ""
	for i := 0; i < len(node.Children)-1; i += 2 {
		varName := node.Children[i].TokenStart.Literal
		//varType := node.Children[1].TokenStart.Literal
		code += varName
		if i+2 < len(node.Children)-1 {
			code += ", "
		}
	}

	varExpr := expr(&node.Children[len(node.Children)-1].Children[0])
	//return varName + " " + varType + " := " + varExpr + "\n"
	return code + " := " + varExpr + "\n"
}

func expr(node *ast.Node) string {
	if node.Type == ast.BINARYOP {
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
		return node.TokenStart.Literal
	}
}

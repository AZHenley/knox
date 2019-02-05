package emitter

import (
	"knox/ast"
)

// Generate outputs code given an AST.
func Generate(node *ast.Node) string {
	return program(node)
}

func header() string {
	return "package main\n\nimport (\n\t\"fmt\"\n)\n\n"
}

func program(node *ast.Node) string {
	var code string
	code += header()

	for _, funcNode := range node.Children {
		code += funcDecl(&funcNode)
	}

	return code
}

func funcDecl(node *ast.Node) string {
	code := "func " + node.Children[0].TokenStart.Literal + "("

	// Parameters.
	for i := 0; i < len(node.Children[1].Children); i += 2 {
		paramName := node.Children[1].Children[i].TokenStart.Literal
		paramType := node.Children[1].Children[i+1].Children[0].TokenStart.Literal
		code += paramName + " " + paramType

		if i+2 < len(node.Children[1].Children) {
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
	code += "{\n"
	for _, s := range node.Children[3].Children {
		code += statement(&s)
	}
	code += "}\n\n"

	return code
}

func statement(node *ast.Node) string {
	var code string
	switch node.Type {
	case ast.VARDECL:
		code = varDecl(node)
	}
	return code
}

func varDecl(node *ast.Node) string {
	varName := node.Children[0].TokenStart.Literal
	varType := node.Children[1].TokenStart.Literal
	varExpr := expr(&node.Children[2].Children[0])
	return "var " + varName + " " + varType + " = " + varExpr + "\n"
}

func expr(node *ast.Node) string {
	if node.Type == ast.BINARYOP {
		return expr(&node.Children[0]) + node.TokenStart.Literal + expr(&node.Children[1])
	} else if node.Type == ast.UNARYOP {
		return ""
	} else if node.Type == ast.FUNCCALL {
		return ""
	} else if node.Type == ast.EXPRESSION {
		return expr(&node.Children[0])
	} else { // Primary.
		return node.TokenStart.Literal
	}
}

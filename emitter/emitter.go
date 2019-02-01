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
	// TODO.
	code += "}\n\n"

	return code
}

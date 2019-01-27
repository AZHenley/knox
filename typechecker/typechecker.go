package typechecker

import (
	"fmt"
	"knox/ast"
)

// Analyze performs type checking on the entire AST.
func Analyze(node *ast.Node) {
	typecheck(node)
}

func typecheck(node *ast.Node) {
	// TODO: Check varassign and vardecl
	for _, child := range node.Children {
		if child.Type == ast.EXPRESSION {
			getType(&child.Children[0])
		} else {
			typecheck(&child)
		}
	}
}

func abortMsg(msg string) {
	fmt.Println("Type error: " + msg)
	panic("Aborted.\n")
}

// TODO: This should probably be split up into methods and a class for each AST node type.
// Get type from expression node.
func getType(node *ast.Node) string {
	fmt.Println(node.Type)
	switch node.Type {
	case ast.BINARYOP:
		left := getType(&node.Children[0])
		right := getType(&node.Children[1])
		if left != right {
			abortMsg("Mismatched types.")
		}
		if left != ast.INT && left != ast.FLOAT && left != ast.STRING {
			abortMsg("Invalid operation.")
		}
		return string(left)
	case ast.INT, ast.FLOAT, ast.STRING, ast.BOOL:
		return string(node.Type)
	}
	return ""
}

// Check if symbol is declared.
func isDeclared(symbol string) bool {
	return false
}

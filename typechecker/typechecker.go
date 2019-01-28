package typechecker

import (
	"fmt"
	"knox/ast"
	"knox/lexer"
	"knox/token"
	"strings"
)

// Analyze performs type checking on the entire AST.
func Analyze(node *ast.Node) {
	typecheck(node)
}

func typecheck(node *ast.Node) {
	// TODO: Check varassign and vardecl
	for _, child := range node.Children {
		if child.Type == ast.EXPRESSION {
			exprType := strings.ToLower(getType(&child.Children[0]))
			// TODO: Handle vardecl, varassign, if, while, for, return
			if node.Type == ast.VARDECL {
				leftType := declType(node)
				if leftType != exprType {
					abortMsg("Mismatched types.")
				}
			}
		} else {
			typecheck(&child)
		}
	}
}

func abortMsg(msg string) {
	fmt.Println("Type error: " + msg)
	panic("Aborted.\n")
}

// Get type from a declaration.
func declType(node *ast.Node) string {
	// Assuming this node is a vardecl
	return node.Children[1].Children[0].TokenStart.Literal
}

// Get type from a symbol.

// Get type from expression node.
func getType(node *ast.Node) string {
	switch node.Type {
	case ast.BINARYOP:
		left := getType(&node.Children[0])
		right := getType(&node.Children[1])
		if left != right {
			abortMsg("Mismatched types.")
		}
		if lexer.IsOperator([]rune(node.TokenStart.Literal)[0]) || node.TokenStart.Literal == ">=" || node.TokenStart.Literal == ">" || node.TokenStart.Literal == "<=" || node.TokenStart.Literal == "<" {
			if left == ast.INT || left == ast.FLOAT {
				return string(left)
			} else if node.TokenStart.Type == token.PLUS && left == ast.STRING {
				return string(left)
			} else {
				abortMsg("Invalid operation.")
			}
		} else if node.TokenStart.Literal == "&&" || node.TokenStart.Literal == "||" {
			if left != ast.BOOL || right != ast.BOOL {
				abortMsg("Invalid operation.")
			}
		}

		return string(left)

	case ast.UNARYOP:
		single := getType(&node.Children[0])
		if node.TokenStart.Type == token.BANG {
			if single != ast.BOOL {
				abortMsg("Invalid operation.")
			}
		}
		if single != ast.INT && single != ast.FLOAT {
			abortMsg("Invalid operation.")
		}
		return string(single)

	case ast.VARREF:
		// TODO
		//name := node.Children[0].TokenStart.Literal
		//isDeclared(name)
		//look up type
		return "INT"

	case ast.FUNCCALL:
		// TODO
		return "INT"

	case ast.INT, ast.FLOAT, ast.STRING, ast.BOOL, ast.NIL:
		return string(node.Type)
	}
	return ""
}

// Check if symbol is declared.
func isDeclared(symbol string) bool {
	return false
}

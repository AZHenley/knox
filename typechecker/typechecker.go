package typechecker

import (
	"fmt"
	"knox/ast"
	"knox/lexer"
	"knox/token"
	"strings"
)

// Internal representation of a type.
type typeObj struct {
	fullName    string    // Name of this type (and all inner types)
	isPrimitive bool      // Is this type a primitive (int, float, string, rune, byte, bool)
	isContainer bool      // Is this type a container (list, map, address)
	isClass     bool      // Is this a user-defined class
	isEnum      bool      // Is this an enum
	isTypedef   bool      // Is this a typedef
	inner       []typeObj // Inner types. TODO: Make this a slice of pointers of typeObj.
}

var typeBOOL *typeObj
var typeINT *typeObj
var typeFLOAT *typeObj
var typeSTRING *typeObj
var typeLIST *typeObj
var typeMAP *typeObj
var typeADDRESS *typeObj

// Analyze performs type checking on the entire AST.
func Analyze(node *ast.Node) {
	setup()
	typecheck(node)
}

func setup() {
	typeBOOL = &typeObj{}
	typeINT = &typeObj{}
	typeFLOAT = &typeObj{}
	typeSTRING = &typeObj{}
	typeBOOL.isPrimitive = true
	typeINT.isPrimitive = true
	typeFLOAT.isPrimitive = true
	typeSTRING.isPrimitive = true
	typeBOOL.fullName = "bool"
	typeINT.fullName = "int"
	typeFLOAT.fullName = "float"
	typeSTRING.fullName = "string"
}

func typecheck(node *ast.Node) {
	for _, child := range node.Children {

		if child.Type == ast.EXPRESSION {
			exprType := getType(&child.Children[0])
			// TODO: Handle for, return
			if node.Type == ast.VARDECL {
				// TODO: Handle multiple assignment.
				leftType := declType2(node)
				if !compareTypes2(leftType, exprType) { // Do the types match?
					abortMsg("Mismatched types.")
				}
			} else if node.Type == ast.VARASSIGN {
				// TODO: Handle multiple assignment.
				decl := child.Symbols.LookupSymbol(node.Children[0].Children[0].TokenStart.Literal)
				if decl == nil {
					abortMsg("Referencing undeclared variable.")
				}
				leftType := declType2(decl)
				if !compareTypes2(leftType, exprType) { // Do the types match?
					abortMsg("Mismatched types.")
				}
			} else if node.Type == ast.IFSTATEMENT || node.Type == ast.WHILESTATEMENT {
				if !compareTypes2(exprType, typeBOOL) {
					abortMsg("Conditionals require boolean expressions.")
				}
			}
		} else if child.Type == ast.FUNCCALL { // Handles funccall outside of an expression.
			// TODO: Are the return values used?
			name := child.Children[0].TokenStart.Literal
			declNode := node.Symbols.LookupSymbol(name)
			if declNode == nil {
				abortMsg("Calling undeclared function.")
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

func compareTypes(a string, b string) bool {
	return strings.ToLower(a) == strings.ToLower(b)
}

func compareTypes2(a *typeObj, b *typeObj) bool {
	return a.fullName == b.fullName
}

func stringToType(prim string) *typeObj {
	primitive := &typeObj{}
	primitive.fullName = prim
	return primitive
}

// Get type from a declaration.
func declType(node *ast.Node) string {
	if node.Type == ast.VARDECL {
		return node.Children[1].Children[0].TokenStart.Literal
	} else if node.Type == ast.FUNCDECL {
		return node.Children[2].Children[0].Children[0].TokenStart.Literal // TODO: Handle multiple return values.
	}
	return "" // Should never happen?
}

// TODO: This needs to build up the typeobj.
// Get type from a declaration.
func declType2(node *ast.Node) *typeObj {
	if node.Type == ast.VARDECL {
		return stringToType(node.Children[1].Children[0].TokenStart.Literal)
	} else if node.Type == ast.FUNCDECL {
		return stringToType(node.Children[2].Children[0].Children[0].TokenStart.Literal) // TODO: Handle multiple return values.
	}
	abortMsg("Unknown type error.")
	return nil
}

// // Get type from a symbol.
// func getSymbolType(symbol string, st *ast.SymTable) string {
// 	node := st.LookupSymbol(symbol)
// 	if node != nil {
// 		return declType(node)
// 	}
// 	return ""
// }

// Get type from expression node.
func getType(node *ast.Node) *typeObj {
	switch node.Type {
	case ast.BINARYOP:
		left := getType(&node.Children[0])
		right := getType(&node.Children[1])
		if !compareTypes2(left, right) {
			abortMsg("Mismatched types.")
		}
		if lexer.IsOperator([]rune(node.TokenStart.Literal)[0]) || node.TokenStart.Literal == ">=" || node.TokenStart.Literal == ">" || node.TokenStart.Literal == "<=" || node.TokenStart.Literal == "<" {
			if compareTypes2(left, typeINT) || compareTypes2(left, typeFLOAT) {
				return left
			} else if node.TokenStart.Type == token.PLUS && compareTypes2(left, typeSTRING) {
				return left
			} else {
				abortMsg("Invalid operation.")
			}
		} else if node.TokenStart.Literal == "&&" || node.TokenStart.Literal == "||" {
			if !compareTypes2(left, typeBOOL) {
				abortMsg("Invalid operation.")
			}
		}

		return left

	case ast.UNARYOP:
		single := getType(&node.Children[0])
		if node.TokenStart.Type == token.BANG {
			if !compareTypes2(single, typeBOOL) {
				abortMsg("Invalid operation.")
			}
		}
		if !compareTypes2(single, typeINT) && !compareTypes2(single, typeFLOAT) {
			abortMsg("Invalid operation.")
		}
		return single

	case ast.VARREF:
		name := node.Children[0].TokenStart.Literal
		declNode := node.Symbols.LookupSymbol(name)
		if declNode == nil {
			abortMsg("Referencing undeclared variable.")
		}
		return declType2(declNode)

	case ast.FUNCCALL:
		name := node.Children[0].TokenStart.Literal
		declNode := node.Symbols.LookupSymbol(name)
		if declNode == nil {
			abortMsg("Calling undeclared function.")
		}
		return declType2(declNode)

	case ast.INT:
		return typeINT
	case ast.FLOAT:
		return typeFLOAT
	case ast.STRING:
		return typeSTRING
	case ast.BOOL:
		return typeBOOL
	//case ast.NIL:
	//return typeNIL

	case ast.EXPRESSION:
		return getType(&node.Children[0])
	}

	return nil
}

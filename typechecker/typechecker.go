package typechecker

import (
	"fmt"
	"knox/ast"
	"knox/lexer"
	"knox/token"
)

// Internal representation of a type.
type typeObj struct {
	fullName    string    // Name of this type (and all inner types)
	isFunction  bool      // Is this a function
	isPrimitive bool      // Is this type a primitive (int, float, string, rune, byte, bool)
	isContainer bool      // Is this type a container (list, map, address, etc.)
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

var currentFunc *ast.Node

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
				leftType := declType(node)
				if !compareTypes(leftType, exprType) { // Do the types match?
					abortMsg("Mismatched types.")
				}
			} else if node.Type == ast.VARASSIGN {
				// TODO: Handle multiple assignment.
				decl := child.Symbols.LookupSymbol(node.Children[0].Children[0].TokenStart.Literal)
				if decl == nil {
					abortMsg("Referencing undeclared variable.")
				}
				leftType := declType(decl)
				if !compareTypes(leftType, exprType) { // Do the types match?
					abortMsg("Mismatched types.")
				}
			} else if node.Type == ast.IFSTATEMENT || node.Type == ast.WHILESTATEMENT {
				if !compareTypes(exprType, typeBOOL) {
					abortMsg("Conditionals require boolean expressions.")
				}
			}
		} else if child.Type == ast.FUNCCALL { // Handles funccall outside of an expression.
			// TODO: Are the arguments the correct type?
			// TODO: Are the return values used?
			name := child.Children[0].TokenStart.Literal
			declNode := node.Symbols.LookupSymbol(name)
			if declNode == nil {
				abortMsg("Calling undeclared function.")
			}
		} else if child.Type == ast.JUMPSTATEMENT {
			if child.TokenStart.Literal == "return" {
				// TODO: Support multiple return types.
				returnType := buildTypeList(&child)
				funcReturnType := buildReturnList(&currentFunc.Children[2])
				if !compareTypes(returnType, funcReturnType) {
					abortMsg("Incorrect return type.")
				}
			}
		} else if child.Type == ast.FUNCDECL {
			currentFunc = &child
			typecheck(&child)
		} else {
			typecheck(&child)
		}
	}
}

func abortMsg(msg string) {
	fmt.Println("Type error: " + msg)
	panic("Aborted.\n")
}

func compareTypes(a *typeObj, b *typeObj) bool {
	return a.fullName == b.fullName
}

func stringToType(prim string) *typeObj {
	primitive := &typeObj{}
	primitive.fullName = prim
	return primitive
}

// Get type from a declaration.
func declType(node *ast.Node) *typeObj {
	if node.Type == ast.VARDECL {
		//return stringToType(node.Children[1].Children[0].TokenStart.Literal)
		return buildTypeObj(&node.Children[1])
	} else if node.Type == ast.FUNCDECL { // TODO: Redo this.
		return stringToType(node.Children[2].Children[0].Children[0].TokenStart.Literal) // TODO: Handle multiple return values.
	}
	abortMsg("Unknown type error.")
	return nil
}

// Builds up a type obj recursively given a varType AST node.
func buildTypeObj(node *ast.Node) *typeObj {
	obj := &typeObj{}

	if isSimple(node) {
		obj.isPrimitive = isPrimitiveType(node)
		obj.isClass = !obj.isPrimitive
		obj.fullName = getName(node)
		return obj
	} else if isList(node) {
		obj.isContainer = true
		obj.inner = append(obj.inner, *buildTypeObj(&node.Children[1]))
		obj.fullName = "[" + obj.inner[0].fullName + "]"
		return obj
	} else { // Complex type
		obj.isContainer = true
		obj.fullName = getName(node) + "["
		for i := 1; i < len(node.Children); i++ {
			obj.inner = append(obj.inner, *buildTypeObj(&node.Children[i]))
			obj.fullName += obj.inner[i-1].fullName
			if i+1 < len(node.Children) {
				obj.fullName += ","
			}
		}
		obj.fullName += "]"
		return obj
	}
}

// Build a list of types from expressions.
func buildTypeList(node *ast.Node) *typeObj {
	obj := &typeObj{}
	for index, expr := range node.Children {
		obj.inner = append(obj.inner, *getType(&expr))
		obj.fullName += obj.inner[index].fullName
		if index+1 < len(node.Children) {
			obj.fullName += ","
		}
	}
	return obj
}

// Build a list of types from func decl return.
func buildReturnList(node *ast.Node) *typeObj {
	obj := &typeObj{}
	for index, ret := range node.Children {
		obj.inner = append(obj.inner, *buildTypeObj(&ret))
		obj.fullName += obj.inner[index].fullName
		if index+1 < len(node.Children) {
			obj.fullName += ","
		}
	}
	return obj
}

func isList(node *ast.Node) bool {
	if len(node.Children) == 2 && getName(node) == "[" {
		return true
	}
	return false
}

func isPrimitiveType(node *ast.Node) bool {
	literal := getName(node)
	return literal == "bool" || literal == "string" || literal == "int" || literal == "float"
}

func isSimple(node *ast.Node) bool {
	if len(node.Children) == 1 {
		return true
	}
	return false
}

func getName(node *ast.Node) string {
	return node.Children[0].TokenStart.Literal
}

// Get type from expression node.
func getType(node *ast.Node) *typeObj {
	switch node.Type {
	case ast.BINARYOP:
		left := getType(&node.Children[0])
		right := getType(&node.Children[1])
		if !compareTypes(left, right) {
			abortMsg("Mismatched types.")
		}
		if lexer.IsOperator([]rune(node.TokenStart.Literal)[0]) || node.TokenStart.Literal == ">=" || node.TokenStart.Literal == ">" || node.TokenStart.Literal == "<=" || node.TokenStart.Literal == "<" {
			if compareTypes(left, typeINT) || compareTypes(left, typeFLOAT) {
				return left
			} else if node.TokenStart.Type == token.PLUS && compareTypes(left, typeSTRING) {
				return left
			} else {
				abortMsg("Invalid operation.")
			}
		} else if node.TokenStart.Literal == "&&" || node.TokenStart.Literal == "||" {
			if !compareTypes(left, typeBOOL) {
				abortMsg("Invalid operation.")
			}
		}

		return left

	case ast.UNARYOP:
		single := getType(&node.Children[0])
		if node.TokenStart.Type == token.BANG {
			if !compareTypes(single, typeBOOL) {
				abortMsg("Invalid operation.")
			}
		} else if node.TokenStart.Type == token.PLUS || node.TokenStart.Type == token.MINUS {
			if !compareTypes(single, typeINT) && !compareTypes(single, typeFLOAT) {
				abortMsg("Invalid operation.")
			}
		} else if node.TokenStart.Type == token.NEW {
			//if !compareTypes(single, typeBOOL) {
			//	abortMsg("Invalid operation.")
			//}
		}
		return single

	case ast.VARREF:
		name := node.Children[0].TokenStart.Literal
		declNode := node.Symbols.LookupSymbol(name)
		if declNode == nil {
			abortMsg("Referencing undeclared variable.")
		}
		return declType(declNode)

	case ast.FUNCCALL:
		// TODO: Are the arguments the correct type?
		name := node.Children[0].TokenStart.Literal
		declNode := node.Symbols.LookupSymbol(name)
		if declNode == nil {
			abortMsg("Calling undeclared function.")
		}
		return declType(declNode)

	case ast.NEW:
		return buildTypeObj(&node.Children[0])

	case ast.LIST:
		obj := &typeObj{}
		obj.isContainer = true
		obj.fullName = "["
		itemType := ""
		for i, item := range node.Children {
			obj.inner = append(obj.inner, *getType(&item))
			if obj.inner[i].fullName != itemType && itemType != "" { // Check if all items are same type.
				abortMsg("Mismatched types in list literal.")
			}
			itemType = obj.inner[i].fullName
		}
		// TODO: Handle if no items.
		obj.fullName += itemType + "]"
		return obj

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

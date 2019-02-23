package typechecker

import (
	"fmt"
	"knox/ast"
	"knox/lexer"
	"knox/token"
)

// Internal representation of a type.
type typeObj struct {
	fullName    string // Name of this type (and all inner types)
	name        string // Name of this outer type
	isFunction  bool   // Is this a function
	isPrimitive bool   // Is this type a primitive (int, float, string, rune, byte, bool)
	isContainer bool   // Is this type a container (list, map, address, etc.)
	isList      bool
	isMap       bool
	isMulti     bool      // Is this a set of types (used for multiple return)
	isClass     bool      // Is this a user-defined class
	isEnum      bool      // Is this an enum
	isTypedef   bool      // Is this a typedef
	inner       []typeObj // Inner types. TODO: Make this a slice of pointers of typeObj.
}

var typeVOID *typeObj
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
	typeVOID = &typeObj{}
	typeBOOL = &typeObj{}
	typeINT = &typeObj{}
	typeFLOAT = &typeObj{}
	typeSTRING = &typeObj{}
	typeBOOL.isPrimitive = true
	typeINT.isPrimitive = true
	typeFLOAT.isPrimitive = true
	typeSTRING.isPrimitive = true
	typeVOID.fullName = "void"
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

				//fmt.Println("Left: " + leftType.fullName)
				//fmt.Println("Right: " + exprType.fullName)

				if !compareTypes(leftType, exprType) { // Do the types match?
					abortMsgf("Mismatched types: %s and %s", leftType.fullName, exprType.fullName)
				}
			} else if node.Type == ast.VARASSIGN {
				// TODO: Handle multiple assignment.
				decl := child.Symbols.LookupSymbol(node.Children[0].Children[0].TokenStart.Literal)
				if decl == nil {
					abortMsgf("Referencing undeclared variable: %s", node.Children[0].Children[0].TokenStart.Literal)
				}
				leftType := declType(decl)
				if !compareTypes(leftType, exprType) { // Do the types match?
					abortMsgf("Mismatched types: %s and %s", leftType.fullName, exprType.fullName)
				}
			} else if node.Type == ast.IFSTATEMENT || node.Type == ast.WHILESTATEMENT {
				if !compareTypes(exprType, typeBOOL) {
					abortMsg("Conditionals require boolean expressions.")
				}
			}
			// } else if child.Type == ast.FUNCCALL { // Handles funccall outside of an expression.
			// 	name := child.Children[0].TokenStart.Literal
			// 	declNode := node.Symbols.LookupSymbol(name)

			// 	// Compare types between args and params.
			// 	checkFuncCall(&child, declNode)

			// 	// Check that nothing is returned.
			// 	if name != "print" && (len(declType(declNode).inner) > 1 || !compareTypes(&declType(declNode).inner[0], typeVOID)) {
			// 		abortMsg("Function call return values must be used.")
			// 	}
		} else if child.Type == ast.JUMPSTATEMENT {
			if child.TokenStart.Literal == "return" {
				// TODO: Support multiple return types.
				returnType := buildTypeList(&child)
				funcReturnType := buildReturnList(&currentFunc.Children[2])
				if compareTypes(funcReturnType, typeVOID) && returnType.fullName == "" { // Check for return; and void type.
				} else if !compareTypes(returnType, funcReturnType) {
					abortMsg("Incorrect return type.")
				}
			}
		} else if child.Type == ast.FUNCDECL {
			currentFunc = &child
			typecheck(&child)
		} else if child.Type == ast.LEFTEXPR {
			only := getType(&child.Children[0])
			if !compareTypes(only, typeVOID) {
				abortMsg("Expression must be of void type.")
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

func abortMsgf(msg string, args ...interface{}) {
	fmt.Printf("Type error: "+msg+"\n", args...)
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
	} else if node.Type == ast.FUNCDECL {
		// Currently this always returns a functions return type
		// TODO: Does this handle multiple return?
		return buildReturnList(&node.Children[2])
	} else if node.Type == ast.CLASS {
		// return
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
		obj.isList = true
		obj.inner = append(obj.inner, *buildTypeObj(&node.Children[1]))
		obj.name = "["
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

func checkFuncCall(node *ast.Node, declNode *ast.Node) {
	name := node.Children[0].TokenStart.Literal
	if name == "" {
		name = node.Children[0].Children[1].TokenStart.Literal
	}
	if name == "print" {
		return
	}

	//declNode := node.Symbols.LookupSymbol(name)
	// TODO: Remove the "print" check after std library is added.
	if declNode == nil {
		abortMsgf("Calling undeclared function: %s", name)
	}

	// Check number of args to number of params.
	if len(node.Children)-1 != len(declNode.Children[1].Children) {
		abortMsg("Incorrect number of arguments.")
	}
	// Check types of args to types of params.
	for i := 1; i < len(node.Children); i++ {
		argType := getType(&node.Children[i])
		expectedType := declType(&declNode.Children[1].Children[i-1])
		if !compareTypes(argType, expectedType) {
			abortMsgf("Mismatched type in function argument.")
		}
	}
}

// Look up a func's declaration given a ref, handling the dot operator.
func lookUpDecl(node ast.Node) *ast.Node {
	if node.Type == ast.IDENT {
		name := node.TokenStart.Literal
		declNode := node.Symbols.LookupSymbol(name)
		return declNode
	} else if node.Type == ast.VARREF {
		name := node.Children[0].TokenStart.Literal
		declNode := node.Symbols.LookupSymbol(name)
		return declNode
	} else if node.Type == ast.DOTOP {
		// TODO: Handle chain of dotops
		left := getType(&node.Children[0])
		var name string
		if left.name == "[" { // Special case for builtin list functions
			name = "list"
		} else {
			name = left.name
		}
		typeDeclNode := node.Symbols.LookupSymbol(name) // Class decl
		if typeDeclNode == nil {
			abortMsgf("Undeclared type: %s", name)
		}
		methodDecl := typeDeclNode.Children[1].Symbols.LookupSymbol(node.Children[1].TokenStart.Literal)
		return methodDecl
	} else if node.Type == ast.EXPRESSION {
		return lookUpDecl(node)
	}
	return nil // Can't happen?
}

// Get type from expression node.
func getType(node *ast.Node) *typeObj {
	switch node.Type {
	case ast.BINARYOP:
		left := getType(&node.Children[0])
		right := getType(&node.Children[1])

		if !compareTypes(left, right) { // All ops require left and right types be same.
			abortMsgf("Mismatched types: %s and %s", left.fullName, right.fullName)
		}
		if lexer.IsOperator([]rune(node.TokenStart.Literal)[0]) {
			if compareTypes(left, typeINT) || compareTypes(left, typeFLOAT) { // Math ops work on numbers.
				return left
			} else if node.TokenStart.Type == token.PLUS && compareTypes(left, typeSTRING) { // + works on strings.
				return left
			} else {
				abortMsg("Invalid operation.")
			}
		} else if node.TokenStart.Literal == ">=" || node.TokenStart.Literal == ">" || node.TokenStart.Literal == "<=" || node.TokenStart.Literal == "<" { // Comparison ops work on numbers, but return a bool.
			if compareTypes(left, typeINT) || compareTypes(left, typeFLOAT) {
				return typeBOOL
			} else {
				abortMsg("Invalid operation.")
			}
		} else if node.TokenStart.Literal == "==" {
			return typeBOOL
		} else if node.TokenStart.Literal == "&&" || node.TokenStart.Literal == "||" {
			if !compareTypes(left, typeBOOL) {
				abortMsg("Invalid operation.")
			}
			return typeBOOL
		}
		return left // Will this ever be reached?

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

	case ast.INDEXOP:
		// Check that left is a list or a map.
		// If list, then right should be int. Return inner type of left.
		// If map, then right should be first inner type of left. Return second inner type of right.
		left := getType(&node.Children[0])
		right := getType(&node.Children[1])

		if left.isList {
			if !compareTypes(right, typeINT) {
				abortMsg("List index must be int.")
			}
			return &left.inner[0]
		} else if left.isMap {

		} else {
			abortMsg("Invalid operation.")
		}

	case ast.VARREF:
		name := node.Children[0].TokenStart.Literal
		declNode := node.Symbols.LookupSymbol(name)
		if declNode == nil {
			abortMsgf("Referencing undeclared variable: %s", name)
		}
		return declType(declNode)

	case ast.FUNCCALL:
		//name := node.Children[0].TokenStart.Literal // TODO: Handle dot op.
		//declNode := node.Symbols.LookupSymbol(name)
		declNode := lookUpDecl(node.Children[0])

		checkFuncCall(node, declNode)

		return &declType(declNode).inner[0] // TODO: This will not work for multiple return...

	case ast.NEW:
		return buildTypeObj(&node.Children[0])

	case ast.LIST:
		obj := &typeObj{}
		obj.isContainer = true
		obj.isList = true
		obj.name = "["
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

package typechecker

import (
	"fmt"
	"knox/ast"
	"knox/lexer"
	"knox/token"
)

var prim primitives // Object holding the primitive types.

var currentFunc *ast.Node  // Keep track of current function to compare return type.
var currentClass *ast.Node // Keep track of current class to check self type.

// Analyze performs type checking on the entire AST.
func Analyze(node *ast.Node) {
	prim.Init()
	typecheck(node)
}

// #137 make sure main has return type void or int

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
					abortMsgf(node, "Mismatched types: %s and %s", leftType.fullName, exprType.fullName)
				}
				if leftType.isClass && !child.Symbols.IsDeclared(leftType.name) {
					abortMsgf(node, "Undeclared type: %s", leftType.name)
				}
			} else if node.Type == ast.VARASSIGN {
				// TODO: Fix member access bug.
				//decl := child.Symbols.LookupSymbol(node.Children[0].Children[0].TokenStart.Literal)
				//if decl == nil {
				//	abortMsgf("Referencing undeclared variable: %s", node.Children[0].Children[0].TokenStart.Literal)
				//}
				//leftType := declType(decl)
				leftType := getType(&node.Children[0])
				if !compareTypes(leftType, exprType) { // Do the types match?
					abortMsgf(node, "Mismatched types: %s and %s", leftType.fullName, exprType.fullName)
				}
			} else if node.Type == ast.IFSTATEMENT || node.Type == ast.WHILESTATEMENT {
				if !compareTypes(exprType, prim.typeBOOL) {
					abortMsg(node, "Conditionals require boolean expressions.")
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
				if compareTypes(funcReturnType, prim.typeVOID) && returnType.fullName == "" { // Check for return; and void type.
				} else if !compareTypes(&returnType.inner[0], &funcReturnType.inner[0]) {
					// TODO: Comparing the inner[0] is correct for single return types, but won't work for multiple. Need to expand compareType to handle this. buildTypeList and buildReturnList should probably not use inner for single return types, which would solve literals and simple types, then set isMulti to true and expand compareTypes to handle recursively comparing inner for multi.
					abortMsgf(node, "Incorrect return type: %v when expecting %v.", returnType.inner[0].fullName, funcReturnType.inner[0].fullName)
				}
			}
		} else if node.Type == ast.FORSTATEMENT {
			// TODO: Right should be a list. Left type should be right inner type.
			// TODO: Is this working?
			left := declType(&node.Children[0])
			right := getType(&node.Children[1])
			fmt.Println("Debugging...", right.fullName, right.isList, right.isClass, right.isPrimitive)
			if !right.isList && !right.isMap {
				abortMsg(node, "For loop requires a list or map")
			}
			if !compareTypes(left, &right.inner[0]) {
				abortMsg(node, "For loop element is incorrect type")
			}

		} else if child.Type == ast.FUNCDECL {
			currentFunc = &child
			typecheck(&child)
		} else if child.Type == ast.CLASS {
			currentClass = &child
			typecheck(&child)
		} else if child.Type == ast.LEFTEXPR {
			only := getType(&child.Children[0])
			if !compareTypes(only, prim.typeVOID) {
				abortMsg(node, "Expression must be of void type, not "+only.fullName)
			}
		} else {
			typecheck(&child)
		}
	}
}

func abortMsg(node *ast.Node, msg string) {
	fmt.Printf("Type error: %v. Line %v.\n", msg, node.TokenStart.Line)
	panic("Aborted.\n")
}

func abortMsgf(node *ast.Node, msg string, args ...interface{}) {
	s := fmt.Sprintf(". Line %v.\n", node.TokenStart.Line)
	fmt.Printf("Type error: "+msg+s, args...)
	panic("Aborted.\n")
}

func compareTypes(a *typeObj, b *typeObj) bool {
	// Handle literals.
	if a.isLiteral || b.isLiteral {
		// Type inference for number literals.
		if a.isNumber && b.isNumber {
			// TODO: Range check the literal.
			return true
		}
	}

	// TODO: Consider adding nil as a subtype of all reference types.
	// Special case for comparing reference types to nil.
	if (a.isClass || a.isContainer) && b.fullName == "nil" {
		return true
	}

	// Recursively check container type.
	if a.isContainer && b.isContainer {
		if len(a.inner) != len(b.inner) {
			return false
		}

		for i := range a.inner {
			if compareTypes(&a.inner[i], &b.inner[i]) == false {
				return false
			}
		}
		return true // All inner types matched.
	}

	// All other cases.
	return a.fullName == b.fullName
}

// Is this still needed?
func stringToType(s string) *typeObj {
	primitive := &typeObj{}
	primitive.fullName = s
	primitive.name = s
	primitive.isPrimitive = prim.IsPrimitiveType(s)
	primitive.isNumber = prim.IsNumberType(s)
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
	} else if node.Type == ast.CLASS { // TODO: Is this code ever used??
		classType := stringToType(node.Children[0].TokenStart.Literal)
		return classType
	}
	abortMsg(node, "Unknown type error.")
	return nil
}

// Builds up a type obj recursively given a varType AST node.
func buildTypeObj(node *ast.Node) *typeObj {
	obj := &typeObj{}

	if isSimple(node) {
		obj.isPrimitive = prim.IsPrimitiveType(getName(node))
		obj.isNumber = prim.IsNumberType(getName(node))
		obj.isClass = !obj.isPrimitive
		obj.fullName = getName(node)
		obj.name = obj.fullName
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
		obj.name = getName(node)
		obj.fullName = obj.name + "["
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

	//declNode := node.Symbols.LookupSymbol(name)
	if declNode == nil {
		abortMsgf(node, "Calling undeclared function: %s", name)
	}

	// Check number of args to number of params.
	if len(node.Children)-1 != len(declNode.Children[1].Children) {
		abortMsg(node, "Incorrect number of arguments.")
	}
	// Check types of args to types of params.
	for i := 1; i < len(node.Children); i++ {
		argType := getType(&node.Children[i])
		expectedType := declType(&declNode.Children[1].Children[i-1])
		if !compareTypes(argType, expectedType) {
			abortMsgf(node, "Mismatched type in function argument.")
		}
	}
}

// TODO: Why doesn't this take a pointer?
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
			abortMsgf(&node, "Undeclared type: %s", name)
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
			abortMsgf(node, "Mismatched types: %s and %s", left.fullName, right.fullName)
		}
		if lexer.IsOperator([]rune(node.TokenStart.Literal)[0]) {
			//if compareTypes(left, prim.typeINT) || compareTypes(left, prim.typeFLOAT) { // Math ops work on numbers.
			if left.isNumber && right.isNumber {
				// TODO: Will this coerce a INTLITERAL to an INT? x + 1 is not the same as 1 + x here.
				return left
			} else if node.TokenStart.Type == token.PLUS && compareTypes(left, prim.typeSTRING) {
				// + works on strings.
				// TODO: Need to properly augment the AST with concat info.
				node.TokenStart.Literal = "concat"
				return left
			} else {
				abortMsg(node, "Invalid operation.") // TODO: Improve this error message.
			}
		} else if node.TokenStart.Literal == ">=" || node.TokenStart.Literal == ">" || node.TokenStart.Literal == "<=" || node.TokenStart.Literal == "<" { // Comparison ops work on numbers, but return a bool.
			if left.isNumber && right.isNumber {
				//if compareTypes(left, prim.typeINT) || compareTypes(left, prim.typeFLOAT) {
				return prim.typeBOOL
			} else {
				abortMsg(node, "Invalid operation.") // TODO: Improve this error message.
			}
		} else if node.TokenStart.Literal == "==" {
			return prim.typeBOOL
		} else if node.TokenStart.Literal == "&&" || node.TokenStart.Literal == "||" {
			if !compareTypes(left, prim.typeBOOL) {
				abortMsg(node, "Invalid operation.") // TODO: Improve this error message.
			}
			return prim.typeBOOL
		}
		return left // Will this ever be reached?

	case ast.UNARYOP:
		single := getType(&node.Children[0])
		if node.TokenStart.Type == token.BANG {
			if !compareTypes(single, prim.typeBOOL) {
				abortMsg(node, "Invalid operation.") // TODO: Improve this error message.
			}
		} else if node.TokenStart.Type == token.PLUS || node.TokenStart.Type == token.MINUS {
			//if !compareTypes(single, prim.typeINT) && !compareTypes(single, prim.typeFLOAT) {
			if !single.isNumber {
				abortMsg(node, "Invalid operation.") // TODO: Improve this error message.
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
			if !compareTypes(right, prim.typeINT) {
				abortMsg(node, "List index must be int.")
			}
			return &left.inner[0]
		} else if left.isMap {
			// TODO
		} else {
			abortMsg(node, "Invalid operation.") // TODO: Improve this error message.
		}

	// Member access
	case ast.DOTOP:
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
			abortMsgf(node, "Undeclared type: %s", name)
		}

		memberDecl := typeDeclNode.Children[1].Symbols.LookupSymbol(node.Children[1].TokenStart.Literal)
		if memberDecl == nil {
			abortMsgf(node, "Referencing undeclared member: %s", node.Children[1].TokenStart.Literal)
		}
		return declType(memberDecl)

	case ast.VARREF:
		name := node.Children[0].TokenStart.Literal
		declNode := node.Symbols.LookupSymbol(name)
		if declNode == nil {
			abortMsgf(node, "Referencing undeclared variable: %s", name)
		}
		return declType(declNode)

	case ast.FUNCCALL:
		//name := node.Children[0].TokenStart.Literal // TODO: Handle dot op.
		//declNode := node.Symbols.LookupSymbol(name)
		declNode := lookUpDecl(node.Children[0])

		checkFuncCall(node, declNode)

		return &declType(declNode).inner[0] // TODO: This will not work for multiple return...

	case ast.CAST:
		// Check that left and right are both primitive.
		// We will rely on C's casting rules for the semantics.
		typeLiteral := node.Children[1].TokenStart.Literal
		left := getType(&node.Children[0])
		isRightPrimitive := prim.IsPrimitiveType(typeLiteral)

		if !left.isPrimitive || !isRightPrimitive {
			abortMsgf(node, "Illegal cast from %s to %s.", node.Children[0].TokenStart.Literal, node.Children[1].TokenStart.Literal)
		}

		return stringToType(typeLiteral)

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
				abortMsg(node, "Mismatched types in list literal.")
			}
			itemType = obj.inner[i].fullName
		}
		// TODO: Handle if no items.
		obj.fullName += itemType + "]"
		return obj

	case ast.SELF:
		return declType(currentClass)
	case ast.INT:
		return prim.typeINTLITERAL
	case ast.FLOAT:
		return prim.typeFLOATLITERAL
	case ast.STRING:
		return prim.typeSTRING
	case ast.BOOL:
		return prim.typeBOOL
	case ast.NIL:
		return prim.typeNIL

	case ast.EXPRESSION:
		return getType(&node.Children[0])
	}

	return nil
}

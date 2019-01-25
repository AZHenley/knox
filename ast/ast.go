package ast

import (
	"fmt"
	"knox/token"
	"strings"
)

// NodeType is a string.
type NodeType string

// Node is for the AST.
type Node struct {
	Type       NodeType
	Children   []Node
	TokenStart token.Token
	//symbols symtable.SymTable // Only blocks get a symbol table.
}

// Predefined AST node types.
const (
	PROGRAM    = "PROGRAM"    // Variable children. One for each funcdecl.
	BLOCK      = "BLOCK"      // Variable children. One for each statement.
	EXPRESSION = "EXPRESSION" // One child. Tree of binary ops, unary ops, and primaries.
	// TODO: Consider removing this.
	BINARYOP = "BINARYOP" // Two children.
	UNARYOP  = "UNARYOP"  // One child.
	INDEXOP  = "INDEXOP"  // Variable children.
	VARDECL  = "VARDECL"  // Two or three children. Name, type, expression for assignment.
	// TODO: Consider making the third child a VARASSIGN.
	VARTYPE = "VARTYPE" // One child. Name.
	// TODO: Support lists.
	VARASSIGN = "VARASSIGN" // Two children. Varref and expression.
	FUNCDECL  = "FUNCDECL"  // Four children. Name, paramlist for params, paramlist for return, block.
	PARAMLIST = "PARAMLIST" // Variable children. Pairs of name and type.
	// TODO: Consider making the pairs a VARDECL node.
	IFSTATEMENT = "IFSTATEMENT" // Three children. Condition, if block, else block (chain elif/else).
	// TODO: Consider changing this.
	FORSTATEMENT   = "FORSTATEMENT"   // Four children. Init, condition, afterthought, block.
	WHILESTATEMENT = "WHILESTATEMENT" // Two children. Condition and block.
	JUMPSTATEMENT  = "JUMPSTATEMENT"  // Leaf.
	VARREF         = "VARREF"         // Variable children. Name and list of expressions for array indices.
	FUNCCALL       = "FUNCCALL"       // Variable children. Name then one expression for each parameter.
	INT            = "INT"            // Leaf.
	FLOAT          = "FLOAT"          // Leaf.
	STRING         = "STRING"         // Leaf.
	BOOL           = "BOOL"           // Leaf.
	NIL            = "NIL"            // Leaf.
	IDENT          = "IDENT"          // Leaf.
)

// Print AST.
func Print(node Node) {
	printUtil(node, 0)
}

func printUtil(node Node, depth int) {
	var prefix = strings.Repeat(">", depth)
	fmt.Printf("%s %s %s\n", prefix, node.Type, node.TokenStart.Literal)
	for _, c := range node.Children {
		printUtil(c, depth+1)
	}
}

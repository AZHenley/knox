package ast

import "knox/token"

// NodeType is a string.
type NodeType string

// Node is for the AST.
type Node struct {
	nodeType   NodeType
	children   []Node
	tokenStart token.Token
	//symbols symtable.SymTable // Only blocks get a symbol table.
}

// Predefined
const (
	EXPRESSION     = "EXPRESSION"     // One child. Tree of binary ops, unary ops, and primaries.
	BINARYOP       = "BINARYOP"       // Two children.
	UNARYOP        = "UNARYOP"        // One child.
	VARDECL        = "VARDECL"        // 2-3 children. Third is optional expression for assignment.
	VARASSIGN      = "VARASSIGN"      // Two children. Name and expression.
	FUNCDECL       = "FUNCDECL"       // Three children. Name, list of param pairs, list of return pairs
	IFSTATEMENT    = "IFSTATEMENT"    // Three children. Condition, if block, else block (chain elif/else).
	FORSTATEMENT   = "FORSTATEMENT"   // Four children. Init, condition, afterthought, block.
	WHILESTATEMENT = "WHILESTATEMENT" // Two children. Condition and block.
	JUMPSTATEMENT  = "JUMPSTATEMENT"  // Leaf.
	VARREF         = "VARREF"         // Leaf.
	FUNCCALL       = "FUNCCALL"       // Leaf.
	INT            = "INT"            // Leaf.
	FLOAT          = "FLOAT"          // Leaf.
	STRING         = "STRING"         // Leaf.
)

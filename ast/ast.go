package ast

import "knox/token"

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
	PROGRAM        = "PROGRAM"        // Variable children. One for each funcdecl.
	BLOCK          = "BLOCK"          // Variable children. One for each statement.
	EXPRESSION     = "EXPRESSION"     // One child. Tree of binary ops, unary ops, and primaries.
	BINARYOP       = "BINARYOP"       // Two children.
	UNARYOP        = "UNARYOP"        // One child.
	VARDECL        = "VARDECL"        // Two or three children. Third is optional expression for assignment.
	VARASSIGN      = "VARASSIGN"      // Two children. Name and expression.
	FUNCDECL       = "FUNCDECL"       // Four children. Name, paramlist for params, paramlist for return, block.
	PARAMLIST      = "PARAMLIST"      // Variable children. Pairs of name and type.
	IFSTATEMENT    = "IFSTATEMENT"    // Three children. Condition, if block, else block (chain elif/else).
	FORSTATEMENT   = "FORSTATEMENT"   // Four children. Init, condition, afterthought, block.
	WHILESTATEMENT = "WHILESTATEMENT" // Two children. Condition and block.
	JUMPSTATEMENT  = "JUMPSTATEMENT"  // Leaf.
	VARREF         = "VARREF"         // Variable children. Variable name and list of expressions for array indices.
	FUNCCALL       = "FUNCCALL"       // Variable children. One expression for each parameter.
	INT            = "INT"            // Leaf.
	FLOAT          = "FLOAT"          // Leaf.
	STRING         = "STRING"         // Leaf.
	IDENT          = "IDENT"          // Leaf.
)

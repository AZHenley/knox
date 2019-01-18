package ast

import "knox/token"

// ASTNodeType is a string.
type ASTNodeType string

// ASTNode is for the AST.
type ASTNode struct {
	nodeType   ASTNodeType
	children   []ASTNode
	tokenStart token.Token
	// symbols symtable.SymTable // Only blocks get a symbol table.
}

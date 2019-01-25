package symtable

import "knox/ast"

// SymTable is a table for symbol entries.
type SymTable struct {
	entries map[string]*ast.Node
	parent  *SymTable
}

// New creates an initialized symbol table.
func New(parentTable *SymTable) *SymTable {
	s := &SymTable{}
	s.entries = make(map[string]*ast.Node)
	return s
}

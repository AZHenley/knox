package ast

// Moved SymTable into AST package to avoid circular depedency.

// SymTable is a table for symbol entries.
type SymTable struct {
	Entries map[string]*Node
	Parent  *SymTable
}

// NewSymTable creates an initialized symbol table.
func NewSymTable() *SymTable {
	s := &SymTable{}
	s.Entries = make(map[string]*Node)
	return s
}

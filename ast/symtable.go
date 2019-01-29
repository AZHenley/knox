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

// InsertSymbol into symbol table, recursively checks if it exists, returns true for success
func (st *SymTable) InsertSymbol(symbol string, node *Node) bool {
	if st.IsDeclared(symbol) {
		return false
	}
	st.Entries[symbol] = node
	return true
}

// IsDeclared returns if symbol is declared, recursively.
func (st *SymTable) IsDeclared(symbol string) bool {
	if _, ok := st.Entries[symbol]; ok {
		return true
	}
	if st.Parent == nil {
		return false
	}
	return st.Parent.IsDeclared(symbol)
}

// LookupSymbol does a recursive lookup.
func (st *SymTable) LookupSymbol(symbol string) *Node {
	if node, ok := st.Entries[symbol]; ok {
		return node
	}
	if st.Parent == nil {
		return nil
	}
	return st.Parent.LookupSymbol(symbol)
}

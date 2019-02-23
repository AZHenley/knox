package builtin

import (
	"io/ioutil"
	"knox/ast"
	"knox/lexer"
	"knox/parser"
)

// import "knox/ast"

// Init setups the builtin.
func Init(node *ast.Node) *ast.Node {

	node = createBuiltin("list", node)
	node = createBuiltin("stl", node)

	return node
}

func createBuiltin(name string, node *ast.Node) *ast.Node {
	code, err := ioutil.ReadFile("builtin/" + name + ".knox")
	if err != nil {
		panic(err)
	}
	l := lexer.New(string(code) + "\n")
	p := parser.New(l)
	a := p.Program()

	// Merge global symtables
	for key, val := range a.Symbols.Entries {
		node.Symbols.InsertSymbol(key, val)
	}
	// Update all the parent nodes
	for _, declaration := range a.Children {
		if declaration.Type == ast.FUNCDECL {
			declaration.Children[3].Symbols = node.Symbols
		} else if declaration.Type == ast.CLASS {
			declaration.Children[0].Symbols = node.Symbols
		}
	}

	node.Children = append(node.Children, a)
	return node
}

func createClass(name string, sym *ast.SymTable) *ast.Node {
	var node ast.Node
	node.Type = ast.CLASS

	var identNode ast.Node
	identNode.Type = ast.IDENT
	identNode.TokenStart.Literal = name

	var blockNode ast.Node
	blockNode.Type = ast.BLOCK
	st := ast.NewSymTable()
	st.Parent = sym
	blockNode.Symbols = st

	node.Children = append(node.Children, identNode, blockNode)

	return &node
}

// // Functions for casting.
// func createCasts() {

// }

// // Placeholder for std lib.
// func createStdLib() {

// }

// // Builtin methods for lists.
// func createListFuncs() {
// 	fLength := "func Length() int {return -1;}"
// 	// Call parser
// 	// Rename node
// }

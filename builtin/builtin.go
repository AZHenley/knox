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

	//listClass := createListFuncs()
	// b := createClass("[", node.Symbols)

	// //bb := createFunction("append")

	// l1 := lexer.New("class list { func append(x : int) void {} func}")
	// p1 := parser.New(l1)
	// a1 := p1.Program().Children[0]

	// ast.Print(a1)

	// l := lexer.New("func append(x : int) void {}" + "\n")
	// p := parser.New(l)
	// a := p.Program().Children[0]

	// ast.Print(*b)
	// ast.Print(a)

	// Create all the functions and attach them to b
	// Attach b to node

	code, err := ioutil.ReadFile("builtin/list.knox") // TODO: Support multiple files.
	//code, err := ioutil.ReadFile("examples/chain.knox")
	if err != nil {
		panic(err)
	}
	l := lexer.New(string(code) + "\n")
	p := parser.New(l)
	a := p.Program()

	// TODO: Build builtin first in case user's program conflicts.

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

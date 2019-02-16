package parser

import (
	"fmt"
	"knox/ast"
	"knox/lexer"

	"knox/token"
)

// Parser object.
type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	//errors    []string
	curSymTable *ast.SymTable
}

// New lexer.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken() // Not sure why these are needed.
	return p
}

func (p *Parser) abort(t token.TokenType) {
	fmt.Printf("Expected %s, got %s instead.\n", t, p.curToken.Type)
	panic("Aborted.\n")
}

func (p *Parser) abortMsg(msg string) {
	fmt.Println(msg)
	panic("Aborted.\n")
}

// forward token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()

	//fmt.Println("Currently at " + string(p.curToken.Literal))
}

// determinate current token is t or not.
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// determinate next token is t or not
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) consume(t token.TokenType) bool {
	if p.curTokenIs(t) {
		p.nextToken()
		return true
	}
	p.abort(t)
	return false // Can't happen.
}

////
// Grammar rules.
////

// Program parses a program.
func (p *Parser) Program() ast.Node {
	var progNode ast.Node
	progNode.Type = ast.PROGRAM

	st := ast.NewSymTable()
	st.Parent = nil
	p.curSymTable = st
	progNode.Symbols = st

	for !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.FUNCTION) {
			progNode.Children = append(progNode.Children, p.funcDecl())
		} else if p.curTokenIs(token.CLASS) {
			progNode.Children = append(progNode.Children, p.classDecl())
		} else {
			p.abortMsg("Expected function or class.")
		}

	}
	return progNode
}

// classDecl = "class" ident classBlock
func (p *Parser) classDecl() ast.Node {
	var classNode ast.Node
	classNode.Type = ast.CLASS
	p.consume(token.CLASS)

	var identNode ast.Node
	identNode.Type = ast.IDENT
	identNode.TokenStart = p.curToken
	classNode.Children = append(classNode.Children, identNode)

	success := p.curSymTable.InsertSymbol(p.curToken.Literal, &classNode)
	if !success {
		p.abortMsg("Class already exists.")
	}
	p.consume(token.IDENT)

	classNode.Children = append(classNode.Children, p.classBlock())
	return classNode
}

// classBlock = "{" {varDecl | funcDecl} "}"
func (p *Parser) classBlock() ast.Node {
	var blockNode ast.Node
	blockNode.Type = ast.BLOCK

	st := ast.NewSymTable()
	st.Parent = p.curSymTable
	p.curSymTable = st
	blockNode.Symbols = st

	p.consume(token.LBRACE)
	for !p.curTokenIs(token.RBRACE) {
		if p.curTokenIs(token.VAR) {
			blockNode.Children = append(blockNode.Children, p.varDecl())
		} else if p.curTokenIs(token.FUNCTION) {
			blockNode.Children = append(blockNode.Children, p.funcDecl())
		}
	}
	p.consume(token.RBRACE)

	p.curSymTable = st.Parent

	return blockNode
}

// funcDecl = "func" ident paramList returnList block
func (p *Parser) funcDecl() ast.Node {
	var funcNode ast.Node
	funcNode.Type = ast.FUNCDECL
	p.consume(token.FUNCTION)

	var identNode ast.Node
	identNode.Type = ast.IDENT
	identNode.TokenStart = p.curToken
	funcNode.Children = append(funcNode.Children, identNode)

	success := p.curSymTable.InsertSymbol(p.curToken.Literal, &funcNode)
	if !success {
		p.abortMsg("Function already exists.")
	}
	p.consume(token.IDENT)

	funcNode.Children = append(funcNode.Children, p.paramList())
	funcNode.Children = append(funcNode.Children, p.returnList())
	funcNode.Children = append(funcNode.Children, p.block())
	return funcNode
}

func (p *Parser) paramList() ast.Node {
	var paramNode ast.Node
	paramNode.Type = ast.PARAMLIST

	p.consume(token.LPAREN)
	for !p.curTokenIs(token.RPAREN) {

		var varNode ast.Node
		varNode.Type = ast.VARDECL
		varNode.Symbols = p.curSymTable

		var identNode ast.Node
		identNode.Type = ast.IDENT
		identNode.TokenStart = p.curToken

		success := p.curSymTable.InsertSymbol(p.curToken.Literal, &varNode)
		if !success {
			p.abortMsg("Variable already exists.")
		}

		p.consume(token.IDENT)
		p.consume(token.COLON)

		varNode.Children = append(varNode.Children, identNode)
		varNode.Children = append(varNode.Children, p.varType())
		paramNode.Children = append(paramNode.Children, varNode)
		if p.curTokenIs(token.COMMA) {
			p.consume(token.COMMA)
			if p.curTokenIs(token.RPAREN) { // No comma before paren.
				p.abort(token.IDENT)
			}
		}
	}
	p.consume(token.RPAREN)
	return paramNode
}

// returnList = varType | "(" varType {"," varType} ")"
func (p *Parser) returnList() ast.Node {
	var returnNode ast.Node
	returnNode.Type = ast.RETURNLIST

	if !p.curTokenIs("(") {
		returnNode.Children = append(returnNode.Children, p.varType())
	} else {
		p.consume(token.LPAREN)
		returnNode.Children = append(returnNode.Children, p.varType())
		for !p.curTokenIs(token.RPAREN) {
			p.consume(token.COMMA)
			returnNode.Children = append(returnNode.Children, p.varType())
		}
		p.nextToken()
	}

	return returnNode
}

func (p *Parser) block() ast.Node {
	var blockNode ast.Node
	blockNode.Type = ast.BLOCK

	st := ast.NewSymTable()
	st.Parent = p.curSymTable
	p.curSymTable = st
	blockNode.Symbols = st

	p.consume(token.LBRACE)
	for !p.curTokenIs(token.RBRACE) {
		//p.statement()
		blockNode.Children = append(blockNode.Children, p.statement())
	}
	p.consume(token.RBRACE)

	p.curSymTable = st.Parent

	return blockNode
}

func (p *Parser) statement() ast.Node {
	var statementNode ast.Node

	if p.curTokenIs(token.VAR) {
		statementNode = p.varDecl()
		p.consume(token.SEMICOLON)
	} else if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.LPAREN) {
		statementNode = p.funcCall()
		p.consume(token.SEMICOLON)
	} else if p.curTokenIs(token.IDENT) {
		statementNode = p.varAssignment()
		p.consume(token.SEMICOLON)
	} else if p.curTokenIs(token.IF) {
		statementNode = p.ifStatement()
	} else if p.curTokenIs(token.FOR) {
		statementNode = p.forStatement()
	} else if p.curTokenIs(token.WHILE) {
		statementNode = p.whileStatement()
	} else if p.curTokenIs(token.RETURN) || p.curTokenIs(token.CONTINUE) || p.curTokenIs(token.BREAK) {
		statementNode = p.jumpStatement()
		p.consume(token.SEMICOLON)
	} else {
		p.abortMsg("Expected statement.")
	}
	return statementNode
}

// varDecl = "var" ident ":" type assignOp expr
func (p *Parser) varDecl() ast.Node {
	var varNode ast.Node
	varNode.Type = ast.VARDECL
	varNode.Symbols = p.curSymTable

	p.consume(token.VAR)

	for p.curTokenIs(token.IDENT) {
		var identNode ast.Node
		identNode.Type = ast.IDENT
		identNode.TokenStart = p.curToken
		success := p.curSymTable.InsertSymbol(p.curToken.Literal, &varNode)
		if !success {
			p.abortMsg("Variable already exists.")
		}
		p.consume(token.IDENT)
		p.consume(token.COLON)

		varNode.Children = append(varNode.Children, identNode)
		varNode.Children = append(varNode.Children, p.varType())

		if !p.curTokenIs(token.COMMA) {
			break
		}
		p.consume(token.COMMA)
	}

	p.consume(token.ASSIGN)
	varNode.Children = append(varNode.Children, p.expr())

	return varNode
}

// funcCall = ident argList
func (p *Parser) funcCall() ast.Node {
	var funcNode ast.Node
	funcNode.Type = ast.FUNCCALL

	var identNode ast.Node
	identNode.Type = ast.IDENT
	identNode.TokenStart = p.curToken
	p.consume(token.IDENT)

	funcNode.Children = append(funcNode.Children, identNode)
	var nodes = p.argList()
	funcNode.Children = append(funcNode.Children, nodes...)

	return funcNode
}

// argList = "(" {expr ","} [expr] ")"
func (p *Parser) argList() []ast.Node {
	var argNodes []ast.Node

	p.consume(token.LPAREN)
	for !p.curTokenIs(token.RPAREN) {
		argNodes = append(argNodes, p.expr())
		if p.curTokenIs(token.COMMA) {
			p.consume(token.COMMA)
			if p.curTokenIs(token.RPAREN) { // No comma before paren.
				p.abort(token.IDENT)
			}
		}
	}
	p.nextToken()
	return argNodes
}

// varType = ident | "[" varType "]" | ident "[" varType {"," varType} "]"
func (p *Parser) varType() ast.Node {
	var typeNode ast.Node
	typeNode.Type = ast.VARTYPE

	var identNode ast.Node
	identNode.Type = ast.IDENT
	identNode.TokenStart = p.curToken
	typeNode.Children = append(typeNode.Children, identNode)

	if p.curTokenIs(token.IDENT) && !p.peekTokenIs(token.LBRACKET) { // simple type
		p.consume(token.IDENT)
	} else if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.LBRACKET) { // container
		p.consume(token.IDENT)
		p.consume(token.LBRACKET)

		typeNode.Children = append(typeNode.Children, p.varType())
		for !p.curTokenIs(token.RBRACKET) {
			p.consume(token.COMMA)
			typeNode.Children = append(typeNode.Children, p.varType())
		}
		p.consume(token.RBRACKET)
	} else if p.curTokenIs(token.LBRACKET) { // list
		p.consume(token.LBRACKET)
		typeNode.Children = append(typeNode.Children, p.varType())
		p.consume(token.RBRACKET)
	}
	return typeNode
}

// varRef = ident {"[" expr "]"}
func (p *Parser) varRef() ast.Node {
	var refNode ast.Node
	refNode.Type = ast.VARREF
	refNode.Symbols = p.curSymTable

	var identNode ast.Node
	identNode.Type = ast.IDENT
	identNode.TokenStart = p.curToken
	refNode.Children = append(refNode.Children, identNode)
	p.consume(token.IDENT)

	for p.curTokenIs(token.LBRACKET) {
		p.nextToken()
		refNode.Children = append(refNode.Children, p.expr())
		p.consume(token.RBRACKET)
	}

	return refNode
}

// varAssignment = varRef assignOp expr
func (p *Parser) varAssignment() ast.Node {
	var assignNode ast.Node
	assignNode.Type = ast.VARASSIGN
	assignNode.Symbols = p.curSymTable

	for p.curTokenIs(token.IDENT) {
		assignNode.Children = append(assignNode.Children, p.varRef())
		if !p.curTokenIs(token.COMMA) {
			break
		}
		p.consume(token.COMMA)
	}
	p.consume(token.ASSIGN)
	assignNode.Children = append(assignNode.Children, p.expr())

	return assignNode
}

// ifStatement = "if" expr block
func (p *Parser) ifStatement() ast.Node {
	var statementNode ast.Node
	statementNode.Type = ast.IFSTATEMENT

	p.consume(token.IF)
	statementNode.Children = append(statementNode.Children, p.expr())
	statementNode.Children = append(statementNode.Children, p.block())

	for p.curTokenIs(token.ELSE) && p.peekTokenIs(token.IF) {
		p.nextToken()
		p.nextToken()
		statementNode.Children = append(statementNode.Children, p.expr())
		statementNode.Children = append(statementNode.Children, p.block())
	}

	if p.curTokenIs(token.ELSE) {
		p.nextToken()
		statementNode.Children = append(statementNode.Children, p.block())
	}

	return statementNode
}

// forStatement = "for" forClause block
func (p *Parser) forStatement() ast.Node {
	var statementNode ast.Node
	statementNode.Type = ast.FORSTATEMENT

	p.consume(token.FOR)
	// TODO: Make this part reuse code from varDecl.
	p.consume(token.VAR)
	p.consume(token.IDENT)
	p.consume(token.COLON)
	p.consume(token.IDENT) // TODO: Vartype.
	p.consume(token.IN)
	statementNode.Children = append(statementNode.Children, p.expr())
	statementNode.Children = append(statementNode.Children, p.block())

	return statementNode
}

func (p *Parser) whileStatement() ast.Node {
	var statementNode ast.Node
	statementNode.Type = ast.WHILESTATEMENT

	p.consume(token.WHILE)
	statementNode.Children = append(statementNode.Children, p.expr())
	statementNode.Children = append(statementNode.Children, p.block())

	return statementNode
}

// jumpStatement = "continue" | "break" | "return" [expr {"," expr}]
func (p *Parser) jumpStatement() ast.Node {
	var statementNode ast.Node
	statementNode.Type = ast.JUMPSTATEMENT
	statementNode.TokenStart = p.curToken
	p.nextToken()

	if !p.curTokenIs(token.SEMICOLON) {
		statementNode.Children = append(statementNode.Children, p.expr())
		for !p.curTokenIs(token.SEMICOLON) {
			p.consume(token.COMMA)
			statementNode.Children = append(statementNode.Children, p.expr())
		}
	}
	return statementNode
}

func (p *Parser) expr() ast.Node {
	var exprNode ast.Node
	exprNode.Type = ast.EXPRESSION
	exprNode.Symbols = p.curSymTable // Make it easier for analyzers to do look ups.

	exprNode.Children = append(exprNode.Children, p.logical())
	return exprNode
}

func (p *Parser) logical() ast.Node {
	var node = p.equality()
	for p.curTokenIs(token.AND) || p.curTokenIs(token.OR) {
		var binaryNode ast.Node
		binaryNode.Type = ast.BINARYOP
		binaryNode.TokenStart = p.curToken

		p.nextToken()
		binaryNode.Children = append(binaryNode.Children, node)
		binaryNode.Children = append(binaryNode.Children, p.equality())
		node = binaryNode
	}
	return node
}

func (p *Parser) equality() ast.Node {
	var node = p.comparison()
	for p.curTokenIs(token.EQ) || p.curTokenIs(token.NOTEQ) {
		var binaryNode ast.Node
		binaryNode.Type = ast.BINARYOP
		binaryNode.TokenStart = p.curToken

		p.nextToken()
		binaryNode.Children = append(binaryNode.Children, node)
		binaryNode.Children = append(binaryNode.Children, p.comparison())
		node = binaryNode
	}
	return node
}

func (p *Parser) comparison() ast.Node {
	var node = p.addition()
	for p.curTokenIs(token.GT) || p.curTokenIs(token.GTEQ) || p.curTokenIs(token.LT) || p.curTokenIs(token.LTEQ) {
		var binaryNode ast.Node
		binaryNode.Type = ast.BINARYOP
		binaryNode.TokenStart = p.curToken

		p.nextToken()
		binaryNode.Children = append(binaryNode.Children, node)
		binaryNode.Children = append(binaryNode.Children, p.addition())
		node = binaryNode
	}
	return node
}

func (p *Parser) addition() ast.Node {
	var node = p.multiplication()
	for p.curTokenIs(token.PLUS) || p.curTokenIs(token.MINUS) {
		var binaryNode ast.Node
		binaryNode.Type = ast.BINARYOP
		binaryNode.TokenStart = p.curToken

		p.nextToken()
		binaryNode.Children = append(binaryNode.Children, node)
		binaryNode.Children = append(binaryNode.Children, p.multiplication())
		node = binaryNode
	}
	return node
}

func (p *Parser) multiplication() ast.Node {
	var node = p.unary()
	for p.curTokenIs(token.ASTERISK) || p.curTokenIs(token.SLASH) {
		var binaryNode ast.Node
		binaryNode.Type = ast.BINARYOP
		binaryNode.TokenStart = p.curToken

		p.nextToken()
		binaryNode.Children = append(binaryNode.Children, node)
		binaryNode.Children = append(binaryNode.Children, p.unary())

		node = binaryNode
	}
	return node
}

func (p *Parser) unary() ast.Node {
	if p.curTokenIs(token.BANG) || p.curTokenIs(token.PLUS) || p.curTokenIs(token.MINUS) {
		var unaryNode ast.Node
		unaryNode.Type = ast.UNARYOP
		unaryNode.TokenStart = p.curToken

		p.nextToken()
		unaryNode.Children = append(unaryNode.Children, p.unary())

		return unaryNode
	}
	return p.postfix()
}

// postfix = paran {"[" expr "]" | argList | "." ident}
func (p *Parser) postfix() ast.Node {
	var node ast.Node
	node = p.paran()

	for p.curTokenIs(token.LBRACKET) || p.curTokenIs(token.LPAREN) || p.curTokenIs(token.DOT) {
		if p.curTokenIs(token.LPAREN) {
			var postNode ast.Node
			postNode.Children = append(postNode.Children, node)
			postNode.Type = ast.FUNCCALL
			postNode.Symbols = p.curSymTable

			var nodes = p.argList()
			postNode.Children = append(postNode.Children, nodes...)

			node = postNode
		} else if p.curTokenIs(token.LBRACKET) {
			var postNode ast.Node
			postNode.Children = append(postNode.Children, node)
			postNode.Type = ast.INDEXOP

			p.nextToken()
			postNode.Children = append(postNode.Children, p.expr())
			p.consume(token.RBRACKET)

			node = postNode
		} else if p.curTokenIs(token.DOT) {
			var postNode ast.Node
			postNode.Children = append(postNode.Children, node)
			postNode.Type = ast.DOTOP

			p.nextToken()
			var identNode ast.Node
			identNode.Type = ast.IDENT
			identNode.TokenStart = p.curToken
			postNode.Children = append(postNode.Children, identNode)
			p.consume(token.IDENT)

			node = postNode
		}
	}
	return node
}

func (p *Parser) paran() ast.Node {
	var paranNode ast.Node

	if p.curTokenIs(token.LPAREN) {
		p.nextToken()
		paranNode = p.expr()
		p.consume(token.RPAREN)
	} else {
		paranNode = p.special()
	}

	return paranNode
}

func (p *Parser) special() ast.Node {
	if p.curTokenIs(token.NEW) {
		p.nextToken()
		var newNode ast.Node
		newNode.Type = ast.NEW
		newNode.Children = append(newNode.Children, p.varType())
		return newNode
	} else {
		return p.primary()
	}
}

func (p *Parser) listLiteral() ast.Node {
	var listNode ast.Node
	listNode.Type = ast.LIST

	p.consume(token.LBRACKET)
	for !p.curTokenIs(token.RBRACKET) {
		listNode.Children = append(listNode.Children, p.expr())
		if p.curTokenIs(token.COMMA) { // TODO: This currently allows a comma proceeded by rbracket.
			p.consume(token.COMMA)
		}
	}
	p.consume(token.RBRACKET)

	return listNode
}

// primary = varRef | INT | FLOAT | STRING | "false" | "true" | "nil" | "(" expr ")" | listLiteral
func (p *Parser) primary() ast.Node {
	var primaryNode ast.Node
	primaryNode.TokenStart = p.curToken

	switch p.curToken.Type {
	case token.INT:
		primaryNode.Type = ast.INT
	case token.FLOAT:
		primaryNode.Type = ast.FLOAT
	case token.STRING:
		primaryNode.Type = ast.STRING
	case token.TRUE, token.FALSE:
		primaryNode.Type = ast.BOOL
	case token.NIL:
		primaryNode.Type = ast.NIL
	case token.IDENT:
		primaryNode.Type = ast.VARREF
		primaryNode.Symbols = p.curSymTable
		var identNode ast.Node
		identNode.Type = ast.IDENT
		identNode.TokenStart = p.curToken
		primaryNode.Children = append(primaryNode.Children, identNode)
	case token.LBRACKET:
		return p.listLiteral()
	}

	p.nextToken()
	return primaryNode

	// if p.curTokenIs(token.INT) || p.curTokenIs(token.FLOAT) || p.curTokenIs(token.STRING) || p.curTokenIs(token.TRUE) || p.curTokenIs(token.FALSE) || p.curTokenIs(token.NIL) {
	// 	p.nextToken()
	// } else if p.curTokenIs(token.LPAREN) {
	// 	p.nextToken()
	// 	p.expr()
	// 	p.consume(token.RPAREN)
	// } else if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.LPAREN) {
	// 	p.funcCall()
	// } else {
	// 	p.varRef()
	// }
}

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
	//ast ast.ASTNode
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

	for !p.curTokenIs(token.EOF) {
		progNode.Children = append(progNode.Children, p.funcDecl())
	}
	return progNode
}

func (p *Parser) funcDecl() ast.Node {
	var funcNode ast.Node
	funcNode.Type = ast.FUNCDECL

	p.consume(token.FUNCTION)

	var identNode ast.Node
	identNode.Type = ast.IDENT
	identNode.TokenStart = p.curToken
	funcNode.Children = append(funcNode.Children, identNode)
	p.consume(token.IDENT)

	funcNode.Children = append(funcNode.Children, p.paramList())
	funcNode.Children = append(funcNode.Children, p.block())
	return funcNode
}

func (p *Parser) paramList() ast.Node {
	var paramNode ast.Node
	paramNode.Type = ast.PARAMLIST

	p.consume(token.LPAREN)
	for !p.curTokenIs(token.RPAREN) {
		var identNode ast.Node
		identNode.Type = ast.IDENT
		identNode.TokenStart = p.curToken
		p.consume(token.IDENT)
		p.consume(token.COLON)

		paramNode.Children = append(paramNode.Children, identNode)
		paramNode.Children = append(paramNode.Children, p.varType())
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

func (p *Parser) block() ast.Node {
	var blockNode ast.Node
	blockNode.Type = ast.BLOCK

	p.consume(token.LBRACE)
	for !p.curTokenIs(token.RBRACE) {
		p.statement()
		blockNode.Children = append(blockNode.Children, p.statement())
	}
	p.consume(token.RBRACE)

	return blockNode
}

func (p *Parser) statement() ast.Node {
	var statementNode ast.Node

	if p.curTokenIs(token.VAR) {
		statementNode = p.varDecl()
	} else if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.LPAREN) {
		statementNode = p.funcCall()
	} else if p.curTokenIs(token.IDENT) {
		statementNode = p.varAssignment()
	} else if p.curTokenIs(token.IF) {
		statementNode = p.ifStatement()
	} else if p.curTokenIs(token.FOR) {
		statementNode = p.forStatement()
	} else if p.curTokenIs(token.WHILE) {
		statementNode = p.whileStatement()
	} else if p.curTokenIs(token.RETURN) || p.curTokenIs(token.CONTINUE) || p.curTokenIs(token.BREAK) {
		statementNode = p.jumpStatement()
	} else {
		p.abortMsg("Expected statement.")
	}
	return statementNode
}

// varDecl = "let" ident ":" type [assignOp expr]
func (p *Parser) varDecl() {
	var varNode ast.Node
	varNode.Type = ast.VARDECL

	p.consume(token.VAR)

	var identNode ast.Node
	identNode.Type = ast.IDENT
	identNode.TokenStart = p.curToken
	p.consume(token.IDENT)
	p.consume(token.COLON)

	varNode.Children = append(varNode.Children, identNode)
	varNode.Children = append(varNode.Children, p.varType())

	if p.curTokenIs(token.ASSIGN) {
		p.consume(token.ASSIGN)
		varNode.Children = append(varNode.Children, p.expr())
		// Does not handle array literal.
	}
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

// varType = ident {"[" [expr] "]"}
func (p *Parser) varType() ast.Node {
	var typeNode ast.Node
	typeNode.Type = ast.VARTYPE

	var identNode ast.Node
	identNode.Type = ast.IDENT
	identNode.TokenStart = p.curToken
	typeNode.Children = append(typeNode.Children, identNode)
	p.consume(token.IDENT)

	for p.curTokenIs(token.LBRACKET) {
		p.nextToken()
		if !p.curTokenIs(token.RBRACKET) {
			p.expr()
		}
		p.consume(token.RBRACKET)
	}

	return typeNode
}

// varRef = ident {"[" expr "]"}
func (p *Parser) varRef() ast.Node {
	var refNode ast.Node
	refNode.Type = ast.VARREF

	var identNode ast.Node
	identNode.Type = ast.IDENT
	identNode.TokenStart = p.curToken
	refNode.Children = append(refNode.Children, identNode)
	p.consume(token.IDENT)

	for p.curTokenIs(token.LBRACKET) {
		p.nextToken()
		p.expr()
		p.consume(token.RBRACKET)
	}

	return refNode
}

// varAssignment = varRef assignOp expr
func (p *Parser) varAssignment() ast.Node {
	var assignNode ast.Node
	assignNode.Type = ast.VARASSIGN

	assignNode.Children = append(assignNode.Children, p.varRef())
	p.consume(token.ASSIGN)
	assignNode.Children = append(assignNode.Children, p.expr())

	return assignNode
}

// ifStatement = "if" expr block
func (p *Parser) ifStatement() {
	var statementNode ast.Node
	statementNode.Type = ast.IFSTATEMENT

	p.consume(token.IF)
	p.expr()
	p.block()
}

// forStatement = "for" forClause block
func (p *Parser) forStatement() {
	p.consume(token.FOR)
	p.forClause()
	p.block()
}

// forClause = [statement] ";" [expr] ";" [statement]
func (p *Parser) forClause() {
	if !p.curTokenIs(token.SEMICOLON) {
		p.statement()
	}
	p.consume(token.SEMICOLON)
	if !p.curTokenIs(token.SEMICOLON) {
		p.expr()
	}
	p.consume(token.SEMICOLON)
	if !p.curTokenIs(token.LBRACE) {
		p.statement()
	}
}

func (p *Parser) whileStatement() {
	p.consume(token.WHILE)
	p.expr()
	p.block()
}

func (p *Parser) jumpStatement() {
	p.nextToken()
}

// Expression grammar based on Crafting Interpreters
func (p *Parser) expr() {
	p.logical()
	fmt.Println("End of expr.")
}

func (p *Parser) logical() {
	p.equality()
	for p.curTokenIs(token.AND) || p.curTokenIs(token.OR) {
		p.nextToken()
		p.equality()
	}
}

func (p *Parser) equality() {
	p.comparison()
	for p.curTokenIs(token.EQ) || p.curTokenIs(token.NOTEQ) {
		p.nextToken()
		p.comparison()
	}
}

func (p *Parser) comparison() {
	p.addition()
	for p.curTokenIs(token.GT) || p.curTokenIs(token.GTEQ) || p.curTokenIs(token.LT) || p.curTokenIs(token.LTEQ) {
		p.nextToken()
		p.addition()
	}
}

func (p *Parser) addition() {
	p.multiplication()
	for p.curTokenIs(token.PLUS) || p.curTokenIs(token.MINUS) {
		p.nextToken()
		p.multiplication()
	}
}

func (p *Parser) multiplication() {
	p.unary()
	for p.curTokenIs(token.ASTERISK) || p.curTokenIs(token.SLASH) {
		p.nextToken()
		p.unary()
	}
}

func (p *Parser) unary() {
	if p.curTokenIs(token.BANG) || p.curTokenIs(token.PLUS) || p.curTokenIs(token.MINUS) {
		p.nextToken()
		p.unary()
	} else {
		p.primary()
	}
}

// primary = funcCall | varRef | INT | FLOAT | STRING | "false" | "true" | "nil" | "(" expr ")"
func (p *Parser) primary() {
	if p.curTokenIs(token.INT) || p.curTokenIs(token.FLOAT) || p.curTokenIs(token.STRING) || p.curTokenIs(token.TRUE) || p.curTokenIs(token.FALSE) || p.curTokenIs(token.NIL) {
		p.nextToken()
	} else if p.curTokenIs(token.LPAREN) {
		p.nextToken()
		p.expr()
		p.consume(token.RPAREN)
	} else if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.LPAREN) {
		p.funcCall()
	} else {
		p.varRef()
	}
}

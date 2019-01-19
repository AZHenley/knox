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
	errors    []string
	ast       ast.ASTNode
}

// New lexer.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.nextToken()
	p.nextToken() // Not sure why these are needed.
	return p
}

// Errors return stored errors
func (p *Parser) Errors() []string {
	return p.errors
}

// if peek token occurs error
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.curToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) consumeError(t token.TokenType) {
	fmt.Printf("Expected next token to be %s, got %s instead.\n", t, p.curToken.Type)
	//p.errors = append(p.errors, msg)
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

// expect next token is t
// succeed: return true and forward token
// failed: return false and store error
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) consume(t token.TokenType) bool {
	if p.curTokenIs(t) {
		p.nextToken()
		return true
	}

	p.consumeError(t)
	panic("Aborted.")
	//return false
}

// Program parses a program.
func (p *Parser) Program() {
	for !p.curTokenIs(token.EOF) {
		p.funcDecl()
	}
}

func (p *Parser) funcDecl() {
	p.consume(token.FUNCTION)
	p.consume(token.IDENT)
	p.paramList()
	p.body()
	fmt.Println("End of func decl.")
}

func (p *Parser) paramList() {
	p.consume(token.LPAREN)
	for !p.curTokenIs(token.RPAREN) {
		p.consume(token.IDENT)
		p.consume(token.COLON)
		p.consume(token.IDENT) // Replace this with varType()
		if p.curTokenIs(token.COMMA) {
			p.consume(token.COMMA)
			if p.curTokenIs(token.RPAREN) {
				p.consumeError(token.IDENT)
			}
		}
	}
	p.consume(token.RPAREN)
}

func (p *Parser) body() {
	p.consume(token.LBRACE)
	for !p.curTokenIs(token.RBRACE) {
		p.statement()
	}
	p.consume(token.RBRACE)
}

func (p *Parser) statement() {
	if p.curTokenIs(token.VAR) {
		p.varDecl()
	} else if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.LPAREN) {
		p.funcCall()
	} else if p.curTokenIs(token.IDENT) {
		p.varAssignment()
	}
}

// varDecl = "let" ident ":" type [assignOp expr]
func (p *Parser) varDecl() {
	p.consume(token.VAR)
	p.consume(token.IDENT)
	p.consume(token.COLON)
	p.consume(token.IDENT) // Replace this with varType()
	if p.curTokenIs(token.ASSIGN) {
		p.consume(token.ASSIGN)
		p.expr() // Does not handle array literal.
	}
}

// funcCall = ident argList
func (p *Parser) funcCall() {
	p.consume(token.IDENT)
	p.argList()
}

// argList = "(" {expr ","} [expr] ")"
func (p *Parser) argList() {
	p.consume(token.LPAREN)
	for !p.curTokenIs(token.RPAREN) {
		p.expr()
		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}
	p.nextToken()
}

// varRef = ident {"[" expr "]"}
func (p *Parser) varRef() {
	p.consume(token.IDENT)
	for p.curTokenIs(token.LBRACKET) {
		p.nextToken()
		p.expr()
		p.consume(token.RBRACKET)
	}
}

func (p *Parser) varAssignment() {

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

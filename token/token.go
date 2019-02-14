// Based on The Interpreter Book.

package token

// TokenType is a string
type TokenType string

// Token struct represent the lexer token
type Token struct {
	Type    TokenType
	Literal string
}

// pre-defined TokenType
const (
	ILLEGAL   = "ILLEGAL"
	EOF       = "EOF"
	IDENT     = "IDENT"
	INT       = "INT"
	FLOAT     = "FLOAT"
	ASSIGN    = "="
	PLUS      = "+"
	COMMA     = ","
	SEMICOLON = ";"
	MINUS     = "-"
	BANG      = "!"
	ASTERISK  = "*"
	SLASH     = "/"
	POWER     = "^"
	LT        = "<"
	LTEQ      = "<="
	GT        = ">"
	GTEQ      = ">="
	AND       = "&&"
	OR        = "||"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	FUNCTION  = "FUNCTION"
	VAR       = "VAR"
	TRUE      = "TRUE"
	FALSE     = "FALSE"
	IF        = "IF"
	ELSE      = "ELSE"
	RETURN    = "RETURN"
	BREAK     = "BREAK"
	CONTINUE  = "CONTINUE"
	FOR       = "FOR"
	IN        = "IN"
	WHILE     = "WHILE"
	NEW       = "NEW"
	EQ        = "=="
	NOTEQ     = "!="
	STRING    = "STRING"
	LBRACKET  = "["
	RBRACKET  = "]"
	COLON     = ":"
	NIL       = "NIL"
)

// reversed keywords
var keywords = map[string]TokenType{
	"func":     FUNCTION,
	"var":      VAR,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"break":    BREAK,
	"continue": CONTINUE,
	"for":      FOR,
	"in":       IN,
	"while":    WHILE,
	"new":      NEW,
	"nil":      NIL,
}

// LookupIdentifier used to determinate whether identifier is keyword nor not
func LookupIdentifier(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENT
}

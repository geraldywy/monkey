package token

import (
	"github.com/geraldywy/monkey/utils"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 1343456

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"
	EQ       = "=="
	NEQ      = "!="
	GTE      = ">="
	LTE      = "<="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
)

var SingleToken = map[byte]TokenType{
	'=': ASSIGN,
	';': SEMICOLON,
	'(': LPAREN,
	')': RPAREN,
	',': COMMA,
	'+': PLUS,
	'-': MINUS,
	'{': LBRACE,
	'}': RBRACE,
	'!': BANG,
	'*': ASTERISK,
	'/': SLASH,
	'<': LT,
	'>': GT,
}

var DoubleToken = map[string]TokenType{
	"==": EQ,
	"!=": NEQ,
	">=": GTE,
	"<=": LTE,
}

var reservedKeywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

func LookupTType(literal string) TokenType {
	if literal == "" {
		return EOF
	}
	if ttype, ok := reservedKeywords[literal]; ok {
		return ttype
	}

	// if identifier starts with a digit, its guaranteed to be a number (we only support integer)
	// bad variables names should already be caught before calling this function
	if utils.IsDigit(literal[0]) {
		return INT
	}

	return IDENT
}

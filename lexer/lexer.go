package lexer

import (
	"github.com/geraldywy/monkey/logger"

	"github.com/geraldywy/monkey/token"
	"github.com/geraldywy/monkey/utils"
)

type Lexer struct {
	input    string
	position int  // current position in input (points to current char
	ch       byte // current char under examination

	// metadata
	fileName string
	lineNum  int
	linePos  int // line position is 1-indexed
}

func New(input string, fileName string) *Lexer {
	l := &Lexer{input: input, fileName: fileName}
	return l
}

func (l *Lexer) peekCurrent() byte {
	if l.position == len(l.input) {
		return 0
	}

	return l.input[l.position]
}

func (l *Lexer) readChar() {
	if l.position == len(l.input) {
		l.ch = 0
		return
	}
	if l.ch == '\n' {
		l.lineNum++
		l.linePos = 0
	}
	l.ch = l.input[l.position]
	l.position++
	l.linePos++
}

func (l *Lexer) NextToken() (*token.Token, error) {
	l.eatWhitespaces()
	l.readChar()
	if tt, exist := token.SingleToken[l.ch]; exist {
		// special case for double tokens
		cand := string(l.ch) + string(l.peekCurrent())
		if dblTt, candExist := token.DoubleToken[cand]; candExist {
			l.readChar()
			return newToken(dblTt, cand), nil
		}

		return newToken(tt, string(l.ch)), nil
	} else if l.ch == 0 {
		return newToken(token.EOF, ""), nil
	}

	// handle all keywords/identifiers/numbers (really, just integers)
	if isSupportedChar(l.ch) {
		literal, err := l.readIdentLiteral()
		if err != nil {
			return nil, err
		}
		return &token.Token{
			Type:    token.LookupTType(literal),
			Literal: literal,
		}, nil
	}

	return newToken(token.ILLEGAL, string(l.ch)), nil
}

func (l *Lexer) eatWhitespaces() {
	for utils.IsWhitespace(l.peekCurrent()) {
		l.readChar()
	}
}

func (l *Lexer) readIdentLiteral() (string, error) {
	start := l.position - 1
	if utils.IsDigit(l.ch) { // is a number
		for utils.IsDigit(l.peekCurrent()) {
			l.readChar()
		}
		if utils.IsAlphaOrUnderscore(l.peekCurrent()) {
			// a variable cannot start with a number
			logger.PrettyPrintErr(l.fileName, l.lineNum, l.linePos, ErrBadVariableName)
			return "", ErrBadVariableName
		}
	} else {
		for utils.IsAlphaOrUnderscore(l.peekCurrent()) || utils.IsDigit(l.peekCurrent()) {
			l.readChar()
		}
	}

	return l.input[start:l.position], nil
}

func isSupportedChar(ch byte) bool {
	return utils.IsDigit(ch) || utils.IsAlphaOrUnderscore(ch)
}

func newToken(tokenType token.TokenType, literal string) *token.Token {
	return &token.Token{
		Type:    tokenType,
		Literal: literal,
	}
}

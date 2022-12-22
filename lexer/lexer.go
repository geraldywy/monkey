package lexer

import (
	"github.com/geraldywy/monkey/logger"

	"github.com/geraldywy/monkey/token"
	"github.com/geraldywy/monkey/utils"
)

type Lexer struct {
	input    string
	position int  // position in input to resume reading (also read as, not read in yet)
	ch       byte // prev char read in

	// metadata
	FileName string
	LineNum  int // line num is 1-indexed
	LinePos  int // line position is 1-indexed
}

func New(input string, fileName string) *Lexer {
	l := &Lexer{input: input, FileName: fileName, LineNum: 1, LinePos: 1}
	return l
}

func (l *Lexer) Debug() (string, int, string) {
	return string(l.ch), l.position, "->" + string(l.input[l.position]) + "<-"
}

func (l *Lexer) peekNext() byte {
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
		l.LineNum++
		l.LinePos = 0
	}
	l.ch = l.input[l.position]
	l.position++
	l.LinePos++
}

func (l *Lexer) byte2Token(ch byte, isPeek bool) (*token.Token, error) {
	startPos := l.position
	defer func() {
		// restore for peeks
		if isPeek {
			l.position = startPos
		}
	}()

	l.readChar()
	if tt, exist := token.SingleToken[ch]; exist {
		// special case for double tokens
		cand := string(ch) + string(l.peekNext())
		if dblTt, candExist := token.DoubleToken[cand]; candExist {
			l.readChar()
			return newToken(dblTt, cand), nil
		}

		return newToken(tt, string(ch)), nil
	} else if ch == 0 {
		return newToken(token.EOF, ""), nil
	}

	// handle all keywords/identifiers/numbers (really, just integers)
	if isSupportedChar(ch) {
		literal, err := l.readIdentLiteral()
		if err != nil {
			return nil, err
		}
		return &token.Token{
			Type:    token.LookupTType(literal),
			Literal: literal,
		}, nil
	}

	return newToken(token.ILLEGAL, string(ch)), nil
}

func (l *Lexer) NextToken() (*token.Token, error) {
	l.eatWhitespaces()
	return l.byte2Token(l.peekNext(), false)
}

func (l *Lexer) PeekToken() (*token.Token, error) {
	l.eatWhitespaces()
	return l.byte2Token(l.peekNext(), true)
}

func (l *Lexer) eatWhitespaces() {
	for utils.IsWhitespace(l.peekNext()) {
		l.readChar()
	}
}

func (l *Lexer) readIdentLiteral() (string, error) {
	start := l.position - 1
	if utils.IsDigit(l.ch) { // is a number
		for utils.IsDigit(l.peekNext()) {
			l.readChar()
		}
		if utils.IsAlphaOrUnderscore(l.peekNext()) {
			// a variable cannot start with a number
			logger.PrettyPrintErr(l.FileName, l.LineNum, l.LinePos, ErrBadVariableName)
			return "", ErrBadVariableName
		}
	} else {
		for utils.IsAlphaOrUnderscore(l.peekNext()) || utils.IsDigit(l.peekNext()) {
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

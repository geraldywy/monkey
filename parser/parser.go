package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/geraldywy/monkey/ast"
	"github.com/geraldywy/monkey/lexer"
	"github.com/geraldywy/monkey/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         //+
	PRODUCT     //*
	PREFIX      //-Xor!X
	CALL        // myFunction(X)
)

type Parser struct {
	l      *lexer.Lexer
	Errors []error

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
		prefixParseFns: map[token.TokenType]prefixParseFn{
			token.IDENT: parseIdentifier,
			token.INT:   parseIntegerLiteral,
		},
	}

	return p
}

func parseIdentifier(tkn *token.Token) (ast.Expression, error) {
	return &ast.Identifier{
		Token: tkn,
		Value: tkn.Literal,
	}, nil
}

func parseIntegerLiteral(tkn *token.Token) (ast.Expression, error) {
	exp := &ast.IntegerLiteral{
		Token: tkn,
	}
	var err error
	exp.Value, err = strconv.ParseInt(tkn.Literal, 10, 64)
	if err != nil {
		return nil, err
	}

	return exp, nil
}

func (p *Parser) ParseProgram() *ast.Program {
	prog := new(ast.Program)
	prog.Statements = make([]ast.Statement, 0)
	var tkn *token.Token
	var nxtErr error

	for tkn, nxtErr = p.l.NextToken(); tkn.Type != token.EOF; tkn, nxtErr = p.l.NextToken() {
		if nxtErr != nil {
			p.Errors = append(p.Errors, nxtErr)
			break
		}
		stmt, err := p.parseStatement(tkn)
		if err != nil {
			p.Errors = append(p.Errors, err)
			continue // redundant, but just leaving it in here for clarity
		}
		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}
	}

	return prog
}

func (p *Parser) parseStatement(startToken *token.Token) (ast.Statement, error) {
	switch startToken.Type {
	case token.LET:
		return p.parseLetStatement(startToken)
	case token.RETURN:
		return p.parseReturnStatement(startToken)
	default:
		return p.parseExpressionStatement(startToken)
	}
}

func (p *Parser) parseLetStatement(startToken *token.Token) (*ast.LetStatement, error) {
	stmt := &ast.LetStatement{
		Token: startToken,
	}
	if err := p.assertPeek(token.IDENT); err != nil {
		return nil, err
	}

	nameToken, err := p.l.NextToken()
	if err != nil {
		return nil, err
	}
	stmt.Name = &ast.Identifier{
		Token: nameToken,
		Value: nameToken.Literal,
	}

	if err := p.assertPeek(token.ASSIGN); err != nil {
		return nil, err
	}
	if _, err := p.l.NextToken(); err != nil {
		return nil, err
	}

	// TODO, skip value
	for err := p.assertPeek(token.SEMICOLON); err != nil; _, err = p.l.NextToken() {
	}
	if err := p.assertPeek(token.SEMICOLON); err != nil {
		return nil, err
	}
	p.l.NextToken() // skip semi

	return stmt, nil
}

func (p *Parser) parseReturnStatement(startToken *token.Token) (*ast.ReturnStatement, error) {
	stmt := &ast.ReturnStatement{
		Token: startToken,
	}

	// TODO, skip value
	for err := p.assertPeek(token.SEMICOLON); err != nil; _, err = p.l.NextToken() {
	}
	if err := p.assertPeek(token.SEMICOLON); err != nil {
		return nil, err
	}
	p.l.NextToken() // skip semi

	return stmt, nil
}

func (p *Parser) parseExpressionStatement(startToken *token.Token) (*ast.ExpressionStatement, error) {
	exp, err := p.parseExpression(startToken, LOWEST)
	if err != nil {
		p.Errors = append(p.Errors, err)
		return nil, err
	}
	stmt := &ast.ExpressionStatement{
		Token:      startToken,
		Expression: exp,
	}

	// advance if is semicolon
	if err := p.assertPeek(token.SEMICOLON); err == nil {
		p.l.NextToken()
	}

	return stmt, nil
}

func (p *Parser) parseExpression(startToken *token.Token, precedence int) (ast.Expression, error) {
	if prefix, exist := p.prefixParseFns[startToken.Type]; exist {
		return prefix(startToken)
	}

	return nil, nil
}

func (p *Parser) assertPeek(wantTkn token.TokenType) error {
	tkn, err := p.l.PeekToken()
	if err != nil {
		return err
	}

	if tkn.Type != wantTkn {
		return errors.New(fmt.Sprintf(
			"%s line: %d col: %d expected token: %s, got %s",
			p.l.FileName,
			p.l.LineNum,
			p.l.LinePos,
			wantTkn,
			tkn.Type,
		))
	}

	return nil
}

type (
	prefixParseFn func(tkn *token.Token) (ast.Expression, error)
	infixParseFn  func(expression ast.Expression) (ast.Expression, error)
)

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

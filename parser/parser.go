package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/geraldywy/monkey/utils"

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
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

type (
	prefixParseFn func(tkn *token.Token) (ast.Expression, error)
	infixParseFn  func(expression ast.Expression, tkn *token.Token) (ast.Expression, error)
)

type Parser struct {
	l      *lexer.Lexer
	Errors []error

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.prefixParseFns = map[token.TokenType]prefixParseFn{
		token.IDENT:      p.parseIdentifier,
		token.INT:        p.parseIntegerLiteral,
		token.TRUE:       p.parseBooleanLiteral,
		token.FALSE:      p.parseBooleanLiteral,
		token.BANG:       p.parsePrefixExpression,
		token.MINUS:      p.parsePrefixExpression,
		token.MINUSMINUS: p.parsePrefixExpression,
		token.PLUSPLUS:   p.parsePrefixExpression,
		token.LPAREN:     p.parseGroupedExpression,
		token.IF:         p.parseIfExpression,
		token.FUNCTION:   p.parseFunctionLiteral,
	}
	p.infixParseFns = map[token.TokenType]infixParseFn{
		token.PLUS:     p.parseInfixExpression,
		token.MINUS:    p.parseInfixExpression,
		token.SLASH:    p.parseInfixExpression,
		token.ASTERISK: p.parseInfixExpression,
		token.EQ:       p.parseInfixExpression,
		token.NEQ:      p.parseInfixExpression,
		token.LT:       p.parseInfixExpression,
		token.GT:       p.parseInfixExpression,
		token.LPAREN:   p.parseCallExpression,
	}

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	prog := new(ast.Program)
	prog.Statements = make([]ast.Statement, 0)
	var tkn *token.Token
	var nxtErr error

	for tkn, nxtErr = p.nextToken(); tkn.Type != token.EOF; tkn, nxtErr = p.nextToken() {
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

func (p *Parser) parseIdentifier(tkn *token.Token) (ast.Expression, error) {
	return &ast.Identifier{
		Token: tkn,
		Value: tkn.Literal,
	}, nil
}

func (p *Parser) parseIntegerLiteral(tkn *token.Token) (ast.Expression, error) {
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

func (p *Parser) parseBooleanLiteral(tkn *token.Token) (ast.Expression, error) {
	exp := &ast.Boolean{
		Token: tkn,
	}
	var err error
	exp.Value, err = strconv.ParseBool(tkn.Literal)
	if err != nil {
		return nil, err
	}

	return exp, nil
}

func (p *Parser) parseFunctionLiteral(tkn *token.Token) (ast.Expression, error) {
	fn := &ast.FunctionLiteral{Token: tkn}

	if _, err := p.assertAndAdvanceTkn(token.LPAREN); err != nil {
		return nil, err
	}

	params, err := p.parseFunctionParams()
	if err != nil {
		return nil, err
	}
	fn.Parameters = params

	lBraceTkn, err := p.assertAndAdvanceTkn(token.LBRACE)
	if err != nil {
		return nil, err
	}

	blockStmt, err := p.parseBlockStatement(lBraceTkn)
	if err != nil {
		return nil, err
	}
	fn.Body = blockStmt

	return fn, nil
}

func (p *Parser) parseFunctionParams() ([]*ast.Identifier, error) {
	idents := make([]*ast.Identifier, 0)
	// scan till rbrace
	for p.assertPeek(token.RPAREN) != nil {
		nxtToken, err := p.nextToken()
		if err != nil {
			return nil, err
		}

		idents = append(idents, &ast.Identifier{
			Token: nxtToken,
			Value: nxtToken.Literal,
		})

		// assert and skip the comma between identifiers
		if _, err := p.assertAndAdvanceTkn(token.COMMA); err != nil {
			break
		}
	}

	if _, err := p.assertAndAdvanceTkn(token.RPAREN); err != nil {
		return nil, err
	}

	return idents, nil
}

func (p *Parser) parsePrefixExpression(tkn *token.Token) (ast.Expression, error) {
	expression := &ast.PrefixExpression{
		Token:    tkn,
		Operator: tkn.Literal,
	}
	var err error

	tkn, err = p.nextToken()
	if err != nil {
		return nil, err
	}

	expression.Right, err = p.parseExpression(tkn, PREFIX)
	if err != nil {
		return nil, err
	}

	return expression, nil
}

func (p *Parser) parseGroupedExpression(_ *token.Token) (ast.Expression, error) {
	nxtToken, err := p.nextToken()
	if err != nil {
		return nil, err
	}

	exp, err := p.parseExpression(nxtToken, LOWEST)
	if err != nil {
		return nil, err
	}
	// advance past ')'
	if _, err := p.assertAndAdvanceTkn(token.RPAREN); err != nil {
		return nil, err
	}

	return exp, nil
}

func (p *Parser) parseIfExpression(tkn *token.Token) (ast.Expression, error) {
	exp := &ast.IfExpression{Token: tkn}

	if _, err := p.assertAndAdvanceTkn(token.LPAREN); err != nil {
		return nil, err
	}

	nxtTkn, err := p.nextToken()
	if err != nil {
		return nil, err
	}

	// Question: Won't this fail? parse expression will read and attempt to evaluate the ')'
	// Answer: No, it is convenient that the default precedence returned for any uncalibrated
	// token is LOWEST, thus, ')' will have a LOWEST precedence, causing the evaluation of the
	// expression to be terminated just right before the ')',
	// See the IF precedence condition in the loop in parser.parseExpression
	exp.Condition, err = p.parseExpression(nxtTkn, LOWEST)
	if err != nil {
		return nil, err
	}

	if _, err := p.assertAndAdvanceTkn(token.RPAREN); err != nil {
		return nil, err
	}

	lBraceTkn, err := p.assertAndAdvanceTkn(token.LBRACE)
	if err != nil {
		return nil, err
	}

	if exp.Consequence, err = p.parseBlockStatement(lBraceTkn); err != nil {
		return nil, err
	}

	// no ELSE to evaluate
	if _, err := p.assertAndAdvanceTkn(token.ELSE); err != nil {
		return exp, nil
	}

	lBraceTkn, err = p.assertAndAdvanceTkn(token.LBRACE)
	if err != nil {
		return nil, err
	}
	if exp.Alternative, err = p.parseBlockStatement(lBraceTkn); err != nil {
		return nil, err
	}

	return exp, nil
}

func (p *Parser) parseBlockStatement(tkn *token.Token) (*ast.BlockStatement, error) {
	block := &ast.BlockStatement{
		Token:      tkn,
		Statements: make([]ast.Statement, 0),
	}

	for p.assertPeek(token.RBRACE, token.EOF) != nil {
		nxtToken, err := p.nextToken()
		if err != nil {
			return nil, err
		}
		stmt, err := p.parseStatement(nxtToken)
		if err != nil {
			return nil, err
		}
		block.Statements = append(block.Statements, stmt)
	}
	// advance past '}'
	if _, err := p.assertAndAdvanceTkn(token.RBRACE); err != nil {
		return nil, err
	}

	return block, nil
}

func (p *Parser) parseInfixExpression(left ast.Expression, tkn *token.Token) (ast.Expression, error) {
	expression := &ast.InfixExpression{
		Token:    tkn,
		Left:     left,
		Operator: tkn.Literal,
	}

	pred := p.getPrecedence(tkn)
	nxtToken, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	exp, err := p.parseExpression(nxtToken, pred)
	if err != nil {
		return nil, err
	}
	expression.Right = exp

	return expression, nil
}

func (p *Parser) parseCallExpression(fn ast.Expression, tkn *token.Token) (ast.Expression, error) {
	exp := &ast.CallExpression{
		Token:    tkn,
		Function: fn,
	}
	arg, err := p.parseCallArgs()
	if err != nil {
		return nil, err
	}
	exp.Arguments = arg

	return exp, nil
}

func (p *Parser) parseCallArgs() ([]ast.Expression, error) {
	args := make([]ast.Expression, 0)

	// scan till rbrace
	for p.assertPeek(token.RPAREN) != nil {
		nxtToken, err := p.nextToken()
		if err != nil {
			return nil, err
		}

		// Question: Why does this work? Wouldn't parseExpression consume the ',' separators
		// as part of the expressions? Why don't we need to manually intervene?
		// Answer: Similar to the previous question I had, the precedence of an uninitialized
		// token ',' is LOWEST, causing the loop in parseExpression to terminate early
		// as the precedence of LOWEST is always in the worst case equal to the caller.
		// This way, the expression will always be evaluated up till the ',' OR ')' OR EOF token.
		arg, err := p.parseExpression(nxtToken, LOWEST)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)

		// assert and skip the comma between expressions
		if _, err := p.assertAndAdvanceTkn(token.COMMA); err != nil {
			break
		}
	}

	if _, err := p.assertAndAdvanceTkn(token.RPAREN); err != nil {
		return nil, err
	}

	return args, nil
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

	nameToken, err := p.assertAndAdvanceTkn(token.IDENT)
	if err != nil {
		return nil, err
	}
	stmt.Name = &ast.Identifier{
		Token: nameToken,
		Value: nameToken.Literal,
	}

	if _, err := p.assertAndAdvanceTkn(token.ASSIGN); err != nil {
		return nil, err
	}

	nxtTkn, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	exp, err := p.parseExpression(nxtTkn, LOWEST)
	stmt.Value = exp

	// assert is semicolon
	if _, err := p.assertAndAdvanceTkn(token.SEMICOLON); err != nil {
		return nil, err
	}

	return stmt, nil
}

func (p *Parser) parseReturnStatement(startToken *token.Token) (*ast.ReturnStatement, error) {
	stmt := &ast.ReturnStatement{
		Token: startToken,
	}

	nxtTkn, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	rv, err := p.parseExpression(nxtTkn, LOWEST)
	if err != nil {
		return nil, err
	}
	stmt.ReturnValue = rv

	// assert is semicolon
	if _, err := p.assertAndAdvanceTkn(token.SEMICOLON); err != nil {
		return nil, err
	}

	return stmt, nil
}

func (p *Parser) parseExpressionStatement(startToken *token.Token) (*ast.ExpressionStatement, error) {
	exp, err := p.parseExpression(startToken, LOWEST)
	if err != nil {
		return nil, err
	}
	stmt := &ast.ExpressionStatement{
		Token:      startToken,
		Expression: exp,
	}

	// advance if is semicolon, intentionally ignoring the assertion error
	p.assertAndAdvanceTkn(token.SEMICOLON)

	return stmt, nil
}

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

func (p *Parser) getPrecedence(tkn *token.Token) int {
	if p, exist := precedences[tkn.Type]; exist {
		return p
	}

	return LOWEST
}

func (p *Parser) parseExpression(startToken *token.Token, precedence int) (ast.Expression, error) {
	prefix, exist := p.prefixParseFns[startToken.Type]
	if !exist {
		return nil, errors.New(fmt.Sprintf(
			"no prefix parse function for %s found",
			startToken.Literal,
		))
	}

	leftExp, err := prefix(startToken)
	if err != nil {
		return nil, err
	}

	for p.assertPeek(token.SEMICOLON, token.EOF) != nil { // peek until next is semicolon or EOF
		nxtToken, err := p.peekToken()
		if err != nil {
			return nil, err
		}
		// current precedence is higher or equal, should be evaluated by itself first
		if precedence >= p.getPrecedence(nxtToken) {
			break
		}

		nxtToken, err = p.nextToken()
		if err != nil {
			return nil, err
		}
		infix, exist := p.infixParseFns[nxtToken.Type]
		if !exist {
			return nil, errors.New(fmt.Sprintf(
				"no infix parse function for %s found",
				nxtToken.Literal,
			))
		}
		leftExp, err = infix(leftExp, nxtToken)
		if err != nil {
			return nil, err
		}
	}

	return leftExp, nil
}

func (p *Parser) nextToken() (*token.Token, error) {
	return p.l.NextToken()
}

func (p *Parser) peekToken() (*token.Token, error) {
	return p.l.PeekToken()
}

func (p *Parser) assertPeek(wantTkns ...token.TokenType) error {
	tkn, err := p.l.PeekToken()
	if err != nil {
		return err
	}

	if !utils.Contains(wantTkns, tkn.Type) {
		return errors.New(fmt.Sprintf(
			"%s line: %d col: %d expected one of tokens: %s, got %s",
			p.l.FileName,
			p.l.LineNum,
			p.l.LinePos,
			wantTkns,
			tkn.Type,
		))
	}

	return nil
}

func (p *Parser) assertAndAdvanceTkn(wantTkns ...token.TokenType) (*token.Token, error) {
	if err := p.assertPeek(wantTkns...); err != nil {
		return nil, err
	}

	return p.nextToken()
}

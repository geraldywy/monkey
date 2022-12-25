package lexer_test

import (
	"testing"

	"github.com/geraldywy/monkey/lexer"

	"github.com/geraldywy/monkey/token"
)

type ts struct {
	name  string
	in    string
	wants []tsWants
}

type tsWants struct {
	wantType    token.TokenType
	wantLiteral string
	wantErr     error
}

func TestNextToken(t *testing.T) {
	tests := []ts{
		{
			name: "basic symbols",
			in:   "=+(){},;",
			wants: []tsWants{
				{token.ASSIGN, "=", nil},
				{token.PLUS, "+", nil},
				{token.LPAREN, "(", nil},
				{token.RPAREN, ")", nil},
				{token.LBRACE, "{", nil},
				{token.RBRACE, "}", nil},
				{token.COMMA, ",", nil},
				{token.SEMICOLON, ";", nil},
				{token.EOF, "", nil},
			},
		},
		{
			name: "simple statements",
			in: `let five = 5;
				let ten = 10;
				   let add = fn(x, y) {
					 x + y;
				};
				   let result = add(five, ten);
				   `,
			wants: []tsWants{
				{token.LET, "let", nil},
				{token.IDENT, "five", nil},
				{token.ASSIGN, "=", nil},
				{token.INT, "5", nil},
				{token.SEMICOLON, ";", nil},
				{token.LET, "let", nil},
				{token.IDENT, "ten", nil},
				{token.ASSIGN, "=", nil},
				{token.INT, "10", nil},
				{token.SEMICOLON, ";", nil},
				{token.LET, "let", nil},
				{token.IDENT, "add", nil},
				{token.ASSIGN, "=", nil},
				{token.FUNCTION, "fn", nil},
				{token.LPAREN, "(", nil},
				{token.IDENT, "x", nil},
				{token.COMMA, ",", nil},
				{token.IDENT, "y", nil},
				{token.RPAREN, ")", nil},
				{token.LBRACE, "{", nil},
				{token.IDENT, "x", nil},
				{token.PLUS, "+", nil},
				{token.IDENT, "y", nil},
				{token.SEMICOLON, ";", nil},
				{token.RBRACE, "}", nil},
				{token.SEMICOLON, ";", nil},
				{token.LET, "let", nil},
				{token.IDENT, "result", nil},
				{token.ASSIGN, "=", nil},
				{token.IDENT, "add", nil},
				{token.LPAREN, "(", nil},
				{token.IDENT, "five", nil},
				{token.COMMA, ",", nil},
				{token.IDENT, "ten", nil},
				{token.RPAREN, ")", nil},
				{token.SEMICOLON, ";", nil},
				{token.EOF, "", nil},
			},
		},
		{
			name: "extended lexer set",
			in: `let five = 5;
				let ten = 10;
				   let add = fn(x, y) {
					 x + y;
				};
				   let result = add(five, ten);
				   !-/*5;
				   5 < 10 > 5;
				   if (5 < 10) {
					   return true;
				   } else {
					   return false;
				}
				10 == 10;
				10 != 9;
				4 >= 2;
				1 <= 2;
				++2;
				--5;
				`,
			wants: []tsWants{
				{token.LET, "let", nil},
				{token.IDENT, "five", nil},
				{token.ASSIGN, "=", nil},
				{token.INT, "5", nil},
				{token.SEMICOLON, ";", nil},
				{token.LET, "let", nil},
				{token.IDENT, "ten", nil},
				{token.ASSIGN, "=", nil},
				{token.INT, "10", nil},
				{token.SEMICOLON, ";", nil},
				{token.LET, "let", nil},
				{token.IDENT, "add", nil},
				{token.ASSIGN, "=", nil},
				{token.FUNCTION, "fn", nil},
				{token.LPAREN, "(", nil},
				{token.IDENT, "x", nil},
				{token.COMMA, ",", nil},
				{token.IDENT, "y", nil},
				{token.RPAREN, ")", nil},
				{token.LBRACE, "{", nil},
				{token.IDENT, "x", nil},
				{token.PLUS, "+", nil},
				{token.IDENT, "y", nil},
				{token.SEMICOLON, ";", nil},
				{token.RBRACE, "}", nil},
				{token.SEMICOLON, ";", nil},
				{token.LET, "let", nil},
				{token.IDENT, "result", nil},
				{token.ASSIGN, "=", nil},
				{token.IDENT, "add", nil},
				{token.LPAREN, "(", nil},
				{token.IDENT, "five", nil},
				{token.COMMA, ",", nil},
				{token.IDENT, "ten", nil},
				{token.RPAREN, ")", nil},
				{token.SEMICOLON, ";", nil},
				{token.BANG, "!", nil},
				{token.MINUS, "-", nil},
				{token.SLASH, "/", nil},
				{token.ASTERISK, "*", nil},
				{token.INT, "5", nil},
				{token.SEMICOLON, ";", nil},
				{token.INT, "5", nil},
				{token.LT, "<", nil},
				{token.INT, "10", nil},
				{token.GT, ">", nil},
				{token.INT, "5", nil},
				{token.SEMICOLON, ";", nil},
				{token.IF, "if", nil},
				{token.LPAREN, "(", nil},
				{token.INT, "5", nil},
				{token.LT, "<", nil},
				{token.INT, "10", nil},
				{token.RPAREN, ")", nil},
				{token.LBRACE, "{", nil},
				{token.RETURN, "return", nil},
				{token.TRUE, "true", nil},
				{token.SEMICOLON, ";", nil},
				{token.RBRACE, "}", nil},
				{token.ELSE, "else", nil},
				{token.LBRACE, "{", nil},
				{token.RETURN, "return", nil},
				{token.FALSE, "false", nil},
				{token.SEMICOLON, ";", nil},
				{token.RBRACE, "}", nil},
				{token.INT, "10", nil},
				{token.EQ, "==", nil},
				{token.INT, "10", nil},
				{token.SEMICOLON, ";", nil},
				{token.INT, "10", nil},
				{token.NEQ, "!=", nil},
				{token.INT, "9", nil},
				{token.SEMICOLON, ";", nil},
				{token.INT, "4", nil},
				{token.GTE, ">=", nil},
				{token.INT, "2", nil},
				{token.SEMICOLON, ";", nil},
				{token.INT, "1", nil},
				{token.LTE, "<=", nil},
				{token.INT, "2", nil},
				{token.SEMICOLON, ";", nil},
				{token.PLUSPLUS, "++", nil},
				{token.INT, "2", nil},
				{token.SEMICOLON, ";", nil},
				{token.MINUSMINUS, "--", nil},
				{token.INT, "5", nil},
				{token.SEMICOLON, ";", nil},
				{token.EOF, "", nil},
			},
		},
	}

	for _, ts := range tests {
		l := lexer.New(ts.in, "lexer_test.go")

		for i, tw := range ts.wants {
			tok, err := l.NextToken()
			if err != tw.wantErr {
				t.Fatalf("test name: %s, tests[%d] - next token err mismatch. expected=%q, got=%q",
					ts.name, i, tw.wantErr, err)
			}

			if tok.Type != tw.wantType {
				t.Fatalf("test name: %s, tests[%d] - tokentype wrong. expected=%q, got=%q",
					ts.name, i, tw.wantType, tok.Type)
			}

			if tok.Literal != tw.wantLiteral {
				t.Fatalf("test name: %s, tests[%d] - literal wrong. expected=%q, got=%q",
					ts.name, i, tw.wantLiteral, tok.Literal)
			}
		}
	}
}

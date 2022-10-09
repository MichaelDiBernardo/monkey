package lexer

import (
	"testing"

	"github.com/MichaelDiBernardo/monkey/token"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
	}

	lexer := New(input)

	for i, ttest := range tests {
		tok := lexer.NextToken()

		if tok.Type != ttest.expectedType {
			t.Fatalf("tests[%d] - wrong tokentype. expected=%q, got=%q", i, ttest.expectedType, tok.Type)
		}

		if tok.Literal != ttest.expectedLiteral {
			t.Fatalf("tests[%d] - wrong literal. expected=%q, got=%q", i, ttest.expectedLiteral, tok.Literal)
		}
	}
}

package lexer

import (
	"testing"

	"github.com/MichaelDiBernardo/monkey/token"
)

type expectedToken struct {
	expectedType    token.TokenType
	expectedLiteral string
}

func TestNextTokenWithSingleCharTokens(t *testing.T) {
	input := `=+(){},;`

	tests := []expectedToken{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
	}
	compareExpectedTokens(t, input, tests)
}

func TestNextTokenWithExampleProgram(t *testing.T) {
	program := `let five = 5;
let ten =  10;

let add = fn(x, y) {
	x + y;
};

let result = add(five, ten);
`
	expectedTokens := []expectedToken{
		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	compareExpectedTokens(t, program, expectedTokens)
}

func compareExpectedTokens(t *testing.T, input string, expectedTokens []expectedToken) {
	lexer := New(input)

	for i, ttest := range expectedTokens {
		tok := lexer.NextToken()

		if tok.Type != ttest.expectedType {
			t.Fatalf("tests[%d] - wrong tokentype. expected=%q, got=%q", i, ttest.expectedType, tok.Type)
		}

		if tok.Literal != ttest.expectedLiteral {
			t.Fatalf("tests[%d] - wrong literal. expected=%q, got=%q", i, ttest.expectedLiteral, tok.Literal)
		}
	}
}

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
	input := `=+(){},;-!*/<>`

	tests := []expectedToken{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.MINUS, "-"},
		{token.BANG, "!"},
		{token.ASTERISK, "*"},
		{token.RSLASH, "/"},
		{token.LANGLE, "<"},
		{token.RANGLE, ">"},
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
!-/*5;
5 < 10 > 5;

if (5 < 10) {
    return true;
} else {
    return false;
}

10 == 10;
10 != 9;
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
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.RSLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LANGLE, "<"},
		{token.INT, "10"},
		{token.RANGLE, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LANGLE, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NEQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.EOF, string(NUL)},
	}
	compareExpectedTokens(t, program, expectedTokens)
}

func TestTokenLocations(t *testing.T) {
	program := `let one = 1;




let two = 2;


two == one;
two != one;
`

	lexer := NewFromString(program)

	expectedLocations := []token.Location{
		{Path: token.NO_FILEPATH, LineN: 1, CharN: 1},
		{Path: token.NO_FILEPATH, LineN: 1, CharN: 5},
		{Path: token.NO_FILEPATH, LineN: 1, CharN: 9},
		{Path: token.NO_FILEPATH, LineN: 1, CharN: 11},
		{Path: token.NO_FILEPATH, LineN: 1, CharN: 12},
		{Path: token.NO_FILEPATH, LineN: 6, CharN: 1},
		{Path: token.NO_FILEPATH, LineN: 6, CharN: 5},
		{Path: token.NO_FILEPATH, LineN: 6, CharN: 9},
		{Path: token.NO_FILEPATH, LineN: 6, CharN: 11},
		{Path: token.NO_FILEPATH, LineN: 6, CharN: 12},
		{Path: token.NO_FILEPATH, LineN: 9, CharN: 1},
		{Path: token.NO_FILEPATH, LineN: 9, CharN: 5},
		{Path: token.NO_FILEPATH, LineN: 9, CharN: 8},
		{Path: token.NO_FILEPATH, LineN: 9, CharN: 11},
		{Path: token.NO_FILEPATH, LineN: 10, CharN: 1},
		{Path: token.NO_FILEPATH, LineN: 10, CharN: 5},
		{Path: token.NO_FILEPATH, LineN: 10, CharN: 8},
		{Path: token.NO_FILEPATH, LineN: 10, CharN: 11},
	}

	for i, expected := range expectedLocations {
		tok := lexer.NextToken()
		if tok.Location != expected {
			t.Fatalf("tests[%d] - wrong location. expected=%q, got=%q", i, expected, tok.Location)
		}
	}

}

func compareExpectedTokens(t *testing.T, input string, expectedTokens []expectedToken) {
	lexer := NewFromString(input)

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

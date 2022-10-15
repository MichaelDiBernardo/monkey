package parser

import (
	"testing"

	"github.com/MichaelDiBernardo/monkey/ast"
	"github.com/MichaelDiBernardo/monkey/lexer"
	"github.com/MichaelDiBernardo/monkey/token"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`

	l := lexer.NewFromString(input)
	p := New(l)

	program := p.ParseProgram()

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if nstmts := len(program.Statements); nstmts != 3 {
		t.Fatalf("expected: len(program.Statements) = 3, got %d", nstmts)
	}

	expected := []string{
		"x",
		"y",
		"foobar",
	}

	for i, e := range expected {
		stmt := program.Statements[i]
		if !testLetStatement(t, e, stmt) {
			return
		}
	}
}

func testLetStatement(t *testing.T, expectedLiteral string, stmt ast.Statement) bool {
	if stype := stmt.Token().Type; stype != token.LET {
		t.Errorf("expected let literal for stmt, got %q", stype)
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt not *ast.LetStatement, got %T", stmt)
		return false
	}

	if v := letStmt.Name.Value; v != expectedLiteral {
		t.Errorf("expected stmt identifier %s, got %s", expectedLiteral, v)
		return false
	}

	if v := letStmt.Name.Token().Literal; v != expectedLiteral {
		t.Errorf("expected stmt identifier %s, got %s", expectedLiteral, v)
		return false
	}

	if tt := letStmt.Name.Token().Type; tt != token.IDENTIFIER {
		t.Errorf("expected stmt identifier to be type %s, got %s", token.IDENTIFIER, tt)
		return false
	}

	return true
}

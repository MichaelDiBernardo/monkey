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
	failIfParserHasErrors(t, p)

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

func TestReturnStatements(t *testing.T) {
	input := `
return 6;
return 11;
`

	l := lexer.NewFromString(input)
	p := New(l)

	program := p.ParseProgram()
	failIfParserHasErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if nstmts := len(program.Statements); nstmts != 2 {
		t.Fatalf("expected: len(program.Statements) = 2, got %d", nstmts)
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt: expected *ast.ReturnStatement, got %T", stmt)
			continue
		}
		if tok := returnStmt.Token(); !tok.Is(token.RETURN) {
			t.Errorf("stmt: expected token to be RETURN, got %q", tok.Type)
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.NewFromString(input)
	p := New(l)

	program := p.ParseProgram()
	failIfParserHasErrors(t, p)

	if nstmts := len(program.Statements); nstmts != 1 {
		t.Fatalf("expected 1 statement, got %d", nstmts)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("stmt was bad type %T", program.Statements[0])
	}

	ident, ok := stmt.Value.(*ast.Identifier)

	if !ok {
		t.Fatalf("ident was bad type %T", stmt.Value)
	}

	if v := ident.Value; v != "foobar" {
		t.Errorf("ident.Value: expected foobar, got %s", v)
	}

	if v := ident.Token().Literal; v != "foobar" {
		t.Errorf("ident.Token().Literal: expected foobar, got %s", v)
	}
}

func TestIntegerExpression(t *testing.T) {
	input := "52;"

	l := lexer.NewFromString(input)
	p := New(l)

	program := p.ParseProgram()
	failIfParserHasErrors(t, p)

	if nstmts := len(program.Statements); nstmts != 1 {
		t.Fatalf("expected 1 statement, got %d", nstmts)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("stmt was bad type %T", program.Statements[0])
	}

	intl, ok := stmt.Value.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("expected *ast.IntegerLiteral, got type %T", stmt.Value)
	}

	if v := intl.Value; v != 52 {
		t.Errorf("intl.Value: expected 52, got %d", v)
	}

	if v := intl.Token().Literal; v != "52" {
		t.Errorf("intl.Token().Literal: expected %q, got %q", "52", v)
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

func failIfParserHasErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, perr := range errors {
		t.Errorf("parser error: %s", perr.String())
	}
	t.FailNow()
}

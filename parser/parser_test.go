package parser

import (
	"fmt"
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
	program := checkParseProgram(t, input, 3)

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
	program := checkParseProgram(t, input, 2)

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
	program := checkParseProgram(t, input, 1)

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
	program := checkParseProgram(t, input, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("stmt was bad type %T", program.Statements[0])
	}

	testIntegerLiteral(t, stmt.Value, 52)
}

func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    int64
	}{
		{"!43", "!", 43},
		{"-92", "-", 92},
	}

	for i, pt := range prefixTests {
		program := checkParseProgram(t, pt.input, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("[%d]: expected ast.ExpressionStatement, got %T", i, program.Statements[0])
		}

		exp, ok := stmt.Value.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("[%d]: expected Value type *ast.PrefixExpression, got %T", i, stmt.Value)
		}

		if exp.Operator != pt.operator {
			t.Fatalf("[%d]: expected operator %q, got %q", i, pt.operator, exp.Operator)
		}

		testIntegerLiteral(t, exp.RHS, pt.value)
	}
}

func checkParseProgram(t *testing.T, input string, expectednstmts int) *ast.Program {
	l := lexer.NewFromString(input)
	p := New(l)

	program := p.ParseProgram()
	failIfParserHasErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if nstmts := len(program.Statements); nstmts != expectednstmts {
		t.Fatalf("expected %d statement(s), got %d", expectednstmts, nstmts)
	}
	return program
}

func testIntegerLiteral(t *testing.T, il ast.Expression, expected int64) {
	intl, ok := il.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("expected *ast.IntegerLiteral, got type %T", il)
	}

	if v := intl.Value; v != expected {
		t.Errorf("intl.Value: expected %d, got %d", expected, v)
	}

	if exp, act := fmt.Sprintf("%d", expected), intl.Token().Literal; exp != act {
		t.Errorf("intl.Token().Literal: expected %q, got %q", exp, act)
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

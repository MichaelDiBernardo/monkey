package parser

import (
	"fmt"
	"testing"

	"github.com/MichaelDiBernardo/monkey/ast"
	"github.com/MichaelDiBernardo/monkey/lexer"
	"github.com/MichaelDiBernardo/monkey/token"
)

func TestLetStatements(t *testing.T) {
	expected := []struct {
		input      string
		identifier string
		assigned   interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = 10", "y", 10},
		{"let foobar = 838383;", "foobar", 838383},
		{"let z = false", "z", false},
		{"let w = true;", "w", true},
	}

	for _, e := range expected {
		program := checkParseProgram(t, e.input, 1)
		stmt := program.Statements[0]
		testLetStatement(t, stmt, e.identifier, e.assigned)
	}
}

func TestReturnStatements(t *testing.T) {
	expected := []struct {
		input     string
		returnval interface{}
	}{
		{"return 6;", 6},
		{"return 11", 11},
		{"return true", true},
		{"return false;", false},
	}

	for _, e := range expected {
		program := checkParseProgram(t, e.input, 1)
		stmt := program.Statements[0]

		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt: expected *ast.ReturnStatement, got %T", stmt)
			continue
		}
		if tok := returnStmt.Token(); !tok.Is(token.RETURN) {
			t.Errorf("stmt: expected token to be RETURN, got %q", tok.Type)
		}

		testLiteralExpression(t, returnStmt.Value, e.returnval)
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	program := checkParseProgram(t, input, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("stmt was bad type %T", program.Statements[0])
	}

	testIdentifier(t, stmt.Value, "foobar")
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

func TestBooleanExpression(t *testing.T) {
	input := "true;"
	program := checkParseProgram(t, input, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("stmt was bad type %T", program.Statements[0])
	}

	testBooleanLiteral(t, stmt.Value, true)
}

func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!43", "!", 43},
		{"-92", "-", 92},
		{"!foobar", "!", "foobar"},
		{"-foobar", "-", "foobar"},
		{"!true", "!", true},
		{"!false", "!", false},
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

		testLiteralExpression(t, exp.RHS, pt.value)
	}
}

func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input    string
		LHS      interface{}
		operator string
		RHS      interface{}
	}{
		{"12 + 13;", 12, "+", 13},
		{"12 - 13;", 12, "-", 13},
		{"12 * 13;", 12, "*", 13},
		{"12 / 13;", 12, "/", 13},
		{"12 > 13;", 12, ">", 13},
		{"12 < 13;", 12, "<", 13},
		{"12 == 13;", 12, "==", 13},
		{"12 != 13;", 12, "!=", 13},
		{"x + 13;", "x", "+", 13},
		{"y - 13;", "y", "-", 13},
		{"12 * z;", 12, "*", "z"},
		{"12 / w;", 12, "/", "w"},
		{"foobar > 13;", "foobar", ">", 13},
		{"12 < x;", 12, "<", "x"},
		{"x == y;", "x", "==", "y"},
		{"foobar != x;", "foobar", "!=", "x"},
		{"x == true;", "x", "==", true},
		{"false < z;", false, "<", "z"},
	}

	for i, it := range infixTests {
		program := checkParseProgram(t, it.input, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("[%d]: expected ast.ExpressionStatement, got %T", i, program.Statements[0])
		}

		exp, ok := stmt.Value.(*ast.InfixExpression)

		if !ok {
			t.Fatalf("[%d]: expected Value type *ast.InfixExpression, got %T", i, stmt.Value)
		}

		if exp.Operator != it.operator {
			t.Fatalf("[%d]: expected operator %q, got %q", i, it.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.LHS, it.LHS) {
			t.Errorf("(failed test [%d]", i)
		}
		if !testLiteralExpression(t, exp.RHS, it.RHS) {
			t.Errorf("(failed test [%d]", i)
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"!!!!true == true",
			"((!(!(!(!true)))) == true)",
		},
	}

	for i, ot := range tests {
		program := checkParseProgram(t, ot.input, -1)

		actual := program.String()

		if actual != ot.expected {
			t.Errorf("[%d] expected: %q, got %q", i, ot.expected, actual)
		}
	}
}

/* test helpers start here */

func checkParseProgram(t *testing.T, input string, expectednstmts int) *ast.Program {
	l := lexer.NewFromString(input)
	p := New(l)

	program := p.ParseProgram()
	failIfParserHasErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if nstmts := len(program.Statements); expectednstmts > 0 && nstmts != expectednstmts {
		t.Fatalf("expected %d statement(s), got %d", expectednstmts, nstmts)
	}
	return program
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

func testLetStatement(t *testing.T, stmt ast.Statement, expectedLiteral string, expectedAssignment interface{}) {
	if stype := stmt.Token().Type; stype != token.LET {
		t.Fatalf("expected let literal for stmt, got %q", stype)
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Fatalf("stmt not *ast.LetStatement, got %T", stmt)
	}

	testIdentifier(t, letStmt.Name, expectedLiteral)
	testLiteralExpression(t, letStmt.Value, expectedAssignment)
}

func testIntegerLiteral(t *testing.T, il ast.Expression, expected int64) bool {
	intl, ok := il.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("expected *ast.IntegerLiteral, got type %T", il)
		return false
	}

	if v := intl.Value; v != expected {
		t.Errorf("intl.Value: expected %d, got %d", expected, v)
		return false
	}

	if exp, act := fmt.Sprintf("%d", expected), intl.Token().Literal; exp != act {
		t.Errorf("intl.Token().Literal: expected %q, got %q", exp, act)
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, expected bool) bool {
	bl, ok := exp.(*ast.BooleanLiteral)

	if !ok {
		t.Errorf("expected *ast.BooleanLiteral, got type %T", bl)
		return false
	}

	if v := bl.Value; v != expected {
		t.Errorf("bl.Value: expected %t, got %t", expected, v)
		return false
	}

	if exp, act := fmt.Sprintf("%t", expected), bl.Token().Literal; exp != act {
		t.Errorf("bl.Token().Literal: expected %q, got %q", exp, act)
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)

	if !ok {
		t.Errorf("expected *ast.Identifier, got %T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("expected literal %s, got %s", value, ident.Value)
		return false
	}

	if act := ident.IdentToken.Literal; value != act {
		t.Errorf("expected literal %s, got %s", value, act)
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch val := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(val))
	case int64:
		return testIntegerLiteral(t, exp, val)
	case string:
		return testIdentifier(t, exp, val)
	case bool:
		return testBooleanLiteral(t, exp, val)
	}
	t.Errorf("expression type %T not handled", expected)
	return false
}

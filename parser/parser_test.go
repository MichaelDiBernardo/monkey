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

func TestIfExpression(t *testing.T) {
	input := "if (x < y) { x }"
	program := checkParseProgram(t, input, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("stmt was bad type %T", program.Statements[0])
	}

	ifexp, ok := stmt.Value.(*ast.IfExpression)

	if !ok {
		t.Fatalf("stmt exp was bad type %T", stmt.Value)
	}

	if exp, act := token.IF, ifexp.Token().Type; exp != act {
		t.Errorf("Expected if expression token to be %q, got %q", exp, act)
	}

	if !testInfixExpression(t, ifexp.Condition, "x", "<", "y") {
		return
	}

	if nstmts := len(ifexp.Consequence.Statements); nstmts != 1 {
		t.Fatalf("expected 1 statement in consequence, got %d", nstmts)
	}

	cexp, ok := ifexp.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("consequence stmt was bad type %T", ifexp.Consequence.Statements[0])
	}

	if !testIdentifier(t, cexp.Value, "x") {
		return
	}

	if ifexp.Alternative != nil {
		t.Errorf("expected alternative to be nil, got %v", ifexp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := "if (x < y) { x } else { y }"
	program := checkParseProgram(t, input, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("stmt was bad type %T", program.Statements[0])
	}

	ifexp, ok := stmt.Value.(*ast.IfExpression)

	if !ok {
		t.Fatalf("stmt exp was bad type %T", stmt.Value)
	}

	if exp, act := token.IF, ifexp.Token().Type; exp != act {
		t.Errorf("Expected if expression token to be %q, got %q", exp, act)
	}

	if !testInfixExpression(t, ifexp.Condition, "x", "<", "y") {
		return
	}

	if nstmts := len(ifexp.Consequence.Statements); nstmts != 1 {
		t.Fatalf("expected 1 statement in consequence, got %d", nstmts)
	}

	cexp, ok := ifexp.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("consequence stmt was bad type %T", ifexp.Consequence.Statements[0])
	}

	if !testIdentifier(t, cexp.Value, "x") {
		return
	}

	if nstmts := len(ifexp.Alternative.Statements); nstmts != 1 {
		t.Fatalf("expected 1 statement in alternative, got %d", nstmts)
	}

	aexp, ok := ifexp.Alternative.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("alternative stmt was bad type %T", ifexp.Alternative.Statements[0])
	}

	if !testIdentifier(t, aexp.Value, "y") {
		return
	}

}

func TestParseFunctionLiteral(t *testing.T) {
	input := "fn(x, y) { x + y;}"
	program := checkParseProgram(t, input, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("stmt was bad type %T", program.Statements[0])
	}

	fl, ok := stmt.Value.(*ast.FunctionLiteral)

	if !ok {
		t.Fatalf("expected *ast.FunctionLiteral, got %T", stmt.Value)
	}

	if fl.FnToken.Type != token.FUNCTION {
		t.Errorf("expected FnToken type %s, got %s", token.FUNCTION, fl.FnToken.Type)
	}

	if exp, act := 2, len(fl.Parameters); exp != act {
		t.Errorf("expected %d params, got %d", exp, act)
	}

	testIdentifier(t, fl.Parameters[0], "x")
	testIdentifier(t, fl.Parameters[1], "y")

	if exp, act := 1, len(fl.Body.Statements); exp != act {
		t.Errorf("expected %d statements in func body, got %d", exp, act)
	}

	bodystmt, ok := fl.Body.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("bodystmt was bad type %T", fl.Body.Statements[0])
	}

	testInfixExpression(t, bodystmt.Value, "x", "+", "y")
}

func TestParseCallExpression(t *testing.T) {
	input := "sum(1, 2 * 3, 4 + 5);"
	program := checkParseProgram(t, input, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("stmt was bad type %T", program.Statements[0])
	}

	callexp, ok := stmt.Value.(*ast.CallExpression)

	if !ok {
		t.Fatalf("stmt was bad type %T", stmt.Value)

	}

	if exp, act := token.LPAREN, callexp.Token().Type; exp != act {
		t.Errorf("expected token type %s, got %s", exp, act)
	}

	testIdentifier(t, callexp.Function, "sum")

	if exp, act := 3, len(callexp.Arguments); exp != act {
		t.Fatalf("expected %d argument expressions, got %d", exp, act)
	}

	testLiteralExpression(t, callexp.Arguments[0], 1)
	testInfixExpression(t, callexp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, callexp.Arguments[2], 4, "+", 5)
}

func TestParseCallExpressionArguments(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		program := checkParseProgram(t, tt.input, 1)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		exp, ok := stmt.Value.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
				stmt.Value)
		}

		if !testIdentifier(t, exp.Function, tt.expectedIdent) {
			return
		}

		if len(exp.Arguments) != len(tt.expectedArgs) {
			t.Fatalf("wrong number of arguments. want=%d, got=%d",
				len(tt.expectedArgs), len(exp.Arguments))
		}

		for i, arg := range tt.expectedArgs {
			if exp.Arguments[i].String() != arg {
				t.Errorf("argument %d wrong. want=%q, got=%q", i,
					arg, exp.Arguments[i].String())
			}
		}
	}
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
		lhs      interface{}
		operator string
		rhs      interface{}
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

		if !testInfixExpression(t, stmt.Value, it.lhs, it.operator, it.rhs) {
			t.Errorf("[%d]: testInfixExpression failed.", i)
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
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
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

func testInfixExpression(t *testing.T, exp ast.Expression, lhs interface{}, operator string, rhs interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.LHS, lhs) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.RHS, rhs) {
		return false
	}
	return true
}

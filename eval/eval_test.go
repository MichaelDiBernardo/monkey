package eval

import (
	"testing"

	"github.com/MichaelDiBernardo/monkey/lexer"
	"github.com/MichaelDiBernardo/monkey/object"
	"github.com/MichaelDiBernardo/monkey/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"98", 98},
	}

	for i, tt := range tests {
		result := evalProgram(t, tt.input)

		if !testIntegerResult(t, result, tt.expected) {
			t.Errorf("[%d] failed testing integer result", i)
		}
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"!true", false},
		{"!false", true},
		{"!!false", false},
		{"!!true", true},
		{"!81", false},
		{"!!81", true},
	}

	for i, tt := range tests {
		result := evalProgram(t, tt.input)

		if !testBooleanResult(t, result, tt.expected) {
			t.Errorf("[%d] failed testing boolean result", i)
		}
	}
}

func testIntegerResult(t *testing.T, result object.Object, expected int64) bool {
	intobj, ok := result.(*object.Integer)

	if !ok {
		t.Fatalf("expected *object.Integer, got %T", result)
		return false
	}

	if act := intobj.Value; expected != act {
		t.Errorf("expected %d, got %d", expected, act)
		return false
	}
	return true
}

func testBooleanResult(t *testing.T, result object.Object, expected bool) bool {
	boolobj, ok := result.(*object.Boolean)

	if !ok {
		t.Fatalf("expected *object.Boolean, got %T", result)
		return false
	}

	if act := boolobj.Value; expected != act {
		t.Errorf("expected %t, got %t", expected, act)
		return false
	}
	return true
}

func evalProgram(t *testing.T, input string) object.Object {
	l := lexer.NewFromString(input)
	p := parser.New(l)

	program := p.ParseProgram()
	failIfParserHasErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	return Eval(program)
}

func failIfParserHasErrors(t *testing.T, p *parser.Parser) {
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

package eval

import (
	"testing"

	"github.com/MichaelDiBernardo/monkey/ast"
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
		program := parseProgram(t, tt.input)

		result := Eval(program)

		if !testIntegerLiteral(result, t, tt.expected) {
			t.Errorf("[%d] failed testing integer literal", i)
		}
	}
}

func testIntegerLiteral(result object.Object, t *testing.T, expected int64) bool {
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

func parseProgram(t *testing.T, input string) *ast.Program {
	l := lexer.NewFromString(input)
	p := parser.New(l)

	program := p.ParseProgram()
	failIfParserHasErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	return program
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

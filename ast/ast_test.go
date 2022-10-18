package ast

import (
	"testing"

	"github.com/MichaelDiBernardo/monkey/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				LetToken: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					IdentToken: token.Token{Type: token.IDENTIFIER, Literal: "myVar"},
					Value:      "myVar",
				},
				Value: &Identifier{
					IdentToken: token.Token{Type: token.IDENTIFIER, Literal: "anotherVar"},
					Value:      "anotherVar",
				},
			},
		},
	}

	if s := program.String(); s != "let myVar = anotherVar;" {
		t.Errorf("program.String(): got %q", s)
	}
}

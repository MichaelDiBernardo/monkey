package ast

import (
	"bytes"
	"fmt"

	"github.com/MichaelDiBernardo/monkey/token"
)

// Node is any node in a Monkey AST. There are two general kinds of nodes;
// statements, which represent sequential instructions in the program, and
// expressions, which are contained in statements and yield values.
type Node interface {
	Token() token.Token // Token is the token that starts this node.
	String() string     // Render this node as Monkey code.
}

// Statement is a single logical instruction in a Monkey program. A Monkey
// program is a series of statements.
type Statement interface {
	Node
	statementNode()
}

// Expression is a node that represents a segment of Monkey code that yields a
// value.
type Expression interface {
	Node
	expressionNode()
}

// Program is a sequence of Monkey statements.
type Program struct {
	Statements []Statement
}

func (p *Program) Token() token.Token {
	if len(p.Statements) > 0 {
		return p.Statements[0].Token()
	} else {
		return token.Token{}
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, stmt := range p.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

type LetStatement struct {
	LetToken token.Token
	Name     *Identifier // Name is 'x' in 'let x = 24'
	Value    Expression  // Value is 24 in 'let x = 24'
}

func (ls *LetStatement) statementNode()     {}
func (ls *LetStatement) Token() token.Token { return ls.LetToken }

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.LetToken.Literal)
	out.WriteString(" ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	// TODO: Remove once we have expressions.
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

type ReturnStatement struct {
	ReturnToken token.Token
	Value       Expression
}

func (rs *ReturnStatement) statementNode()     {}
func (rs *ReturnStatement) Token() token.Token { return rs.ReturnToken }

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.ReturnToken.Literal)
	out.WriteString(" ")

	// TODO: Remove once we have expressions.
	if rs.Value != nil {
		out.WriteString(rs.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

// ExpressionStatement represents a bare expression e.g. that is typed into the
// REPL. For example, '>> 24 + 3'
type ExpressionStatement struct {
	FirstToken token.Token // First token in the expression statement.
	Value      Expression
}

func (es *ExpressionStatement) statementNode()     {}
func (es *ExpressionStatement) Token() token.Token { return es.FirstToken }

func (es *ExpressionStatement) String() string {
	// TODO: Remove once we have expressions.
	if es.Value != nil {
		return es.Value.String()
	}
	return ""
}

// Identifier is an expression composed of a single identifier.
type Identifier struct {
	IdentToken token.Token
	Value      string // Same as IdentToken.Literal
}

func (i *Identifier) expressionNode()    {}
func (i *Identifier) Token() token.Token { return i.IdentToken }

func (i *Identifier) String() string {
	return i.Value
}

// IntegerLiteral is an expression composed of an integer literal.
type IntegerLiteral struct {
	IntToken token.Token
	Value    int64
}

func (il *IntegerLiteral) expressionNode()    {}
func (il *IntegerLiteral) Token() token.Token { return il.IntToken }

func (il *IntegerLiteral) String() string {
	return fmt.Sprintf("%d", il.Value)
}

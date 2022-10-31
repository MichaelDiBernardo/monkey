package ast

import (
	"bytes"
	"fmt"
	"strings"

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

// BlockStatement is an aggregate of statements contained within curly braces.
type BlockStatement struct {
	StartToken token.Token // The LPAREN token that starts the block.
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

func (bs *BlockStatement) Token() token.Token {
	return bs.StartToken
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	out.WriteString("{")
	for _, stmt := range bs.Statements {
		out.WriteString(stmt.String())
	}
	out.WriteString("}")

	return out.String()
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

// BooleanLiteral is an expression composed of an integer literal.
type BooleanLiteral struct {
	BoolToken token.Token
	Value     bool
}

func (bl *BooleanLiteral) expressionNode()    {}
func (bl *BooleanLiteral) Token() token.Token { return bl.BoolToken }

func (bl *BooleanLiteral) String() string {
	return fmt.Sprintf("%t", bl.Value)
}

type PrefixExpression struct {
	OperatorToken token.Token
	Operator      string
	RHS           Expression
}

func (pe *PrefixExpression) expressionNode()    {}
func (pe *PrefixExpression) Token() token.Token { return pe.OperatorToken }

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.RHS.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	OperatorToken token.Token
	Operator      string
	LHS           Expression
	RHS           Expression
}

func (pe *InfixExpression) expressionNode()    {}
func (pe *InfixExpression) Token() token.Token { return pe.OperatorToken }

func (pe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.LHS.String())
	out.WriteString(" ")
	out.WriteString(pe.Operator)
	out.WriteString(" ")
	out.WriteString(pe.RHS.String())
	out.WriteString(")")

	return out.String()
}

type IfExpression struct {
	IfToken     token.Token // IF token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()    {}
func (ie *IfExpression) Token() token.Token { return ie.IfToken }

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	FnToken    token.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()    {}
func (fl *FunctionLiteral) Token() token.Token { return fl.FnToken }

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")

	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	LPToken   token.Token // The lparen before the args.
	Function  Expression  // Expression that evaluates to func literal
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()    {}
func (ce *CallExpression) Token() token.Token { return ce.LPToken }

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

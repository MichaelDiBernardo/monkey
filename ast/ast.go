package ast

import "github.com/MichaelDiBernardo/monkey/token"

type Node interface {
	Token() token.Token
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

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

type LetStatement struct {
	LetToken token.Token
	Name     *Identifier
	Value    Expression
}

func (ls *LetStatement) statementNode()     {}
func (ls *LetStatement) Token() token.Token { return ls.LetToken }

type ReturnStatement struct {
	ReturnToken token.Token
	Value       Expression
}

func (rs *ReturnStatement) statementNode()     {}
func (rs *ReturnStatement) Token() token.Token { return rs.ReturnToken }

type Identifier struct {
	IdentToken token.Token
	Value      string
}

func (i *Identifier) expressionNode()    {}
func (i *Identifier) Token() token.Token { return i.IdentToken }

package parser

import (
	"github.com/MichaelDiBernardo/monkey/ast"
	"github.com/MichaelDiBernardo/monkey/lexer"
	"github.com/MichaelDiBernardo/monkey/token"
)

type Parser struct {
	lexer *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{LetToken: p.curToken}

	if !p.advanceIfPeekTokenIs(token.IDENTIFIER) {
		return nil
	}
	ident := p.curToken

	if !p.peekToken.Is(token.ASSIGN) {
		return nil
	}

	stmt.Name = &ast.Identifier{IdentToken: ident, Value: ident.Literal}

	// Skip expression for now.
	for p.curToken.Is(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) advanceIfPeekTokenIs(ttype token.TokenType) bool {
	if p.peekToken.Is(ttype) {
		p.nextToken()
		return true
	}
	return false
}

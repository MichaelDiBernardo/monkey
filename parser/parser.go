package parser

import (
	"fmt"

	"github.com/MichaelDiBernardo/monkey/ast"
	"github.com/MichaelDiBernardo/monkey/lexer"
	"github.com/MichaelDiBernardo/monkey/token"
)

type ParseError struct {
	Message  string
	Location token.Location
}

func (pe *ParseError) String() string {
	msg := "%s (at line %d, col %d)"
	return fmt.Sprintf(msg, pe.Message, pe.Location.LineN, pe.Location.CharN)
}

type Parser struct {
	lexer *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []ParseError
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l, errors: []ParseError{}}
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

func (p *Parser) Errors() []ParseError {
	return p.errors
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
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

	if !p.advanceIfPeekTokenIs(token.ASSIGN) {
		return nil
	}

	stmt.Name = &ast.Identifier{IdentToken: ident, Value: ident.Literal}

	// Skip expression for now.
	for p.curToken.Is(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{ReturnToken: p.curToken}

	// Skip expression for now.
	for p.curToken.Is(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// advanceIfPeekTokenIs is a predicate that also does a bunch of work. This is
// generally frowned upon, but for this specific operation is pretty handy.
//
// It:
// - calls nextToken() and returns true if ttype matches the peekToken's type
// - adds a parse error and returns false otherwise
func (p *Parser) advanceIfPeekTokenIs(ttype token.TokenType) bool {
	if p.peekToken.Is(ttype) {
		p.nextToken()
		return true
	} else {
		p.addErrorForMismatchedPeekToken(ttype)
		return false
	}
}

// addErrorForMismatchedPeekToken adds an appropriate error to the errors
// collection when the peekToken's type doesn't match the given expectedType.
func (p *Parser) addErrorForMismatchedPeekToken(expectedType token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s '%s' instead", expectedType, p.peekToken.Type, p.peekToken.Literal)
	p.errors = append(p.errors, ParseError{Message: msg, Location: p.peekToken.Location})
}

package parser

import (
	"fmt"
	"strconv"

	"github.com/MichaelDiBernardo/monkey/ast"
	"github.com/MichaelDiBernardo/monkey/lexer"
	"github.com/MichaelDiBernardo/monkey/token"
)

// ParseError is a message that is emitted by the parser when there is an error.
// Message should be legible by the program author if emitted by the parser.
type ParseError struct {
	Message  string
	Location token.Location
}

// String() should be legible by the program author if emitted by the parser.
func (pe *ParseError) String() string {
	msg := "%s (at line %d, col %d)"
	return fmt.Sprintf(msg, pe.Message, pe.Location.LineN, pe.Location.CharN)
}

type (
	// prefixParseFn is a Pratt prefix-parse function. It is called when parsing
	// a prefix operation.
	prefixParseFn func() ast.Expression
	// infixParseFn is a Pratt infix-parse function. It is called when parsing
	// an infix operation.
	infixParseFn func(ast.Expression) ast.Expression
)

// The precedence an operation takes over another.
type Precedence int

const (
	_ Precedence = iota
	P_LOWEST
	P_EQUALS
	P_LESSGREATER
	P_SUM
	P_PRODUCT
	P_PREFIX
	P_CALL
)

type Parser struct {
	lexer  *lexer.Lexer
	errors []ParseError

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l, errors: []ParseError{}}
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		program.Statements = append(program.Statements, stmt)
		p.nextToken()
	}
	return program
}

func (p *Parser) Errors() []ParseError {
	return p.errors
}

func (p *Parser) registerPrefix(tt token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tt] = fn
}

func (p *Parser) registerInfix(tt token.TokenType, fn infixParseFn) {
	p.infixParseFns[tt] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
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
	for !p.curToken.Is(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{ReturnToken: p.curToken}

	// Skip expression for now.
	for !p.curToken.Is(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		FirstToken: p.curToken,
		Value:      p.parseExpression(P_LOWEST),
	}

	// Semicolons are optional in expression statements.
	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence Precedence) ast.Expression {
	pfn := p.prefixParseFns[p.curToken.Type]

	if pfn == nil {
		p.addErrorForMissingPrefixFn(p.curToken.Type)
		return nil
	}
	lhs := pfn()

	return lhs
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

func (p *Parser) addErrorForMissingPrefixFn(tt token.TokenType) {
	msg := fmt.Sprintf("no prefix parse fn for tokentype %s", tt)
	p.errors = append(p.errors, ParseError{Message: msg, Location: p.curToken.Location})
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{IdentToken: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	intval, err := strconv.ParseInt(p.curToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %s as integer literal", p.curToken.Literal)
		p.errors = append(p.errors, ParseError{Message: msg, Location: p.curToken.Location})
		return nil
	}

	return &ast.IntegerLiteral{IntToken: p.curToken, Value: intval}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		OperatorToken: p.curToken,
		Operator:      p.curToken.Literal,
	}

	p.nextToken()

	expression.RHS = p.parseExpression(P_PREFIX)
	return expression
}

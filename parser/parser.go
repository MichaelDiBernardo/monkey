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

var precedences = map[token.TokenType]Precedence{
	token.EQ:       P_EQUALS,
	token.NEQ:      P_EQUALS,
	token.LANGLE:   P_LESSGREATER,
	token.RANGLE:   P_LESSGREATER,
	token.PLUS:     P_SUM,
	token.MINUS:    P_SUM,
	token.RSLASH:   P_PRODUCT,
	token.ASTERISK: P_PRODUCT,
}

func precedenceOfTokenType(tt token.TokenType) Precedence {
	if prec, ok := precedences[tt]; ok {
		return prec
	}

	return P_LOWEST
}

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
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.RSLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.LANGLE, p.parseInfixExpression)
	p.registerInfix(token.RANGLE, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curToken.Is(token.EOF) {
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

// advanceIfPeekTokenIs is a predicate that also does a bunch of work. This is
// generally frowned upon, but for this specific operation is pretty handy.
//
// It:
// - calls nextToken() and returns true if ttype matches the peekToken's type
// - returns false otherwise
func (p *Parser) advanceIfPeekTokenIs(ttype token.TokenType) bool {
	if p.peekToken.Is(ttype) {
		p.nextToken()
		return true
	} else {
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
	msg := fmt.Sprintf("unexpected token type %s while parsing prefix expression", tt)
	p.errors = append(p.errors, ParseError{Message: msg, Location: p.curToken.Location})
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
		p.addErrorForMismatchedPeekToken(token.IDENTIFIER)
		return nil
	}

	ident := p.curToken
	stmt.Name = &ast.Identifier{IdentToken: ident, Value: ident.Literal}

	if !p.advanceIfPeekTokenIs(token.ASSIGN) {
		p.addErrorForMismatchedPeekToken(token.ASSIGN)
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(P_LOWEST)

	if p.peekToken.Is(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{ReturnToken: p.curToken}

	p.nextToken()
	stmt.Value = p.parseExpression(P_LOWEST)

	if p.peekToken.Is(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// Assign these outside the struct declaration, because parseExpression
	// advances curToken.
	firstToken := p.curToken
	value := p.parseExpression(P_LOWEST)

	stmt := &ast.ExpressionStatement{
		FirstToken: firstToken,
		Value:      value,
	}

	// Semicolons are optional in expression statements.
	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{StartToken: p.curToken}
	p.nextToken()

	block.Statements = []ast.Statement{}

	for !p.curToken.Is(token.RBRACE) && !p.curToken.Is(token.EOF) {
		stmt := p.parseStatement()
		block.Statements = append(block.Statements, stmt)
		p.nextToken()
	}

	return block
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{IdentToken: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	boolval, err := strconv.ParseBool(p.curToken.Literal)

	if err != nil {
		msg := fmt.Sprintf("could not parse %s as bool literal", p.curToken.Literal)
		p.errors = append(p.errors, ParseError{Message: msg, Location: p.curToken.Location})
		return nil
	}

	return &ast.BooleanLiteral{BoolToken: p.curToken, Value: boolval}
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

func (p *Parser) parseInfixExpression(lhs ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		LHS:           lhs,
		OperatorToken: p.curToken,
		Operator:      p.curToken.Literal,
	}

	prec := precedenceOfTokenType(p.curToken.Type)
	p.nextToken()
	expression.RHS = p.parseExpression(prec)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(P_LOWEST)

	if !p.advanceIfPeekTokenIs(token.RPAREN) {
		p.addErrorForMissingPrefixFn(token.RPAREN)
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	iftok := p.curToken

	if !p.advanceIfPeekTokenIs(token.LPAREN) {
		p.addErrorForMismatchedPeekToken(token.LPAREN)
		return nil
	}

	p.nextToken()

	condition := p.parseExpression(P_LOWEST)

	if !p.advanceIfPeekTokenIs(token.RPAREN) {
		p.addErrorForMismatchedPeekToken(token.RPAREN)
		return nil
	}

	if !p.advanceIfPeekTokenIs(token.LBRACE) {
		p.addErrorForMismatchedPeekToken(token.LBRACE)
		return nil
	}

	consequence := p.parseBlockStatement()

	if consequence == nil {
		return nil
	}

	var alternative *ast.BlockStatement = nil

	if p.advanceIfPeekTokenIs(token.ELSE) {
		if !p.advanceIfPeekTokenIs(token.LBRACE) {
			p.addErrorForMismatchedPeekToken(token.LBRACE)
			return nil
		}
		alternative = p.parseBlockStatement()

		if alternative == nil {
			p.errors = append(p.errors, ParseError{Message: "couldn't parse else clause", Location: iftok.Location})
			return nil
		}
	}

	return &ast.IfExpression{IfToken: iftok, Condition: condition, Consequence: consequence, Alternative: alternative}

}
func (p *Parser) parseExpression(precedence Precedence) ast.Expression {
	pfn := p.prefixParseFns[p.curToken.Type]

	if pfn == nil {
		p.addErrorForMissingPrefixFn(p.curToken.Type)
		return nil
	}
	lhs := pfn()

	for !p.peekToken.Is(token.SEMICOLON) && precedence < precedenceOfTokenType(p.peekToken.Type) {
		infix := p.infixParseFns[p.peekToken.Type]

		if infix == nil {
			return lhs
		}

		p.nextToken()

		lhs = infix(lhs)
	}

	return lhs
}

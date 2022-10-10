package lexer

import "github.com/MichaelDiBernardo/monkey/token"

type Lexer struct {
	input      string
	currentPos int  // current position in input
	peekPos    int  // lookahead position
	ch         byte // char being inspected
}

const NUL byte = 0

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.eatWhitespace()

	switch l.ch {
	case '=':
		if l.peek() == '=' {
			tok = token.Token{Type: token.EQ, Literal: "=="}
			l.readChar()
		} else {
			tok = token.NewFromByte(token.ASSIGN, l.ch)
		}
	case '!':
		if l.peek() == '=' {
			tok = token.Token{Type: token.NEQ, Literal: "!="}
			l.readChar()
		} else {
			tok = token.NewFromByte(token.BANG, l.ch)
		}
	case '+':
		tok = token.NewFromByte(token.PLUS, l.ch)
	case '-':
		tok = token.NewFromByte(token.MINUS, l.ch)
	case '*':
		tok = token.NewFromByte(token.ASTERISK, l.ch)
	case '/':
		tok = token.NewFromByte(token.RSLASH, l.ch)
	case ',':
		tok = token.NewFromByte(token.COMMA, l.ch)
	case ';':
		tok = token.NewFromByte(token.SEMICOLON, l.ch)
	case '(':
		tok = token.NewFromByte(token.LPAREN, l.ch)
	case ')':
		tok = token.NewFromByte(token.RPAREN, l.ch)
	case '{':
		tok = token.NewFromByte(token.LBRACE, l.ch)
	case '}':
		tok = token.NewFromByte(token.RBRACE, l.ch)
	case '<':
		tok = token.NewFromByte(token.LANGLE, l.ch)
	case '>':
		tok = token.NewFromByte(token.RANGLE, l.ch)
	case NUL:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isIdentifierChar(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupMulticharTokenType(tok.Literal)
			return tok
		}
		if isNumericChar(l.ch) {
			tok.Literal = l.readNumericLiteral()
			tok.Type = token.INT
			return tok
		}
		tok = token.NewFromByte(token.ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.peekPos >= len(l.input) {
		l.ch = NUL
	} else {
		l.ch = l.input[l.peekPos]
	}
	l.currentPos = l.peekPos
	l.peekPos += 1
}

func (l *Lexer) readIdentifier() string {
	pos := l.currentPos
	for isIdentifierChar(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.currentPos]
}

func (l *Lexer) readNumericLiteral() string {
	pos := l.currentPos
	for isNumericChar(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.currentPos]
}

func (l *Lexer) eatWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) peek() byte {
	if l.peekPos >= len(l.input) {
		return NUL
	}
	return l.input[l.peekPos]
}

func isIdentifierChar(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isNumericChar(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

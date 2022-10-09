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

	switch l.ch {
	case '=':
		tok = token.NewFromByte(token.ASSIGN, l.ch)
	case '+':
		tok = token.NewFromByte(token.PLUS, l.ch)
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
	case NUL:
		tok.Literal = ""
		tok.Type = token.EOF
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

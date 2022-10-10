package lexer

import "github.com/MichaelDiBernardo/monkey/token"

type Lexer struct {
	input      string
	currentPos int            // current position in input
	peekPos    int            // lookahead position
	ch         byte           // char being inspected
	currentLoc token.Location // Location of current token
}

const NUL byte = 0

func New(input string) *Lexer {
	l := &Lexer{input: input, currentLoc: token.Location{Path: token.NO_FILEPATH, CharN: 0, LineN: 1}}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.eatWhitespace()

	switch l.ch {
	case '=':
		if l.peek() == '=' {
			tok = token.NewTwoCharToken(token.EQ, "==", l.currentLoc)
			l.readChar()
		} else {
			tok = token.NewOneCharToken(token.ASSIGN, l.ch, l.currentLoc)
		}
	case '!':
		if l.peek() == '=' {
			tok = token.NewTwoCharToken(token.NEQ, "!=", l.currentLoc)
			l.readChar()
		} else {
			tok = token.NewOneCharToken(token.BANG, l.ch, l.currentLoc)
		}
	case '+':
		tok = token.NewOneCharToken(token.PLUS, l.ch, l.currentLoc)
	case '-':
		tok = token.NewOneCharToken(token.MINUS, l.ch, l.currentLoc)
	case '*':
		tok = token.NewOneCharToken(token.ASTERISK, l.ch, l.currentLoc)
	case '/':
		tok = token.NewOneCharToken(token.RSLASH, l.ch, l.currentLoc)
	case ',':
		tok = token.NewOneCharToken(token.COMMA, l.ch, l.currentLoc)
	case ';':
		tok = token.NewOneCharToken(token.SEMICOLON, l.ch, l.currentLoc)
	case '(':
		tok = token.NewOneCharToken(token.LPAREN, l.ch, l.currentLoc)
	case ')':
		tok = token.NewOneCharToken(token.RPAREN, l.ch, l.currentLoc)
	case '{':
		tok = token.NewOneCharToken(token.LBRACE, l.ch, l.currentLoc)
	case '}':
		tok = token.NewOneCharToken(token.RBRACE, l.ch, l.currentLoc)
	case '<':
		tok = token.NewOneCharToken(token.LANGLE, l.ch, l.currentLoc)
	case '>':
		tok = token.NewOneCharToken(token.RANGLE, l.ch, l.currentLoc)
	case NUL:
		tok = token.NewOneCharToken(token.EOF, NUL, l.currentLoc)
	default:
		if isIdentifierChar(l.ch) {
			literal := l.readIdentifier()
			tok = token.NewMultiCharToken(token.LookupMulticharTokenType(literal), literal, l.currentLoc)
			return tok
		}
		if isNumericChar(l.ch) {
			tok = token.NewMultiCharToken(token.INT, l.readNumericLiteral(), l.currentLoc)
			return tok
		}
		tok = token.NewOneCharToken(token.ILLEGAL, l.ch, l.currentLoc)
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

	if l.ch == '\n' {
		l.currentLoc.NextLine()
	} else {
		l.currentLoc.NextChar()
	}
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

package lexer

import "github.com/MichaelDiBernardo/monkey/token"

type Lexer struct {
	input        string
	position     int  // current position in input
	readPosition int  // lookahead position
	ch           byte // char being inspected
}

func New(input string) *Lexer {
	return &Lexer{input: input}
}

func (l *Lexer) NextToken() token.Token {
	return token.Token{}
}

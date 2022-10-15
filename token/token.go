// There are three broad types of tokens in monkey: Single-character tokens,
// double-character tokens, and multi-character tokens. These are all scanned
// slightly differently, and so these token types are made explicit here by
// giving them all unique constructors.
package token

import "fmt"

type TokenType string

const NO_FILEPATH = "<input>"

// A location in a Monkey program text.
type Location struct {
	Path  string // Full path to filename of program.
	LineN uint   // 1-indexed line number
	CharN uint   // 1-indexed character number in line
}

// Increment the line number this location is tracking.
func (l *Location) NextLine() {
	l.LineN++
	l.CharN = 0
}

// Increment the character this location is tracking on the current line.
func (l *Location) NextChar() {
	l.CharN++
}

// A Monkey-language token.
type Token struct {
	Type     TokenType
	Literal  string
	Location Location
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENTIFIER = "IDENTIFIER"
	INT        = "INT"

	ASSIGN   = "="
	EQ       = "=="
	NEQ      = "!="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	BANG     = "!"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	LANGLE = "<"
	RANGLE = ">"
	RSLASH = "/"

	FUNCTION = "FUNCTION"
	LET      = "LET"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

// One char literals are detected at the read head, so location is not modified
// for these tokens.
func NewOneCharToken(tokenType TokenType, ch byte, location Location) Token {
	return Token{Type: tokenType, Literal: string(ch), Location: location}
}

// Two char tokens are detected on the first char of the literal (e.g. on the
// bang of "!=") so location is not modified for these tokens.
func NewTwoCharToken(tokenType TokenType, literal string, location Location) Token {
	if len(literal) != 2 {
		panic(fmt.Sprintf("NewTwoCharToken given literal '%s' at loc %q", literal, location))
	}
	return Token{Type: tokenType, Literal: literal, Location: location}
}

// Given a token type, a literal, and a location from the lexer, this will
// return a Token composed of those bits. This function expects a location that
// has CharN pointing to the whitespace at the _end_ of the multichar literal;
// it will adjust the location to point to its start.
func NewMultiCharToken(tokenType TokenType, literal string, location Location) Token {
	tok := Token{Type: tokenType, Literal: literal, Location: location}
	tok.Location.CharN -= uint(len(tok.Literal))
	return tok
}

// Given a literal, checks if it is a keyword. Otherwise, it is an identifier.
func LookupMulticharTokenType(literal string) TokenType {
	if ttype, ok := keywords[literal]; ok {
		return ttype
	}
	return IDENTIFIER
}

func (t *Token) Is(ttype TokenType) bool {
	return t.Type == ttype
}

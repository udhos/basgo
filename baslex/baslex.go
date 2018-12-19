package baslex

import (
	"io"
)

// Tokens
const (
	TkEOF   = iota // EOF
	TkFIXME = iota // FIXME

	TkLineNumber = iota // Line number
	TkString     = iota // String
	TkEqual      = iota // Equal
	TkUnequal    = iota // Unequal

	TkKeywordCls = iota // CLS

	TkIdentifier = iota // Identifier
)

// Token is a lexical token
type Token struct {
	ID     int
	Value  string
	Offset int
}

// IsEOF checks for EOF token
func (t Token) IsEOF() bool {
	return t.ID == TkEOF
}

// Lex is a full lexer object
type Lex struct {
	r   io.Reader
	eof bool
}

// New creates a Lex object
func New(input io.Reader) *Lex {
	return &Lex{r: input}
}

var tokenEOF = Token{ID: TkEOF, Value: "EOF"}

func (l *Lex) returnEOF() Token {
	l.eof = true // EOF sent
	return tokenEOF
}

// Next gets next token
func (l *Lex) Next() Token {
	if l.eof {
		// will send EOF forever
		return l.returnEOF()
	}
	return Token{ID: TkFIXME, Value: "FIXME-WRITEME"}
}

// HasToken checks if there are more tokens
func (l *Lex) HasToken() bool {
	return !l.eof
}

package baslex

import (
	"io"
)

const (
	TkEOF = iota
	TkError = iota

	TkLineNumber = iota
	TkString = iota
	TkEqual = iota
	TkUnequal = iota

	TkKeywordCls = iota

	TkIdentifier = iota
)

type Token struct {
	Id int
	Value string
	Offset int
}

func (t Token) IsEOF() bool {
	return t.ID == TkEOF
}

type Lex struct {
	r io.Reader
	eof bool
}

func New(input io.Reader) *Lex {
	return &Lex{r: input}
}

func (l *Lex) returnEOF() Token {
	l.eof = true
	return Token{Id: TkEOF, Value: "EOF"}
}

func (l *Lex) Next() Token {
	if eof {
		return l.returnEOF()
	}
	return Token{Id: TkError, Value: "FIXME-WRITEME"}

func (l *Lex) HasToken() bool {
	return l.eof
}

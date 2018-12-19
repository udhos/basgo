package baslex

import (
	"io"
	"log"
)

// Tokens
const (
	TkNull  = iota // Null token should never be seen
	TkEOF   = iota // EOF
	TkFIXME = iota // FIXME

	TkErrInput    = iota // Input error -- first error
	TkErrInternal = iota // Internal error
	TkErrLarge    = iota // Large token -- last error

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

// IsError checks for error token
func (t Token) IsError() bool {
	return t.ID >= TkErrInput && t.ID <= TkErrLarge
}

// Lex is a full lexer object
type Lex struct {
	r           io.Reader
	eof         bool // has sent EOF?
	broken      bool // hit error?
	buf         []byte
	tokenOffset int
	tokenSize   int
}

// New creates a Lex object
func New(input io.Reader) *Lex {
	return &Lex{r: input, buf: make([]byte, 0, 10)}
}

var tokenNull = Token{}
var tokenEOF = Token{ID: TkEOF, Value: "EOF"}
var tokenFIXME = Token{ID: TkFIXME, Value: "FIXME-WRITEME"}
var tokenErrInput = Token{ID: TkErrInput, Value: "ERROR-INPUT"}
var tokenErrInternal = Token{ID: TkErrInternal, Value: "ERROR-INTERNAL"}
var tokenErrLarge = Token{ID: TkErrLarge, Value: "ERROR-LARGE-TOKEN"}

func (l *Lex) returnEOF() Token {
	l.eof = true // set EOF, no more tokens
	return tokenEOF
}

// Next gets next token
func (l *Lex) Next() Token {
	if !l.HasToken() {
		// will send EOF forever
		return l.returnEOF()
	}

	t := l.findToken()
	if t.IsError() {
		l.broken = true // set fail state, no more tokens
	}

	return t
}

// HasToken checks if there are more tokens
func (l *Lex) HasToken() bool {
	return !l.eof && !l.broken
}

func (l *Lex) findToken() Token {

	log.Printf("findToken: len=%d cap=%d", len(l.buf), cap(l.buf))

	for {
		if consume := l.tokenOffset + l.tokenSize; consume > 0 {
			// shift last token
			log.Printf("findToken: consume=%d", consume)
			l.buf = append(l.buf[:0], l.buf[consume:]...)
			l.tokenOffset = 0
			l.tokenSize = 0
		}

		size := len(l.buf)
		if size >= cap(l.buf) {
			return tokenErrLarge // no room for more data
		}

		// grab more data
		n, errRead := l.r.Read(l.buf[size:cap(l.buf)])
		l.buf = l.buf[:size+n]
		log.Printf("findToken: size=%d read=[%s]", n, string(l.buf[size:]))
		switch errRead {
		case nil:
			if n < 1 {
				return tokenErrInternal // ugh should not happen
			}
			l.buf = l.buf[:size+n]
			log.Printf("findToken: buf=[%s]", string(l.buf))
		case io.EOF:
			return l.returnEOF() // EOF
		default:
			return tokenErrInput // unexpected input error
		}

		if len(l.buf) < 1 {
			return tokenErrInternal // ugh should not happen
		}

		if t, found := l.match(); found {
			return t
		}

		log.Printf("findToken: FIXME WRITEME")
		return tokenFIXME // stop loop
	}

	// NOT REACHED
}

func (l *Lex) match() (Token, bool) {
	log.Printf("match: FIXME WRITEME")
	return tokenFIXME, true
}

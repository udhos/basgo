package baslex

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	//"log"
	"strings"
)

// Tokens
const (
	TkNull  = iota // Null token should never be seen
	TkEOF   = iota // EOF
	TkEOL   = iota // EOL
	TkFIXME = iota // FIXME

	TkErrInput    = iota // Input error -- first error
	TkErrInternal = iota // Internal error
	TkErrLarge    = iota // Large token -- last error

	TkColon      = iota // Colon :
	TkCommentQ   = iota // Comment '
	TkLineNumber = iota // Line number
	TkString     = iota // String
	TkEqual      = iota // Equal
	TkUnequal    = iota // Unequal

	TkKeywordCls = iota // CLS

	TkIdentifier = iota // Identifier (variable)
)

const (
	stBlank    = iota
	stCommentQ = iota
	//stName        = iota
	//stNumber      = iota
	//stString      = iota
)

// Token is a lexical token
type Token struct {
	ID    int
	Value string
	//	Offset int
}

// IsEOF checks for EOF token
func (t Token) IsEOF() bool {
	return t.ID == TkEOF
}

// IsError checks for error token
func (t Token) IsError() bool {
	return t.ID >= TkErrInput && t.ID <= TkErrLarge || t.ID == TkFIXME
}

// Lex is a full lexer object
type Lex struct {
	r      *bufio.Reader
	eof    bool // hit EOF?
	broken bool // hit error?
	buf    bytes.Buffer
	state  int
}

// New creates a Lex object
func New(input *bufio.Reader) *Lex {
	//return &Lex{r: input, buf: make([]byte, 0, 10)}
	return &Lex{r: input}
}

// NewStr creates a Lex object from string
func NewStr(s string) *Lex {
	return New(bufio.NewReader(strings.NewReader(s)))
}

var tokenNull = Token{}
var tokenEOF = Token{ID: TkEOF, Value: "EOF"}
var tokenFIXME = Token{ID: TkFIXME, Value: "FIXME-WRITEME"}

//var tokenErrInput = Token{ID: TkErrInput, Value: "ERROR-INPUT"}
//var tokenErrInternal = Token{ID: TkErrInternal, Value: "ERROR-INTERNAL"}
//var tokenErrLarge = Token{ID: TkErrLarge, Value: "ERROR-LARGE-TOKEN"}

func (l *Lex) returnTokenEOF() Token {
	//l.eof = true // set EOF, no more tokens
	return tokenEOF
}

// Next gets next token
func (l *Lex) Next() Token {
	if !l.HasToken() {
		// will send EOF forever
		return l.returnTokenEOF()
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

	for {
		b, errRead := l.r.ReadByte()
		switch errRead {
		case nil:
		case io.EOF:
			return l.foundEOF()
		default:
			return l.foundErrorInput(errRead)
		}

		t := l.match(b)
		if t.ID != TkNull {
			return t
		}
	}

	// never reached
}

func (l *Lex) foundErrorInput(err error) Token {
	return Token{ID: TkErrInput, Value: fmt.Sprintf("ERROR-INPUT: after [%s]: %s", l.buf.String(), err.Error())}
}

/*
func (l *Lex) findToken() Token {

	log.Printf("findToken: len=%d cap=%d", len(l.buf), cap(l.buf))

	for {
		size := len(l.buf)
		if size >= cap(l.buf) {
			return tokenErrLarge // no room for more data
		}

		if size < 1 {
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
				return l.foundEOF() // EOF
			default:
				return tokenErrInput // unexpected input error
			}
		}

		if len(l.buf) < 1 {
			return tokenErrInternal // ugh should not happen
		}

		if t, found := l.match(); found {
			return t
		}
	}

	// NOT REACHED
}

func (l *Lex) consumeToken(size int) {
	log.Printf("consume: size=%d [%s]", size, string(l.buf[:size]))
	l.buf = append(l.buf[:0], l.buf[size:]...)
}

func (l *Lex) foundEOF() Token {
	switch l.state {
	case stCommentQ:
		return Token{ID: TkCommentQ, Value: "'", Offset: l.offset}
	}
	return l.returnEOF() // EOF
}

func (l *Lex) match() (Token, bool) {

	for i := 0; i < len(l.buf); i++ {
		b := l.buf[i]
		switch l.state {
		case stBegin:
			switch {
			case b == '\r':
				// Search for LF
				l.state = stBeginCR
			case b == '\n':
				// LF
				l.consumeToken(i + 1)
				return Token{ID: TkEOL, Value: "EOL", Offset: l.offset}, true
			case b == ' ':
				l.state = stBlank
			case b == '\'':
				l.state = stCommentQ
				l.consumeToken(i + 1)
				return Token{ID: TkCommentQ, Value: "'", Offset: l.offset}, true
			}
		case stCommentQ:
			switch {
			case b == '\r':
				// Search for LF
				l.state = stCommentQ_CR
			case b == '\n':
				// LF
				l.consumeToken(i + 1)
				return Token{ID: TkCommentQ, Value: "'", Offset: l.offset}, true
			}
		case stCommentQ_CR:
			if b == '\n' {
				// CR LF
				l.state = stBegin
				l.consumeToken(i + 1)
				return Token{ID: TkCommentQ, Value: "'", Offset: l.offset}, true
			}
			l.state = stBegin // restart
			//l.consumeToken(0) // do not consume = push back
			return Token{ID: TkCommentQ, Value: "'", Offset: l.offset}, true
		case stBeginCR:
			if b == '\n' {
				// CR LF
				l.consumeToken(i + 1)
				return Token{ID: TkEOL, Value: "EOL", Offset: l.offset}, true
			}
			l.state = stBegin // restart
			//l.consumeToken(0) // do not consume = push back
			return Token{ID: TkEOL, Value: "EOL", Offset: l.offset}, true
		case stBlank:
		case stName:
		case stNumber:
		case stString:
		default:
			return tokenErrInternal, true // ugh should not happen
		}
		l.offset++
	}

	return tokenNull, false
}
*/

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

	TkColon    = iota // Colon :
	TkCommentQ = iota // Comment '
	TkString   = iota // String "
	TkNumber   = iota // Number [0-9]+

	TkEqual   = iota // Equal
	TkLT      = iota // <
	TkGT      = iota // >
	TkUnequal = iota // Unequal <>
	TkLE      = iota // <=
	TkGE      = iota // >=

	TkPlus  = iota // +
	TkMinus = iota // -
	TkMult  = iota // *
	TkDiv   = iota // /

	TkKeywordCls   = iota // CLS
	TkKeywordEnd   = iota // END
	TkKeywordPrint = iota // PRINT
	TkKeywordTime  = iota // TIME$

	TkIdentifier = iota // Identifier (variable)
)

var tabType = []string{
	"NULL",
	"EOF",
	"EOL",
	"FIXME",

	"ERROR-INPUT",
	"ERROR-INTERNAL",
	"ERROR-LARGE",

	"COLON",
	"COMMENT-Q",
	"STRING",
	"NUMBER",

	"EQUAL",
	"LT",
	"GT",
	"UNEQUAL",
	"LE",
	"GE",

	"PLUS",
	"MINUS",
	"MULT",
	"DIV",

	"CLS",
	"END",
	"PRINT",
	"TIME",

	"IDENTIFIER",
}

// Token is a lexical token
type Token struct {
	ID        int
	Value     string
	LineCount int
}

// Type gets the token type
func (t Token) Type() string {
	if t.ID < 0 || t.ID >= len(tabType) {
		return fmt.Sprintf("TYPE-ERROR:%d", t.ID)
	}
	return tabType[t.ID]
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
	r       io.ByteScanner
	eofSeen bool // hit EOF?
	eofSent bool // delivered EOF?
	broken  bool // hit error?
	buf     bytes.Buffer
	state   int
}

// New creates a Lex object
func New(input io.ByteScanner) *Lex {
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
	l.eofSent = true
	return tokenEOF
}

// Next gets next token
func (l *Lex) Next() Token {
	if !l.HasToken() {
		return l.returnTokenEOF()
	}

	if l.eofSeen {
		// deliver pending EOF
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
	return !l.eofSent && !l.broken
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

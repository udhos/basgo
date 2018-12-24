package baslex

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	//"log"
	"strings"
)

/*
Keep tokens in sync:

(A) const tokens
(B) var tabKeywords
(C) var tabType
*/

// (A) const tokens
const (
	TkNull = iota // Null token should never be seen
	TkEOF  = iota // EOF
	TkEOL  = iota // EOL

	TkErrInput    = iota // Input error -- first error
	TkErrInternal = iota // Internal error
	TkErrInvalid  = iota // Invalid, unexpected token found
	TkErrLarge    = iota // Large token -- last error

	TkColon        = iota // Colon :
	TkComma        = iota // Comma ,
	TkSemicolon    = iota // Semicolon ; (newline suppressor)
	TkParLeft      = iota // (
	TkParRight     = iota // )
	TkBracketLeft  = iota // [
	TkBracketRight = iota // ]
	TkCommentQ     = iota // Comment '
	TkString       = iota // String "
	TkNumber       = iota // Number [0-9]+

	TkEqual   = iota // Equal
	TkLT      = iota // <
	TkGT      = iota // >
	TkUnequal = iota // Unequal <>
	TkLE      = iota // <=
	TkGE      = iota // >=

	TkPlus      = iota // +
	TkMinus     = iota // -
	TkMult      = iota // *
	TkDiv       = iota // /
	TkBackSlash = iota // \

	TkKeywordCls    = iota // CLS
	TkKeywordCont   = iota // CONT
	TkKeywordElse   = iota // ELSE
	TkKeywordEnd    = iota // END
	TkKeywordFor    = iota // FOR
	TkKeywordGosub  = iota // GOSUB
	TkKeywordGoto   = iota // GOTO
	TkKeywordInput  = iota // INPUT
	TkKeywordIf     = iota // IF
	TkKeywordLet    = iota // LET
	TkKeywordList   = iota // LIST
	TkKeywordLoad   = iota // LOAD
	TkKeywordNext   = iota // NEXT
	TkKeywordPrint  = iota // PRINT
	TkKeywordRem    = iota // REM
	TkKeywordReturn = iota // RETURN
	TkKeywordRun    = iota // RUN
	TkKeywordSave   = iota // SAVE
	TkKeywordStep   = iota // STEP
	TkKeywordStop   = iota // STOP
	TkKeywordSystem = iota // SYSTEM
	TkKeywordThen   = iota // THEN
	TkKeywordTime   = iota // TIME$
	TkKeywordTo     = iota // TO

	TkIdentifier = iota // Identifier (variable)
)

// Token ID marks
const (
	TokenIDFirst = TkNull
	TokenIDLast  = TkIdentifier
)

// (B) var tabKeywords
var tabKeywords = []struct {
	TokenID int
	Name    string
}{
	{TkKeywordCls, "CLS"},
	{TkKeywordCont, "CONT"},
	{TkKeywordElse, "ELSE"},
	{TkKeywordEnd, "END"},
	{TkKeywordFor, "FOR"},
	{TkKeywordGosub, "GOSUB"},
	{TkKeywordGoto, "GOTO"},
	{TkKeywordIf, "IF"},
	{TkKeywordInput, "INPUT"},
	{TkKeywordLet, "LET"},
	{TkKeywordList, "LIST"},
	{TkKeywordLoad, "LOAD"},
	{TkKeywordNext, "NEXT"},
	{TkKeywordPrint, "PRINT"},
	{TkKeywordRem, "REM"},
	{TkKeywordReturn, "RETURN"},
	{TkKeywordRun, "RUN"},
	{TkKeywordSave, "SAVE"},
	{TkKeywordStep, "STEP"},
	{TkKeywordStop, "STOP"},
	{TkKeywordSystem, "SYSTEM"},
	{TkKeywordThen, "THEN"},
	{TkKeywordTime, "TIME$"},
	{TkKeywordTo, "TO"},
}

// (C) var tabType
var tabType = []string{
	"NULL",
	"EOF",
	"EOL",

	"ERROR-INPUT",
	"ERROR-INTERNAL",
	"ERROR-INVALID",
	"ERROR-LARGE",

	"COLON",
	"COMMA",
	"SEMICOLON",
	"ROUND-BRACKET-LEFT",
	"ROUND-BRACKET-RIGHT",
	"SQUARE-BRACKET-LEFT",
	"SQUARE-BRACKET-RIGHT",
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
	"BACK-SLASH",

	"CLS",
	"CONT",
	"ELSE",
	"END",
	"FOR",
	"GOSUB",
	"GOTO",
	"IF",
	"INPUT",
	"LET",
	"LIST",
	"LOAD",
	"NEXT",
	"PRINT",
	"REM",
	"RETURN",
	"RUN",
	"SAVE",
	"STEP",
	"STOP",
	"SYSTEM",
	"THEN",
	"TIME$",
	"TO",

	"IDENTIFIER",
}

// Token is a lexical token
type Token struct {
	ID         int
	Value      string
	LineCount  int
	LineOffset int
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
	return t.ID >= TkErrInput && t.ID <= TkErrLarge
}

// Lex is a full lexer object
type Lex struct {
	r          io.ByteScanner
	eofSeen    bool // hit EOF?
	eofSent    bool // delivered EOF?
	broken     bool // hit error?
	buf        bytes.Buffer
	state      int
	lineCount  int
	lineOffset int
}

// Line returns current line count.
func (l *Lex) Line() int {
	return l.lineCount
}

// Column returns current column offset within line.
func (l *Lex) Column() int {
	return l.lineOffset
}

// New creates a Lex object
func New(input io.ByteScanner) *Lex {
	return &Lex{r: input, lineCount: 1}
}

// NewStr creates a Lex object from string
func NewStr(s string) *Lex {
	return New(bufio.NewReader(strings.NewReader(s)))
}

var tokenNull = Token{}
var tokenEOF = Token{ID: TkEOF, Value: "EOF"}

//var tokenFIXME = Token{ID: TkFIXME, Value: "FIXME-WRITEME"}
//var tokenErrInput = Token{ID: TkErrInput, Value: "ERROR-INPUT"}
//var tokenErrInternal = Token{ID: TkErrInternal, Value: "ERROR-INTERNAL"}
//var tokenErrLarge = Token{ID: TkErrLarge, Value: "ERROR-LARGE-TOKEN"}

func (l *Lex) returnTokenEOF() Token {
	l.eofSent = true
	return tokenEOF
}

// Next gets next token.
// Will return EOF token unless Lex.HasToken() is true.
// Check for EOF token with Token.IsEOF() method.
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

		l.lineOffset++

		t := l.match(b)
		switch t.ID {
		case TkNull:
			continue
		case TkEOL:
			l.lineOffset = 0
			l.lineCount++
		}

		return t
	}

	// never reached
}

func (l *Lex) foundErrorInput(err error) Token {
	return l.saveLocationEmpty(Token{ID: TkErrInput, Value: fmt.Sprintf("ERROR-INPUT: after [%s]: %s", l.buf.String(), err.Error())})
}

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
	TkFloat        = iota // .digits | digits. | digits.digits

	TkKeywordImp = iota // IMP
	TkKeywordEqv = iota // EQV
	TkKeywordXor = iota // XOR
	TkKeywordOr  = iota // OR
	TkKeywordAnd = iota // AND
	TkKeywordNot = iota // NOT

	TkEqual   = iota // Equal
	TkUnequal = iota // Unequal <>
	TkLT      = iota // <
	TkGT      = iota // >
	TkLE      = iota // <=
	TkGE      = iota // >=

	TkPlus       = iota // +
	TkMinus      = iota // -
	TkKeywordMod = iota // MOD
	TkBackSlash  = iota // \
	TkMult       = iota // *
	TkDiv        = iota // /
	TkPow        = iota // ^
	UnaryPlus    = iota // fictitious
	UnaryMinus   = iota // fictitious

	TkKeywordCls    = iota // CLS
	TkKeywordCont   = iota // CONT
	TkKeywordElse   = iota // ELSE
	TkKeywordEnd    = iota // END
	TkKeywordFor    = iota // FOR
	TkKeywordGosub  = iota // GOSUB
	TkKeywordGoto   = iota // GOTO
	TkKeywordIf     = iota // IF
	TkKeywordInput  = iota // INPUT
	TkKeywordInt    = iota // INT
	TkKeywordLeft   = iota // LEFT$
	TkKeywordLen    = iota // LEN
	TkKeywordLet    = iota // LET
	TkKeywordList   = iota // LIST
	TkKeywordLoad   = iota // LOAD
	TkKeywordNext   = iota // NEXT
	TkKeywordOn     = iota // ON
	TkKeywordPrint  = iota // PRINT
	TkKeywordRem    = iota // REM
	TkKeywordReturn = iota // RETURN
	TkKeywordRnd    = iota // RND
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
	{TkKeywordAnd, "AND"},
	{TkKeywordCls, "CLS"},
	{TkKeywordCont, "CONT"},
	{TkKeywordElse, "ELSE"},
	{TkKeywordEnd, "END"},
	{TkKeywordEqv, "EQV"},
	{TkKeywordFor, "FOR"},
	{TkKeywordGosub, "GOSUB"},
	{TkKeywordGoto, "GOTO"},
	{TkKeywordIf, "IF"},
	{TkKeywordImp, "IMP"},
	{TkKeywordInput, "INPUT"},
	{TkKeywordInt, "INT"},
	{TkKeywordLeft, "LEFT$"},
	{TkKeywordLen, "LEN"},
	{TkKeywordLet, "LET"},
	{TkKeywordList, "LIST"},
	{TkKeywordLoad, "LOAD"},
	{TkKeywordMod, "MOD"},
	{TkKeywordNext, "NEXT"},
	{TkKeywordNot, "NOT"},
	{TkKeywordOn, "ON"},
	{TkKeywordOr, "OR"},
	{TkKeywordPrint, "PRINT"},
	{TkKeywordRem, "REM"},
	{TkKeywordReturn, "RETURN"},
	{TkKeywordRnd, "RND"},
	{TkKeywordRun, "RUN"},
	{TkKeywordSave, "SAVE"},
	{TkKeywordStep, "STEP"},
	{TkKeywordStop, "STOP"},
	{TkKeywordSystem, "SYSTEM"},
	{TkKeywordThen, "THEN"},
	{TkKeywordTime, "TIME$"},
	{TkKeywordTo, "TO"},
	{TkKeywordXor, "XOR"},
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
	"FLOAT",

	"IMP",
	"EQV",
	"XOR",
	"OR",
	"AND",
	"NOT",

	"EQUAL",
	"UNEQUAL",
	"LT",
	"GT",
	"LE",
	"GE",

	"PLUS",
	"MINUS",
	"MOD",
	"BACK-SLASH",
	"MULT",
	"DIV",
	"POW",
	"UNARY-PLUS",
	"UNARY-MINUS",

	"CLS",
	"CONT",
	"ELSE",
	"END",
	"FOR",
	"GOSUB",
	"GOTO",
	"IF",
	"INPUT",
	"INT",
	"LEFT$",
	"LEN",
	"LET",
	"LIST",
	"LOAD",
	"NEXT",
	"ON",
	"PRINT",
	"REM",
	"RETURN",
	"RND",
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
	r                io.ByteScanner
	eofSeen          bool         // hit EOF?
	eofSent          bool         // delivered EOF?
	broken           bool         // hit error?
	buf              bytes.Buffer // current token
	rawLine          bytes.Buffer // current raw line
	pendingLineReset bool
	state            int
	lineCount        int
	lineOffset       int
}

// Line returns current line count.
func (l *Lex) Line() int {
	return l.lineCount
}

// Column returns current column offset within line.
func (l *Lex) Column() int {
	return l.lineOffset
}

// RawLine returns current raw line.
func (l *Lex) RawLine() string {
	return l.rawLine.String()
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

	if l.pendingLineReset {
		// after returning EOL, on reentrance, we ought to discard full raw line
		l.rawLine.Reset()
		l.pendingLineReset = false
	}

	for {
		b, errRead := l.r.ReadByte()
		switch errRead {
		case nil:
		case io.EOF:
			return l.foundEOF()
		default:
			return l.saveLocationEmpty(Token{ID: TkErrInput, Value: fmt.Sprintf("ERROR-INPUT: after [%s]: %v", l.buf.String(), errRead)})
		}

		if errRaw := l.rawLine.WriteByte(b); errRaw != nil {
			return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: %v", errRaw)})
		}

		l.lineOffset++

		t := l.match(b)
		switch t.ID {
		case TkNull:
			continue
		case TkEOL:
			l.lineOffset = 0
			l.lineCount++
			l.pendingLineReset = true
		}

		return t
	}

	// never reached
}

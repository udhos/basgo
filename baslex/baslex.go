package baslex

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
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

	TkHash         = iota // Colon #
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
	TkNumberHex    = iota // Number &H[0-9]+
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

	TkKeywordAbs         = iota // ABS
	TkKeywordAsc         = iota // ASC
	TkKeywordBeep        = iota // BEEP
	TkKeywordChain       = iota // CHAIN
	TkKeywordChr         = iota // CHR$
	TkKeywordClear       = iota // CLEAR
	TkKeywordClose       = iota // CLOSE
	TkKeywordCls         = iota // CLS
	TkKeywordColor       = iota // COLOR
	TkKeywordCommon      = iota // COMMON
	TkKeywordCont        = iota // CONT
	TkKeywordCos         = iota // COS
	TkKeywordCsrlin      = iota // CSRLIN
	TkKeywordData        = iota // DATA
	TkKeywordDate        = iota // DATE$
	TkKeywordDef         = iota // DEF
	TkKeywordDefdbl      = iota // DEFDBL
	TkKeywordDefint      = iota // DEFINT
	TkKeywordDefsng      = iota // DEFSNG
	TkKeywordDefstr      = iota // DEFSTR
	TkKeywordDim         = iota // DIM
	TkKeywordElse        = iota // ELSE
	TkKeywordEnd         = iota // END
	TkKeywordEnviron     = iota // ENVIRON
	TkKeywordEnvironFunc = iota // ENVIRON$
	TkKeywordEof         = iota // EOF
	TkKeywordErase       = iota // ERASE
	TkKeywordError       = iota // ERROR
	TkKeywordFiles       = iota // FILES
	TkKeywordFix         = iota // FIX
	TkKeywordFor         = iota // FOR
	TkKeywordGodecl      = iota // _GODECL
	TkKeywordGofunc      = iota // _GOFUNC
	TkKeywordGoimport    = iota // _GOIMPORT
	TkKeywordGoproc      = iota // _GOPROC
	TkKeywordGosub       = iota // GOSUB
	TkKeywordGoto        = iota // GOTO
	TkKeywordIf          = iota // IF
	TkKeywordInkey       = iota // INKEY$
	TkKeywordInput       = iota // INPUT
	TkKeywordInputFunc   = iota // INPUT$
	TkKeywordInstr       = iota // INSTR
	TkKeywordInt         = iota // INT
	TkKeywordKey         = iota // KEY
	TkKeywordKill        = iota // KILL
	TkKeywordLeft        = iota // LEFT$
	TkKeywordLen         = iota // LEN
	TkKeywordLet         = iota // LET
	TkKeywordLine        = iota // LINE
	TkKeywordList        = iota // LIST
	TkKeywordLoad        = iota // LOAD
	TkKeywordLocate      = iota // LOCATE
	TkKeywordLof         = iota // LOF
	TkKeywordMid         = iota // MID$
	TkKeywordNext        = iota // NEXT
	TkKeywordNew         = iota // NEW
	TkKeywordOff         = iota // OFF
	TkKeywordOn          = iota // ON
	TkKeywordOpen        = iota // OPEN
	TkKeywordOption      = iota // OPTION
	TkKeywordPeek        = iota // PEEK
	TkKeywordPlay        = iota // PLAY
	TkKeywordPoke        = iota // POKE
	TkKeywordPos         = iota // POS
	TkKeywordPrint       = iota // PRINT
	TkKeywordRandomize   = iota // RANDOMIZE
	TkKeywordRead        = iota // READ
	TkKeywordRem         = iota // REM
	TkKeywordReset       = iota // RESET
	TkKeywordRestore     = iota // RESTORE
	TkKeywordResume      = iota // RESUME
	TkKeywordReturn      = iota // RETURN
	TkKeywordRight       = iota // RIGHT$
	TkKeywordRnd         = iota // RND
	TkKeywordRun         = iota // RUN
	TkKeywordSave        = iota // SAVE
	TkKeywordScreen      = iota // SCREEN
	TkKeywordSgn         = iota // SGN
	TkKeywordSin         = iota // SIN
	TkKeywordSound       = iota // SOUND
	TkKeywordSpace       = iota // SPACE$
	TkKeywordSpc         = iota // SPC
	TkKeywordSqr         = iota // SQR
	TkKeywordStep        = iota // STEP
	TkKeywordStop        = iota // STOP
	TkKeywordStr         = iota // STR$
	TkKeywordString      = iota // STRING$
	TkKeywordSwap        = iota // SWAP
	TkKeywordSystem      = iota // SYSTEM
	TkKeywordTab         = iota // TAB
	TkKeywordTan         = iota // TAN
	TkKeywordThen        = iota // THEN
	TkKeywordTime        = iota // TIME$
	TkKeywordTimer       = iota // TIMER
	TkKeywordTo          = iota // TO
	TkKeywordUsing       = iota // USING
	TkKeywordVal         = iota // VAL
	TkKeywordView        = iota // VIEW
	TkKeywordWend        = iota // WEND
	TkKeywordWhile       = iota // WHILE
	TkKeywordWidth       = iota // WIDTH

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
	{TkKeywordAbs, "ABS"},
	{TkKeywordAnd, "AND"},
	{TkKeywordAsc, "ASC"},
	{TkKeywordBeep, "BEEP"},
	{TkKeywordChain, "CHAIN"},
	{TkKeywordChr, "CHR$"},
	{TkKeywordClear, "CLEAR"},
	{TkKeywordClose, "CLOSE"},
	{TkKeywordCls, "CLS"},
	{TkKeywordColor, "COLOR"},
	{TkKeywordCommon, "COMMON"},
	{TkKeywordCont, "CONT"},
	{TkKeywordCos, "COS"},
	{TkKeywordCsrlin, "CSRLIN"},
	{TkKeywordData, "DATA"},
	{TkKeywordDate, "DATE$"},
	{TkKeywordDef, "DEF"},
	{TkKeywordDefdbl, "DEFDBL"},
	{TkKeywordDefint, "DEFINT"},
	{TkKeywordDefsng, "DEFSNG"},
	{TkKeywordDefstr, "DEFSTR"},
	{TkKeywordDim, "DIM"},
	{TkKeywordElse, "ELSE"},
	{TkKeywordEnd, "END"},
	{TkKeywordEnviron, "ENVIRON"},
	{TkKeywordEnvironFunc, "ENVIRON$"},
	{TkKeywordEof, "EOF"},
	{TkKeywordEqv, "EQV"},
	{TkKeywordErase, "ERASE"},
	{TkKeywordError, "ERROR"},
	{TkKeywordFiles, "FILES"},
	{TkKeywordFix, "FIX"},
	{TkKeywordFor, "FOR"},
	{TkKeywordGodecl, "_GODECL"},
	{TkKeywordGofunc, "_GOFUNC"},
	{TkKeywordGoimport, "_GOIMPORT"},
	{TkKeywordGoproc, "_GOPROC"},
	{TkKeywordGosub, "GOSUB"},
	{TkKeywordGoto, "GOTO"},
	{TkKeywordIf, "IF"},
	{TkKeywordImp, "IMP"},
	{TkKeywordInkey, "INKEY$"},
	{TkKeywordInput, "INPUT"},
	{TkKeywordInputFunc, "INPUT$"},
	{TkKeywordInstr, "INSTR"},
	{TkKeywordInt, "INT"},
	{TkKeywordKey, "KEY"},
	{TkKeywordKill, "KILL"},
	{TkKeywordLeft, "LEFT$"},
	{TkKeywordLen, "LEN"},
	{TkKeywordLet, "LET"},
	{TkKeywordLine, "LINE"},
	{TkKeywordList, "LIST"},
	{TkKeywordLoad, "LOAD"},
	{TkKeywordLocate, "LOCATE"},
	{TkKeywordLof, "LOF"},
	{TkKeywordMid, "MID$"},
	{TkKeywordMod, "MOD"},
	{TkKeywordNext, "NEXT"},
	{TkKeywordNew, "NEW"},
	{TkKeywordNot, "NOT"},
	{TkKeywordOff, "OFF"},
	{TkKeywordOn, "ON"},
	{TkKeywordOpen, "OPEN"},
	{TkKeywordOption, "OPTION"},
	{TkKeywordOr, "OR"},
	{TkKeywordPeek, "PEEK"},
	{TkKeywordPlay, "PLAY"},
	{TkKeywordPoke, "POKE"},
	{TkKeywordPos, "POS"},
	{TkKeywordPrint, "PRINT"},
	{TkKeywordRandomize, "RANDOMIZE"},
	{TkKeywordRead, "READ"},
	{TkKeywordRem, "REM"},
	{TkKeywordReset, "RESET"},
	{TkKeywordRestore, "RESTORE"},
	{TkKeywordResume, "RESUME"},
	{TkKeywordReturn, "RETURN"},
	{TkKeywordRight, "RIGHT$"},
	{TkKeywordRnd, "RND"},
	{TkKeywordRun, "RUN"},
	{TkKeywordSave, "SAVE"},
	{TkKeywordScreen, "SCREEN"},
	{TkKeywordSgn, "SGN"},
	{TkKeywordSin, "SIN"},
	{TkKeywordSound, "SOUND"},
	{TkKeywordSpace, "SPACE$"},
	{TkKeywordSpc, "SPC"},
	{TkKeywordSqr, "SQR"},
	{TkKeywordStep, "STEP"},
	{TkKeywordStop, "STOP"},
	{TkKeywordStr, "STR$"},
	{TkKeywordString, "STRING$"},
	{TkKeywordSwap, "SWAP"},
	{TkKeywordSystem, "SYSTEM"},
	{TkKeywordTab, "TAB"},
	{TkKeywordTan, "TAN"},
	{TkKeywordThen, "THEN"},
	{TkKeywordTime, "TIME$"},
	{TkKeywordTimer, "TIMER"},
	{TkKeywordTo, "TO"},
	{TkKeywordUsing, "USING"},
	{TkKeywordVal, "VAL"},
	{TkKeywordView, "VIEW"},
	{TkKeywordWend, "WEND"},
	{TkKeywordWhile, "WHILE"},
	{TkKeywordWidth, "WIDTH"},
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

	"HASH",
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
	"NUMBER-HEX",
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

	"ABS",
	"ASC",
	"BEEP",
	"CHAIN",
	"CHR$",
	"CLEAR",
	"CLOSE",
	"CLS",
	"COLOR",
	"COMMON",
	"CONT",
	"COS",
	"CSRLIN",
	"DATA",
	"DATE",
	"DEF",
	"DEFDBL",
	"DEFINT",
	"DEFSNG",
	"DEFSTR",
	"DIM",
	"ELSE",
	"END",
	"ENVIRON",
	"ENVIRON$",
	"EOF",
	"ERASE",
	"ERROR",
	"FILES",
	"FIX",
	"FOR",
	"_GODECL",
	"_GOFUNC",
	"_GOIMPORT",
	"_GOPROC",
	"GOSUB",
	"GOTO",
	"IF",
	"INKEY$",
	"INPUT",
	"INPUT$",
	"INSTR",
	"INT",
	"KEY",
	"KILL",
	"LEFT$",
	"LEN",
	"LET",
	"LINE",
	"LIST",
	"LOAD",
	"LOCATE",
	"LOF",
	"MID$",
	"NEXT",
	"NEW",
	"OFF",
	"ON",
	"OPEN",
	"OPTION",
	"PEEK",
	"PLAY",
	"POKE",
	"POS",
	"PRINT",
	"RANDOMIZE",
	"READ",
	"REM",
	"RESET",
	"RESTORE",
	"RESUME",
	"RETURN",
	"RIGHT$",
	"RND",
	"RUN",
	"SAVE",
	"SCREEN",
	"SGN",
	"SIN",
	"SOUND",
	"SPACE$",
	"SPC",
	"SQR",
	"STEP",
	"STOP",
	"STR$",
	"STRING$",
	"SWAP",
	"SYSTEM",
	"TAB",
	"TAN",
	"THEN",
	"TIME$",
	"TIMER",
	"TO",
	"USING",
	"VAL",
	"VIEW",
	"WEND",
	"WHILE",
	"WIDTH",

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
	data             int // support DATA unquoted string
	countRead        int
	debug            bool
	countToken       int
}

const (
	DataOff         = iota
	DataBeforeValue = iota
	DataAfterValue  = iota
)

// Line returns current line count.
func (l *Lex) Line() int {
	return l.lineCount
}

// Column returns current column offset within line.
func (l *Lex) Column() int {
	return l.lineOffset
}

// Offset returns total read bytes.
func (l *Lex) Offset() int {
	return l.countRead
}

// RawLine returns current raw line.
func (l *Lex) RawLine() string {
	return l.rawLine.String()
}

// New creates a Lex object
func New(input io.ByteScanner, debug bool) *Lex {
	return &Lex{r: input, lineCount: 1, debug: debug}
}

// NewStr creates a Lex object from string
func NewStr(s string, debug bool) *Lex {
	return New(bufio.NewReader(strings.NewReader(s)), debug)
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
	t := l.nextToken()

	if l.debug {
		l.countToken++
		s := TokenString("Lex.Next", t, l.Offset())
		log.Printf(s+" countToken=%d", l.countToken)
	}

	return t
}

func TokenString(label string, t Token, offset int) string {
	return fmt.Sprintf("%-12s: line=%02d col=%02d offset=%d id=%03d %-s [%-s]", label, t.LineCount, t.LineOffset, offset, t.ID, t.Type(), t.Value)
}

func (l *Lex) nextToken() Token {
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
			if l.debug {
				log.Printf("Lex.findToken: io.EOF from ReadByte()")
			}
			return l.foundEOF()
		default:
			return l.saveLocationEmpty(Token{ID: TkErrInput, Value: fmt.Sprintf("ERROR-INPUT: after [%s]: %v", l.buf.String(), errRead)})
		}

		l.countRead++

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
			l.data = DataOff // disable DATA unquoted string
		case TkColon:
			l.data = DataOff // disable DATA unquoted string
		}

		if t.ID != TkKeywordData && t.ID != TkComma && l.data > DataOff {
			// data value returned
			//log.Printf("token[%s] switch to DataAfterValue", t.Value)
			l.data = DataAfterValue
		}

		return t
	}

	// never reached
}

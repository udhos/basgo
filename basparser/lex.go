package basparser

import (
	"io"
	"log"

	"github.com/udhos/basgo/baslex"
)

func NewInputLex(input io.ByteScanner, debug bool) *InputLex {
	return &InputLex{lex: baslex.New(input), debug: debug}
}

type InputLex struct {
	lex              *baslex.Lex
	debug            bool
	syntaxErrorCount int
	lastToken        baslex.Token // save token for parser error reporting
}

// Reduced is hook for recording a reduction.
// https://godoc.org/modernc.org/goyacc
// Optionally the argument to yyParse may implement the following interface:
// Reduced(rule, state int, lval *yySymType) (stop bool) // Client should copy *lval.
func (l *InputLex) Reduced(rule, state int, lval *InputSymType) (stop bool) {
	if !l.debug {
		return false
	}
	log.Printf("Reduced: rule=%d state=%d", rule, state)
	return false
}

func (l *InputLex) Errors() int {
	return l.syntaxErrorCount
}

func (l *InputLex) Lex(lval *InputSymType) int {

	if !l.lex.HasToken() {
		return 0 // 0 means real EOF for the parser
	}

	t := l.lex.Next()

	l.lastToken = t // save token for parser error reporting

	// ATTENTION: t.ID is in lex token space

	id := parserToken(t.ID) // convert lex ID to parser ID

	// ATTENTION: id is in parser token space

	if l.debug {
		log.Printf("InputLex.Lex: %s [%s] basicLine=%s line=%d col=%d offset=%d\n", t.Type(), t.Value, lastLineNum, l.lex.Line(), l.lex.Column(), l.lex.Offset())
	}

	// need to store values only for some terminals
	// when a parser rule action need to consume the value
	// for example: ident, literals (number, string)
	switch id {
	case TkKeywordRem:
		lval.typeRem = t.Value
	case TkCommentQ:
		lval.typeRem = t.Value
	case TkString:
		lval.typeString = t.Value
	case TkNumber:
		lval.typeNumber = t.Value
	case TkNumberHex:
		lval.typeNumber = t.Value
	case TkFloat:
		lval.typeFloat = t.Value
	case TkIdentifier:
		lval.typeIdentifier = t.Value
	case TkEOL:
		lval.typeRawLine = l.lex.RawLine()
	case TkEOF:
		lval.typeRawLine = l.lex.RawLine()
	}

	return id
}

func (l *InputLex) Error(s string) {
	l.syntaxErrorCount++
	log.Printf("InputLex.Error: PARSER: %s", s)
	log.Printf("InputLex.Error: PARSER: last token: %s [%s]", l.lastToken.Type(), l.lastToken.Value)
	log.Printf("InputLex.Error: PARSER: basicLine=%s inputLine=%d column=%d totalOffset=%d", lastLineNum, l.lex.Line(), l.lex.Column(), l.lex.Offset())
	log.Printf("InputLex.Error: PARSER: errors=%d", l.syntaxErrorCount)
}

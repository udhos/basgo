package baslex

import (
	"fmt"
	//"log"
	"strings"
)

/*
Keep states in sync:

(1) const states
(2) var tabState
(3) func foundEOF()
*/

// (1) const states
const (
	stBlank      = iota
	stCommentQ   = iota
	stCommentRem = iota
	stString     = iota
	stNumber     = iota
	stFloat      = iota
	stName       = iota
	stLT         = iota
	stGT         = iota
)

// "CLS" => TkKeywordCls
func findKeyword(name string) int {
	nameUp := strings.ToUpper(name)
	for _, k := range tabKeywords {
		if nameUp == k.Name {
			return k.TokenID
		}
	}
	return TkIdentifier
}

type funcState func(l *Lex, b byte) Token

// (2) var tabState
var tabState = []funcState{
	matchBlank,
	matchCommentQ,
	matchCommentRem,
	matchString,
	matchNumber,
	matchFloat,
	matchName,
	matchLT,
	matchGT,
}

func (l *Lex) saveLocation(t Token, size int) Token {
	t.LineCount = l.lineCount          // save line
	t.LineOffset = l.lineOffset - size // save offset
	return t
}

func (l *Lex) saveLocationEmpty(t Token) Token {
	return l.saveLocation(t, 0)
}

func (l *Lex) saveLocationValue(t Token) Token {
	return l.saveLocation(t, len(t.Value))
}

func (l *Lex) consume(t Token) Token {
	t.Value = l.buf.String() // save value

	t = l.saveLocationValue(t)

	//log.Printf("consume: [%s]", t.Value)

	l.buf.Reset()
	return t
}

func (l *Lex) consumeName() Token {
	name := l.buf.String()
	id := findKeyword(name)
	return l.consume(Token{ID: id})
}

// (3) func foundEOF()
func (l *Lex) foundEOF() Token {

	l.eofSeen = true

	switch l.state {
	case stBlank:
		return l.saveLocationEmpty(l.returnTokenEOF())
	case stCommentQ:
		return l.consume(Token{ID: TkCommentQ})
	case stCommentRem:
		return l.consume(Token{ID: TkKeywordRem})
	case stString:
		return l.consume(Token{ID: TkString})
	case stNumber:
		return l.consume(Token{ID: TkNumber})
	case stFloat:
		return l.consume(Token{ID: TkFloat})
	case stName:
		return l.consumeName()
	case stLT:
		return l.consume(Token{ID: TkLT})
	case stGT:
		return l.consume(Token{ID: TkGT})
	}

	return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL:foundEOF: bad state=%d", l.state)})
}

func (l *Lex) match(b byte) Token {

	if l.state < 0 || l.state >= len(tabState) {
		return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: match bad state=%d", l.state)})
	}

	return tabState[l.state](l, b)
}

func (l *Lex) save(b byte) Token {
	if errSave := l.buf.WriteByte(b); errSave != nil {
		return l.saveLocationEmpty(Token{ID: TkErrLarge, Value: fmt.Sprintf("ERROR-LARGE-TOKEN: %s", errSave)})
	}
	return tokenNull
}

func matchBlank(l *Lex, b byte) Token {

	switch {
	case eol(b):
		return l.saveLocationEmpty(Token{ID: TkEOL, Value: "EOL"})
	case blank(b):
		return tokenNull
	case b == '\'':
		l.state = stCommentQ
		return l.save(b)
	case b == '"':
		l.state = stString
		return l.save(b)
	case b == '+':
		return l.saveLocationValue(Token{ID: TkPlus, Value: "+"})
	case b == '-':
		return l.saveLocationValue(Token{ID: TkMinus, Value: "-"})
	case b == '*':
		return l.saveLocationValue(Token{ID: TkMult, Value: "*"})
	case b == '/':
		return l.saveLocationValue(Token{ID: TkDiv, Value: "/"})
	case b == '\\':
		return l.saveLocationValue(Token{ID: TkBackSlash, Value: "\\"})
	case b == ':':
		return l.saveLocationValue(Token{ID: TkColon, Value: ":"})
	case b == '=':
		return l.saveLocationValue(Token{ID: TkEqual, Value: "="})
	case b == ',':
		return l.saveLocationValue(Token{ID: TkComma, Value: ","})
	case b == ';':
		return l.saveLocationValue(Token{ID: TkSemicolon, Value: ";"})
	case b == '(':
		return l.saveLocationValue(Token{ID: TkParLeft, Value: "("})
	case b == ')':
		return l.saveLocationValue(Token{ID: TkParRight, Value: ")"})
	case b == '[':
		return l.saveLocationValue(Token{ID: TkBracketLeft, Value: "["})
	case b == ']':
		return l.saveLocationValue(Token{ID: TkBracketRight, Value: "]"})
	case b == '<':
		l.state = stLT
		return l.save(b)
	case b == '>':
		l.state = stGT
		return l.save(b)
	case digit(b):
		l.state = stNumber
		return l.save(b)
	case b == '.':
		l.state = stFloat
		return l.save(b)
	case letter(b):
		l.state = stName
		return l.save(b)
	}

	invalid := fmt.Sprintf("INVALID: byte=%d: '%c'", b, b)
	//log.Printf("matchBlank: %s", invalid)
	return l.saveLocationEmpty(Token{ID: TkErrInvalid, Value: invalid})
}

func letter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func digit(b byte) bool {
	return b >= '0' && b <= '9'
}

func blank(b byte) bool {
	return b == ' ' || b == '\t'
}

func eol(b byte) bool {
	return b == '\r' || b == '\n'
}

// push back byte
func unread(l *Lex) error {
	l.rawLine.Truncate(l.rawLine.Len() - 1) // unwrite byte from raw line buf

	errInputUnread := l.r.UnreadByte()
	if errInputUnread == nil {
		l.lineOffset--
	}
	return errInputUnread
}

func matchCommentQ(l *Lex, b byte) Token {

	switch {
	case eol(b):
		// push back EOL
		if errUnread := unread(l); errUnread != nil {
			return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)})
		}
		l.state = stBlank // blank state will deliver EOL

		return l.consume(Token{ID: TkCommentQ})
	}

	return l.save(b)
}

func matchCommentRem(l *Lex, b byte) Token {

	switch {
	case eol(b):
		// push back EOL
		if errUnread := unread(l); errUnread != nil {
			return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)})
		}
		l.state = stBlank // blank state will deliver EOL

		return l.consume(Token{ID: TkKeywordRem})
	}

	return l.save(b)
}

func matchString(l *Lex, b byte) Token {

	switch {
	case b == '"':
		l.state = stBlank
		// attention: must save byte before extracting value for new token
		if t := l.save(b); t.ID != TkNull {
			return t
		}
		return l.consume(Token{ID: TkString})
	case eol(b):
		// push back EOL
		if errUnread := unread(l); errUnread != nil {
			return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)})
		}
		l.state = stBlank // blank state will deliver EOL

		return l.consume(Token{ID: TkString})
	}

	return l.save(b)
}

func matchNumber(l *Lex, b byte) Token {

	switch {
	case digit(b):
		return l.save(b)
	case b == '.':
		l.state = stFloat // switch from number to float
		return l.save(b)
	}

	// push back non-digit
	if errUnread := unread(l); errUnread != nil {
		return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)})
	}
	l.state = stBlank // blank state will deliver next token

	return l.consume(Token{ID: TkNumber})
}

func matchFloat(l *Lex, b byte) Token {

	if digit(b) {
		return l.save(b)
	}

	// push back non-digit
	if errUnread := unread(l); errUnread != nil {
		return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)})
	}
	l.state = stBlank // blank state will deliver next token

	return l.consume(Token{ID: TkFloat})
}

func matchName(l *Lex, b byte) Token {

	switch {

	case letter(b) || digit(b) || b == '.':
		return l.save(b)

	case b == '$' || b == '%' || b == '!' || b == '#':
		l.state = stBlank
		// attention: must save byte before extracting value for new token
		if t := l.save(b); t.ID != TkNull {
			return t
		}
		return l.consumeName()
	}

	// found name

	// push back non-name byte
	if errUnread := unread(l); errUnread != nil {
		return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)})
	}
	l.state = stBlank // blank state will deliver next token

	// trap special keyword REM

	name := l.buf.String()
	id := findKeyword(name)
	if id == TkKeywordRem {
		l.state = stCommentRem
		return tokenNull // keep matching REM
	}

	return l.consumeName()
}

func matchLT(l *Lex, b byte) Token {

	switch b {
	case '>':
		l.state = stBlank
		// attention: must save byte before extracting value for new token
		if t := l.save(b); t.ID != TkNull {
			return t
		}
		return l.consume(Token{ID: TkUnequal})
	case '=':
		l.state = stBlank
		// attention: must save byte before extracting value for new token
		if t := l.save(b); t.ID != TkNull {
			return t
		}
		return l.consume(Token{ID: TkLE})
	}

	// push back
	if errUnread := unread(l); errUnread != nil {
		return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)})
	}
	l.state = stBlank // blank state will deliver next token

	return l.consume(Token{ID: TkLT})
}

func matchGT(l *Lex, b byte) Token {

	switch b {
	case '=':
		l.state = stBlank
		// attention: must save byte before extracting value for new token
		if t := l.save(b); t.ID != TkNull {
			return t
		}
		return l.consume(Token{ID: TkGE})
	}

	// push back
	if errUnread := unread(l); errUnread != nil {
		return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)})
	}
	l.state = stBlank // blank state will deliver next token

	return l.consume(Token{ID: TkGT})
}

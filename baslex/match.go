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
	stBlank                = iota
	stCommentQ             = iota
	stCommentRem           = iota
	stString               = iota
	stStringUnquoted       = iota
	stStringUnquotedFinish = iota
	stNumber               = iota
	stFloat                = iota
	stFloatE               = iota
	stFloatEE              = iota
	stFloatEEE             = iota
	stName                 = iota
	stLT                   = iota
	stGT                   = iota
	stEqual                = iota
	stAmpersand            = iota
	stAmperH               = iota
	stHex                  = iota
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
	matchStringUnquoted,
	matchStringUnquotedFinish,
	matchNumber,
	matchFloat,
	matchFloatE,
	matchFloatEE,
	matchFloatEEE,
	matchName,
	matchLT,
	matchGT,
	matchEqual,
	matchAmpersand,
	matchAmperH,
	matchHex,
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

func quoteString(t Token) Token {
	t.Value = `"` + t.Value + `"`
	return t
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
	case stStringUnquoted:
		return quoteString(l.consume(Token{ID: TkString}))
	case stNumber:
		return l.consume(Token{ID: TkNumber})
	case stFloat, stFloatE, stFloatEEE:
		return l.consume(Token{ID: TkFloat})
	case stName:
		return l.consumeName()
	case stLT:
		return l.consume(Token{ID: TkLT})
	case stGT:
		return l.consume(Token{ID: TkGT})
	case stEqual:
		return l.consume(Token{ID: TkEqual})
	case stHex:
		return l.consume(Token{ID: TkNumberHex})

		/*
			// Some states do not support EOF:

			case stFloatEE:
			case stAmpersend:
			case stAmperH:
		*/
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

/*
func finishStringUnquoted(l *Lex) Token {
	// push back
	if errUnread := unread(l); errUnread != nil {
		return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)})
	}

	return quoteString(l.consume(Token{ID: TkString}))
}
*/

func matchBlankData(l *Lex, b byte) Token {

	//log.Printf("matchBlankData: byte=%c", b)

	switch {
	case eol(b):
		l.data = false
		// push back
		if errUnread := unread(l); errUnread != nil {
			return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)})
		}
		l.state = stStringUnquoted
		return tokenNull
	case blank(b):
		return tokenNull
	case b == ',':
		// push back
		if errUnread := unread(l); errUnread != nil {
			return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)})
		}
		l.state = stStringUnquoted
		return tokenNull
	case b == ':':
		l.data = false
		// push back
		if errUnread := unread(l); errUnread != nil {
			return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)})
		}
		l.state = stStringUnquoted
		return tokenNull
	case b == '"':
		l.state = stString
		return l.save(b)
	case digit(b):
		l.state = stNumber
		return l.save(b)
	}

	// anything else is unquoted string

	l.state = stStringUnquoted
	return l.save(b)
}

func matchBlank(l *Lex, b byte) Token {

	if l.data {
		return matchBlankData(l, b)
	}

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
	case b == '^':
		return l.saveLocationValue(Token{ID: TkPow, Value: "^"})
	case b == '\\':
		return l.saveLocationValue(Token{ID: TkBackSlash, Value: "\\"})
	case b == '#':
		return l.saveLocationValue(Token{ID: TkHash, Value: "#"})
	case b == ':':
		return l.saveLocationValue(Token{ID: TkColon, Value: ":"})
	case b == '=':
		l.state = stEqual
		return l.save(b)
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
	case b == '_', letter(b): // support keyword with '_' prefix for _GOFUNC
		l.state = stName
		return l.save(b)
	case b == '&':
		l.state = stAmpersand
		return l.save(b)
	}

	invalid := fmt.Sprintf("matchBlank: INVALID: byte=%d: '%c'", b, b)
	return l.saveLocationEmpty(Token{ID: TkErrInvalid, Value: invalid})
}

func letter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func digit(b byte) bool {
	return b >= '0' && b <= '9'
}

func hexdigit(b byte) bool {
	return digit(b) || (b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F')
}

func blank(b byte) bool {
	return b == ' ' || b == '\t' || b == '\r'
}

func eol(b byte) bool {
	return b == '\n'
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

func matchStringUnquoted(l *Lex, b byte) Token {

	switch {
	case b == ',', b == ':', eol(b):
		// push back , : EOL
		if errUnread := unread(l); errUnread != nil {
			return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)})
		}
		l.state = stStringUnquotedFinish // will deliver , : EOL

		return quoteString(l.consume(Token{ID: TkString}))
	}

	return l.save(b)
}

func matchStringUnquotedFinish(l *Lex, b byte) Token {

	switch {
	case eol(b):
		l.state = stBlank
		return l.saveLocationEmpty(Token{ID: TkEOL, Value: "EOL"})
	case b == ':':
		l.state = stBlank
		return l.saveLocationValue(Token{ID: TkColon, Value: ":"})
	case b == ',':
		l.state = stBlank
		return l.saveLocationValue(Token{ID: TkComma, Value: ","})
	}

	invalid := fmt.Sprintf("matchStringUnquotedFinish: INVALID: byte=%d: '%c'", b, b)
	return l.saveLocationEmpty(Token{ID: TkErrInvalid, Value: invalid})
}

func matchNumber(l *Lex, b byte) Token {

	switch {
	case digit(b):
		return l.save(b)
	case b == '.':
		l.state = stFloat // switch from number to float
		return l.save(b)
	case b == 'e', b == 'E':
		l.state = stFloatE // switch from number to floatE
		return l.save(b)
	case b == '!':
		// force float
		l.state = stBlank // blank state will deliver next token
		return l.consume(Token{ID: TkFloat})
	}

	// push back non-digit
	if errUnread := unread(l); errUnread != nil {
		return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)})
	}
	l.state = stBlank // blank state will deliver next token

	return l.consume(Token{ID: TkNumber})
}

func matchFloat(l *Lex, b byte) Token {

	switch {
	case digit(b):
		return l.save(b)
	case b == 'e', b == 'E':
		l.state = stFloatE // switch from float to floatE
		return l.save(b)
	}

	// push back non-digit
	if errUnread := unread(l); errUnread != nil {
		return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)})
	}
	l.state = stBlank // blank state will deliver next token

	return l.consume(Token{ID: TkFloat})
}

// expect digit, -, +
func matchFloatE(l *Lex, b byte) Token {

	switch {
	case digit(b):
		l.state = stFloatEEE // switch from floatE to floatEEE
		return l.save(b)
	case b == '-', b == '+':
		l.state = stFloatEE // switch from floatE to floatEE
		return l.save(b)
	}

	// push back non-digit
	if errUnread := unread(l); errUnread != nil {
		return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)})
	}
	l.state = stBlank // blank state will deliver next token

	return l.consume(Token{ID: TkFloat})
}

// only digit is accepted, non-digit is error
func matchFloatEE(l *Lex, b byte) Token {

	if digit(b) {
		l.state = stFloatEEE // switch to floatEEE
		return l.save(b)
	}

	invalid := fmt.Sprintf("matchFloatEE: INVALID: byte=%d: '%c'", b, b)
	return l.saveLocationEmpty(Token{ID: TkErrInvalid, Value: invalid})
}

// only digit is accepted, non-digit finishes
func matchFloatEEE(l *Lex, b byte) Token {

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

	case b == '#':

		name := l.buf.String()
		id := findKeyword(name)
		switch id {
		case TkKeywordInput, TkKeywordPrint:
			// push back hash
			if errUnread := unread(l); errUnread != nil {
				return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)})
			}
			l.state = stBlank // blank state will deliver next token #

			return l.consumeName() // consume INPUT, PRINT
		}

		l.state = stBlank
		// attention: must save byte before extracting value for new token
		if t := l.save(b); t.ID != TkNull {
			return t
		}
		return l.consumeName() // consume IDENT#

	case b == '$' || b == '%' || b == '!':
		l.state = stBlank
		// attention: must save byte before extracting value for new token
		if t := l.save(b); t.ID != TkNull {
			return t
		}
		return l.consumeName() // consume IDENTx (x: $,%,!)
	}

	// found name

	// push back non-name byte
	if errUnread := unread(l); errUnread != nil {
		return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)})
	}
	l.state = stBlank // blank state will deliver next token

	// trap special keywords REM, DATA

	name := l.buf.String()
	id := findKeyword(name)
	switch id {
	case TkKeywordRem:
		l.state = stCommentRem
		return tokenNull // keep matching REM
	case TkKeywordData:
		l.data = true // support DATA unquoted string
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
	case '<':
		l.state = stBlank
		// attention: must save byte before extracting value for new token
		if t := l.save(b); t.ID != TkNull {
			return t
		}
		return l.consume(Token{ID: TkUnequal})
	}

	// push back
	if errUnread := unread(l); errUnread != nil {
		return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)})
	}
	l.state = stBlank // blank state will deliver next token

	return l.consume(Token{ID: TkGT})
}

func matchEqual(l *Lex, b byte) Token {

	switch b {
	case '>':
		l.state = stBlank
		// attention: must save byte before extracting value for new token
		if t := l.save(b); t.ID != TkNull {
			return t
		}
		return l.consume(Token{ID: TkGE})
	case '<':
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

	return l.consume(Token{ID: TkEqual})
}

// expect only H
func matchAmpersand(l *Lex, b byte) Token {

	switch b {
	case 'h', 'H':
		l.state = stAmperH
		return l.save(b)
	}

	invalid := fmt.Sprintf("matchAmpersand: INVALID: byte=%d: '%c'", b, b)
	return l.saveLocationEmpty(Token{ID: TkErrInvalid, Value: invalid})
}

// expect hexdigit
func matchAmperH(l *Lex, b byte) Token {

	if hexdigit(b) {
		l.state = stHex
		return l.save(b)
	}

	invalid := fmt.Sprintf("matchAmperH: INVALID: byte=%d: '%c'", b, b)
	return l.saveLocationEmpty(Token{ID: TkErrInvalid, Value: invalid})
}

func matchHex(l *Lex, b byte) Token {

	if hexdigit(b) {
		l.state = stHex
		return l.save(b)
	}

	// push back non-digit
	if errUnread := unread(l); errUnread != nil {
		return l.saveLocationEmpty(Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)})
	}
	l.state = stBlank // blank state will deliver next token

	return l.consume(Token{ID: TkNumberHex})
}

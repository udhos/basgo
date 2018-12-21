package baslex

import (
	"fmt"
	"log"
)

const (
	stBlank    = iota
	stCommentQ = iota
	stString   = iota
	stNumber   = iota
	//stName        = iota
)

type funcState func(l *Lex, b byte) Token

var tabState = []funcState{
	matchBlank,
	matchCommentQ,
	matchString,
	matchNumber,
}

func (l *Lex) consume(t Token) Token {
	t.Value = l.buf.String()
	//log.Printf("consume: [%s]", t.Value)
	l.buf.Reset()
	return t
}

func (l *Lex) foundEOF() Token {

	l.eofSeen = true

	switch l.state {
	case stBlank:
		return l.returnTokenEOF()
	case stCommentQ:
		return l.consume(Token{ID: TkCommentQ})
	case stString:
		return l.consume(Token{ID: TkString})
	case stNumber:
		return l.consume(Token{ID: TkNumber})
	}

	return Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: foundEOF bad state=%d", l.state)}
}

func (l *Lex) match(b byte) Token {

	if l.state < 0 || l.state >= len(tabState) {
		return Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: match bad state=%d", l.state)}
	}

	return tabState[l.state](l, b)
}

func (l *Lex) save(b byte) Token {
	if errSave := l.buf.WriteByte(b); errSave != nil {
		return Token{ID: TkErrLarge, Value: fmt.Sprintf("ERROR-LARGE-TOKEN: %s", errSave)}
	}
	return tokenNull
}

func matchBlank(l *Lex, b byte) Token {

	switch {
	case eol(b):
		return Token{ID: TkEOL, Value: "EOL"}
	case blank(b):
		return tokenNull
	case b == '\'':
		l.state = stCommentQ
		return l.save(b)
	case b == '"':
		l.state = stString
		return l.save(b)
	case b == ':':
		return Token{ID: TkColon, Value: ":"}
	case digit(b):
		l.state = stNumber
		return l.save(b)
	}

	log.Printf("matchBlank: FIXME-WRITEME")
	return tokenFIXME
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

func matchCommentQ(l *Lex, b byte) Token {

	switch {
	case eol(b):
		// push back EOL
		if errUnread := l.r.UnreadByte(); errUnread != nil {
			return Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)}
		}
		l.state = stBlank // blank state will deliver EOL

		return l.consume(Token{ID: TkCommentQ})
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
		if errUnread := l.r.UnreadByte(); errUnread != nil {
			return Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)}
		}
		l.state = stBlank // blank state will deliver EOL

		return l.consume(Token{ID: TkString})
	}

	return l.save(b)
}

func matchNumber(l *Lex, b byte) Token {

	if digit(b) {
		return l.save(b)
	}

	// push back non-digit
	if errUnread := l.r.UnreadByte(); errUnread != nil {
		return Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)}
	}
	l.state = stBlank // blank state will deliver next token

	return l.consume(Token{ID: TkNumber})
}

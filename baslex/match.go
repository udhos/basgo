package baslex

import (
	"fmt"
	"log"
)

type funcState func(l *Lex, b byte) (Token, bool)

var tabState = []funcState{
	matchBlank,
	matchCommentQ,
}

func (l *Lex) foundEOF() Token {

	l.eof = true // set EOF, no more tokens

	switch l.state {
	case stBlank:
		return l.returnTokenEOF() // EOF
	case stCommentQ:
		return Token{ID: TkCommentQ, Value: l.buf.String()} // deliver comment q
	}

	return Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: foundEOF bad state=%d", l.state)}
}

func (l *Lex) match(b byte) (Token, bool) {

	if l.state < 0 || l.state >= len(tabState) {
		return Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: match bad state=%d", l.state)}, true
	}

	return tabState[l.state](l, b)
}

func (l *Lex) save(b byte) (Token, bool) {
	if errSave := l.buf.WriteByte(b); errSave != nil {
		return Token{ID: TkErrLarge, Value: fmt.Sprintf("ERROR-LARGE-TOKEN: %s", errSave)}, true
	}
	return tokenNull, false
}

func matchBlank(l *Lex, b byte) (Token, bool) {

	switch {
	case eol(b):
		return Token{ID: TkEOL, Value: "EOL"}, true
	case blank(b):
		return tokenNull, false
	case b == '\'':
		l.state = stCommentQ
		return l.save(b)
	case b == ':':
		return Token{ID: TkColon, Value: ":"}, true
	}

	log.Printf("matchBlank: FIXME-WRITEME")
	return tokenFIXME, true
}

func blank(b byte) bool {
	return b == ' ' || b == '\t'
}

func eol(b byte) bool {
	return b == '\r' || b == '\n'
}

func matchCommentQ(l *Lex, b byte) (Token, bool) {

	switch {
	case eol(b):
		// push back EOL
		if errUnread := l.r.UnreadByte(); errUnread != nil {
			return Token{ID: TkErrInternal, Value: fmt.Sprintf("ERROR-INTERNAL: unread: %s", errUnread)}, true
		}
		l.state = stBlank // blank state will deliver EOL

		return Token{ID: TkCommentQ, Value: l.buf.String()}, true // deliver comment q
	}

	return l.save(b)
}

package baslex

import (
	"testing"
)

func TestEOF(t *testing.T) {
	compareID(t, "eof", "", []int{TkEOF})
}

func TestColon(t *testing.T) {
	compareID(t, "colon-empty", "", []int{TkEOF})
	compareID(t, "colon-1space", " ", []int{TkEOF})
	compareID(t, "colon-1colon", ":", []int{TkColon, TkEOF})
	compareID(t, "colon-1colon-spaces", "  :  ", []int{TkColon, TkEOF})
	compareID(t, "colon-2colon", "::", []int{TkColon, TkColon, TkEOF})
	compareID(t, "colon-2colon-spaces", "  ::  ", []int{TkColon, TkColon, TkEOF})
	compareID(t, "colon-2colon-spaces-between", "  :  :  ", []int{TkColon, TkColon, TkEOF})
}

func TestCommentQ(t *testing.T) {
	compareID(t, "commentq-empty", "'", []int{TkCommentQ, TkEOF})
	compareID(t, "commentq-hi", "' hi", []int{TkCommentQ, TkEOF})
	compareID(t, "commentq-hi-comment", "' hi '", []int{TkCommentQ, TkEOF})
	compareID(t, "commentq-colon-after", "' hi :", []int{TkCommentQ, TkEOF})
	compareID(t, "commentq-colon-before", ":' hi", []int{TkColon, TkCommentQ, TkEOF})
	compareID(t, "commentq-colon-before-spaces", " : ' hi", []int{TkColon, TkCommentQ, TkEOF})
}

func TestString(t *testing.T) {
	compareValue(t, "string-empty", `"`, []Token{{ID: TkString, Value: `"`}, {ID: TkEOF, Value: "EOF"}})
	compareValue(t, "string-hi", `" hi`, []Token{{ID: TkString, Value: `" hi`}, {ID: TkEOF, Value: "EOF"}})
	compareValue(t, "string-hi-comment", `" hi '`, []Token{{ID: TkString, Value: `" hi '`}, {ID: TkEOF, Value: "EOF"}})
	compareValue(t, "string-colon-after", `" hi :`, []Token{{ID: TkString, Value: `" hi :`}, {ID: TkEOF, Value: "EOF"}})
	compareValue(t, "string-colon-before", `:" hi`, []Token{{ID: TkColon, Value: ":"}, {ID: TkString, Value: `" hi`}, {ID: TkEOF, Value: "EOF"}})
	compareValue(t, "string-colon-before-spaces", ` : " hi`, []Token{{ID: TkColon, Value: ":"}, {ID: TkString, Value: `" hi`}, {ID: TkEOF, Value: "EOF"}})

	compareValue(t, "string2-empty", `""`, []Token{{ID: TkString, Value: `""`}, {ID: TkEOF, Value: "EOF"}})
	compareValue(t, "string2-hi", `" hi"`, []Token{{ID: TkString, Value: `" hi"`}, {ID: TkEOF, Value: "EOF"}})
	compareValue(t, "string2-hi-comment", `" hi "'`, []Token{{ID: TkString, Value: `" hi "`}, {ID: TkCommentQ, Value: "'"}, {ID: TkEOF, Value: "EOF"}})
	compareValue(t, "string2-colon-after", `" hi :"`, []Token{{ID: TkString, Value: `" hi :"`}, {ID: TkEOF, Value: "EOF"}})
	compareValue(t, "string2-colon-after", `" hi ":`, []Token{{ID: TkString, Value: `" hi "`}, {ID: TkColon, Value: ":"}, {ID: TkEOF, Value: "EOF"}})
	compareValue(t, "string2-colon-before", `:" hi"`, []Token{{ID: TkColon, Value: ":"}, {ID: TkString, Value: `" hi"`}, {ID: TkEOF, Value: "EOF"}})
	compareValue(t, "string2-colon-before-spaces", ` : " hi"`, []Token{{ID: TkColon, Value: ":"}, {ID: TkString, Value: `" hi"`}, {ID: TkEOF, Value: "EOF"}})
}

func compareValue(t *testing.T, label, str string, tokens []Token) {

	lex := NewStr(str)
	var i int
	for ; lex.HasToken(); i++ {
		tok := lex.Next()
		if tok.ID != tokens[i].ID {
			t.Errorf("compareValue: %s: bad id: found %v expected: %v", label, tok, tokens[i])
			return
		}
		if tok.Value != tokens[i].Value {
			t.Errorf("compareValue: %s: bad value: found %v expected: %v", label, tok, tokens[i])
			return
		}
	}

	if i != len(tokens) {
		t.Errorf("compareValue: %s: bad token count: found %d expected: %d", label, i, len(tokens))
	}
}

func compareID(t *testing.T, label, str string, tokens []int) {

	lex := NewStr(str)
	var i int
	for ; lex.HasToken(); i++ {
		tok := lex.Next()
		if tok.ID != tokens[i] {
			t.Errorf("compareID: %s: bad token: found %v expected: %v", label, tok, tokens[i])
			return
		}
	}

	if i != len(tokens) {
		t.Errorf("compareID: %s: bad token count: found %d expected: %d", label, i, len(tokens))
	}
}

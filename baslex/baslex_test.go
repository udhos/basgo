package baslex

import (
	"testing"
)

func TestEOF(t *testing.T) {
	compare(t, "eof", "", []int{TkEOF})
}

func TestColon(t *testing.T) {
	compare(t, "colon-empty", "", []int{TkEOF})
	compare(t, "colon-1space", " ", []int{TkEOF})
	compare(t, "colon-1colon", ":", []int{TkColon, TkEOF})
	compare(t, "colon-1colon-spaces", "  :  ", []int{TkColon, TkEOF})
	compare(t, "colon-2colon", "::", []int{TkColon, TkColon, TkEOF})
	compare(t, "colon-2colon-spaces", "  ::  ", []int{TkColon, TkColon, TkEOF})
	compare(t, "colon-2colon-spaces-between", "  :  :  ", []int{TkColon, TkColon, TkEOF})
}

func TestCommentQ(t *testing.T) {
	compare(t, "commentq-empty", "'", []int{TkCommentQ, TkEOF})
	compare(t, "commentq-hi", "' hi", []int{TkCommentQ, TkEOF})
	compare(t, "commentq-hi-comment", "' hi '", []int{TkCommentQ, TkEOF})
	compare(t, "commentq-colon-after", "' hi :", []int{TkCommentQ, TkEOF})
	compare(t, "commentq-colon-before", ":' hi", []int{TkColon, TkCommentQ, TkEOF})
	compare(t, "commentq-colon-before-spaces", " : ' hi", []int{TkColon, TkCommentQ, TkEOF})
}

func compare(t *testing.T, label, str string, tokens []int) {

	lex := NewStr(str)
	var i int
	for ; lex.HasToken(); i++ {
		tok := lex.Next()
		if tok.ID != tokens[i] {
			t.Errorf("compare: %s: bad token: found %d expected: %v", label, tok.ID, tokens[i])
			return
		}
	}

	if i != len(tokens) {
		t.Errorf("compare: %s: bad token count: found %d expected: %d", label, i, len(tokens))
	}
}

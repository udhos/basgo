package baslex

import (
	"testing"
)

func TestEOF(t *testing.T) {
	lex := NewStr("")
	if !lex.HasToken() {
		t.Errorf("could not find any token")
	}
	tok := lex.Next()
	if !tok.IsEOF() {
		t.Errorf("non-EOF token: %v", tok)
	}
}

func TestCommentQ1(t *testing.T) {
	lex := NewStr("'")
	if !lex.HasToken() {
		t.Errorf("could not find any token")
	}
	tok := lex.Next()
	if tok.ID != TkCommentQ {
		t.Errorf("non-comment-q token: %v", tok)
	}
	if lex.HasToken() {
		t.Errorf("non-expected token after comment-q: %v", lex.Next())
	}
}

func TestCommentQ2(t *testing.T) {
	lex := NewStr(" '")
	if !lex.HasToken() {
		t.Errorf("could not find any token")
	}
	tok := lex.Next()
	if tok.ID != TkCommentQ {
		t.Errorf("non-comment-q token: %v", tok)
	}
	if lex.HasToken() {
		t.Errorf("non-expected token after comment-q")
	}
}

func TestCommentQ3(t *testing.T) {
	lex := NewStr("' ")
	if !lex.HasToken() {
		t.Errorf("could not find any token")
	}
	tok := lex.Next()
	if tok.ID != TkCommentQ {
		t.Errorf("non-comment-q token: %v", tok)
	}
	if lex.HasToken() {
		t.Errorf("non-expected token after comment-q")
	}
}

func TestCommentQ4(t *testing.T) {
	lex := NewStr(" ' ")
	if !lex.HasToken() {
		t.Errorf("could not find any token")
	}
	tok := lex.Next()
	if tok.ID != TkCommentQ {
		t.Errorf("non-comment-q token: %v", tok)
	}
	if lex.HasToken() {
		t.Errorf("non-expected token after comment-q")
	}
}

func TestCls(t *testing.T) {
	lex := NewStr(" 10cls ")

	if !lex.HasToken() {
		t.Errorf("could not find 1st token")
	}
	if tok := lex.Next(); tok.ID != TkLineNumber {
		t.Errorf("non-line-number token: %v", tok)
	}

	if !lex.HasToken() {
		t.Errorf("could not find 2nd token")
	}
	if tok := lex.Next(); tok.ID != TkKeywordCls {
		t.Errorf("non-cls token: %v", tok)
	}

	if !lex.HasToken() {
		t.Errorf("could not find 3rd token")
	}
	if tok := lex.Next(); !tok.IsEOF() {
		t.Errorf("non-EOF token: %v", tok)
	}

	if lex.HasToken() {
		t.Errorf("non-expected token after EOF")
	}
}

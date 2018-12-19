package baslex

import (
	"strings"
	"testing"
)

func TestEOF(t *testing.T) {
	lex := New(strings.NewReader(""))
	if !lex.HasToken() {
		t.Errorf("could not find any token")
	}
	tok := lex.Next()
	if !tok.IsEOF() {
		t.Errorf("non-EOF token: %v", tok)
	}
}

func TestCls(t *testing.T) {
	lex := New(strings.NewReader(" 10cls "))

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

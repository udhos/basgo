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

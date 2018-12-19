package baslex

import (
	"bytes"
	"testing"

	"github.com/udhos/basgo/basgo"
)

func TestEOF(t *testing.T) {
	lex := New(strings.NewReader(""))
	if !lex.HasToken() {
		t.Errorf("could not find any token")
	}
	tok := lex.Next()
	if !tok.isEOF() {
		t.Errorf("non-EOF token: %v", tok)
	}
}

package basparser

import (
	"testing"

	"github.com/udhos/basgo/baslex"
)

func TestTokenTable(t *testing.T) {

	parserTokenFirst := parserToken(baslex.TokenIDFirst)
	if parserTokenFirst != parserTokenIDFirst {
		t.Errorf("bad first parser token: lexIndex=%d found=%d expected=%d", baslex.TokenIDFirst, parserTokenFirst, parserTokenIDFirst)
	}

	parserTokenLast := parserToken(baslex.TokenIDLast)
	if parserTokenLast != parserTokenIDLast {
		t.Errorf("bad last parser token: lexIndex=%d found=%d expected=%d", baslex.TokenIDLast, parserTokenLast, parserTokenIDLast)
	}
}

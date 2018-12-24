package basparser

import (
	"github.com/udhos/basgo/baslex"
)

const (
	parserTokenIDFirst = TkNull
	parserTokenIDLast  = TkIdentifier
)

func parserToken(lexToken int) int {
	if lexToken < baslex.TokenIDFirst || lexToken > baslex.TokenIDLast {
		return 0
	}
	return lexToken + parserTokenIDFirst
}

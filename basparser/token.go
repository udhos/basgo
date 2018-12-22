package basparser

// index: lex token    (baslex.Tk...)
// value: parser token (basparser.Tk...)

var tabLexToken = []int{
	TkNull, // TkNull
	TkEOF,  // TkEOF
	TkEOL,  // TkEOL

	TkErrInput,    // TkErrInput
	TkErrInternal, // TkErrInternal
	TkErrInvalid,  // TkErrInvalid
	TkErrLarge,    // TkErrLarge

	TkColon,        // TkColon
	TkComma,        // TkComma
	TkSemicolon,    // TkSemicolon
	TkParLeft,      // TkParLeft
	TkParRight,     // TkParRight
	TkBracketLeft,  // TkBracketLeft
	TkBracketRight, // TkBracketRight
	TkCommentQ,     // TkCommentQ
	TkString,       // TkString
	TkNumber,       // TkNumber

	TkEqual,   // TkEqual
	TkLT,      // TkLT
	TkGT,      // TkGT
	TkUnequal, // TkUnequal
	TkLE,      // TkLE
	TkGE,      // TkGE

	TkPlus,      // TkPlus
	TkMinus,     // TkMinus
	TkMult,      // TkMult
	TkDiv,       // TkDiv
	TkBackSlash, // TkBackSlash

	TkKeywordCls,    // TkKeywordCls
	TkKeywordCont,   // TkKeywordCont
	TkKeywordElse,   // TkKeywordElse
	TkKeywordEnd,    // TkKeywordEnd
	TkKeywordGoto,   // TkKeywordGoto
	TkKeywordInput,  // TkKeywordInput
	TkKeywordIf,     // TkKeywordIf
	TkKeywordLet,    // TkKeywordLet
	TkKeywordList,   // TkKeywordList
	TkKeywordLoad,   // TkKeywordLoad
	TkKeywordPrint,  // TkKeywordPrint
	TkKeywordRun,    // TkKeywordRun
	TkKeywordSave,   // TkKeywordSave
	TkKeywordStop,   // TkKeywordStop
	TkKeywordSystem, // TkKeywordSystem
	TkKeywordThen,   // TkKeywordThen
	TkKeywordTime,   // TkKeywordTime

	TkIdentifier, // TkIdentifier
}

const (
	parserTokenIDFirst = TkNull
	parserTokenIDLast  = TkIdentifier
)

func parserToken(lexToken int) int {
	if lexToken < parserTokenIDFirst || lexToken > parserTokenIDLast {
		return 0
	}
	return tabLexToken[lexToken]
}

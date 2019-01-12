package baslex

import (
	"testing"
)

func TestTabKeywords(t *testing.T) {
	// should be +1
	// however these 7 keywords don't belong to keyword block:
	// REM EQV IMP NOT AND OR XOR
	size := TkKeywordTo - TkKeywordCls + 8
	if len(tabKeywords) != size {
		t.Errorf("mismatch keywords table size: table=%d tokens=%d", len(tabKeywords), size)
	}
}

func TestTokenTable(t *testing.T) {

	if TokenIDFirst != 0 {
		t.Errorf("Bad first token ID: %d", TokenIDFirst)
	}

	if len(tabType) != (TokenIDLast + 1) {
		t.Errorf("Token table size=%d MISMATCHES last token id=%d", len(tabType), TokenIDLast+1)
	}
}

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
	expectTokenEOF := tokenEOF
	expectTokenColon := Token{ID: TkColon, Value: ":"}

	compareValue(t, "string-empty", `"`, []Token{{ID: TkString, Value: `"`}, expectTokenEOF})
	compareValue(t, "string-hi", `" hi`, []Token{{ID: TkString, Value: `" hi`}, expectTokenEOF})
	compareValue(t, "string-hi-comment", `" hi '`, []Token{{ID: TkString, Value: `" hi '`}, expectTokenEOF})
	compareValue(t, "string-colon-after", `" hi :`, []Token{{ID: TkString, Value: `" hi :`}, expectTokenEOF})
	compareValue(t, "string-colon-before", `:" hi`, []Token{expectTokenColon, {ID: TkString, Value: `" hi`}, expectTokenEOF})
	compareValue(t, "string-colon-before-spaces", ` : " hi`, []Token{expectTokenColon, {ID: TkString, Value: `" hi`}, expectTokenEOF})

	compareValue(t, "string2-empty", `""`, []Token{{ID: TkString, Value: `""`}, expectTokenEOF})
	compareValue(t, "string2-hi", `" hi"`, []Token{{ID: TkString, Value: `" hi"`}, expectTokenEOF})
	compareValue(t, "string2-hi-comment", `" hi "'`, []Token{{ID: TkString, Value: `" hi "`}, {ID: TkCommentQ, Value: "'"}, expectTokenEOF})
	compareValue(t, "string2-colon-after", `" hi :"`, []Token{{ID: TkString, Value: `" hi :"`}, expectTokenEOF})
	compareValue(t, "string2-colon-after", `" hi ":`, []Token{{ID: TkString, Value: `" hi "`}, expectTokenColon, expectTokenEOF})
	compareValue(t, "string2-colon-before", `:" hi"`, []Token{expectTokenColon, {ID: TkString, Value: `" hi"`}, expectTokenEOF})
	compareValue(t, "string2-colon-before-spaces", ` : " hi"`, []Token{expectTokenColon, {ID: TkString, Value: `" hi"`}, expectTokenEOF})
}

func TestNumber(t *testing.T) {
	expectTokenEOF := tokenEOF
	expectTokenColon := Token{ID: TkColon, Value: ":"}
	num67 := Token{ID: TkNumber, Value: `67`}
	num345 := Token{ID: TkNumber, Value: `345`}
	strHi := Token{ID: TkString, Value: `" hi "`}

	compareValue(t, "number-simple", `0`, []Token{{ID: TkNumber, Value: `0`}, expectTokenEOF})
	compareValue(t, "number-simple-spaces", ` 1 `, []Token{{ID: TkNumber, Value: `1`}, expectTokenEOF})
	compareValue(t, "number-simple2", `20`, []Token{{ID: TkNumber, Value: `20`}, expectTokenEOF})
	compareValue(t, "number-simple2-spaces", ` 345 `, []Token{num345, expectTokenEOF})
	compareValue(t, "number-two-spaces", ` 345 67 `, []Token{num345, num67, expectTokenEOF})
	compareValue(t, "number-two-colon", ` 345:67 `, []Token{num345, expectTokenColon, num67, expectTokenEOF})
	compareValue(t, "number-two-color-spc", ` 345 : 67 `, []Token{num345, expectTokenColon, num67, expectTokenEOF})
	compareValue(t, "number-two-string", ` 345" hi "67 `, []Token{num345, strHi, num67, expectTokenEOF})
	compareValue(t, "number-two-string-spc", ` 345  " hi "  67 `, []Token{num345, strHi, num67, expectTokenEOF})
}

func TestName(t *testing.T) {
	expectTokenEOF := tokenEOF
	expectTokenColon := Token{ID: TkColon, Value: ":"}
	num67 := Token{ID: TkNumber, Value: `67`}
	strHi := Token{ID: TkString, Value: `" hi "`}

	compareValue(t, "name", `a`, []Token{{ID: TkIdentifier, Value: `a`}, expectTokenEOF})
	compareValue(t, "name", `a$`, []Token{{ID: TkIdentifier, Value: `a$`}, expectTokenEOF})
	compareValue(t, "name", `a!`, []Token{{ID: TkIdentifier, Value: `a!`}, expectTokenEOF})
	compareValue(t, "name", `a%`, []Token{{ID: TkIdentifier, Value: `a%`}, expectTokenEOF})
	compareValue(t, "name", `a#`, []Token{{ID: TkIdentifier, Value: `a#`}, expectTokenEOF})

	compareValue(t, "name", ` a.2 `, []Token{{ID: TkIdentifier, Value: `a.2`}, expectTokenEOF})
	compareValue(t, "name", ` abc `, []Token{{ID: TkIdentifier, Value: `abc`}, expectTokenEOF})
	compareValue(t, "name", ` TIME 67 " hi "`, []Token{{ID: TkIdentifier, Value: `TIME`}, num67, strHi, expectTokenEOF})
	compareValue(t, "name", ` TIME$ 67 " hi "`, []Token{{ID: TkKeywordTime, Value: `TIME$`}, num67, strHi, expectTokenEOF})
	compareValue(t, "name", ` : CLS  67 " hi "`, []Token{expectTokenColon, {ID: TkKeywordCls, Value: `CLS`}, num67, strHi, expectTokenEOF})
}

func TestEqual(t *testing.T) {
	expectTokenEOF := tokenEOF
	eq := Token{ID: TkEqual, Value: `=`}
	un := Token{ID: TkUnequal, Value: `<>`}
	lt := Token{ID: TkLT, Value: `<`}
	gt := Token{ID: TkGT, Value: `>`}
	le := Token{ID: TkLE, Value: `<=`}
	ge := Token{ID: TkGE, Value: `>=`}

	compareValue(t, "equal-lt", `<`, []Token{lt, expectTokenEOF})
	compareValue(t, "equal-lt", ` < `, []Token{lt, expectTokenEOF})
	compareValue(t, "equal-lt2", ` << `, []Token{lt, lt, expectTokenEOF})
	compareValue(t, "equal-lt2", ` < < `, []Token{lt, lt, expectTokenEOF})

	compareValue(t, "equal-gt", `>`, []Token{gt, expectTokenEOF})
	compareValue(t, "equal-gt", ` > `, []Token{gt, expectTokenEOF})
	compareValue(t, "equal-gt2", ` >> `, []Token{gt, gt, expectTokenEOF})
	compareValue(t, "equal-gt2", ` > > `, []Token{gt, gt, expectTokenEOF})

	compareValue(t, "equal-le", `<=`, []Token{le, expectTokenEOF})
	compareValue(t, "equal-le", ` <= `, []Token{le, expectTokenEOF})
	compareValue(t, "equal-le2", ` <=<= `, []Token{le, le, expectTokenEOF})
	compareValue(t, "equal-le2", ` <= <= `, []Token{le, le, expectTokenEOF})

	compareValue(t, "equal-ge", `>=`, []Token{ge, expectTokenEOF})
	compareValue(t, "equal-ge", ` >= `, []Token{ge, expectTokenEOF})
	compareValue(t, "equal-ge2", ` >=>= `, []Token{ge, ge, expectTokenEOF})
	compareValue(t, "equal-ge2", ` >= >= `, []Token{ge, ge, expectTokenEOF})

	compareValue(t, "equal-eq", `=`, []Token{eq, expectTokenEOF})
	compareValue(t, "equal-eq", ` = `, []Token{eq, expectTokenEOF})
	compareValue(t, "equal-un", `<>`, []Token{un, expectTokenEOF})
	compareValue(t, "equal-un", ` <> `, []Token{un, expectTokenEOF})
	compareValue(t, "equal-lg", ` < > `, []Token{lt, gt, expectTokenEOF})
	compareValue(t, "equal-eq-un", `=<>`, []Token{eq, un, expectTokenEOF})
	compareValue(t, "equal-eq-un", ` = <> `, []Token{eq, un, expectTokenEOF})
	compareValue(t, "equal-un-eq", ` <>= `, []Token{un, eq, expectTokenEOF})
	compareValue(t, "equal-un-eq", ` <> = `, []Token{un, eq, expectTokenEOF})
}

func TestArith(t *testing.T) {

	expectTokenEOF := tokenEOF
	plus := Token{ID: TkPlus, Value: `+`}
	minus := Token{ID: TkMinus, Value: `-`}
	mult := Token{ID: TkMult, Value: `*`}
	div := Token{ID: TkDiv, Value: `/`}
	bs := Token{ID: TkBackSlash, Value: `\`}
	ident := Token{ID: TkIdentifier, Value: `a`}

	seq1 := []Token{plus, minus, mult, div, bs, expectTokenEOF}
	seq2 := []Token{ident, plus, ident, minus, ident, mult, ident, div, ident, bs, ident, expectTokenEOF}

	compareValue(t, "arith", `+-*/\`, seq1)
	compareValue(t, "arith", ` + - * / \ `, seq1)
	compareValue(t, "arith", `a+a-a*a/a\a`, seq2)
	compareValue(t, "arith", ` a + a - a * a / a \ a `, seq2)
}

func TestMarks(t *testing.T) {

	expectTokenEOF := tokenEOF
	comma := Token{ID: TkComma, Value: `,`}
	semi := Token{ID: TkSemicolon, Value: `;`}
	lp := Token{ID: TkParLeft, Value: `(`}
	rp := Token{ID: TkParRight, Value: `)`}
	end := Token{ID: TkKeywordEnd, Value: `end`}

	seq1 := []Token{comma, semi, lp, rp, expectTokenEOF}
	seq2 := []Token{end, comma, end, semi, end, lp, end, rp, end, expectTokenEOF}

	compareValue(t, "mark", `,;()`, seq1)
	compareValue(t, "mark", ` , ; ( ) `, seq1)
	compareValue(t, "mark", `end,end;end(end)end`, seq2)
	compareValue(t, "mark", ` end , end ; end ( end ) end `, seq2)
}

func TestBrackets(t *testing.T) {

	expectTokenEOF := tokenEOF
	lb := Token{ID: TkBracketLeft, Value: `[`}
	rb := Token{ID: TkBracketRight, Value: `]`}
	let := Token{ID: TkKeywordLet, Value: `let`}

	seq1 := []Token{lb, rb, expectTokenEOF}
	seq2 := []Token{let, lb, let, rb, let, expectTokenEOF}

	compareValue(t, "bracket", `[]`, seq1)
	compareValue(t, "bracket", ` [ ] `, seq1)
	compareValue(t, "bracket", `let[let]let`, seq2)
	compareValue(t, "bracket", ` let [ let ] let `, seq2)
}

func TestKeywords(t *testing.T) {

	expectTokenEOF := tokenEOF
	kwIf := Token{ID: TkKeywordIf, Value: `if`}
	kwThen := Token{ID: TkKeywordThen, Value: `then`}
	kwElse := Token{ID: TkKeywordElse, Value: `else`}
	kwStop := Token{ID: TkKeywordStop, Value: `stop`}
	kwSystem := Token{ID: TkKeywordSystem, Value: `system`}
	kwCont := Token{ID: TkKeywordCont, Value: `cont`}

	seq := []Token{kwIf, kwThen, kwElse, kwStop, kwSystem, kwCont, expectTokenEOF}

	compareValue(t, "keywords", ` if then else stop system cont `, seq)
}

func TestRem(t *testing.T) {

	expectTokenEOF := tokenEOF
	rem := Token{ID: TkKeywordRem, Value: `rem`}
	remSpc := Token{ID: TkKeywordRem, Value: `rem `}
	remBig := Token{ID: TkKeywordRem, Value: `rem  this is a comment! : print`}

	seq := []Token{rem, expectTokenEOF}
	seqSpc := []Token{remSpc, expectTokenEOF}
	seqBig := []Token{remBig, expectTokenEOF}

	compareValue(t, "rem", `rem`, seq)
	compareValue(t, "rem", ` rem`, seq)
	compareValue(t, "rem", ` rem `, seqSpc)
	compareValue(t, "rem", ` rem  this is a comment! : print`, seqBig)
}

func TestLoop(t *testing.T) {

	expectTokenEOF := tokenEOF
	kwFor := Token{ID: TkKeywordFor, Value: `for`}
	kwTo := Token{ID: TkKeywordTo, Value: `to`}
	kwStep := Token{ID: TkKeywordStep, Value: `step`}
	kwNext := Token{ID: TkKeywordNext, Value: `next`}
	kwGosub := Token{ID: TkKeywordGosub, Value: `gosub`}
	kwReturn := Token{ID: TkKeywordReturn, Value: `return`}

	seq := []Token{kwFor, kwTo, kwStep, kwNext, kwGosub, kwReturn, expectTokenEOF}

	compareValue(t, "loop", " for to step next gosub return", seq)
}

func TestLogical(t *testing.T) {

	expectTokenEOF := tokenEOF
	kwAnd := Token{ID: TkKeywordAnd, Value: `and`}
	kwNot := Token{ID: TkKeywordNot, Value: `not`}
	kwImp := Token{ID: TkKeywordImp, Value: `imp`}
	kwEqv := Token{ID: TkKeywordEqv, Value: `eqv`}
	kwOr := Token{ID: TkKeywordOr, Value: `or`}
	kwXor := Token{ID: TkKeywordXor, Value: `xor`}

	seq := []Token{kwAnd, kwNot, kwImp, kwEqv, kwOr, kwXor, expectTokenEOF}

	compareValue(t, "logical", " and not imp eqv or xor", seq)
}

func compareValue(t *testing.T, label, str string, tokens []Token) {

	lex := NewStr(str)
	var i int
	for ; lex.HasToken(); i++ {
		tok := lex.Next()
		if tok.ID != tokens[i].ID {
			t.Errorf("compareValue: %s: bad id: found %v id=%d expected: tok=%v id=%d", label, tok, tok.ID, tokens[i], tokens[i].ID)
			return
		}
		if tok.Value != tokens[i].Value {
			t.Errorf("compareValue: %s: bad value: found [%s] expected: [%s]", label, tok.Value, tokens[i].Value)
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

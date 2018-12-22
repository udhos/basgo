%{

package basparser

import (
	//"bufio"
	"fmt"
	//"os"
	"unicode"
)

%}

// fields inside this union end up as the fields in a structure known
// as ${PREFIX}SymType, of which a reference is passed to the lexer.
%union{
	val string
}

// any non-terminal which returns a value needs a type, which is
// really a field name in the above union struct
%type <val> input

// same for terminals
%token <val> CHARACTER

%token <val> TkNull
%token <val> TkEOF
%token <val> TkEOL

%token <val> TkErrInput
%token <val> TkErrInternal
%token <val> TkErrInvalid
%token <val> TkErrLarge

%token <val> TkColon
%token <val> TkComma
%token <val> TkSemicolon
%token <val> TkParLeft
%token <val> TkParRight
%token <val> TkBracketLeft
%token <val> TkBracketRight
%token <val> TkCommentQ
%token <val> TkString
%token <val> TkNumber

%token <val> TkEqual
%token <val> TkLT
%token <val> TkGT
%token <val> TkUnequal
%token <val> TkLE
%token <val> TkGE

%token <val> TkPlus
%token <val> TkMinus
%token <val> TkMult
%token <val> TkDiv
%token <val> TkBackSlash

%token <val> TkKeywordCls
%token <val> TkKeywordCont
%token <val> TkKeywordElse
%token <val> TkKeywordEnd
%token <val> TkKeywordGoto
%token <val> TkKeywordInput
%token <val> TkKeywordIf
%token <val> TkKeywordLet
%token <val> TkKeywordList
%token <val> TkKeywordLoad
%token <val> TkKeywordPrint
%token <val> TkKeywordRem
%token <val> TkKeywordRun
%token <val> TkKeywordSave
%token <val> TkKeywordStop
%token <val> TkKeywordSystem
%token <val> TkKeywordThen
%token <val> TkKeywordTime

%token <val> TkIdentifier

%%

in : /* empty */
  | in input '\n'
     { fmt.Printf("Read character: %s\n", $2) }
  ;

input : CHARACTER
  | input CHARACTER
      { $$ = $1 + $2 }
  ;

%%

func NewInputLex(line string) *InputLex {
 	return &InputLex{s: line}
}

type InputLex struct {
        // contains one complete input string (with the trailing \n)
        s string
        // used to keep track of parser position along the above input string
        pos int
}

func (l *InputLex) Lex(lval *InputSymType) int {
	var c rune = ' '

        // skip through all the spaces, both at the ends and in between
	for c == ' ' {
		if l.pos == len(l.s) {
			return 0
		}
		c = rune(l.s[l.pos])
		l.pos += 1
	}

        // only look for input characters that are either digits or lower case
	if unicode.IsDigit(c) || unicode.IsLower(c) {
	    lval.val = string(c)
	    return CHARACTER
	}

        // do not return any token in case of unrecognized grammar
        // this results in syntax error
	return int(c)
}

func (l *InputLex) Error(s string) {
	fmt.Printf("syntax error: %s\n", s)
}


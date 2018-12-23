%{

package basparser

import (
	//"bufio"
	"fmt"
	//"os"
	//"unicode"
	"io"

	"github.com/udhos/basgo/baslex"
)

%}

// fields inside this union end up as the fields in a structure known
// as ${PREFIX}SymType, of which a reference is passed to the lexer.
%union{
	typeProg int
	typeLine int
	typeStmtList int
	typeStmt int
	tok int
}

// any non-terminal which returns a value needs a type, which is
// really a field name in the above union struct
//%type <val> input
%type <typeProg> prog
%type <typeLine> line
%type <typeStmtList> statements
%type <typeStmt> stmt

// same for terminals
//%token <val> CHARACTER

%token <tok> TkNull
%token <tok> TkEOF
%token <tok> TkEOL

%token <tok> TkErrInput
%token <tok> TkErrInternal
%token <tok> TkErrInvalid
%token <tok> TkErrLarge

%token <tok> TkColon
%token <tok> TkComma
%token <tok> TkSemicolon
%token <tok> TkParLeft
%token <tok> TkParRight
%token <tok> TkBracketLeft
%token <tok> TkBracketRight
%token <tok> TkCommentQ
%token <tok> TkString
%token <tok> TkNumber

%token <tok> TkEqual
%token <tok> TkLT
%token <tok> TkGT
%token <tok> TkUnequal
%token <tok> TkLE
%token <tok> TkGE

%token <tok> TkPlus
%token <tok> TkMinus
%token <tok> TkMult
%token <tok> TkDiv
%token <tok> TkBackSlash

%token <tok> TkKeywordCls
%token <tok> TkKeywordCont
%token <tok> TkKeywordElse
%token <tok> TkKeywordEnd
%token <tok> TkKeywordGoto
%token <tok> TkKeywordInput
%token <tok> TkKeywordIf
%token <tok> TkKeywordLet
%token <tok> TkKeywordList
%token <tok> TkKeywordLoad
%token <tok> TkKeywordPrint
%token <tok> TkKeywordRem
%token <tok> TkKeywordRun
%token <tok> TkKeywordSave
%token <tok> TkKeywordStop
%token <tok> TkKeywordSystem
%token <tok> TkKeywordThen
%token <tok> TkKeywordTime

%token <tok> TkIdentifier

%%

prog: line_list TkEOF
     {
	 fmt.Printf("parser action - full prog?\n")
	 $$ = 1
     }
  ;

line_list: line
  | line_list TkEOL line
  ;

line: statements
     { $$ = 4 /* statements */ }
  | TkNumber statements
     { $$ = 5 /* number statements */ }
  ;

statements: stmt
     { $$ = 6 }
  | statements TkColon stmt
     { $$ = 7 /* stmt */ }
  ;

stmt: /* empty */
     { $$ = 8  }
  | TkKeywordEnd
     { $$ = 9 /* end */ }
  | TkKeywordPrint
     { $$ = 10 /* print */ }
  ;

//in : /* empty */
//  | in input '\n'
//     { fmt.Printf("Read character: %s\n", $2) }
//  ;
//
//input : CHARACTER
//  | input CHARACTER
//      { $$ = $1 + $2 }
//  ;

%%

func NewInputLex(input io.ByteScanner) *InputLex {
 	return &InputLex{lex: baslex.New(input)}
}

type InputLex struct {
	lex *baslex.Lex
}

func (l *InputLex) Lex(lval *InputSymType) int {

	if !l.lex.HasToken() {
		return 0
	}

	t := l.lex.Next()

	id := parserToken(t.ID) // convert lex ID to parser ID

	fmt.Printf("InputLex.Lex: lex=%v parser=%d\n", t, id)

	return id
}

func (l *InputLex) Error(s string) {
	fmt.Printf("InputLex.Error: syntax error: %s\n", s)
}


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

// parser auxiliary variables
var (
	Root []Node
	lineList []Node
	nodeList []Node
)

%}

// fields inside this union end up as the fields in a structure known
// as ${PREFIX}SymType, of which a reference is passed to the lexer.
%union{
	//typeProg []Node
	typeLineList []Node
	typeLine Node
	typeStmtList []Node
	typeStmt Node
	tok int
}

// any non-terminal which returns a value needs a type, which is
// really a field name in the above union struct

//%type <typeProg> prog
%type <typeLineList> line_list
%type <typeLine> line
%type <typeStmtList> statements
%type <typeStmt> stmt

// same for terminals

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
%token <tok> TkKeywordFor
%token <tok> TkKeywordGosub
%token <tok> TkKeywordGoto
%token <tok> TkKeywordInput
%token <tok> TkKeywordIf
%token <tok> TkKeywordLet
%token <tok> TkKeywordList
%token <tok> TkKeywordLoad
%token <tok> TkKeywordNext
%token <tok> TkKeywordPrint
%token <tok> TkKeywordRem
%token <tok> TkKeywordReturn
%token <tok> TkKeywordRun
%token <tok> TkKeywordSave
%token <tok> TkKeywordStep
%token <tok> TkKeywordStop
%token <tok> TkKeywordSystem
%token <tok> TkKeywordThen
%token <tok> TkKeywordTime
%token <tok> TkKeywordTo

%token <tok> TkIdentifier

%%

prog: line_list TkEOF
     { Root = $1 }
  ;

line_list: line
     {
        lineList = []Node{$1} // reset line list
	$$ = lineList
     }
  | line_list TkEOL line
     {
        lineList = append(lineList, $3)
        $$ = lineList
     }
  ;

line: statements
     { $$ = &LineImmediate{Nodes:$1} }
  | TkNumber statements
     { $$ = &LineNumbered{LineNumber:$1, Nodes:$2} }
  ;

statements: stmt
     {
        nodeList = []Node{$1} // reset node list
	$$ = nodeList
     }
  | statements TkColon stmt
     {
        nodeList = append(nodeList, $3)
        $$ = nodeList
     }
  ;

stmt: /* empty */
     { $$ = &NodeEmpty{} }
  | TkKeywordEnd
     { $$ = &NodeEnd{} }
  | TkKeywordPrint
     { $$ = &NodePrint{} }
  ;

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


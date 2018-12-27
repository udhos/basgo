%{

package basparser

import (
	//"bufio"
	"fmt"
	//"os"
	//"unicode"
	"io"
	//"strconv"
        "log"

	"github.com/udhos/basgo/baslex"
	"github.com/udhos/basgo/node"
)

// parser auxiliary variables
var (
	Root []node.Node
	lineList []node.Node
	nodeList []node.Node
)

%}

// fields inside this union end up as the fields in a structure known
// as ${PREFIX}SymType, of which a reference is passed to the lexer.
%union{
	typeLineList []node.Node
	typeLine node.Node
	typeStmtList []node.Node
	typeStmt node.Node

	typeNumber string
	typeRawLine string

	tok int
}

// any non-terminal which returns a value needs a type, which is
// really a field name in the above union struct

%type <typeLineList> line_list
%type <typeLine> line
%type <typeStmtList> statements
%type <typeStmt> stmt

// same for terminals

%token <tok> TkNull
%token <typeRawLine> TkEOF
%token <typeRawLine> TkEOL

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
%token <typeNumber> TkNumber

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
     {
         list := $1
         captureRawLine("EOF", list, $2) // only last line
         
	 Root = $1 // save for caller
     }
  ;

line_list: line
     {
        lineList = []node.Node{$1} // reset line list
	$$ = lineList
     }
  | line_list TkEOL line
     {
        captureRawLine("EOL", lineList, $2) // all lines except last

        lineList = append(lineList, $3)
        $$ = lineList
     }
  ;

line: statements
     {
        $$ = &node.LineImmediate{Nodes:$1}
     }
  | TkNumber statements
     {
       $$ = &node.LineNumbered{LineNumber:$1, Nodes:$2}
     }
  ;

statements: stmt
     {
        nodeList = []node.Node{$1} // reset node list
	$$ = nodeList
     }
  | statements TkColon stmt
     {
        nodeList = append(nodeList, $3)
        $$ = nodeList
     }
  ;

stmt: /* empty */
     { $$ = &node.NodeEmpty{} }
  | TkKeywordEnd
     { $$ = &node.NodeEnd{} }
  | TkKeywordPrint
     { $$ = &node.NodePrint{} }
  ;

%%

func captureRawLine(label string, list []node.Node, rawLine string) {
	last := len(list) - 1
	if last < 0 {
		log.Printf("captureRawLine: %s last line index=%d < 0", label, last)
		return
	}

	switch n := list[last].(type) {
		case *node.LineNumbered:
			n.RawLine = rawLine
			list[last] = n	
             		//log.Printf("captureRawLine: %s numbered index=%d raw=[%s]", label, last, n.RawLine)
		case *node.LineImmediate:
			n.RawLine = rawLine
			list[last] = n	
             		//log.Printf("captureRawLine: %s immediate index=%d raw=[%s]", label, last, n.RawLine)
		default:
			log.Printf("captureRawLine: %s non-line node: %v", label, list[last])
	}
}

func NewInputLex(input io.ByteScanner, debug bool) *InputLex {
 	return &InputLex{lex: baslex.New(input), debug:debug}
}

type InputLex struct {
	lex *baslex.Lex
	debug bool
}

func (l *InputLex) Lex(lval *InputSymType) int {

	if !l.lex.HasToken() {
		return 0
	}

	t := l.lex.Next()

	// ATTENTION: t.ID is in lex token space

	id := parserToken(t.ID) // convert lex ID to parser ID

	// ATTENTION: id is in parser token space

	if l.debug {
		fmt.Printf("InputLex.Lex: %s [%s]\n", t.Type(), t.Value)
	}

	// need to store values only for some terminals
        // when a parser rule action need to consume the value
	// for example: ident, literals (number, string)
	switch id {
		case TkNumber:
			lval.typeNumber = t.Value
		case TkEOL:
			lval.typeRawLine = l.lex.RawLine()
		case TkEOF:
			lval.typeRawLine = l.lex.RawLine()
		case TkColon: // do not store
		case TkKeywordEnd: // do not store
		case TkKeywordPrint: // do not store
		default:
			fmt.Printf("InputLex.Lex: FIXME token value [%s] not stored for parser actions\n", t.Value)
	}

	return id
}

func (l *InputLex) Error(s string) {
	fmt.Printf("InputLex.Error: line=%d column=%d: %s\n", l.lex.Line(), l.lex.Column(), s)
}


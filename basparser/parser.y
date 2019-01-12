%{

package basparser

import (
	//"bufio"
	"fmt"
	//"os"
	//"unicode"
	"io"
	"strconv"
        "log"

	"github.com/udhos/basgo/baslex"
	"github.com/udhos/basgo/node"
)

// parser auxiliary variables
var (
	Root []node.Node
	lineList []node.Node
	nodeList []node.Node
	expList []node.NodeExp
)

%}

// fields inside this union end up as the fields in a structure known
// as ${PREFIX}SymType, of which a reference is passed to the lexer.
%union{
	typeLineList []node.Node
	typeLine node.Node
	typeStmtList []node.Node
	typeStmt node.Node

	typeExpressions []node.NodeExp
	typeExp node.NodeExp

	typeRem string
	typeNumber string
	typeFloat string
	typeString string
	typeIdentifier string
	typeRawLine string

	tok int
}

// any non-terminal which returns a value needs a type, which is
// really a field name in the above union struct

%type <typeLineList> line_list
%type <typeLine> line
%type <typeStmtList> statements
%type <typeStmt> stmt
%type <typeStmt> assign
%type <typeExpressions> expressions
%type <typeExp> exp

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
%token <typeString> TkString
%token <typeNumber> TkNumber
%token <typeFloat> TkFloat

//%token <tok> TkEqual
//%token <tok> TkLT
//%token <tok> TkGT
//%token <tok> TkUnequal
//%token <tok> TkLE
//%token <tok> TkGE

//%left <tok> TkKeywordImp
//%left <tok> TkKeywordEqv
//%left <tok> TkKeywordXor
//%left <tok> TkKeywordOr
//%left <tok> TkKeywordAnd
//%left <tok> TkKeywordNot

%left <tok> TkEqual TkUnequal TkLT TkGT TkLE TkGE
%left <tok> TkPlus TkMinus
%left <tok> TkKeywordMod
%left <tok> TkBackSlash
%left <tok> TkMult TkDiv
%left <tok> TkPow
%precedence UnaryPlus // fictitious
%precedence UnaryMinus // fictitious

%token <tok> TkKeywordCls
%token <tok> TkKeywordCont
%token <tok> TkKeywordElse
%token <tok> TkKeywordEnd
%token <tok> TkKeywordFor
%token <tok> TkKeywordGosub
%token <tok> TkKeywordGoto
%token <tok> TkKeywordInput
%token <tok> TkKeywordIf
%token <tok> TkKeywordLen
%token <tok> TkKeywordLet
%token <tok> TkKeywordList
%token <tok> TkKeywordLoad
%token <tok> TkKeywordNext
%token <tok> TkKeywordPrint
%token <typeRem> TkKeywordRem
%token <tok> TkKeywordReturn
%token <tok> TkKeywordRun
%token <tok> TkKeywordSave
%token <tok> TkKeywordStep
%token <tok> TkKeywordStop
%token <tok> TkKeywordSystem
%token <tok> TkKeywordThen
%token <tok> TkKeywordTime
%token <tok> TkKeywordTo

%token <typeIdentifier> TkIdentifier

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
  | TkKeywordLet assign
     { $$ = $2 }
  | assign
     { $$ = $1 }
  | TkKeywordList
     { $$ = &node.NodeList{} }
  | TkKeywordPrint
     { 
        $$ = &node.NodePrint{}
     }
  | TkKeywordPrint expressions
     {
        $$ = &node.NodePrint{Expressions: $2}
     }
  | TkKeywordRem
     { $$ = &node.NodeRem{Value: $1} }
  ;

assign: TkIdentifier TkEqual exp
     {
        $$ = &node.NodeAssign{Left: $1, Right: $3}
     }
  ;

expressions: exp
	{
        	expList = []node.NodeExp{$1} // reset
	        $$ = expList
	}
    |
        expressions exp
        {
		expList = append(expList, $2)
		$$ = expList
	}
    |
        expressions TkComma exp
        {
		expList = append(expList, $3)
		$$ = expList
	}
    |
        expressions TkSemicolon exp
        {
		expList = append(expList, $3)
		$$ = expList
	}
    ;

exp: TkNumber { $$ = &node.NodeExpNumber{Value:$1} }
   | TkFloat 
     {
       n := &node.NodeExpFloat{}
       v := $1
       if v != "." {
         var errParse error
         n.Value, errParse = strconv.ParseFloat(v, 64)
         if errParse != nil {
           msg := fmt.Sprintf("TkFloat action syntax error: %v", errParse)

           // Code inside the grammar actions may refer to the variable yylex,
           // which holds the yyLexer passed to yyParse.
           yylex.Error(msg)
         }
       }
       $$ = n
     }
   | TkString { $$ = &node.NodeExpString{Value:$1} }
   | TkIdentifier { $$ = &node.NodeExpIdentifier{Value:$1} }
   | exp TkPlus exp { $$ = &node.NodeExpPlus{Left: $1, Right: $3} }
   | exp TkMinus exp { $$ = &node.NodeExpMinus{Left: $1, Right: $3} }
   | exp TkKeywordMod exp { $$ = &node.NodeExpMod{Left: $1, Right: $3} }
   | exp TkBackSlash exp { $$ = &node.NodeExpDivInt{Left: $1, Right: $3} }
   | exp TkMult exp { $$ = &node.NodeExpMult{Left: $1, Right: $3} }
   | exp TkDiv exp { $$ = &node.NodeExpDiv{Left: $1, Right: $3} }
   | exp TkPow exp { $$ = &node.NodeExpPow{Left: $1, Right: $3} }
   | TkPlus exp %prec UnaryPlus { $$ = &node.NodeExpUnaryPlus{Value:$2} }
   | TkMinus exp %prec UnaryMinus { $$ = &node.NodeExpUnaryMinus{Value:$2} }
   | TkParLeft exp TkParRight { $$ = &node.NodeExpGroup{Value:$2} }
   | TkKeywordLen exp { $$ = &node.NodeExpLen{Value:$2} }
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
	syntaxErrorCount int
}

func (l *InputLex) Errors() int {
	return l.syntaxErrorCount
}

func (l *InputLex) Lex(lval *InputSymType) int {

	if !l.lex.HasToken() {
		return 0 // 0 means real EOF for the parser
	}

	t := l.lex.Next()

	// ATTENTION: t.ID is in lex token space

	id := parserToken(t.ID) // convert lex ID to parser ID

	// ATTENTION: id is in parser token space

	if l.debug {
		log.Printf("InputLex.Lex: %s [%s]\n", t.Type(), t.Value)
	}

	// need to store values only for some terminals
        // when a parser rule action need to consume the value
	// for example: ident, literals (number, string)
	switch id {
		case TkKeywordRem:
			lval.typeRem = t.Value
		case TkString:
			lval.typeString = t.Value
		case TkNumber:
			lval.typeNumber = t.Value
		case TkFloat:
			lval.typeFloat = t.Value
		case TkIdentifier:
			lval.typeIdentifier = t.Value
		case TkEOL:
			lval.typeRawLine = l.lex.RawLine()
		case TkEOF:
			lval.typeRawLine = l.lex.RawLine()
		case TkEqual: // do not store
		case TkParLeft: // do not store
		case TkParRight: // do not store
		case TkColon: // do not store
		case TkComma: // do not store
		case TkSemicolon: // do not store
		case TkPlus: // do not store
		case TkMinus: // do not store
		case TkMult: // do not store
		case TkDiv: // do not store
		case TkBackSlash: // do not store
		case TkPow: // do not store
		case TkKeywordEnd: // do not store
		case TkKeywordLen: // do not store
		case TkKeywordList: // do not store
		case TkKeywordMod: // do not store
		case TkKeywordPrint: // do not store
		default:
			log.Printf("InputLex.Lex: FIXME token value [%s] not stored for parser actions\n", t.Value)
	}

	return id
}

func (l *InputLex) Error(s string) {
	l.syntaxErrorCount++
	log.Printf("InputLex.Error: count=%d line=%d column=%d: %s\n", l.syntaxErrorCount, l.lex.Line(), l.lex.Column(), s)
}


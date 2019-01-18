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
	LineNumbers = map[string]node.LineNumber{} // used by GOTO GOSUB etc
)

func Reset() {
	Root = []node.Node{}
	LineNumbers = map[string]node.LineNumber{} // used by GOTO GOSUB etc
}

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

%left <tok> TkKeywordImp
%left <tok> TkKeywordEqv
%left <tok> TkKeywordXor
%left <tok> TkKeywordOr
%left <tok> TkKeywordAnd
%left <tok> TkKeywordNot

%left <tok> TkEqual TkUnequal TkLT TkGT TkLE TkGE
%left <tok> TkPlus TkMinus
%left <tok> TkKeywordMod
%left <tok> TkBackSlash
%left <tok> TkMult TkDiv
%right <tok> TkPow
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
       n := $1
       ln, found := LineNumbers[n]
       if found {
         // set defined, keep used unchanged
         ln.Defined = true
         LineNumbers[n] = ln
       } else {
         // set defined, unset used
         LineNumbers[n] = node.LineNumber{Defined: true}
       }
       $$ = &node.LineNumbered{LineNumber:n, Nodes:$2}
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
  | TkKeywordGoto TkNumber
     { 
       n := $2
       ln, found := LineNumbers[n]
       if found {
         // set used, keep defined unchanged
         ln.Used = true
         LineNumbers[n] = ln
       } else {
         // set used, unset defined
         LineNumbers[n] = node.LineNumber{Used: true}
       }
       $$ = &node.NodeGoto{Line: n}
     }
  | TkKeywordLet assign
     { $$ = $2 }
  | assign
     { $$ = $1 }
  | TkKeywordList
     { $$ = &node.NodeList{} }
  | TkKeywordPrint
     { 
        $$ = &node.NodePrint{Newline: true}
     }
  | TkKeywordPrint expressions
     {
        $$ = &node.NodePrint{Expressions: $2, Newline: true}
     }
  | TkKeywordPrint expressions TkSemicolon
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
   | exp TkPlus exp
     {
       if $1.Type() == node.TypeString && $3.Type() != node.TypeString {
           yylex.Error("TkPlus string and non-string")
       }
       if $1.Type() != node.TypeString && $3.Type() == node.TypeString {
           yylex.Error("TkPlus non-string and string")
       }
       n := &node.NodeExpPlus{Left: $1, Right: $3}
       if n.Type() == node.TypeUnknown {
           yylex.Error("TkPlus produces unknown type")
       }
       $$ = n
     }
   | exp TkMinus exp
     {
       switch $1.Type() {
       case node.TypeString:
           yylex.Error("TkMinus left value has string type")
       case node.TypeUnknown:
           yylex.Error("TkMinus left value has unknown type")
       }
       switch $3.Type() {
       case node.TypeString:
           yylex.Error("TkMinus right value has string type")
       case node.TypeUnknown:
           yylex.Error("TkMinus right value has unknown type")
       }
       n := &node.NodeExpMinus{Left: $1, Right: $3}
       switch n.Type() {
       case node.TypeString:
           yylex.Error("TkMinus produces string type")
       case node.TypeUnknown:
           yylex.Error("TkMinus produces unknown type")
       }
       $$ = n
     }
   | exp TkKeywordMod exp
     {
       switch $1.Type() {
       case node.TypeString:
           yylex.Error("TkMod left value has string type")
       case node.TypeUnknown:
           yylex.Error("TkMod left value has unknown type")
       }
       switch $3.Type() {
       case node.TypeString:
           yylex.Error("TkMod right value has string type")
       case node.TypeUnknown:
           yylex.Error("TkMod right value has unknown type")
       }
       n := &node.NodeExpMod{Left: $1, Right: $3}
       if  n.Type() != node.TypeInteger {
           yylex.Error("TkMod produces non-integer type")
       }
       $$ = n
     }
   | exp TkBackSlash exp
     {
       switch $1.Type() {
       case node.TypeString:
           yylex.Error("Integer division left value has string type")
       case node.TypeUnknown:
           yylex.Error("Integer division left value has unknown type")
       }
       switch $3.Type() {
       case node.TypeString:
           yylex.Error("Integer division right value has string type")
       case node.TypeUnknown:
           yylex.Error("Integer division right value has unknown type")
       }
       n := &node.NodeExpDivInt{Left: $1, Right: $3}
       if  n.Type() != node.TypeInteger {
           yylex.Error("Integer division produces non-integer type")
       }
       $$ = n
     }
   | exp TkMult exp
     {
       switch $1.Type() {
       case node.TypeString:
           yylex.Error("TkMult left value has string type")
       case node.TypeUnknown:
           yylex.Error("TkMult left value has unknown type")
       }
       switch $3.Type() {
       case node.TypeString:
           yylex.Error("TkMult right value has string type")
       case node.TypeUnknown:
           yylex.Error("TkMult right value has unknown type")
       }
       n := &node.NodeExpMult{Left: $1, Right: $3}
       switch n.Type() {
       case node.TypeString:
           yylex.Error("TkMult produces string type")
       case node.TypeUnknown:
           yylex.Error("TkMult produces unknown type")
       }
       $$ = n
     }
   | exp TkDiv exp
     {
       switch $1.Type() {
       case node.TypeString:
           yylex.Error("TkDiv left value has string type")
       case node.TypeUnknown:
           yylex.Error("TkDiv left value has unknown type")
       }
       switch $3.Type() {
       case node.TypeString:
           yylex.Error("TkDiv right value has string type")
       case node.TypeUnknown:
           yylex.Error("TkDiv right value has unknown type")
       }
       n := &node.NodeExpDiv{Left: $1, Right: $3}
       if  n.Type() != node.TypeFloat {
           yylex.Error("TkDiv produces non-float type")
       }
       $$ = n
     }
   | exp TkPow exp
     {
       switch $1.Type() {
       case node.TypeString:
           yylex.Error("TkPow left value has string type")
       case node.TypeUnknown:
           yylex.Error("TkPow left value has unknown type")
       }
       switch $3.Type() {
       case node.TypeString:
           yylex.Error("TkPow right value has string type")
       case node.TypeUnknown:
           yylex.Error("TkPow right value has unknown type")
       }
       n := &node.NodeExpPow{Left: $1, Right: $3}
       if  n.Type() != node.TypeFloat {
           yylex.Error("TkPow produces non-float type")
       }
       $$ = n
     }
   | TkPlus exp %prec UnaryPlus
     {
       switch $2.Type() {
       case node.TypeString:
           yylex.Error("Unary plus has string type")
       case node.TypeUnknown:
           yylex.Error("Unary plus has unknown type")
       }
       $$ = &node.NodeExpUnaryPlus{Value:$2}
     }
   | TkMinus exp %prec UnaryMinus
     {
       switch $2.Type() {
       case node.TypeString:
           yylex.Error("Unary minus has string type")
       case node.TypeUnknown:
           yylex.Error("Unary minus has unknown type")
       }
       $$ = &node.NodeExpUnaryMinus{Value:$2}
     }
   | TkParLeft exp TkParRight { $$ = &node.NodeExpGroup{Value:$2} }
   | TkKeywordNot exp
     {
       switch $2.Type() {
       case node.TypeString:
           yylex.Error("Not has string type")
       case node.TypeUnknown:
           yylex.Error("Not has unknown type")
       }
       $$ = &node.NodeExpNot{Value:$2}
     }
   | exp TkKeywordAnd exp
     {
       switch $1.Type() {
       case node.TypeString:
           yylex.Error("TkAnd left value has string type")
       case node.TypeUnknown:
           yylex.Error("TkAnd left value has unknown type")
       }
       switch $3.Type() {
       case node.TypeString:
           yylex.Error("TkAnd right value has string type")
       case node.TypeUnknown:
           yylex.Error("TkAnd right value has unknown type")
       }
       n := &node.NodeExpAnd{Left:$1, Right:$3}
       if  n.Type() != node.TypeInteger {
           yylex.Error("TkAnd produces non-integer type")
       }
       $$ = n
     }
   | exp TkKeywordEqv exp
     {
       switch $1.Type() {
       case node.TypeString:
           yylex.Error("TkEqv left value has string type")
       case node.TypeUnknown:
           yylex.Error("TkEqv left value has unknown type")
       }
       switch $3.Type() {
       case node.TypeString:
           yylex.Error("TkEqv right value has string type")
       case node.TypeUnknown:
           yylex.Error("TkEqv right value has unknown type")
       }
       n := &node.NodeExpEqv{Left:$1, Right:$3}
       if  n.Type() != node.TypeInteger {
           yylex.Error("TkEqv produces non-integer type")
       }
       $$ = n
     }
   | exp TkKeywordImp exp
     {
       switch $1.Type() {
       case node.TypeString:
           yylex.Error("TkImp left value has string type")
       case node.TypeUnknown:
           yylex.Error("TkImp left value has unknown type")
       }
       switch $3.Type() {
       case node.TypeString:
           yylex.Error("TkImp right value has string type")
       case node.TypeUnknown:
           yylex.Error("TkImp right value has unknown type")
       }
       n := &node.NodeExpImp{Left:$1, Right:$3}
       if  n.Type() != node.TypeInteger {
           yylex.Error("TkImp produces non-integer type")
       }
       $$ = n
     }
   | exp TkKeywordOr exp
     {
       switch $1.Type() {
       case node.TypeString:
           yylex.Error("TkOr left value has string type")
       case node.TypeUnknown:
           yylex.Error("TkOr left value has unknown type")
       }
       switch $3.Type() {
       case node.TypeString:
           yylex.Error("TkOr right value has string type")
       case node.TypeUnknown:
           yylex.Error("TkOr right value has unknown type")
       }
       n := &node.NodeExpOr{Left:$1, Right:$3}
       if  n.Type() != node.TypeInteger {
           yylex.Error("TkOr produces non-integer type")
       }
       $$ = n
     }
   | exp TkKeywordXor exp
     {
       switch $1.Type() {
       case node.TypeString:
           yylex.Error("TkXor left value has string type")
       case node.TypeUnknown:
           yylex.Error("TkXor left value has unknown type")
       }
       switch $3.Type() {
       case node.TypeString:
           yylex.Error("TkXor right value has string type")
       case node.TypeUnknown:
           yylex.Error("TkXor right value has unknown type")
       }
       n := &node.NodeExpXor{Left:$1, Right:$3}
       if  n.Type() != node.TypeInteger {
           yylex.Error("TkXor produces non-integer type")
       }
       $$ = n
     }
   | exp TkEqual exp
     {
       if !node.TypeCompare($1.Type(), $3.Type()) {
           yylex.Error("TkEqual type mismatch")
       }
       $$ = &node.NodeExpEqual{Left:$1, Right:$3}
     }
   | exp TkUnequal exp
     {
       if !node.TypeCompare($1.Type(), $3.Type()) {
           yylex.Error("TkUnequal type mismatch")
       }
       $$ = &node.NodeExpUnequal{Left:$1, Right:$3}
     }
   | exp TkGT exp
     {
       if !node.TypeCompare($1.Type(), $3.Type()) {
           yylex.Error("TkGT type mismatch")
       }
       $$ = &node.NodeExpGT{Left:$1, Right:$3}
     }
   | exp TkLT exp
     {
       if !node.TypeCompare($1.Type(), $3.Type()) {
           yylex.Error("TkLT type mismatch")
       }
       $$ = &node.NodeExpLT{Left:$1, Right:$3}
     }
   | exp TkGE exp
     {
       if !node.TypeCompare($1.Type(), $3.Type()) {
           yylex.Error("TkGE type mismatch")
       }
       $$ = &node.NodeExpGE{Left:$1, Right:$3}
     }
   | exp TkLE exp
     {
       if !node.TypeCompare($1.Type(), $3.Type()) {
           yylex.Error("TkLE type mismatch")
       }
       $$ = &node.NodeExpLE{Left:$1, Right:$3}
     }
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
		case TkUnequal: // do not store
		case TkLT: // do not store
		case TkGT: // do not store
		case TkLE: // do not store
		case TkGE: // do not store
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
		case TkKeywordGoto: // do not store
		case TkKeywordLen: // do not store
		case TkKeywordLet: // do not store
		case TkKeywordList: // do not store
		case TkKeywordMod: // do not store
		case TkKeywordPrint: // do not store
		case TkKeywordNot: // do not store
		case TkKeywordAnd: // do not store
		case TkKeywordEqv: // do not store
		case TkKeywordImp: // do not store
		case TkKeywordOr: // do not store
		case TkKeywordXor: // do not store
		default:
			log.Printf("InputLex.Lex: FIXME token value [%s] not stored for parser actions\n", t.Value)
	}

	return id
}

func (l *InputLex) Error(s string) {
	l.syntaxErrorCount++
	log.Printf("InputLex.Error: count=%d line=%d column=%d: %s\n", l.syntaxErrorCount, l.lex.Line(), l.lex.Column(), s)
}


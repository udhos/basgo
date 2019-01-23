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
	//"strings"

	"github.com/udhos/basgo/baslex"
	"github.com/udhos/basgo/node"
)

type ParserResult struct {
	Root []node.Node
	LineNumbers map[string]node.LineNumber // used by GOTO GOSUB etc
	LibInput bool
	LibReadData bool
	ForStack []*node.NodeFor
	CountFor int
	CountNext int
	ArrayTable map[string]int
}

// parser auxiliary variables
var (
	Result = ParserResult{
		LineNumbers: map[string]node.LineNumber{},
		ArrayTable: map[string]int{},
	}

	nodeListStack [][]node.Node // support nested node lists (1)
	lineList []node.Node
	expListStack [][]node.NodeExp // support nested exp lists (2)
	constList []node.NodeExp
	numberList []string
	identList []string
	lastLineNum string // basic line number for parser error reporting

	// (1) stmt IF-THEN can nest list of stmt: THEN CLS:IF:CLS
	// (2) exp can nest list of exp: array(exp,exp,exp)
)

func Reset() {
	Result = ParserResult{
		LineNumbers: map[string]node.LineNumber{},
		ArrayTable: map[string]int{},
	}

	nodeListStack = [][]node.Node{}
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
	typeExpArray *node.NodeExpArray

	typeRem string
	typeNumber string
	typeFloat string
	typeString string
	typeIdentifier string
	typeRawLine string
	typeNumberList []string
	typeLineNumber string
	typeIdentList []string

	tok int
}

// any non-terminal which returns a value needs a type, which is
// really a field name in the above union struct

%type <typeLineList> line_list
%type <typeLine> line
%type <typeNumber> line_num
%type <typeStmtList> statements
%type <typeStmtList> statements_aux
%type <typeStmt> stmt
%type <typeStmt> stmt_goto
%type <typeStmt> assign
%type <typeExpressions> print_expressions
%type <typeExpressions> array_index_exp_list
%type <typeExp> exp
%type <typeExp> one_const
%type <typeExp> one_const_num
%type <typeExpArray> array_exp
%type <typeExpArray> one_dim
%type <typeNumberList> number_list
%type <typeLineNumber> use_line_number
%type <typeExpressions> const_list
%type <typeExpressions> const_list_num
%type <typeExpressions> dim_list
%type <typeExp> one_var
%type <typeExpressions> var_list

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
%token <tok> TkKeywordData
%token <tok> TkKeywordDim
%token <tok> TkKeywordElse
%token <tok> TkKeywordEnd
%token <tok> TkKeywordFor
%token <tok> TkKeywordGosub
%token <tok> TkKeywordGoto
%token <tok> TkKeywordIf
%token <tok> TkKeywordInput
%token <tok> TkKeywordInt
%token <tok> TkKeywordLeft
%token <tok> TkKeywordLen
%token <tok> TkKeywordLet
%token <tok> TkKeywordList
%token <tok> TkKeywordLoad
%token <tok> TkKeywordNext
%token <tok> TkKeywordOn
%token <tok> TkKeywordPrint
%token <tok> TkKeywordRead
%token <typeRem> TkKeywordRem
%token <tok> TkKeywordReturn
%token <tok> TkKeywordRnd
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
         
	 Result.Root = $1 // save for caller
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

statements_push:
     {
        // create new nested node list
        // because an IF node can hold a nested list of nodes
        nodeListStack = append(nodeListStack, []node.Node{})
     }
     ;

statements_pop:
     {
        // drop nested node list
	last := len(nodeListStack) - 1
	nodeListStack = nodeListStack[:last]
     }
     ;

statements_aux: statements_push statements statements_pop
     {
        $$ = $2
     }
     ;

line_num: TkNumber
     {
       lastLineNum = $1 // save for parser error reporting
       $$ = $1
     };

line: statements_aux
     {
	$$ = &node.LineImmediate{Nodes:$1}
     }
  | line_num statements_aux
     {
       n := $1
       ln, found := Result.LineNumbers[n]
       if found {
         // set defined, keep used unchanged
         ln.Defined = true
         Result.LineNumbers[n] = ln
       } else {
         // set defined, unset used
         Result.LineNumbers[n] = node.LineNumber{Defined: true}
       }
       $$ = &node.LineNumbered{LineNumber:n, Nodes:$2}
     }
  ;

statements: stmt
     {
	last := len(nodeListStack) - 1
	nodeListStack[last] = []node.Node{$1} // init node list
	$$ = nodeListStack[last]
     }
  | statements TkColon stmt
     {
	last := len(nodeListStack) - 1
	nodeListStack[last] = append(nodeListStack[last], $3)
        $$ = nodeListStack[last]
     }
  ;

stmt_goto: use_line_number
    {
       $$ = &node.NodeGoto{Line: $1}
    }
  ;

then_or_goto: TkKeywordThen
           |
           TkKeywordGoto
           ;

one_dim: TkIdentifier bracket_left const_list_num bracket_right
	{
        	name := $1
        	indices := $3
		/*
        	err := node.ArraySetDeclared(Result.ArrayTable, name, len(indices))
        	if err != nil {
	           yylex.Error("error declaring array: " + err.Error())
        	}
		*/
      		$$ = &node.NodeExpArray{Name: name, Indices: indices}
        }
        ;

dim_list: one_dim
	{
		last := len(expListStack) - 1
        	expListStack[last] = []node.NodeExp{$1} // reset dim_list
	        $$ = expListStack[last]
	}
        | dim_list TkComma one_dim
        {
		last := len(expListStack) - 1
		expListStack[last] = append(expListStack[last], $3)
	        $$ = expListStack[last]
	}
        ;

stmt: /* empty */
     { $$ = &node.NodeEmpty{} }
  | TkKeywordEnd
     { $$ = &node.NodeEnd{} }
  | TkKeywordStop
     { $$ = &node.NodeEnd{} }
  | TkKeywordData const_list
     {
        Result.LibReadData = true
        $$ = &node.NodeData{Expressions: $2}
     }
  | TkKeywordDim expressions_push dim_list expressions_pop
     {
        //$$ = &node.NodeDim{Variables: $3}
        $$ = &node.NodeEmpty{}
     }
  | TkKeywordFor one_var TkEqual exp TkKeywordTo exp
     {
	ident := $2
	first := $4
	last := $6
	if !node.TypeNumeric(ident.Type()) {
           yylex.Error("FOR variable must be numeric")
	}
        if !node.TypeNumeric(first.Type()) {
           yylex.Error("FOR first value must be numeric")
        }
        if !node.TypeNumeric(last.Type()) {
           yylex.Error("FOR last value must be numeric")
        }
        f := &node.NodeFor{Index: Result.CountFor, Variable: ident, First: first, Last: last, Step: &node.NodeExpNumber{Value: "1"}}
	Result.CountFor++
	Result.ForStack = append(Result.ForStack, f) // push
        $$ = f
     }
  | TkKeywordFor one_var TkEqual exp TkKeywordTo exp TkKeywordStep exp
     {
	ident := $2
	first := $4
	last := $6
	step := $8
	if !node.TypeNumeric(ident.Type()) {
           yylex.Error("FOR variable must be numeric")
	}
        if !node.TypeNumeric(first.Type()) {
           yylex.Error("FOR first value must be numeric")
        }
        if !node.TypeNumeric(last.Type()) {
           yylex.Error("FOR last value must be numeric")
        }
        if !node.TypeNumeric(step.Type()) {
           yylex.Error("FOR step value must be numeric")
        }
        f := &node.NodeFor{Index: Result.CountFor, Variable: ident, First: first, Last: last, Step: step}
	Result.CountFor++
	Result.ForStack = append(Result.ForStack, f) // push
        $$ = f
     }
  | TkKeywordNext
     {
	var f *node.NodeFor
	stackTop := len(Result.ForStack)-1
	if stackTop < 0 {
           yylex.Error("NEXT without FOR")
	} else {
           f = Result.ForStack[stackTop]
	   Result.ForStack = Result.ForStack[:stackTop] // pop
	}
	Result.CountNext++
        $$ = &node.NodeNext{Fors: []*node.NodeFor{f}}
     }
  | TkKeywordNext expressions_push var_list expressions_pop
     {
        list := $3
        forList := []*node.NodeFor{}
	for _, ident := range list {
	   if !node.TypeNumeric(ident.Type()) {
              yylex.Error("NEXT variable must be numeric: "+ident.String())
              continue
	   }

           stackTop := len(Result.ForStack)-1
           if stackTop < 0 {
              yylex.Error(fmt.Sprintf("NEXT '%s' without FOR", ident.String()))
              continue
           }

           f := Result.ForStack[stackTop]
           forList = append(forList,f)
           Result.ForStack = Result.ForStack[:stackTop] // pop

           if !node.VarMatch(f.Variable.String(), ident.String()) {
              yylex.Error(fmt.Sprintf("FOR var %s mismatches NEXT var %s", f.Variable.String(), ident.String()))
              continue
           }

	   Result.CountNext++
	}
        $$ = &node.NodeNext{Variables: list, Fors: forList}
     }
  | TkKeywordIf exp then_or_goto stmt_goto
     {
       cond := $2
       if !node.TypeNumeric(cond.Type()) {
           yylex.Error("IF condition must be boolean")
       }
       $$ = &node.NodeIf{Cond: cond, Then: []node.Node{$4}, Else: []node.Node{&node.NodeEmpty{}}}
     }
  | TkKeywordIf exp then_or_goto stmt_goto TkKeywordElse stmt_goto
     {
       cond := $2
       if !node.TypeNumeric(cond.Type()) {
           yylex.Error("IF condition must be boolean")
       }
       $$ = &node.NodeIf{Cond: cond, Then: []node.Node{$4}, Else: []node.Node{$6}}
     }
  | TkKeywordIf exp then_or_goto stmt_goto TkKeywordElse statements_aux
     {
       cond := $2
       if !node.TypeNumeric(cond.Type()) {
           yylex.Error("IF condition must be boolean")
       }
       $$ = &node.NodeIf{Cond: cond, Then: []node.Node{$4}, Else: $6}
     }
  | TkKeywordIf exp TkKeywordThen statements_aux
     {
       cond := $2
       if !node.TypeNumeric(cond.Type()) {
           yylex.Error("IF condition must be boolean")
       }
       $$ = &node.NodeIf{Cond: cond, Then: $4, Else: []node.Node{&node.NodeEmpty{}}}
     }
  | TkKeywordIf exp TkKeywordThen statements_aux TkKeywordElse stmt_goto
     {
       cond := $2
       if !node.TypeNumeric(cond.Type()) {
           yylex.Error("IF condition must be boolean")
       }
       $$ = &node.NodeIf{Cond: cond, Then: $4, Else: []node.Node{$6}}
     }
  | TkKeywordIf exp TkKeywordThen statements_aux TkKeywordElse statements_aux
     {
       cond := $2
       if !node.TypeNumeric(cond.Type()) {
           yylex.Error("IF condition must be boolean")
       }
       $$ = &node.NodeIf{Cond: cond, Then: $4, Else: $6}
     }
  | TkKeywordInput TkIdentifier
     {
        Result.LibInput = true
        $$ = &node.NodeInput{Variable: $2}
     }
  | TkKeywordGoto stmt_goto
     { $$ = $2 }
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
  | TkKeywordPrint expressions_push print_expressions expressions_pop
     {
        $$ = &node.NodePrint{Expressions: $3, Newline: true}
     }
  | TkKeywordPrint expressions_push print_expressions TkSemicolon expressions_pop
     {
        $$ = &node.NodePrint{Expressions: $3}
     }
  | TkKeywordPrint expressions_push print_expressions TkComma expressions_pop
     {
        $$ = &node.NodePrint{Expressions: $3, Tab: true}
     }
  | TkKeywordRead expressions_push var_list expressions_pop
     {
        Result.LibReadData = true
        $$ = &node.NodeRead{Variables: $3}
     }
  | TkKeywordRem
     { $$ = &node.NodeRem{Value: $1} }
  | TkKeywordOn exp TkKeywordGoto number_list
     {
       cond := $2
       lines := $4
       if !node.TypeNumeric(cond.Type()) {
           yylex.Error("ON-GOTO condition must be numeric")
       }
       $$ = &node.NodeOnGoto{Cond: cond, Lines: lines}
     }
  ;

expressions_push:
     {
        // create new nested exp list
        // because an exp can hold a nested list of exp
        expListStack = append(expListStack, []node.NodeExp{})
     }
     ;

expressions_pop:
     {
        // drop nested exp list
	last := len(expListStack) - 1
	expListStack = expListStack[:last]
     }
     ;

use_line_number: TkNumber
    {
       n := $1
       ln, found := Result.LineNumbers[n]
       if found {
         // set used, keep defined unchanged
         ln.Used = true
         Result.LineNumbers[n] = ln
       } else {
         // set used, unset defined
         Result.LineNumbers[n] = node.LineNumber{Used: true}
       }
       $$ = n
    }
  ;

number_list: use_line_number
     {
        numberList = []string{$1} // reset list
	$$ = numberList
     }
  | number_list TkComma use_line_number
     {
        numberList = append(numberList, $3)
        $$ = numberList
     }
  ;

one_var: TkIdentifier
     {
        log.Printf("parser.y one_var FIXME LHS identifier should NOT be marked as used in its FindVars method")
        $$ = &node.NodeExpIdentifier{Value:$1}
     }
     | array_exp
     {
        $$ = $1 // node.NodeExpArray
     }
     ;

var_list: one_var
	{
		last := len(expListStack) - 1
        	expListStack[last] = []node.NodeExp{$1} // reset var_list
	        $$ = expListStack[last]
	}
        | var_list TkComma one_var
        {
		last := len(expListStack) - 1
		expListStack[last] = append(expListStack[last], $3)
	        $$ = expListStack[last]
	}
        ;

const_list: one_const
     {
        constList = []node.NodeExp{$1} // reset list
	$$ = constList
     }
  | const_list TkComma one_const
     {
        constList = append(constList, $3)
        $$ = constList
     }
  ;

const_list_num: one_const_num
     {
        constList = []node.NodeExp{$1} // reset list
	$$ = constList
     }
  | const_list_num TkComma one_const_num
     {
        constList = append(constList, $3)
        $$ = constList
     }
  ;

assign: TkIdentifier TkEqual exp
     {
	i := $1
	e := $3
	ti := node.VarType(i)
	te := e.Type()
	if !node.TypeCompare(ti, te) {
           yylex.Error("Assignment type mismatch")
	}
        $$ = &node.NodeAssign{Left: i, Right: e}
     }
  | array_exp TkEqual exp
     {
	a := $1
	e := $3
	ta := a.Type()
	te := e.Type()
	if !node.TypeCompare(ta, te) {
           yylex.Error("Array assignment type mismatch")
	}
        $$ = &node.NodeAssignArray{Left: a, Right: e}
     }
  ;

array_index_exp_list: exp
	{
		e := $1
		if !node.TypeNumeric(e.Type()) {
			yylex.Error("Array index must be numeric")
		}
		last := len(expListStack) - 1
        	expListStack[last] = []node.NodeExp{e} // reset array_index_exp_list
	        $$ = expListStack[last]
	}
    |
        array_index_exp_list TkComma exp
        {
		e := $3
		if !node.TypeNumeric(e.Type()) {
			yylex.Error("Array index must be numeric")
		}
		last := len(expListStack) - 1
		expListStack[last] = append(expListStack[last], e)
	        $$ = expListStack[last]
	}
    ;

print_expressions: exp
	{
		last := len(expListStack) - 1
        	expListStack[last] = []node.NodeExp{$1} // reset print_expressions
	        $$ = expListStack[last]
	}
    |
        print_expressions exp
        {
		last := len(expListStack) - 1
		expListStack[last] = append(expListStack[last], $2)
	        $$ = expListStack[last]
	}
    |
        print_expressions TkComma exp
        {
		last := len(expListStack) - 1
		expListStack[last] = append(expListStack[last], $3)
	        $$ = expListStack[last]
	}
    |
        print_expressions TkSemicolon exp
        {
		last := len(expListStack) - 1
		expListStack[last] = append(expListStack[last], $3)
	        $$ = expListStack[last]
	}
    ;

one_const_num: TkNumber { $$ = &node.NodeExpNumber{Value:$1} }
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
   ;

one_const: one_const_num { $$ = $1 }
   | TkString { $$ = &node.NodeExpString{Value:$1} }
   ;

bracket_left: TkParLeft
              |
              TkBracketLeft
              ;

bracket_right: TkParRight
              |
              TkBracketRight
              ;

array_exp: TkIdentifier bracket_left expressions_push array_index_exp_list expressions_pop bracket_right
   {
      name := $1
      indices := $4
      err := node.ArraySetUsed(Result.ArrayTable, name, len(indices))
      if err != nil {
         yylex.Error("error using array: " + err.Error())
      }
      $$ = &node.NodeExpArray{Name: name,Indices: indices}
   }
   ;

exp: one_const
      { $$ = $1 }
   | TkIdentifier { $$ = &node.NodeExpIdentifier{Value:$1} }
   | array_exp
     {
        $$ = $1
     }
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
   | TkKeywordInt exp
     {
       e := $2
       if !node.TypeNumeric(e.Type()) {
           yylex.Error("INT expression must be numeric")
       }
       $$ = &node.NodeExpInt{Value:e}
     }
   | TkKeywordLeft TkParLeft exp TkComma exp TkParRight
     {
       e1 := $3
       e2 := $5
       if e1.Type() != node.TypeString {
           yylex.Error("LEFT$ value must be string")
       }
       if !node.TypeNumeric(e2.Type()) {
           yylex.Error("LEFT$ size must be numeric")
       }
       $$ = &node.NodeExpLeft{Value:e1, Size:e2}
     }
   | TkKeywordLen exp { $$ = &node.NodeExpLen{Value:$2} }
   | TkKeywordRnd { $$ = &node.NodeExpRnd{Value:&node.NodeExpNumber{Value:"1"}} }
   | TkKeywordRnd exp
     {
       e := $2
       if !node.TypeNumeric(e.Type()) {
           yylex.Error("RND expression must be numeric")
       }
       $$ = &node.NodeExpRnd{Value:e}
     }
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
	lastToken baslex.Token // save token for parser error reporting
}

func (l *InputLex) Errors() int {
	return l.syntaxErrorCount
}

func (l *InputLex) Lex(lval *InputSymType) int {

	if !l.lex.HasToken() {
		return 0 // 0 means real EOF for the parser
	}

	t := l.lex.Next()

	l.lastToken = t // save token for parser error reporting

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
		case TkKeywordData: // do not store
		case TkKeywordDim: // do not store
		case TkKeywordEnd: // do not store
		case TkKeywordElse: // do not store
		case TkKeywordFor: // do not store
		case TkKeywordGoto: // do not store
		case TkKeywordIf: // do not store
		case TkKeywordInput: // do not store
		case TkKeywordInt: // do not store
		case TkKeywordLeft: // do not store
		case TkKeywordLen: // do not store
		case TkKeywordLet: // do not store
		case TkKeywordList: // do not store
		case TkKeywordMod: // do not store
		case TkKeywordNext: // do not store
		case TkKeywordOn: // do not store
		case TkKeywordPrint: // do not store
		case TkKeywordNot: // do not store
		case TkKeywordAnd: // do not store
		case TkKeywordEqv: // do not store
		case TkKeywordImp: // do not store
		case TkKeywordOr: // do not store
		case TkKeywordXor: // do not store
		case TkKeywordRead: // do not store
		case TkKeywordRnd: // do not store
		case TkKeywordStep: // do not store
		case TkKeywordThen: // do not store
		case TkKeywordTo: // do not store
		default:
			log.Printf("InputLex.Lex: FIXME token value [%s] not stored for parser actions\n", t.Value)
	}

	return id
}

func (l *InputLex) Error(s string) {
	l.syntaxErrorCount++
	log.Printf("InputLex.Error: %s", s)
	log.Printf("InputLex.Error: last token: %s [%s]", l.lastToken.Type(), l.lastToken.Value)
	log.Printf("InputLex.Error: basicLine=%s inputLine=%d column=%d", lastLineNum, l.lex.Line(), l.lex.Column())
	log.Printf("InputLex.Error: errors=%d", l.syntaxErrorCount)
}


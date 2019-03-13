%{

package basparser

// header of parser.y -- see lex.go and misc.go

import (
        "fmt"
        //"log"
        "strings"
	"strconv"

        "github.com/udhos/basgo/baslib/file"
        "github.com/udhos/basgo/node"
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
	typeExpArray *node.NodeExpArray
	typeExpString *node.NodeExpString

	typeRem string
	typeNumber string
	typeFloat string
	typeString string
	typeIdentifier string
	typeRawLine string
	typeNumberList []string
	typeLineNumber string
	typeIdentList []string
	typeRangeList [][]string
	typeBool bool

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
%type <typeExpressions> printfile_expressions
%type <typeExpressions> array_index_exp_list
%type <typeExpressions> call_exp_list
%type <typeExpressions> null_exp_list
%type <typeExp> exp
%type <typeExp> null_exp
%type <typeExp> one_const_any
%type <typeExp> one_const_noneg
%type <typeExp> one_const_num_any
%type <typeExp> one_const_num_noneg
%type <typeExp> one_const_int
%type <typeExp> one_const_float
%type <typeExpString> one_const_str
%type <typeExpArray> array_exp
%type <typeExp> array_or_call
%type <typeExpArray> one_dim
%type <typeLineNumber> jump_number
%type <typeNumberList> jump_list
%type <typeLineNumber> use_line_number
%type <typeLineNumber> restore_line_number
%type <typeExpressions> const_list_any
%type <typeExpressions> const_list_num_noneg
%type <typeExpressions> dim_list
%type <typeExpressions> const_list_int
%type <typeExp> one_var
%type <typeExp> single_var
%type <typeExpressions> var_list
%type <typeExpressions> single_var_list
%type <typeExp> file_num
%type <typeExpressions> file_num_list
%type <typeIdentifier> letter_single
%type <typeIdentList> letter_range
%type <typeRangeList> letter_range_list
%type <typeBool> semicolon_optional

// same for terminals

%token <tok> TkNull
%token <typeRawLine> TkEOF
%token <typeRawLine> TkEOL

%token <tok> TkErrInput
%token <tok> TkErrInternal
%token <tok> TkErrInvalid
%token <tok> TkErrLarge

%token <tok> TkHash
%token <tok> TkColon
%token <tok> TkComma
%token <tok> TkSemicolon
%token <tok> TkParLeft
%token <tok> TkParRight
%token <tok> TkBracketLeft
%token <tok> TkBracketRight
%token <typeRem> TkCommentQ
%token <typeString> TkString
%token <typeNumber> TkNumber
%token <typeNumber> TkNumberHex
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

%token <tok> TkKeywordAbs
%token <tok> TkKeywordAsc
%token <tok> TkKeywordBeep
%token <tok> TkKeywordChain
%token <tok> TkKeywordChdir
%token <tok> TkKeywordChr
%token <tok> TkKeywordClear
%token <tok> TkKeywordClose
%token <tok> TkKeywordCls
%token <tok> TkKeywordColor
%token <tok> TkKeywordCommon
%token <tok> TkKeywordCont
%token <tok> TkKeywordCos
%token <tok> TkKeywordCsrlin
%token <tok> TkKeywordData
%token <tok> TkKeywordDate
%token <tok> TkKeywordDef
%token <tok> TkKeywordDefdbl
%token <tok> TkKeywordDefint
%token <tok> TkKeywordDefsng
%token <tok> TkKeywordDefstr
%token <tok> TkKeywordDim
%token <tok> TkKeywordElse
%token <tok> TkKeywordEnd
%token <tok> TkKeywordEnviron
%token <tok> TkKeywordEnvironFunc
%token <tok> TkKeywordEof
%token <tok> TkKeywordErase
%token <tok> TkKeywordError
%token <tok> TkKeywordFiles
%token <tok> TkKeywordFix
%token <tok> TkKeywordFor
%token <tok> TkKeywordGodecl
%token <tok> TkKeywordGofunc
%token <tok> TkKeywordGoimport
%token <tok> TkKeywordGoproc
%token <tok> TkKeywordGosub
%token <tok> TkKeywordGoto
%token <tok> TkKeywordHex
%token <tok> TkKeywordIf
%token <tok> TkKeywordInkey
%token <tok> TkKeywordInput
%token <tok> TkKeywordInputFunc
%token <tok> TkKeywordInstr
%token <tok> TkKeywordInt
%token <tok> TkKeywordKey
%token <tok> TkKeywordKill
%token <tok> TkKeywordLeft
%token <tok> TkKeywordLen
%token <tok> TkKeywordLet
%token <tok> TkKeywordLine
%token <tok> TkKeywordList
%token <tok> TkKeywordLoad
%token <tok> TkKeywordLocate
%token <tok> TkKeywordLof
%token <tok> TkKeywordMid
%token <tok> TkKeywordMkdir
%token <tok> TkKeywordName
%token <tok> TkKeywordNext
%token <tok> TkKeywordNew
%token <tok> TkKeywordOct
%token <tok> TkKeywordOff
%token <tok> TkKeywordOn
%token <tok> TkKeywordOpen
%token <tok> TkKeywordOption
%token <tok> TkKeywordPeek
%token <tok> TkKeywordPlay
%token <tok> TkKeywordPoke
%token <tok> TkKeywordPos
%token <tok> TkKeywordPrint
%token <tok> TkKeywordRandomize
%token <tok> TkKeywordRead
%token <typeRem> TkKeywordRem
%token <tok> TkKeywordReset
%token <tok> TkKeywordRestore
%token <tok> TkKeywordResume
%token <tok> TkKeywordReturn
%token <tok> TkKeywordRight
%token <tok> TkKeywordRmdir
%token <tok> TkKeywordRnd
%token <tok> TkKeywordRun
%token <tok> TkKeywordSave
%token <tok> TkKeywordScreen
%token <tok> TkKeywordSgn
%token <tok> TkKeywordSin
%token <tok> TkKeywordSound
%token <tok> TkKeywordSpace
%token <tok> TkKeywordSpc
%token <tok> TkKeywordSqr
%token <tok> TkKeywordStep
%token <tok> TkKeywordStop
%token <tok> TkKeywordStr
%token <tok> TkKeywordString
%token <tok> TkKeywordSwap
%token <tok> TkKeywordSystem
%token <tok> TkKeywordTab
%token <tok> TkKeywordTan
%token <tok> TkKeywordThen
%token <tok> TkKeywordTime
%token <tok> TkKeywordTimer
%token <tok> TkKeywordTo
%token <tok> TkKeywordUsing
%token <tok> TkKeywordVal
%token <tok> TkKeywordView
%token <tok> TkKeywordWend
%token <tok> TkKeywordWhile
%token <tok> TkKeywordWidth

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

statements_aux: statements_push statements statements_pop { $$ = $2 } ;

line_num: TkNumber
     {
       lastLineNum = $1 // save for parser error reporting
       $$ = $1
     };

comment_q: /* empty */
         | TkCommentQ
         ;

line: statements_aux comment_q
     {
	$$ = &node.LineImmediate{Nodes:$1}
     }
  | line_num statements_aux comment_q
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
           | TkKeywordGoto
           ;

one_dim: TkIdentifier bracket_left const_list_num_noneg bracket_right
	{
        	name := $1
        	indices := $3

		if node.IsFuncName(name) {
	           yylex.Error("DIM array name can't start with DEF FN prefix: " + name)
		}

		strList := []string{}
		for _, e := range indices {
			strList = append(strList, e.String())
		}
        	err := node.ArraySetDeclared(Result.ArrayTable, name, strList)
        	if err != nil {
	           yylex.Error("error declaring array: " + err.Error())
        	}
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

file_num: exp { $$ = $1 }
	| TkHash exp { $$ = $2 }
        ;

file_num_list: file_num
        {
                last := len(expListStack) - 1
                expListStack[last] = []node.NodeExp{$1} // reset file_num_list
                $$ = expListStack[last]
        }
    |
        file_num_list TkComma file_num
        {
                last := len(expListStack) - 1
                expListStack[last] = append(expListStack[last], $3)
                $$ = expListStack[last]
        }
    ;

letter_single: TkIdentifier
	{
		s := strings.ToLower($1)
		if len(s) != 1 {
			yylex.Error("range must be a single letter")
		}
		$$ = s
	}
	;

letter_range: letter_single
	{
		s := $1
		$$ = []string{s,s}
	}
	| letter_single TkMinus letter_single
	{
		s1 := $1
		s2 := $3
		if s1 > s2 {
			yylex.Error("bad range order: first char greater than last char")
		}
		$$ = []string{s1,s2}
	}
	;

letter_range_list: letter_range
	{
		r := $1
        	rangeList = [][]string{r} // reset range list
        	$$ = rangeList
	}
	| letter_range_list TkComma letter_range
	{
		r := $3
        	rangeList = append(rangeList, r)
		$$ = rangeList
	}
	;

semicolon_optional: { $$ = false }
	| TkSemicolon { $$ = true } ;

stmt: /* empty */
     { $$ = &node.NodeEmpty{} }
  | TkKeywordEnd
     { $$ = createEndNode(&Result, "") }
  | TkKeywordStop
     { $$ = createEndNode(&Result, "") }
  | TkKeywordSystem
     { $$ = createEndNode(&Result, "") }
  | TkKeywordData
     {
        line := lastLineNum
	if _, found := Result.RestoreTable[line]; !found {
		Result.RestoreTable[line] = Result.DataOffset
	}
	Result.DataOffset++
        $$ = &node.NodeData{Expressions: []node.NodeExp{node.NewNodeExpStringEmpty()}}
     }
  | TkKeywordData const_list_any
     {
	list := $2
        line := lastLineNum
	if _, found := Result.RestoreTable[line]; !found {
		Result.RestoreTable[line] = Result.DataOffset
	}
	Result.DataOffset += len(list)
        $$ = &node.NodeData{Expressions: list}
     }
  | TkKeywordDef TkIdentifier
     { 
        ident := $2
        var unsup node.Node
        if isSymbol(ident, "SEG") {
           unsup = unsupportedEmpty("DEF SEG")
        } else {
           unsup = unsupportedEmpty("DEF " + ident)
	}
        $$ = unsup
     }
  | TkKeywordDef TkIdentifier TkEqual exp
     {
        ident := $2
        var unsup node.Node
        if isSymbol(ident, "SEG") {
           unsup = unsupportedEmpty("DEF SEG")
        } else {
           unsup = unsupportedEmpty("DEF " + ident)
	}
        $$ = unsup
     }
  | TkKeywordDef TkIdentifier TkParLeft TkParRight TkEqual exp
     {
        i := $2
        e := $6
	if !node.IsFuncName(i) {
           yylex.Error("DEF FN bad function name: " + i)
	}
	if !node.TypeCompare(node.VarType(Result.TypeTable, i), e.Type(Result.TypeTable)) {
           yylex.Error("DEF FN type mismatch")
	}
        n := &node.NodeDefFn{FuncName: i, Body: e}
       	err := node.FuncSetDeclared(Result.FuncTable, n)
        if err != nil {
           yylex.Error("error declaring DEF FN func: " + err.Error())
        }
        $$ = n
     }
  | TkKeywordDef TkIdentifier TkParLeft single_var_list TkParRight TkEqual exp
     {
        i := $2
        list := $4
        e := $7

	if !node.IsFuncName(i) {
           yylex.Error("DEF FN bad function name: " + i)
	}

        dedupVar := map[string]struct{}{}
	for _, v := range list {
                vName := v.String()
		if _, found := dedupVar[vName]; found {
                        yylex.Error("DEF FN dup var: " + vName)
			break
		}
                dedupVar[vName] = struct{}{}
	}
        
	if !node.TypeCompare(node.VarType(Result.TypeTable, i), e.Type(Result.TypeTable)) {
           yylex.Error("DEF FN type mismatch")
	}
        n := &node.NodeDefFn{FuncName: i, Variables: list, Body: e}
       	err := node.FuncSetDeclared(Result.FuncTable, n)
        if err != nil {
           yylex.Error("error declaring DEF FN func: " + err.Error())
        }
        $$ = n
     }
  | TkKeywordDim expressions_push dim_list expressions_pop
     {
        $$ = &node.NodeDim{Arrays: $3}
     }
  | TkKeywordFor one_var TkEqual exp TkKeywordTo exp
     {
	ident := $2
	first := $4
	last := $6
	if !node.TypeNumeric(ident.Type(Result.TypeTable)) {
           yylex.Error("FOR variable must be numeric")
	}
        if !node.TypeNumeric(first.Type(Result.TypeTable)) {
           yylex.Error("FOR first value must be numeric")
        }
        if !node.TypeNumeric(last.Type(Result.TypeTable)) {
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
	if !node.TypeNumeric(ident.Type(Result.TypeTable)) {
           yylex.Error("FOR variable must be numeric")
	}
        if !node.TypeNumeric(first.Type(Result.TypeTable)) {
           yylex.Error("FOR first value must be numeric")
        }
        if !node.TypeNumeric(last.Type(Result.TypeTable)) {
           yylex.Error("FOR last value must be numeric")
        }
        if !node.TypeNumeric(step.Type(Result.TypeTable)) {
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
	   if !node.TypeNumeric(ident.Type(Result.TypeTable)) {
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
       if !node.TypeNumeric(cond.Type(Result.TypeTable)) {
           yylex.Error("IF condition must be boolean")
       }
       Result.CountIf++
       $$ = &node.NodeIf{Index:Result.CountIf, Cond: cond, Then: []node.Node{$4}, Else: []node.Node{&node.NodeEmpty{}}}
     }
  | TkKeywordIf exp then_or_goto stmt_goto TkKeywordElse stmt_goto
     {
       cond := $2
       if !node.TypeNumeric(cond.Type(Result.TypeTable)) {
           yylex.Error("IF condition must be boolean")
       }
       Result.CountIf++
       $$ = &node.NodeIf{Index:Result.CountIf, Cond: cond, Then: []node.Node{$4}, Else: []node.Node{$6}}
     }
  | TkKeywordIf exp then_or_goto stmt_goto TkKeywordElse statements_aux
     {
       cond := $2
       if !node.TypeNumeric(cond.Type(Result.TypeTable)) {
           yylex.Error("IF condition must be boolean")
       }
       Result.CountIf++
       $$ = &node.NodeIf{Index:Result.CountIf, Cond: cond, Then: []node.Node{$4}, Else: $6}
     }
  | TkKeywordIf exp TkKeywordThen statements_aux
     {
       cond := $2
       if !node.TypeNumeric(cond.Type(Result.TypeTable)) {
           yylex.Error("IF condition must be boolean")
       }
       Result.CountIf++
       $$ = &node.NodeIf{Index:Result.CountIf, Cond: cond, Then: $4, Else: []node.Node{&node.NodeEmpty{}}}
     }
  | TkKeywordIf exp TkKeywordThen statements_aux TkKeywordElse stmt_goto
     {
       cond := $2
       if !node.TypeNumeric(cond.Type(Result.TypeTable)) {
           yylex.Error("IF condition must be boolean")
       }
       Result.CountIf++
       $$ = &node.NodeIf{Index:Result.CountIf, Cond: cond, Then: $4, Else: []node.Node{$6}}
     }
  | TkKeywordIf exp TkKeywordThen statements_aux TkKeywordElse statements_aux
     {
       cond := $2
       if !node.TypeNumeric(cond.Type(Result.TypeTable)) {
           yylex.Error("IF condition must be boolean")
       }
       Result.CountIf++
       $$ = &node.NodeIf{Index:Result.CountIf, Cond: cond, Then: $4, Else: $6}
     }
  | TkKeywordClose
     {
       $$ = &node.NodeClose{}
     }
  | TkKeywordClose expressions_push file_num_list expressions_pop
     {
       list := $3
       for _, num := range list {
          if !node.TypeNumeric(num.Type(Result.TypeTable)) {
              yylex.Error("CLOSE file number must be numeric")
          }
       }
       $$ = &node.NodeClose{Numbers: list}
     }
  | TkKeywordPrint TkHash exp TkComma expressions_push printfile_expressions expressions_pop semicolon_optional
     {
       num := $3
       list := $6
       semi := $8

       if !node.TypeNumeric(num.Type(Result.TypeTable)) {
          yylex.Error("PRINT# file number must be numeric")
       }

       Result.Baslib = true
       $$ = &node.NodePrintFile{Number:num, Expressions:list, Newline:!semi}
     }
  | TkKeywordInput TkHash exp TkComma expressions_push var_list expressions_pop
     {
       num := $3
       list := $6

       if !node.TypeNumeric(num.Type(Result.TypeTable)) {
          yylex.Error("INPUT# file number must be numeric")
       }

       Result.Baslib = true
       $$ = &node.NodeInputFile{Number:num, Variables:list}
     }
  | TkKeywordInput expressions_push var_list expressions_pop
     {
        Result.Baslib = true
        $$ = &node.NodeInput{Variables: $3, AddQuestion: true}
     }
  | TkKeywordInput one_const_str TkComma expressions_push var_list expressions_pop
     {
        Result.Baslib = true
        $$ = &node.NodeInput{PromptString:$2, Variables: $5}
     }
  | TkKeywordInput one_const_str TkSemicolon expressions_push var_list expressions_pop
     {
        Result.Baslib = true
        $$ = &node.NodeInput{PromptString:$2, Variables: $5, AddQuestion: true}
     }
  | TkKeywordLine TkKeywordInput expressions_push one_var expressions_pop
     {
        v := $4
        if v.Type(Result.TypeTable) != node.TypeString {
           yylex.Error("LINE INPUT variable must be string")
        }
        Result.Baslib = true
        $$ = &node.NodeInput{Variables: []node.NodeExp{v}}
     }
  | TkKeywordLine TkKeywordInput one_const_str TkComma expressions_push one_var expressions_pop
     {
        v := $6
        if v.Type(Result.TypeTable) != node.TypeString {
           yylex.Error("LINE INPUT variable must be string")
        }
        Result.Baslib = true
        $$ = &node.NodeInput{PromptString:$3, Variables: []node.NodeExp{v}}
     }
  | TkKeywordLine TkKeywordInput one_const_str TkSemicolon expressions_push one_var expressions_pop
     {
        v := $6
        if v.Type(Result.TypeTable) != node.TypeString {
           yylex.Error("LINE INPUT variable must be string")
        }
        Result.Baslib = true
        $$ = &node.NodeInput{PromptString:$3, Variables: []node.NodeExp{v}}
     }
  | TkKeywordGosub use_line_number
     {
        Result.LibGosubReturn = true
        g := &node.NodeGosub{Index: Result.CountGosub, Line: $2}
        Result.CountGosub++
        $$ = g
     }
  | TkKeywordGoto stmt_goto { $$ = $2 }
  | TkKeywordLet assign { $$ = $2 }
  | assign { $$ = $1 }
  | TkKeywordList { $$ = &node.NodeList{} }
  | TkKeywordName exp TkIdentifier exp
	{
		// NAME "a" AS "b"
		e1 := $2
		labelAs := $3
		e2 := $4
		if e1.Type(Result.TypeTable) != node.TypeString {
			yylex.Error("NAME old filename must be string")
		}
		if e2.Type(Result.TypeTable) != node.TypeString {
			yylex.Error("NAME new filename must be string")
		}
		if !isSymbol(labelAs, "AS") {
			yylex.Error("NAME expecting 'AS', found: " + labelAs)
		}
		$$ = &node.NodeName{From:e1, To:e2}
	}
  | TkKeywordOpen exp TkKeywordFor TkKeywordInput TkIdentifier file_num
     {
        // OPEN "arq" FOR INPUT AS 1
	filename := $2
	labelAs := $5
	num := $6

	if filename.Type(Result.TypeTable) != node.TypeString {
           yylex.Error("OPEN filename must be string")
	}
        if !isSymbol(labelAs, "AS") {
           yylex.Error("OPEN expecting 'AS', found: " + labelAs)
        }
	if !node.TypeNumeric(num.Type(Result.TypeTable)) {
           yylex.Error("OPEN file number must be numeric")
	}
	
        Result.Baslib = true
        $$ = &node.NodeOpen{File:filename, Number:num, Mode:file.OpenInput}
     }
  | TkKeywordOpen exp TkKeywordFor TkIdentifier TkIdentifier file_num
     {
        // OPEN "arq" FOR OUTPUT AS 1
	filename := $2
	mode := $4
	labelAs := $5
	num := $6
        var m int

	if filename.Type(Result.TypeTable) != node.TypeString {
           yylex.Error("OPEN filename must be string")
	}
	switch strings.ToLower(mode) {
           case "output":
              m = file.OpenOutput
           default:
              yylex.Error("OPEN unexpected mode: " + mode)
        }
        if !isSymbol(labelAs, "AS") {
           yylex.Error("OPEN expecting 'AS', found: " + labelAs)
        }
	if !node.TypeNumeric(num.Type(Result.TypeTable)) {
           yylex.Error("OPEN file number must be numeric")
	}
	
        Result.Baslib = true
        $$ = &node.NodeOpen{File:filename, Number:num, Mode:m}
     }
  | TkKeywordPrint
     {
        Result.Baslib = true
        $$ = &node.NodePrint{Newline: true}
     }
  | TkKeywordPrint expressions_push print_expressions expressions_pop
     {
        Result.Baslib = true
        $$ = &node.NodePrint{Expressions: $3, Newline: true}
     }
  | TkKeywordPrint expressions_push print_expressions TkSemicolon expressions_pop
     {
        Result.Baslib = true
        $$ = &node.NodePrint{Expressions: $3}
     }
  | TkKeywordPrint expressions_push print_expressions TkComma expressions_pop
     {
        Result.Baslib = true
        $$ = &node.NodePrint{Expressions: $3, Tab: true}
     }
  | TkKeywordRead expressions_push var_list expressions_pop
     {
        Result.LibReadData = true
        $$ = &node.NodeRead{Variables: $3}
     }
  | TkCommentQ { $$ = &node.NodeRem{Value: $1} }
  | TkKeywordRem { $$ = &node.NodeRem{Value: $1} }
  | TkKeywordRestore
     {
       Result.LibReadData = true
       $$ = &node.NodeRestore{}
     }
  | TkKeywordRestore restore_line_number
     {
       Result.LibReadData = true
       $$ = &node.NodeRestore{Line: $2}
     }
  | TkKeywordReturn
     {
       Result.LibGosubReturn = true
       Result.CountReturn++ // this RETURN jumps to label generated by GOSUB
       $$ = &node.NodeReturn{}
     }
  | TkKeywordReturn use_line_number
     {
       // this return DOES NOT jump to label generated by GOSUB
       Result.LibGosubReturn = true
       $$ = &node.NodeReturn{Line: $2}
     }
  | TkKeywordRun { $$ = unsupportedEnd(&Result, "RUN") }
  | TkKeywordRun use_line_number { $$ = unsupportedEnd(&Result, "RUN") }
  | TkKeywordRun use_line_number TkComma single_var { $$ = unsupportedEnd(&Result, "RUN") }
  | TkKeywordRun one_const_str { $$ = unsupportedEnd(&Result, "RUN") }
  | TkKeywordRun one_const_str TkComma single_var { $$ = unsupportedEnd(&Result, "RUN") }
  | TkKeywordOn exp TkKeywordGosub jump_list
     {
       cond := $2
       lines := $4
       if !node.TypeNumeric(cond.Type(Result.TypeTable)) {
           yylex.Error("ON-GOSUB condition must be numeric")
       }
       g := &node.NodeOnGosub{Index: Result.CountGosub, Cond: cond, Lines: lines}
       Result.CountGosub++
       $$ = g
     }
  | TkKeywordOn exp TkKeywordGoto jump_list
     {
       cond := $2
       lines := $4
       if !node.TypeNumeric(cond.Type(Result.TypeTable)) {
           yylex.Error("ON-GOTO condition must be numeric")
       }
       $$ = &node.NodeOnGoto{Cond: cond, Lines: lines}
     }
  | TkKeywordWhile exp
     {
       cond := $2
       if !node.TypeNumeric(cond.Type(Result.TypeTable)) {
           yylex.Error("WHILE condition must be boolean")
       }
       while := &node.NodeWhile{Cond: cond, Index: Result.CountWhile}
       Result.CountWhile++
       Result.WhileStack = append(Result.WhileStack, while) // push
       $$ = while
     }
  | TkKeywordWend
     {
	var while *node.NodeWhile
	last := len(Result.WhileStack)-1
	if last < 0 {
           yylex.Error("WEND without WHILE")
	} else {
           while = Result.WhileStack[last]
	   Result.WhileStack = Result.WhileStack[:last] // pop
	}
	Result.CountWend++
        $$ = &node.NodeWend{While: while}
     }
  | TkKeywordSwap one_var TkComma one_var
     {
	v1 := $2
	v2 := $4
	if v1.Type(Result.TypeTable) != v2.Type(Result.TypeTable) {
           yylex.Error("SWAP type mismatch")
	}
        $$ = &node.NodeSwap{Var1: v1, Var2: v2}
     }
   | TkKeywordGodecl TkParLeft one_const_str TkParRight
     {
       decl := $3
       Result.Declarations = append(Result.Declarations, decl.Value)
       $$ = &node.NodeGodecl{Value:decl}
     }
   | TkKeywordGoimport TkParLeft one_const_str TkParRight
     {
       imp := $3
       Result.Imports[imp.Value] = struct{}{}       
       $$ = &node.NodeGoimport{Value:imp}
     }
   | TkKeywordGoproc TkParLeft one_const_str TkParRight
     {
       $$ = &node.NodeGoproc{ProcName: $3}
     }
   | TkKeywordGoproc TkParLeft one_const_str TkComma expressions_push call_exp_list expressions_pop TkParRight
     {
       $$ = &node.NodeGoproc{ProcName: $3, Arguments: $6}
     }
   | TkKeywordRandomize
     {
       Result.Baslib = true
       $$ = &node.NodeRandomize{}
     }
   | TkKeywordRandomize exp
     {
       seed := $2
       if !node.TypeNumeric(seed.Type(Result.TypeTable)) {
           yylex.Error("RANDOMIZE seed must be numeric")
       }
       Result.Baslib = true
       $$ = &node.NodeRandomize{Seed: seed}
     }
  | TkKeywordKey TkKeywordOn { $$ = unsupportedEmpty("KEY")  }
  | TkKeywordKey TkKeywordOff { $$ = unsupportedEmpty("KEY") }
  | TkKeywordKey TkKeywordList { $$ = unsupportedEmpty("KEY")  }
  | TkKeywordKey exp TkComma exp { $$ = unsupportedEmpty("KEY")  }
  | TkKeywordBeep { $$ = unsupportedEmpty("BEEP") }
  | TkKeywordCls
	{
		Result.Baslib = true
		$$ = &node.NodeCls{}
	}
  | TkKeywordWidth exp
	{
		w := $2
		if !node.TypeNumeric(w.Type(Result.TypeTable)) {
			yylex.Error("WIDTH value must be numeric")
		}
		Result.Baslib = true
		$$ = &node.NodeWidth{Width: w}
	}
  | TkKeywordDefdbl letter_range_list
	{
		list := $2
		defineTypeRange(&Result, list, node.TypeDouble) // DEFDBL
		$$ = &node.NodeEmpty{Value: "DEFDBL"}
	}
  | TkKeywordDefint letter_range_list
	{
		list := $2
		defineTypeRange(&Result, list, node.TypeInteger) // DEFINT
		$$ = &node.NodeEmpty{Value: "DEFINT"}
	}
  | TkKeywordDefsng letter_range_list
	{
		list := $2
		defineTypeRange(&Result, list, node.TypeFloat) // DEFSNG
		$$ = &node.NodeEmpty{Value: "DEFSNG"}
	}
  | TkKeywordDefstr letter_range_list
	{
		list := $2
		defineTypeRange(&Result, list, node.TypeString) // DEFSTR
		$$ = &node.NodeEmpty{Value: "DEFSTR"}
	}
  | TkKeywordChain expressions_push call_exp_list expressions_pop { $$ = unsupportedEmpty("CHAIN") }
  | TkKeywordClear { $$ = unsupportedEmpty("CLEAR") }
  | TkKeywordClear expressions_push call_exp_list expressions_pop { $$ = unsupportedEmpty("CLEAR") }
  | TkKeywordColor expressions_push null_exp_list expressions_pop
	{
		list := $3
		if len(list) < 1 {
			yylex.Error("Missing COLOR arguments")
		}
		var fg,bg node.NodeExp
		if f := list[0]; f != nodeExpNull {
			if !node.TypeNumeric(f.Type(Result.TypeTable)) {
				yylex.Error("COLOR foreground be numeric")
			}
			fg = f
		}
		if len(list) > 1 {
			if b := list[1]; b != nodeExpNull {
				if !node.TypeNumeric(b.Type(Result.TypeTable)) {
					yylex.Error("COLOR background must be numeric")
				}
				bg = b 
			}
		}
		Result.Baslib = true
		$$ = &node.NodeColor{Foreground: fg, Background: bg}
	}
  | TkKeywordErase single_var_list
	{
		list := $2

		for _, a := range list {
			name := a.String()
			_, found := Result.ArrayTable[strings.ToLower(name)]
			if !found {
				yylex.Error("ERASE array not declared: '" + name + "'")
			}
		}

		$$ = &node.NodeErase{Arrays: list}
	}
  | TkKeywordCommon expressions_push common_var_list expressions_pop { $$ = unsupportedEmpty("COMMON") }
  | TkKeywordLocate expressions_push null_exp_list expressions_pop
	{
		list := $3
		if len(list) < 1 {
			yylex.Error("Missing LOCATE arguments")
		}
		var row,col,cursor node.NodeExp
		if r := list[0]; r != nodeExpNull {
			if !node.TypeNumeric(r.Type(Result.TypeTable)) {
				yylex.Error("LOCATE row must be numeric")
			}
			row = r
		}
		if len(list) > 1 {
			if c := list[1]; c != nodeExpNull {
				if !node.TypeNumeric(c.Type(Result.TypeTable)) {
					yylex.Error("LOCATE column must be numeric")
				}
				col = c 
			}
		}
		if len(list) > 2 {
			if cur := list[2]; cur != nodeExpNull {
				if !node.TypeNumeric(cur.Type(Result.TypeTable)) {
					yylex.Error("LOCATE cursor must be numeric")
				}
				cursor = cur 
			}
		}
		Result.Baslib = true
		$$ = &node.NodeLocate{Row: row, Col: col, Cursor: cursor}
	}
  | TkKeywordNew { $$ = unsupportedEnd(&Result, "NEW") }
  | TkKeywordOn TkKeywordError TkKeywordGoto TkNumber { $$ = unsupportedEmpty("ON-ERROR-GOTO") }
  | TkKeywordPlay exp { $$ = unsupportedEmpty("PLAY") }
  | TkKeywordPoke TkParLeft exp TkComma exp TkParRight { $$ = unsupportedEmpty("POKE") }
  | TkKeywordPoke exp TkComma exp { $$ = unsupportedEmpty("POKE") }
  | TkKeywordReset
	{
		$$ = &node.NodeClose{Label:"RESET"}
	}
  | TkKeywordResume { $$ = unsupportedEmpty("RESUME") }
  | TkKeywordResume TkNumber { $$ = unsupportedEmpty("RESUME") }
  | TkKeywordResume TkKeywordNext { $$ = unsupportedEmpty("RESUME") }
  | TkKeywordScreen const_list_int
	{
		list := $2
		m := list[0].String()
		mode, errInt := strconv.Atoi(m)
		if errInt != nil {
			yylex.Error("Bad screen mode: " + errInt.Error())
		}
		if mode != 0 {
			yylex.Error("Only screen mode 0 is supported")
		}
        	Result.Baslib = true
		$$ = &node.NodeScreen{Mode: mode}
	}
  | TkKeywordSound exp TkComma exp { $$ = unsupportedEmpty("SOUND") }
  | TkKeywordView TkKeywordPrint
	{
        	Result.Baslib = true
		$$ = &node.NodeViewPrint{}
	}
  | TkKeywordView TkKeywordPrint exp TkKeywordTo exp
	{
		top := $3
		bottom := $5
		if !node.TypeNumeric(top.Type(Result.TypeTable)) {
			yylex.Error("VIEW PRINT top line must be numeric")
		}
		if !node.TypeNumeric(bottom.Type(Result.TypeTable)) {
			yylex.Error("VIEW PRINT bottom line must be numeric")
		}
        	Result.Baslib = true
		$$ = &node.NodeViewPrint{Top: top, Bottom: bottom}
	}
  | TkKeywordFiles
	{
        	Result.Baslib = true
		$$ = &node.NodeFiles{Pattern: node.NewNodeExpString("*")}
	}
  | TkKeywordFiles exp
	{
		pattern := $2
		if pattern.Type(Result.TypeTable) != node.TypeString {
			yylex.Error("FILES pattern must be string")
		}
        	Result.Baslib = true
		$$ = &node.NodeFiles{Pattern: pattern}
	}
  | TkKeywordEnviron exp
	{
		str := $2
		if str.Type(Result.TypeTable) != node.TypeString {
			yylex.Error("ENVIRON value must be string")
		}
        	Result.Baslib = true
		$$ = &node.NodeEnviron{Value: str}
	}
  | TkKeywordKill exp
	{
		str := $2
		if str.Type(Result.TypeTable) != node.TypeString {
			yylex.Error("KILL value must be string")
		}
        	Result.Baslib = true
		$$ = &node.NodeKill{Value: str}
	}
  | TkKeywordOption TkIdentifier exp
     { 
        ident := $2
        var unsup node.Node
        if isSymbol(ident, "BASE") {
           unsup = unsupportedEmpty("OPTION BASE")
        } else {
	   yylex.Error("OPTION unknown keyword: " + ident)
           unsup = unsupportedEmpty("OPTION " + ident)
	}
        $$ = unsup
     }
  | TkKeywordChdir exp
	{
		str := $2
		if str.Type(Result.TypeTable) != node.TypeString {
			yylex.Error("CHDIR directory must be string")
		}
        	Result.Baslib = true
		$$ = &node.NodeChdir{Value: str}
	}
  | TkKeywordMkdir exp
	{
		str := $2
		if str.Type(Result.TypeTable) != node.TypeString {
			yylex.Error("MKDIR directory must be string")
		}
        	Result.Baslib = true
		$$ = &node.NodeMkdir{Value: str}
	}
  | TkKeywordRmdir exp
	{
		str := $2
		if str.Type(Result.TypeTable) != node.TypeString {
			yylex.Error("RMDIR directory must be string")
		}
        	Result.Baslib = true
		$$ = &node.NodeRmdir{Value: str}
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
         // set used
         Result.LineNumbers[n] = node.LineNumber{Used: true}
       }
       $$ = n
    }
  ;

restore_line_number: TkNumber
    {
       n := $1
       ln, found := Result.LineNumbers[n]
       if found {
         // set used, keep defined unchanged
         ln.UsedRestore = true
         Result.LineNumbers[n] = ln
       } else {
         // set used
         Result.LineNumbers[n] = node.LineNumber{UsedRestore: true}
       }
       $$ = n
    }
  ;

jump_number: /* empty */ { $$ = "" }
	| use_line_number { $$ = $1 }
	;

jump_list: jump_number
     {
        numberList = []string{$1} // reset jump_list
	$$ = numberList
     }
  | jump_list TkComma jump_number
     {
        numberList = append(numberList, $3)
        $$ = numberList
     }
  ;

common_var: TkIdentifier | TkIdentifier TkParLeft TkParRight ;

common_var_list: common_var | common_var_list TkComma common_var ;

single_var: TkIdentifier
	{
		i := $1
		$$ = node.NewNodeExpIdent(Result.TypeTable, i)
	}
	;

single_var_list: single_var
	{
        	varList = []node.NodeExp{$1} // reset single_var_list
	        $$ = varList
	}
        | single_var_list TkComma single_var
        {
		varList = append(varList, $3)
	        $$ = varList
	}
        ;

one_var: single_var { $$ = $1 }
     | array_exp { $$ = $1 /* node.NodeExpArray */ }
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

const_list_any: one_const_any
     {
        constList = []node.NodeExp{$1} // reset list
	$$ = constList
     }
  | const_list_any TkComma one_const_any
     {
        constList = append(constList, $3)
        $$ = constList
     }
  ;

const_list_num_noneg: one_const_num_noneg
     {
        constList = []node.NodeExp{$1} // reset list
	$$ = constList
     }
  | const_list_num_noneg TkComma one_const_num_noneg
     {
        constList = append(constList, $3)
        $$ = constList
     }
  ;

const_list_int: one_const_int
     {
        constList = []node.NodeExp{$1} // reset list
	$$ = constList
     }
  | const_list_int TkComma one_const_int
     {
        constList = append(constList, $3)
        $$ = constList
     }
  ;

assign: single_var TkEqual exp
     {
	i := $1
	e := $3
	ti := i.Type(Result.TypeTable)
	te := e.Type(Result.TypeTable)
	if !node.TypeCompare(ti, te) {
           yylex.Error("Assignment type mismatch: " + 
		fmt.Sprintf("%s = %s | ", i.String(), e.String()) +
		fmt.Sprintf("%s = %s", node.TypeLabel(ti), node.TypeLabel(te)))
	}
        $$ = &node.NodeAssign{Left: i, Right: e}
     }
  | array_exp TkEqual exp
     {
	a := $1
	e := $3
	ta := a.Type(Result.TypeTable)
	te := e.Type(Result.TypeTable)
	if !node.TypeCompare(ta, te) {
           yylex.Error("Array assignment type mismatch")
	}
        $$ = &node.NodeAssignArray{Left: a, Right: e}
     }
  ;

array_index_exp_list: exp
	{
		e := $1
		if !node.TypeNumeric(e.Type(Result.TypeTable)) {
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
		if !node.TypeNumeric(e.Type(Result.TypeTable)) {
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

printfile_expressions: exp
	{
		last := len(expListStack) - 1
        	expListStack[last] = []node.NodeExp{$1} // reset print_expressions
	        $$ = expListStack[last]
	}
    |
        printfile_expressions TkComma exp
        {
		last := len(expListStack) - 1
		expListStack[last] = append(expListStack[last], $3)
	        $$ = expListStack[last]
	}
    ;

one_const_int: TkNumber
    {
        str := $1

        // str->int->str: make sure it can be used as literal int const in Go source code
        num, err := strconv.Atoi(str)
        if err != nil {
            yylex.Error("error parsing number: "+err.Error())
        }
	str = strconv.Itoa(num)

        $$ = &node.NodeExpNumber{Value:str}
    }
    | TkNumberHex
    {
        str := $1

	if len(str) < 3 {
            yylex.Error("short hex number: "+str)
	} else {
            str = str[2:]
	}

        // str->int->str: make sure it can be used as literal int const in Go source code
        num, err := strconv.ParseInt(str, 16, 64)
        if err != nil {
            yylex.Error("error parsing hex number: "+err.Error())
        }
	maxInt := int64(^uint(0) >> 1)
	if num > maxInt {
            yylex.Error("hex number overflow: " + str)
	}
	str = strconv.Itoa(int(num))

        $$ = &node.NodeExpNumber{Value:str}
    }
    ;

one_const_float: TkFloat
     {
       n := &node.NodeExpFloat{}
       v := $1
       if v != "." {
 
         // 1e => 1e0
         last := len(v) - 1
         if v[last] == 'e' || v[last] == 'E' {
           v += "0"
         }

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

one_const_num_any: one_const_int { $$ = $1 }
   | TkPlus one_const_int { $$ = &node.NodeExpUnaryPlus{Value:$2} }
   | TkMinus one_const_int { $$ = &node.NodeExpUnaryMinus{Value:$2} }
   | one_const_float { $$ = $1 }
   | TkPlus one_const_float { $$ = &node.NodeExpUnaryPlus{Value:$2} }
   | TkMinus one_const_float { $$ = &node.NodeExpUnaryMinus{Value:$2} }
   ;

one_const_num_noneg: one_const_int { $$ = $1 }
   | one_const_float { $$ = $1 }
   ;

one_const_str: TkString { $$ = node.NewNodeExpString($1) } ;

one_const_any: one_const_num_any { $$ = $1 }
   | one_const_str { $$ = $1 }
   ;

one_const_noneg: one_const_num_noneg { $$ = $1 }
   | one_const_str { $$ = $1 }
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

array_or_call: TkIdentifier TkBracketLeft expressions_push array_index_exp_list expressions_pop TkBracketRight
   {
      // square bracket is array-only
      name := $1
      indices := $4
      err := node.ArraySetUsed(Result.ArrayTable, name, len(indices))
      if err != nil {
         yylex.Error("error using array: " + err.Error())
      }
      $$ = &node.NodeExpArray{Name: name,Indices: indices}
   }
   | TkIdentifier TkParLeft TkParRight
   {
      //
      // function call
      //
      name := $1
      err := node.FuncSetUsed(Result.FuncTable, name, nil, Result.TypeTable)
      if err != nil {
         yylex.Error("error using DEF FN: " + err.Error())
      }
      $$ = &node.NodeExpFuncCall{Name: name}
   }
   | TkIdentifier TkParLeft expressions_push call_exp_list expressions_pop TkParRight
   {
      //
      // round bracket is either array or function call
      //
      var n node.NodeExp
      list := $4
      name := $1
      if node.IsFuncName(name) {
         //
         // function call
         //
         err := node.FuncSetUsed(Result.FuncTable, name, list, Result.TypeTable)
         if err != nil {
            yylex.Error("error using DEF FN: " + err.Error())
         }
         n = &node.NodeExpFuncCall{Name: name,Parameters: list}
      } else {
         //
         // array
         //
         indices := $4
         for _, i := range indices {
            if !node.TypeNumeric(i.Type(Result.TypeTable)) {
               yylex.Error("array index must be numeric")
            }
         }
         err := node.ArraySetUsed(Result.ArrayTable, name, len(indices))
         if err != nil {
            yylex.Error("error using array: " + err.Error())
         }
         n = &node.NodeExpArray{Name: name,Indices: list}
      }
      $$ = n
   }
   ;

call_exp_list: exp
	{
		last := len(expListStack) - 1
        	expListStack[last] = []node.NodeExp{$1} // reset call_exp_list
	        $$ = expListStack[last]
	}
    |
        call_exp_list TkComma exp
        {
		last := len(expListStack) - 1
		expListStack[last] = append(expListStack[last], $3)
	        $$ = expListStack[last]
	}
    ;

null_exp: /* empty */
		{ $$ = nodeExpNull }
	| exp
		{ $$ = $1 }
	;

null_exp_list: null_exp
	{
		last := len(expListStack) - 1
        	expListStack[last] = []node.NodeExp{$1} // reset null_exp_list
	        $$ = expListStack[last]
	}
    |
        call_exp_list TkComma null_exp
        {
		last := len(expListStack) - 1
		expListStack[last] = append(expListStack[last], $3)
	        $$ = expListStack[last]
	}
    ;

exp: one_const_noneg { $$ = $1 }
   | TkIdentifier { $$ = node.NewNodeExpIdent(Result.TypeTable, $1) }
   | array_or_call { $$ = $1 }
   | exp TkPlus exp
     {
       e1 := $1
       e2 := $3
       if e1.Type(Result.TypeTable) == node.TypeString && e2.Type(Result.TypeTable) != node.TypeString {
           yylex.Error("TkPlus string and non-string")
       }
       if e1.Type(Result.TypeTable) != node.TypeString && e2.Type(Result.TypeTable) == node.TypeString {
           yylex.Error("TkPlus non-string and string")
       }
       n := &node.NodeExpPlus{Left: e1, Right: e2}
       if n.Type(Result.TypeTable) == node.TypeUnknown {
           yylex.Error("TkPlus produces unknown type")
       }
       $$ = n
     }
   | exp TkMinus exp
     {
       switch $1.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("TkMinus left value has string type")
       case node.TypeUnknown:
           yylex.Error("TkMinus left value has unknown type")
       }
       switch $3.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("TkMinus right value has string type")
       case node.TypeUnknown:
           yylex.Error("TkMinus right value has unknown type")
       }
       n := &node.NodeExpMinus{Left: $1, Right: $3}
       switch n.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("TkMinus produces string type")
       case node.TypeUnknown:
           yylex.Error("TkMinus produces unknown type")
       }
       $$ = n
     }
   | exp TkKeywordMod exp
     {
       switch $1.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("TkMod left value has string type")
       case node.TypeUnknown:
           yylex.Error("TkMod left value has unknown type")
       }
       switch $3.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("TkMod right value has string type")
       case node.TypeUnknown:
           yylex.Error("TkMod right value has unknown type")
       }
       n := &node.NodeExpMod{Left: $1, Right: $3}
       if  n.Type(Result.TypeTable) != node.TypeInteger {
           yylex.Error("TkMod produces non-integer type")
       }
       $$ = n
     }
   | exp TkBackSlash exp
     {
       switch $1.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("Integer division left value has string type")
       case node.TypeUnknown:
           yylex.Error("Integer division left value has unknown type")
       }
       switch $3.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("Integer division right value has string type")
       case node.TypeUnknown:
           yylex.Error("Integer division right value has unknown type")
       }
       n := &node.NodeExpDivInt{Left: $1, Right: $3}
       if  n.Type(Result.TypeTable) != node.TypeInteger {
           yylex.Error("Integer division produces non-integer type")
       }
       $$ = n
     }
   | exp TkMult exp
     {
       switch $1.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("TkMult left value has string type")
       case node.TypeUnknown:
           yylex.Error("TkMult left value has unknown type")
       }
       switch $3.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("TkMult right value has string type")
       case node.TypeUnknown:
           yylex.Error("TkMult right value has unknown type")
       }
       n := &node.NodeExpMult{Left: $1, Right: $3}
       switch n.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("TkMult produces string type")
       case node.TypeUnknown:
           yylex.Error("TkMult produces unknown type")
       }
       $$ = n
     }
   | exp TkDiv exp
     {
       switch $1.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("TkDiv left value has string type")
       case node.TypeUnknown:
           yylex.Error("TkDiv left value has unknown type")
       }
       switch $3.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("TkDiv right value has string type")
       case node.TypeUnknown:
           yylex.Error("TkDiv right value has unknown type")
       }
       n := &node.NodeExpDiv{Left: $1, Right: $3}
       if  n.Type(Result.TypeTable) != node.TypeFloat {
           yylex.Error("TkDiv produces non-float type")
       }
       $$ = n
     }
   | exp TkPow exp
     {
       switch $1.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("TkPow left value has string type")
       case node.TypeUnknown:
           yylex.Error("TkPow left value has unknown type")
       }
       switch $3.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("TkPow right value has string type")
       case node.TypeUnknown:
           yylex.Error("TkPow right value has unknown type")
       }
       n := &node.NodeExpPow{Left: $1, Right: $3}
       if  n.Type(Result.TypeTable) != node.TypeFloat {
           yylex.Error("TkPow produces non-float type")
       }
       $$ = n
     }
   | TkPlus exp %prec UnaryPlus
     {
       switch $2.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("Unary plus has string type")
       case node.TypeUnknown:
           yylex.Error("Unary plus has unknown type")
       }
       $$ = &node.NodeExpUnaryPlus{Value:$2}
     }
   | TkMinus exp %prec UnaryMinus
     {
       switch $2.Type(Result.TypeTable) {
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
       e := $2
       switch e.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("Not has string type")
       case node.TypeUnknown:
           yylex.Error("Not has unknown type")
       }
       $$ = &node.NodeExpNot{Value:e}
     }
   | exp TkKeywordAnd exp
     {
       switch $1.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("TkAnd left value has string type")
       case node.TypeUnknown:
           yylex.Error("TkAnd left value has unknown type")
       }
       switch $3.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("TkAnd right value has string type")
       case node.TypeUnknown:
           yylex.Error("TkAnd right value has unknown type")
       }
       n := &node.NodeExpAnd{Left:$1, Right:$3}
       if  n.Type(Result.TypeTable) != node.TypeInteger {
           yylex.Error("TkAnd produces non-integer type")
       }
       $$ = n
     }
   | exp TkKeywordEqv exp
     {
       switch $1.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("TkEqv left value has string type")
       case node.TypeUnknown:
           yylex.Error("TkEqv left value has unknown type")
       }
       switch $3.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("TkEqv right value has string type")
       case node.TypeUnknown:
           yylex.Error("TkEqv right value has unknown type")
       }
       n := &node.NodeExpEqv{Left:$1, Right:$3}
       if  n.Type(Result.TypeTable) != node.TypeInteger {
           yylex.Error("TkEqv produces non-integer type")
       }
       $$ = n
     }
   | exp TkKeywordImp exp
     {
       switch $1.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("TkImp left value has string type")
       case node.TypeUnknown:
           yylex.Error("TkImp left value has unknown type")
       }
       switch $3.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("TkImp right value has string type")
       case node.TypeUnknown:
           yylex.Error("TkImp right value has unknown type")
       }
       n := &node.NodeExpImp{Left:$1, Right:$3}
       if  n.Type(Result.TypeTable) != node.TypeInteger {
           yylex.Error("TkImp produces non-integer type")
       }
       $$ = n
     }
   | exp TkKeywordOr exp
     {
       switch $1.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("TkOr left value has string type")
       case node.TypeUnknown:
           yylex.Error("TkOr left value has unknown type")
       }
       switch $3.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("TkOr right value has string type")
       case node.TypeUnknown:
           yylex.Error("TkOr right value has unknown type")
       }
       n := &node.NodeExpOr{Left:$1, Right:$3}
       if  n.Type(Result.TypeTable) != node.TypeInteger {
           yylex.Error("TkOr produces non-integer type")
       }
       $$ = n
     }
   | exp TkKeywordXor exp
     {
       switch $1.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("TkXor left value has string type")
       case node.TypeUnknown:
           yylex.Error("TkXor left value has unknown type")
       }
       switch $3.Type(Result.TypeTable) {
       case node.TypeString:
           yylex.Error("TkXor right value has string type")
       case node.TypeUnknown:
           yylex.Error("TkXor right value has unknown type")
       }
       n := &node.NodeExpXor{Left:$1, Right:$3}
       if  n.Type(Result.TypeTable) != node.TypeInteger {
           yylex.Error("TkXor produces non-integer type")
       }
       $$ = n
     }
   | exp TkEqual exp
     {
       e1 := $1
       e2 := $3
       t1 := e1.Type(Result.TypeTable)
       t2 := e2.Type(Result.TypeTable)
       if !node.TypeCompare(t1, t2) {
           yylex.Error("TkEqual type mismatch: " + 
		fmt.Sprintf("%s = %s | ", e1.String(), e2.String()) +
		fmt.Sprintf("%s = %s", node.TypeLabel(t1), node.TypeLabel(t2)))
       }
       Result.Baslib = true // BoolToInt
       $$ = &node.NodeExpEqual{Left:e1, Right:e2}
     }
   | exp TkUnequal exp
     {
       if !node.TypeCompare($1.Type(Result.TypeTable), $3.Type(Result.TypeTable)) {
           yylex.Error("TkUnequal type mismatch")
       }
       Result.Baslib = true // BoolToInt
       $$ = &node.NodeExpUnequal{Left:$1, Right:$3}
     }
   | exp TkGT exp
     {
       if !node.TypeCompare($1.Type(Result.TypeTable), $3.Type(Result.TypeTable)) {
           yylex.Error("TkGT type mismatch")
       }
       Result.Baslib = true // BoolToInt
       $$ = &node.NodeExpGT{Left:$1, Right:$3}
     }
   | exp TkLT exp
     {
       if !node.TypeCompare($1.Type(Result.TypeTable), $3.Type(Result.TypeTable)) {
           yylex.Error("TkLT type mismatch")
       }
       Result.Baslib = true // BoolToInt
       $$ = &node.NodeExpLT{Left:$1, Right:$3}
     }
   | exp TkGE exp
     {
       if !node.TypeCompare($1.Type(Result.TypeTable), $3.Type(Result.TypeTable)) {
           yylex.Error("TkGE type mismatch")
       }
       Result.Baslib = true // BoolToInt
       $$ = &node.NodeExpGE{Left:$1, Right:$3}
     }
   | exp TkLE exp
     {
       if !node.TypeCompare($1.Type(Result.TypeTable), $3.Type(Result.TypeTable)) {
           yylex.Error("TkLE type mismatch")
       }
       Result.Baslib = true // BoolToInt
       $$ = &node.NodeExpLE{Left:$1, Right:$3}
     }
   | TkKeywordFix TkParLeft exp TkParRight
     {
       e := $3
       if !node.TypeNumeric(e.Type(Result.TypeTable)) {
           yylex.Error("FIX expression must be numeric")
       }
       Result.Baslib = true
       $$ = &node.NodeExpFix{Value:e}
     }
   | TkKeywordInt TkParLeft exp TkParRight
     {
       e := $3
       if !node.TypeNumeric(e.Type(Result.TypeTable)) {
           yylex.Error("INT expression must be numeric")
       }
       Result.Baslib = true
       $$ = &node.NodeExpInt{Value:e}
     }
   | TkKeywordLeft TkParLeft exp TkComma exp TkParRight
     {
       e1 := $3
       e2 := $5
       if e1.Type(Result.TypeTable) != node.TypeString {
           yylex.Error("LEFT$ value must be string")
       }
       if !node.TypeNumeric(e2.Type(Result.TypeTable)) {
           yylex.Error("LEFT$ size must be numeric")
       }
       Result.Baslib = true
       $$ = &node.NodeExpLeft{Value:e1, Size:e2}
     }
   | TkKeywordRight TkParLeft exp TkComma exp TkParRight
     {
       e1 := $3
       e2 := $5
       if e1.Type(Result.TypeTable) != node.TypeString {
           yylex.Error("RIGHT$ value must be string")
       }
       if !node.TypeNumeric(e2.Type(Result.TypeTable)) {
           yylex.Error("RIGHT$ size must be numeric")
       }
       Result.Baslib = true
       $$ = &node.NodeExpRight{Value:e1, Size:e2}
     }
   | TkKeywordLen TkParLeft exp TkParRight  { $$ = &node.NodeExpLen{Value:$3} }
   | TkKeywordMid TkParLeft exp TkComma exp TkParRight
     {
       e1 := $3
       e2 := $5
       if e1.Type(Result.TypeTable) != node.TypeString {
           yylex.Error("MID$ value must be string")
       }
       if !node.TypeNumeric(e2.Type(Result.TypeTable)) {
           yylex.Error("MID$ begin must be numeric")
       }
       Result.Baslib = true
       $$ = &node.NodeExpMid{Value:e1, Begin:e2}
     }
   | TkKeywordMid TkParLeft exp TkComma exp TkComma exp TkParRight
     {
       e1 := $3
       e2 := $5
       e3 := $7
       if e1.Type(Result.TypeTable) != node.TypeString {
           yylex.Error("MID$ value must be string")
       }
       if !node.TypeNumeric(e2.Type(Result.TypeTable)) {
           yylex.Error("MID$ begin must be numeric")
       }
       if !node.TypeNumeric(e3.Type(Result.TypeTable)) {
           yylex.Error("MID$ size must be numeric")
       }
       Result.Baslib = true
       $$ = &node.NodeExpMid{Value:e1, Begin:e2, Size:e3}
     }
   | TkKeywordRnd { $$ = &node.NodeExpRnd{Value:&node.NodeExpNumber{Value:"1"}} }
   | TkKeywordRnd TkParLeft exp TkParRight
     {
       e := $3
       if !node.TypeNumeric(e.Type(Result.TypeTable)) {
           yylex.Error("RND expression must be numeric")
       }
       Result.Baslib = true
       $$ = &node.NodeExpRnd{Value:e}
     }
   | TkKeywordStr TkParLeft exp TkParRight
     {
       num := $3
       if !node.TypeNumeric(num.Type(Result.TypeTable)) {
           yylex.Error("STR$ expression must be numeric")
       }
       Result.Baslib = true
       $$ = &node.NodeExpStr{Value:num}
     }
   | TkKeywordVal TkParLeft exp TkParRight
     {
       str := $3
       if str.Type(Result.TypeTable) != node.TypeString {
           yylex.Error("VAL expression must be string")
       }
       Result.Baslib = true
       $$ = &node.NodeExpVal{Value:str}
     }
   | TkKeywordTab TkParLeft exp TkParRight
     {
       num := $3
       if !node.TypeNumeric(num.Type(Result.TypeTable)) {
           yylex.Error("TAB expression must be numeric")
       }
       Result.Baslib = true
       $$ = &node.NodeExpTab{Value:num}
     }
   | TkKeywordSpc TkParLeft exp TkParRight
     {
       num := $3
       if !node.TypeNumeric(num.Type(Result.TypeTable)) {
           yylex.Error("SPC expression must be numeric")
       }
       $$ = &node.NodeExpSpc{Value:num}
       Result.Baslib = true
     }
   | TkKeywordSpace TkParLeft exp TkParRight
     {
       num := $3
       if !node.TypeNumeric(num.Type(Result.TypeTable)) {
           yylex.Error("SPACE$ expression must be numeric")
       }
       Result.Baslib = true
       $$ = &node.NodeExpSpace{Value:num}
     }
   | TkKeywordString TkParLeft exp TkComma exp TkParRight
     {
       num := $3
       char := $5
       if !node.TypeNumeric(num.Type(Result.TypeTable)) {
           yylex.Error("STRING$ expression must be numeric")
       }
       t := char.Type(Result.TypeTable)
       if !node.TypeNumeric(t) && t != node.TypeString  {
           yylex.Error("STRING$ char must be numeric or string")
       }
       Result.Baslib = true
       $$ = &node.NodeExpFuncString{Value:num, Char: char}
     }
   | TkKeywordAsc TkParLeft exp TkParRight
     {
       str := $3
       if str.Type(Result.TypeTable) != node.TypeString {
           yylex.Error("ASC expression must be string")
       }
       Result.Baslib = true
       $$ = &node.NodeExpAsc{Value:str}
     }
   | TkKeywordChr TkParLeft exp TkParRight
     {
       num := $3
       if !node.TypeNumeric(num.Type(Result.TypeTable)) {
           yylex.Error("CHR$ expression must be numeric")
       }
       $$ = &node.NodeExpChr{Value:num}
     }
   | TkKeywordDate {
       Result.Baslib = true
       $$ = &node.NodeExpDate{}
     }
   | TkKeywordTime {
       Result.Baslib = true
       $$ = &node.NodeExpTime{}
     }
   | TkKeywordTimer {
       Result.Baslib = true
       $$ = &node.NodeExpTimer{}
     }
   | TkKeywordAbs TkParLeft exp TkParRight
     {
       num := $3
       if !node.TypeNumeric(num.Type(Result.TypeTable)) {
           yylex.Error("ABS expression must be numeric")
       }
       Result.LibMath = true
       $$ = &node.NodeExpAbs{Value:num}
     }
   | TkKeywordSgn TkParLeft exp TkParRight
     {
       num := $3
       if !node.TypeNumeric(num.Type(Result.TypeTable)) {
           yylex.Error("SGN expression must be numeric")
       }
       Result.Baslib = true
       $$ = &node.NodeExpSgn{Value:num}
     }
   | TkKeywordCos TkParLeft exp TkParRight
     {
       num := $3
       if !node.TypeNumeric(num.Type(Result.TypeTable)) {
           yylex.Error("COS expression must be numeric")
       }
       Result.LibMath = true
       $$ = &node.NodeExpCos{Value:num}
     }
   | TkKeywordSin TkParLeft exp TkParRight
     {
       num := $3
       if !node.TypeNumeric(num.Type(Result.TypeTable)) {
           yylex.Error("SIN expression must be numeric")
       }
       Result.LibMath = true
       $$ = &node.NodeExpSin{Value:num}
     }
   | TkKeywordSqr TkParLeft exp TkParRight
     {
       num := $3
       if !node.TypeNumeric(num.Type(Result.TypeTable)) {
           yylex.Error("SQR expression must be numeric")
       }
       Result.LibMath = true
       $$ = &node.NodeExpSqr{Value:num}
     }
   | TkKeywordTan TkParLeft exp TkParRight
     {
       num := $3
       if !node.TypeNumeric(num.Type(Result.TypeTable)) {
           yylex.Error("TAN expression must be numeric")
       }
       Result.LibMath = true
       $$ = &node.NodeExpTan{Value:num}
     }
   | TkKeywordGofunc TkParLeft one_const_str TkParRight
     {
       $$ = &node.NodeExpGofunc{Name: $3}
     }
   | TkKeywordGofunc TkParLeft one_const_str TkComma expressions_push call_exp_list expressions_pop TkParRight
     {
       $$ = &node.NodeExpGofunc{Name: $3, Arguments: $6}
     }
   | TkKeywordInkey
     {
       Result.Baslib = true
       $$ = &node.NodeExpInkey{}
     }
   | TkKeywordInstr TkParLeft exp TkComma exp TkParRight
     {
       str := $3
       sub := $5
       if str.Type(Result.TypeTable) != node.TypeString {
           yylex.Error("INSTR wrong string type")
       }
       if sub.Type(Result.TypeTable) != node.TypeString {
           yylex.Error("INSTR wrong sub-string type")
       }
       Result.Baslib = true
       $$ = &node.NodeExpInstr{Begin:&node.NodeExpNumber{Value:"1"},Str:str,Sub:sub}
     }
   | TkKeywordInstr TkParLeft exp TkComma exp TkComma exp TkParRight
     {
       begin := $3
       str := $5
       sub := $7
       if !node.TypeNumeric(begin.Type(Result.TypeTable)) {
           yylex.Error("INSTR offset must be numeric")
       }
       if str.Type(Result.TypeTable) != node.TypeString {
           yylex.Error("INSTR wrong string type")
       }
       if sub.Type(Result.TypeTable) != node.TypeString {
           yylex.Error("INSTR wrong sub-string type")
       }
       Result.Baslib = true
       $$ = &node.NodeExpInstr{Begin:begin,Str:str,Sub:sub}
     }
  | TkKeywordPeek TkParLeft exp TkParRight
     {
       $$ = &node.NodeExpPeek{}
     }
  | TkKeywordInputFunc TkParLeft exp TkParRight
     {
       count := $3
       if !node.TypeNumeric(count.Type(Result.TypeTable)) {
           yylex.Error("INPUT$ count must be numeric")
       }	
       Result.Baslib = true
       $$ = &node.NodeExpInput{Count: count}
     }
  | TkKeywordInputFunc TkParLeft exp TkComma file_num TkParRight
     {
       count := $3
       num := $5
       if !node.TypeNumeric(count.Type(Result.TypeTable)) {
           yylex.Error("INPUT$ count must be numeric")
       }	
       if !node.TypeNumeric(num.Type(Result.TypeTable)) {
           yylex.Error("INPUT$ file number must be numeric")
       }	
       Result.Baslib = true
       $$ = &node.NodeExpInputFile{Count: count, Number: num}
     }
  | TkKeywordPos TkParLeft exp TkParRight
     {
        Result.Baslib = true
        $$ = &node.NodeExpPos{}
     }
  | TkKeywordScreen TkParLeft exp TkComma exp TkParRight
     {
	row := $3
	col := $5
        if !node.TypeNumeric(row.Type(Result.TypeTable)) {
           yylex.Error("SCREEN() row must be numeric")
        }	
        if !node.TypeNumeric(col.Type(Result.TypeTable)) {
           yylex.Error("SCREEN() column must be numeric")
        }	
        Result.Baslib = true
        $$ = &node.NodeExpScreen{Row: row, Col: col, ColorFlag: &node.NodeExpNumber{Value:"0"}}
     }
  | TkKeywordScreen TkParLeft exp TkComma exp TkComma exp TkParRight
     {
	row := $3
	col := $5
	colorFlag := $7
        if !node.TypeNumeric(row.Type(Result.TypeTable)) {
           yylex.Error("SCREEN() row must be numeric")
        }	
        if !node.TypeNumeric(col.Type(Result.TypeTable)) {
           yylex.Error("SCREEN() column must be numeric")
        }	
        if !node.TypeNumeric(colorFlag.Type(Result.TypeTable)) {
           yylex.Error("SCREEN() colorflag must be numeric")
        }	
        Result.Baslib = true
        $$ = &node.NodeExpScreen{Row: row, Col: col, ColorFlag: colorFlag}
     }
  | TkKeywordCsrlin 
     {
        Result.Baslib = true
        $$ = &node.NodeExpCsrlin{}
     }
  | TkKeywordEof TkParLeft exp TkParRight
	{
		num := $3
		if !node.TypeNumeric(num.Type(Result.TypeTable)) {
			yylex.Error("EOF file number must be numeric")
		}
		Result.Baslib = true
		$$ = &node.NodeExpEof{Number: num}
	}
  | TkKeywordEnvironFunc TkParLeft exp TkParRight
	{
		key := $3
		Result.Baslib = true
		$$ = &node.NodeExpEnviron{Key: key}
	}
  | TkKeywordLof TkParLeft exp TkParRight
	{
		num := $3
		if !node.TypeNumeric(num.Type(Result.TypeTable)) {
			yylex.Error("LOF file number must be numeric")
		}
		Result.Baslib = true
		$$ = &node.NodeExpLof{Number: num}
	}
  | TkKeywordHex TkParLeft exp TkParRight
	{
		num := $3
		if !node.TypeNumeric(num.Type(Result.TypeTable)) {
			yylex.Error("HEX$ value must be numeric")
		}
		Result.Baslib = true
		$$ = &node.NodeExpHex{Number: num}
	}
  | TkKeywordOct TkParLeft exp TkParRight
	{
		num := $3
		if !node.TypeNumeric(num.Type(Result.TypeTable)) {
			yylex.Error("OCT$ value must be numeric")
		}
		Result.Baslib = true
		$$ = &node.NodeExpOct{Number: num}
	}
   ;

%%

// footer of parser.y -- see lex.go and misc.go


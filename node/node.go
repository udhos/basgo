package node

import (
	"fmt"
	"log"
	//"bufio"
	"strconv"
	"strings"
	"unicode"
)

// LineNumber track used+undefined line numbers
type LineNumber struct {
	Used        bool // GOTO 10, GOSUB 10 etc
	UsedRestore bool // RESTORE 30 -- do NOT create label
	Defined     bool // 10 print
}

// ByLineNumber helper type to sort lines by number
type ByLineNumber []Node

func (a ByLineNumber) Len() int      { return len(a) }
func (a ByLineNumber) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByLineNumber) Less(i, j int) bool {
	n1, n1Numbered := a[i].(*LineNumbered)
	n2, n2Numbered := a[j].(*LineNumbered)
	if n1Numbered && n2Numbered {
		v1, err1 := strconv.Atoi(n1.LineNumber)
		if err1 != nil {
			log.Printf("node sort: bad line number: '%s': %v", n1.LineNumber, err1)
		}
		v2, err2 := strconv.Atoi(n2.LineNumber)
		if err2 != nil {
			log.Printf("node sort: bad line number: '%s': %v", n2.LineNumber, err2)
		}
		return v1 < v2
	}
	if n1Numbered {
		return true
	}
	return false
}

// FuncPrintf is func type for printf
type FuncPrintf func(format string, v ...interface{}) (int, error)

type ArraySymbol struct {
	UsedDimensions     int      // used
	DeclaredDimensions []string // DIM
}

func (a ArraySymbol) ArrayType(table []int, name string) string {
	t := VarType(table, name)
	tt := TypeName(name, t)

	var indices string

	declared := len(a.DeclaredDimensions)
	if declared > 0 {
		for _, d := range a.DeclaredDimensions {
			indices += "[" + d + "+1]"
		}
	} else {
		for i := 0; i < a.UsedDimensions; i++ {
			indices += "[11]"
		}
	}

	arrayType := indices + tt

	return arrayType
}

func TypeName(name string, t int) string {
	var tt string
	switch t {
	case TypeString:
		tt = "string"
	case TypeInteger:
		tt = "int"
	case TypeFloat:
		tt = "float64"
	default:
		log.Printf("node.TypeName: unknown var %s type: %d", name, t)
		tt = "node_TypeName_TYPE_UNKNOWN"
	}
	return tt
}

type VarSymbol struct {
	Name string
	Type int
}

// BuildOptions holds state required for issuing Go code
type BuildOptions struct {
	Headers      map[string]struct{}
	Vars         map[string]VarSymbol
	Arrays       map[string]ArraySymbol
	LineNumbers  map[string]LineNumber // numbers used by GOTO, GOSUB etc
	CountGosub   int
	CountReturn  int
	ReadData     []string // DATA for READ
	RestoreTable map[string]int
	TypeTable    []int
}

// VarSetUsed sets variable as used
func (o *BuildOptions) VarSetUsed(name string, t int) {
	s := RenameVarType(name, t)
	o.Vars[s] = VarSymbol{Name: name, Type: t}
}

// VarIsUsed checks whether variable is used
func (o *BuildOptions) VarIsUsed(name string, t int) bool {
	s := RenameVarType(name, t)
	_, used := o.Vars[s]
	return used
}

type FuncSymbol struct {
	Func *NodeDefFn
	Used bool
}

// fna => true
// fn1 => false
// fn  => false
func IsFuncName(name string) bool {
	if len(name) < 3 {
		return false
	}
	if !strings.HasPrefix(strings.ToLower(name), "fn") {
		return false
	}
	return unicode.IsLetter(rune(name[2]))
}

func FuncSetDeclared(tab map[string]FuncSymbol, f *NodeDefFn) error {
	name := f.FuncName
	_, found := FuncGet(tab, name)
	if found {
		return fmt.Errorf("can't redeclare DEF FN '%s'", name)
	}
	tab[strings.ToLower(name)] = FuncSymbol{Func: f}
	return nil
}

func FuncSetUsed(tab map[string]FuncSymbol, name string, parameters []NodeExp, typeTable []int) error {
	symb, found := FuncGet(tab, name)
	if !found {
		return fmt.Errorf("can't use non-declared DEF FN func '%s'", name)
	}
	if len(symb.Func.Variables) != len(parameters) {
		return fmt.Errorf("wrong arg count for DEF FN '%s': declared=%d used=%d", name, len(symb.Func.Variables), len(parameters))
	}
	for i, p := range parameters {
		s1 := symb.Func.Variables[i].String()
		s2 := p.String()
		t1 := symb.Func.Variables[i].Type(typeTable)
		t2 := p.Type(typeTable)
		if !TypeCompare(t1, t2) {
			return fmt.Errorf("arg type mismatch for DEF FN '%s': arg=%d declaration=%s(%s) call=%s(%s)", name, i, TypeLabel(t1), s1, TypeLabel(t2), s2)
		}
	}
	symb.Used = true
	tab[strings.ToLower(name)] = symb // overwrite
	return nil
}

func FuncGet(tab map[string]FuncSymbol, name string) (FuncSymbol, bool) {
	symb, found := tab[strings.ToLower(name)]
	return symb, found
}

// ArraySetDeclared sets array as declared
func ArraySetDeclared(tab map[string]ArraySymbol, name string, dimensions []string) error {
	low := strings.ToLower(name)

	var used, declared bool
	a, found := tab[low]
	if found {
		used = a.UsedDimensions > 0
		declared = len(a.DeclaredDimensions) > 0
	}
	if used {
		// cannot change used dimensions
		if a.UsedDimensions != len(dimensions) {
			return fmt.Errorf("array '%s' used with new dimensions %d, old ones were %d", name, len(dimensions), a.UsedDimensions)
		}
	}
	if declared {
		// cannot redeclare dimensions
		if len(a.DeclaredDimensions) != len(dimensions) {
			return fmt.Errorf("array '%s' redeclared with new dimensions %d, old ones were %d", name, len(dimensions), len(a.DeclaredDimensions))
		}
		for i, d := range dimensions {
			if d != a.DeclaredDimensions[i] {
				return fmt.Errorf("array '%s' redeclared dimension %d as %s, old one was %s", name, i, d, a.DeclaredDimensions[i])
			}
		}
	}

	a.DeclaredDimensions = dimensions
	tab[low] = a

	return nil
}

// ArraySetUsed sets array as used
func ArraySetUsed(tab map[string]ArraySymbol, name string, dimensions int) error {
	low := strings.ToLower(name)

	var used, declared bool
	a, found := tab[low]
	if found {
		used = a.UsedDimensions > 0
		declared = len(a.DeclaredDimensions) > 0
	}
	if used {
		// cannot change used dimensions
		if a.UsedDimensions != dimensions {
			return fmt.Errorf("array '%s' used with new dimensions %d, old ones were %d", name, dimensions, a.UsedDimensions)
		}
	}
	if declared {
		// cannot change declared dimensions
		d := len(a.DeclaredDimensions)
		if d != dimensions {
			return fmt.Errorf("array '%s' used with new dimensions %d, declared ones were %d", name, dimensions, d)
		}
	}

	a.UsedDimensions = dimensions
	tab[low] = a

	return nil
}

// ArrayIsUsed checks whether array is used
func ArrayIsUsed(tab map[string]ArraySymbol, name string) bool {
	a, found := tab[strings.ToLower(name)]
	if !found {
		return false
	}
	return a.UsedDimensions > 0
}

// VarMatch matches names of two variables
func VarMatch(s1, s2 string) bool {
	return strings.ToLower(s1) == strings.ToLower(s2)
}

// RenameFunc renames a.B$ => func_str_a_b
func RenameFunc(table []int, name string) string {
	return "func_" + RenameVar(table, name)
}

// RenameArray renames a.B$ => array_str_a_b
func RenameArray(table []int, name string) string {
	return "array_" + RenameVar(table, name)
}

// RenameVar renames a.B$ => str_a_b
func RenameVar(table []int, name string) string {
	last := len(name) - 1
	if last < 0 {
		log.Printf("RenameVar: bad var name: [%s]", name)
		return "RenameVar_failure"
	}
	switch name[last] {
	case '$':
		return RenameVarType(name, TypeString)
	case '%':
		return RenameVarType(name, TypeInteger)
	case '!':
		return RenameVarType(name, TypeFloat)
	case '#':
		return RenameVarType(name, TypeDouble)
	}
	t := VarTypeDouble(table, name)
	return RenameVarType(name, t)
}

func RenameVarType(name string, t int) string {
	last := len(name) - 1
	if last < 0 {
		log.Printf("RenameVarType: bad var name: [%s]", name)
		return "RenameVarType_failure"
	}
	switch name[last] {
	case '$', '%', '!', '#':
		name = name[:last]
	}
	name = strings.ToLower(name)
	name = strings.Replace(name, ".", "_", -1)
	switch t {
	case TypeString:
		return "str_" + name
	case TypeInteger:
		return "int_" + name
	case TypeFloat:
		return "sng_" + name
	case TypeDouble:
		return "dbl_" + name
	}
	log.Printf("RenameVarType: bad var type=%d [%s]", t, name)
	return "sng_" + name
}

// VarType finds var type: a$ => string
func VarType(table []int, name string) int {
	t := VarTypeDouble(table, name)
	if t == TypeDouble {
		return TypeFloat
	}
	return t
}

func VarTypeDouble(table []int, name string) int {
	last := len(name) - 1
	if last < 0 {
		log.Printf("VarType: bad var name: len=%d [%s]", len(name), name)
		return TypeFloat
	}
	switch name[last] {
	case '$':
		return TypeString
	case '%':
		return TypeInteger
	case '!':
		return TypeFloat
	case '#':
		return TypeDouble
	}
	if len(table) != 26 {
		log.Printf("VarType: unexpected type table size: %d", len(table))
		return TypeFloat
	}
	i := unicode.ToLower(rune(name[0])) - 'a'
	if i < 0 && i > 25 {
		log.Printf("VarType: bad var first letter index: letter=%c index=%d", name[0], i)
		return TypeFloat
	}
	//log.Printf("VarType: %s = %s", name, TypeLabel(table[i]))
	return table[i]
}

// Node is element for syntax tree
type Node interface {
	Show(printf FuncPrintf)
	Name() string
	Build(options *BuildOptions, outputf FuncPrintf)
	FindUsedVars(options *BuildOptions)
}

// LineNumbered is numbered line
type LineNumbered struct {
	LineNumber string
	Nodes      []Node
	RawLine    string
}

// Show displays the node
func (n *LineNumbered) Show(printf FuncPrintf) {
	printf("line[%s]: ", n.LineNumber)
	for _, n := range n.Nodes {
		n.Show(printf)
	}
	printf("\n")
}

// Name returns the name of the node
func (n *LineNumbered) Name() string {
	return "NUMBERED-LINE:" + n.LineNumber
}

// Build generates code
func (n *LineNumbered) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// line %s\n", n.LineNumber)
	if ln, found := options.LineNumbers[n.LineNumber]; found && ln.Used {
		// generate label for GOTO GOSUB etc
		outputf("line%s:\n", n.LineNumber)
	}
	for _, n := range n.Nodes {
		n.Build(options, outputf)
	}
}

// FindUsedVars finds used vars
func (n *LineNumbered) FindUsedVars(options *BuildOptions) {
	for _, n := range n.Nodes {
		n.FindUsedVars(options)
	}
}

// LineImmediate is unnumbered line
type LineImmediate struct {
	Nodes   []Node
	RawLine string
}

// Show displays the node
func (n *LineImmediate) Show(printf FuncPrintf) {
	printf("immediate: ")
	for _, n := range n.Nodes {
		n.Show(printf)
	}
	printf("\n")
}

// Name returns the name of the node
func (n *LineImmediate) Name() string {
	return "UNNUMBERED-LINE"
}

// Build generates code
func (n *LineImmediate) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// unnumbered line ignored: '%s'\n", strings.TrimSpace(n.RawLine))
}

// FindUsedVars finds used vars
func (n *LineImmediate) FindUsedVars(options *BuildOptions) {
	// do nothing
}

// NodeEmpty is empty
type NodeEmpty struct {
	Value string
}

// Name returns the name of the node
func (n *NodeEmpty) Name() string {
	s := "EMPTY-NODE"
	if n.Value != "" {
		s += ":" + n.Value
	}
	return s
}

// Show displays the node
func (n *NodeEmpty) Show(printf FuncPrintf) {
	printf("[" + n.Name() + "]")
}

// Build generates code
func (n *NodeEmpty) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// empty node ignored: ")
	n.Show(outputf)
	outputf("\n")
}

// FindUsedVars finds used vars
func (n *NodeEmpty) FindUsedVars(options *BuildOptions) {
	// do nothing
}

// NodeAssign is assignment
type NodeAssign struct {
	Left  NodeExp
	Right NodeExp
}

// Name returns the name of the node
func (n *NodeAssign) Name() string {
	return "LET"
}

// Show displays the node
func (n *NodeAssign) Show(printf FuncPrintf) {
	printf("[")
	printf(n.Name())
	printf(" ")
	printf(n.Left.String())
	printf("=")
	printf(n.Right.String())
	printf("]")
}

func assignCode(options *BuildOptions, op string, v1, v2 string, t1, t2 int) string {
	switch {
	case t1 == TypeFloat && t2 == TypeInteger:
		v2 = toFloat(v2)
	case t1 == TypeInteger && t2 == TypeFloat:
		options.Headers["math"] = struct{}{}
		v2 = toInt("math.Round(" + v2 + ")")
	}

	code := fmt.Sprintf("%s %s %s", v1, op, v2)

	return code
}

// Build generates code
func (n *NodeAssign) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	tv := n.Left.Type(options.TypeTable)
	te := n.Right.Type(options.TypeTable)

	v := n.Left.Exp(options)
	e := n.Right.Exp(options)

	code := assignCode(options, "=", v, e, tv, te)

	if options.VarIsUsed(n.Left.String(), tv) {
		outputf(code + "\n")
	} else {
		outputf("// %s // suppressed: '%s' not used\n", code, n.Left.String())
	}
}

// FindUsedVars finds used vars
func (n *NodeAssign) FindUsedVars(options *BuildOptions) {
	n.Right.FindUsedVars(options)
}

// NodeAssignArray is array assignment
type NodeAssignArray struct {
	Left  *NodeExpArray
	Right NodeExp
}

// Name returns the name of the node
func (n *NodeAssignArray) Name() string {
	return "LET"
}

// Show displays the node
func (n *NodeAssignArray) Show(printf FuncPrintf) {
	printf("[")
	printf(n.Name())
	printf(" ")
	printf(n.Left.String())
	printf("=")
	printf(n.Right.String())
	printf("]")
}

// Build generates code
func (n *NodeAssignArray) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	ta := n.Left.Type(options.TypeTable)
	te := n.Right.Type(options.TypeTable)

	a := n.Left.Exp(options)
	e := n.Right.Exp(options)

	code := assignCode(options, "=", a, e, ta, te)

	used := ArrayIsUsed(options.Arrays, n.Left.Name)

	if used {
		outputf(code + "\n")
	} else {
		outputf("// %s // suppressed: array '%s' not used\n", code, n.Left)
	}
}

// FindUsedVars finds used vars
func (n *NodeAssignArray) FindUsedVars(options *BuildOptions) {
	n.Left.FindUsedVars(options)
	n.Right.FindUsedVars(options)
}

// NodeData is data
type NodeData struct {
	Expressions []NodeExp
}

// Name returns the name of the node
func (n *NodeData) Name() string {
	return "DATA"
}

// Show displays the node
func (n *NodeData) Show(printf FuncPrintf) {
	printf("[%s %q]", n.Name(), n.Expressions)
}

// Build generates code
func (n *NodeData) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	for _, e := range n.Expressions {
		s := e.Exp(options)
		//log.Printf("NodeData.Build: '%s' %s", s, TypeLabel(e.Type(options.TypeTable)))
		options.ReadData = append(options.ReadData, s)
	}
}

// FindUsedVars finds used vars
func (n *NodeData) FindUsedVars(options *BuildOptions) {
	// DATA allows only constant expressions - no vars
}

// NodeDefFn is DEF FN
type NodeDefFn struct {
	FuncName  string
	Variables []NodeExp
	Body      NodeExp
}

// Name returns the name of the node
func (n *NodeDefFn) Name() string {
	return "DEF FN"
}

// Show displays the node
func (n *NodeDefFn) Show(printf FuncPrintf) {
	printf("[%s %s %q = ", n.Name(), n.FuncName, n.Variables)
	printf(n.Body.String())
	printf("]")
}

func FuncBuildType(options *BuildOptions, name string, variables []NodeExp) string {

	var varList string
	if len(variables) > 0 {
		v0 := variables[0]
		vtt0 := TypeName(v0.String(), v0.Type(options.TypeTable))
		varList = v0.Exp(options) + " " + vtt0
		for i := 1; i < len(variables); i++ {
			v := variables[i]
			vtt := TypeName(v.String(), v.Type(options.TypeTable))
			varList += "," + v.Exp(options) + " " + vtt
		}
	}

	t := VarType(options.TypeTable, name)
	tt := TypeName(name, t)

	return "func(" + varList + ") " + tt
}

// Build generates code
func (n *NodeDefFn) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	name := RenameFunc(options.TypeTable, n.FuncName)

	t := VarType(options.TypeTable, n.FuncName)

	var body string
	bodyType := n.Body.Type(options.TypeTable)
	switch {
	case t == TypeInteger && bodyType == TypeFloat:
		body = forceInt(options, n.Body)
	case t == TypeFloat && bodyType == TypeInteger:
		body = forceFloat(options, n.Body)
	default:
		body = n.Body.Exp(options)
	}

	outputf(name)
	outputf(" = ")
	outputf(FuncBuildType(options, n.FuncName, n.Variables))
	outputf(" {\n")
	outputf(" return ")
	outputf(body)
	outputf("\n")
	outputf("}\n")
}

// FindUsedVars finds used vars
func (n *NodeDefFn) FindUsedVars(options *BuildOptions) {

	// find used vars

	tmp := &BuildOptions{
		Vars:   map[string]VarSymbol{},
		Arrays: map[string]ArraySymbol{},
	}

	n.Body.FindUsedVars(tmp)

	// remove func args from used vars

	for _, arg := range n.Variables {
		//log.Printf("NodeDefFn.FindUsedVar: arg: %s/%s", arg.String(), arg.Exp(options))
		delete(tmp.Vars, arg.Exp(options))
	}

	// return used vars

	for k, v := range tmp.Vars {
		options.Vars[k] = v
		//log.Printf("NodeDefFn.FindUsedVar: used var: %s/%s/%s", v.Name, k, TypeLabel(v.Type))
	}

	for k, v := range tmp.Arrays {
		options.Arrays[k] = v
	}
}

// NodeDim is dim
type NodeDim struct {
	Arrays []NodeExp
}

// Name returns the name of the node
func (n *NodeDim) Name() string {
	return "DIM"
}

// Show displays the node
func (n *NodeDim) Show(printf FuncPrintf) {
	printf("[%s ", n.Name())
	for _, a := range n.Arrays {
		printf(a.String())
	}
	printf("]")
}

// Build generates code
func (n *NodeDim) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	for _, e := range n.Arrays {
		arrayExp, isArray := e.(*NodeExpArray)
		if !isArray {
			msg := fmt.Sprintf("NodeDim.Build: unexpected non-array: %v %s", e, e.String())
			log.Printf(msg)
			outputf("// %s\n", msg)
			continue
		}
		v := arrayExp.Name
		arr, found := options.Arrays[strings.ToLower(v)]
		if !found {
			msg := fmt.Sprintf("NodeDim.Build: array not found: %s", v)
			log.Printf(msg)
			outputf("// %s\n", msg)
			continue
		}
		name := RenameArray(options.TypeTable, v)
		arrayType := arr.ArrayType(options.TypeTable, v)
		outputf("%s = %s{} // DIM reset array [%s]\n", name, arrayType, v)
	}
}

// FindUsedVars finds used vars
func (n *NodeDim) FindUsedVars(options *BuildOptions) {
	// DIM allows only constant expressions - no vars
}

// NodeErase is erase
type NodeErase struct {
	Arrays []NodeExp
}

// Name returns the name of the node
func (n *NodeErase) Name() string {
	return "ERASE"
}

// Show displays the node
func (n *NodeErase) Show(printf FuncPrintf) {
	printf("[%s ", n.Name())
	for _, a := range n.Arrays {
		printf(a.String())
	}
	printf("]")
}

// Build generates code
func (n *NodeErase) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	for _, e := range n.Arrays {
		v := e.String()
		arr, found := options.Arrays[strings.ToLower(v)]
		if !found {
			msg := fmt.Sprintf("NodeErase.Build: array not found: '%s'", v)
			log.Printf(msg)
			outputf("// %s\n", msg)
			continue
		}
		name := RenameArray(options.TypeTable, v)
		arrayType := arr.ArrayType(options.TypeTable, v)
		outputf("%s = %s{} // ERASE reset array [%s]\n", name, arrayType, v)
	}
}

// FindUsedVars finds used vars
func (n *NodeErase) FindUsedVars(options *BuildOptions) {
}

// NodeOnGosub is ongosub
type NodeOnGosub struct {
	Cond  NodeExp
	Lines []string
	Index int
}

// Name returns the name of the node
func (n *NodeOnGosub) Name() string {
	return "ON-GOSUB"
}

// Show displays the node
func (n *NodeOnGosub) Show(printf FuncPrintf) {
	printf("[" + n.Name() + " ")
	printf(n.Cond.String())
	printf(" %q]", n.Lines)
}

// Build generates code
func (n *NodeOnGosub) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ON %s GOSUB %q\n", n.Cond.String(), n.Lines)

	outputf("switch %s {\n", forceInt(options, n.Cond))
	for i, num := range n.Lines {
		if num == "" {
			continue
		}
		outputf("case %d:\n", i+1)
		buildOneGosub(options, outputf, num, n.Index)
	}
	outputf("}\n")

	if options.CountReturn > 0 {
		outputf("gosub_return_%d:\n", n.Index)
	}
}

// FindUsedVars finds used vars
func (n *NodeOnGosub) FindUsedVars(options *BuildOptions) {
	n.Cond.FindUsedVars(options)
}

// NodeOnGoto is ongoto
type NodeOnGoto struct {
	Cond  NodeExp
	Lines []string
}

// Name returns the name of the node
func (n *NodeOnGoto) Name() string {
	return "ON-GOTO"
}

// Show displays the node
func (n *NodeOnGoto) Show(printf FuncPrintf) {
	printf("[" + n.Name() + " ")
	printf(n.Cond.String())
	printf(" %q]", n.Lines)
}

// Build generates code
func (n *NodeOnGoto) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ON %s GOTO %q\n", n.Cond.String(), n.Lines)

	outputf("switch %s {\n", forceInt(options, n.Cond))
	for i, num := range n.Lines {
		if num == "" {
			continue
		}
		outputf("case %d: goto line%s\n", i+1, num)
	}
	outputf("}\n")
}

// FindUsedVars finds used vars
func (n *NodeOnGoto) FindUsedVars(options *BuildOptions) {
	n.Cond.FindUsedVars(options)
}

// NodePrint is print
type NodePrint struct {
	Newline     bool
	Tab         bool
	Expressions []NodeExp
}

// Name returns the name of the node
func (n *NodePrint) Name() string {
	return "PRINT"
}

// Show displays the node
func (n *NodePrint) Show(printf FuncPrintf) {
	printf("[" + n.Name())
	for _, e := range n.Expressions {
		printf(" <")
		printf(e.String())
		printf(">")
	}
	if n.Tab {
		printf(" TAB")
	}
	if n.Newline {
		printf(" NEWLINE")
	}
	printf("]")
}

// Build generates code
func (n *NodePrint) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	for _, e := range n.Expressions {
		t := e.Type(options.TypeTable)
		switch t {
		case TypeString:
			outputf("baslib.Print(%s)\n", e.Exp(options))
		case TypeInteger:
			outputf("baslib.PrintInt(%s)\n", e.Exp(options))
		case TypeFloat:
			outputf("baslib.PrintFloat(%s)\n", e.Exp(options))
		default:
			log.Printf("NodePrint.Build: unsupported type: %d", t)
		}

		/*
			numeric := TypeNumeric(e.Type(options.TypeTable))
			if numeric {
				outputf("baslib.Print(` `) // PRINT space before number\n")
			}
			outputf("baslib.Print(%s)\n", e.Exp(options))
			if numeric {
				outputf("baslib.Print(` `) // PRINT space after number\n")
			}
		*/
	}

	if n.Tab {
		outputf(`baslib.Print("        ") // PRINT tab due to ending comma\n`)
	}

	if n.Newline {
		outputf("baslib.Println(``) // PRINT newline not suppressed\n")
	}

	//options.Headers["fmt"] = struct{}{} // used package
}

// FindUsedVars finds used vars
func (n *NodePrint) FindUsedVars(options *BuildOptions) {
	for _, e := range n.Expressions {
		e.FindUsedVars(options)
	}
}

// NodeEnd is end
type NodeEnd struct {
	Message string
}

// Name returns the name of the node
func (n *NodeEnd) Name() string {
	return "END"
}

// Show displays the node
func (n *NodeEnd) Show(printf FuncPrintf) {
	printf("[" + n.Name() + "]")
}

// Build generates code
func (n *NodeEnd) Build(options *BuildOptions, outputf FuncPrintf) {
	if n.Message != "" {
		outputf("log.Printf(\"END: %s\")\n", n.Message)
	}
	outputf("baslib.End()\n")
	outputf("os.Exit(0) // %s\n", n.Name())
	options.Headers["os"] = struct{}{} // used package
}

// FindUsedVars finds used vars
func (n *NodeEnd) FindUsedVars(options *BuildOptions) {
	// do nothing
}

// NodeFor is for
type NodeFor struct {
	Index    int // FOR and NEXT are linked thru the same index
	Variable NodeExp
	First    NodeExp
	Last     NodeExp
	Step     NodeExp
}

// Name returns the name of the node
func (n *NodeFor) Name() string {
	return "FOR"
}

// Show displays the node
func (n *NodeFor) Show(printf FuncPrintf) {
	printf("[" + n.Name())
	printf(" " + n.Variable.String())
	printf(" = " + n.First.String())
	printf(" TO " + n.Last.String())
	printf(" STEP " + n.Step.String())
	printf(" Index=%d", n.Index)
	printf("]")
}

// Build generates code
func (n *NodeFor) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	v := n.Variable.Exp(options)
	first := n.First.Exp(options)
	typeV := n.Variable.Type(options.TypeTable)
	typeFirst := n.First.Type(options.TypeTable)
	code := assignCode(options, "=", v, first, typeV, typeFirst)
	outputf("%s // FOR %d initialization\n", code, n.Index)

	last := n.Last.Exp(options)
	typeLast := n.Last.Type(options.TypeTable)
	codeGT := assignCode(options, ">", v, last, typeV, typeLast)
	codeLT := assignCode(options, "<", v, last, typeV, typeLast)

	outputf("for_loop_%d:\n", n.Index)
	outputf("if (%s) >= 0 { // FOR step non-negative?\n", n.Step.Exp(options))
	outputf("  if %s {\n", codeGT)
	outputf("    goto for_exit_%d\n", n.Index)
	outputf("  }\n")
	outputf("} else {\n")
	outputf("  if %s {\n", codeLT)
	outputf("    goto for_exit_%d\n", n.Index)
	outputf("  }\n")
	outputf("}\n")
}

// FindUsedVars finds used vars
func (n *NodeFor) FindUsedVars(options *BuildOptions) {

	switch v := n.Variable.(type) {
	case *NodeExpIdent:
		options.VarSetUsed(v.Value, v.Type(options.TypeTable))
	case *NodeExpArray:
		err := ArraySetUsed(options.Arrays, v.Name, len(v.Indices))
		if err != nil {
			log.Printf("NodeFor.FindUsedVars: ArraySetUsed: %s: %v", v.String(), err)
		}
	default:
		log.Printf("NodeFor.FindUsedVars: unexpected %s node: %v", v.String(), n.Variable)
	}

	n.First.FindUsedVars(options)
	n.Last.FindUsedVars(options)
	n.Step.FindUsedVars(options)
}

// NodeNext is next
type NodeNext struct {
	Variables []NodeExp
	Fors      []*NodeFor
}

// Name returns the name of the node
func (n *NodeNext) Name() string {
	return "NEXT"
}

// Show displays the node
func (n *NodeNext) Show(printf FuncPrintf) {
	printf("[%s vars_size=%d fors_size=%d]", n.Name(), len(n.Variables), len(n.Fors))
}

// Build generates code
func (n *NodeNext) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	for _, f := range n.Fors {

		v := f.Variable.Exp(options)
		step := f.Step.Exp(options)
		typeV := f.Variable.Type(options.TypeTable)
		typeStep := f.Step.Type(options.TypeTable)
		code := assignCode(options, "+=", v, step, typeV, typeStep)
		outputf("%s // FOR %d step\n", code, f.Index)

		outputf("goto for_loop_%d\n", f.Index)
		outputf("for_exit_%d:\n", f.Index)
	}
}

// FindUsedVars finds used vars
func (n *NodeNext) FindUsedVars(options *BuildOptions) {
	for _, i := range n.Variables {
		switch v := i.(type) {
		case *NodeExpIdent:
			options.VarSetUsed(v.Value, v.Type(options.TypeTable))
		case *NodeExpArray:
			err := ArraySetUsed(options.Arrays, v.Name, len(v.Indices))
			if err != nil {
				log.Printf("NodeFor.FindUsedVars: ArraySetUsed: %s: %v", v.String(), err)
			}
		default:
			log.Printf("NodeFor.FindUsedVars: unexpected %s node: %v", v.String(), i)
		}
	}
}

// NodeGosub is gosub
type NodeGosub struct {
	Index int
	Line  string
}

// Name returns the name of the node
func (n *NodeGosub) Name() string {
	return "GOSUB"
}

// Show displays the node
func (n *NodeGosub) Show(printf FuncPrintf) {
	printf("[" + n.Name() + " " + n.Line + "]")
}

func buildOneGosub(options *BuildOptions, outputf FuncPrintf, line string, index int) {
	outputf("gosubStack = append(gosubStack, %d) // gosub push return index\n", index)
	outputf("goto line%s // Index=%d\n", line, index)
}

// Build generates code
func (n *NodeGosub) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	buildOneGosub(options, outputf, n.Line, n.Index)
	if options.CountReturn > 0 {
		outputf("gosub_return_%d:\n", n.Index)
	}
}

// FindUsedVars finds used vars
func (n *NodeGosub) FindUsedVars(options *BuildOptions) {
	// do nothing
}

// NodeGoto is goto
type NodeGoto struct {
	Line string
}

// Name returns the name of the node
func (n *NodeGoto) Name() string {
	return "GOTO"
}

// Show displays the node
func (n *NodeGoto) Show(printf FuncPrintf) {
	printf("[" + n.Name() + " " + n.Line + "]")
}

// Build generates code
func (n *NodeGoto) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("goto line%s // %s %s\n", n.Line, n.Name(), n.Line)
}

// FindUsedVars finds used vars
func (n *NodeGoto) FindUsedVars(options *BuildOptions) {
	// do nothing
}

// NodeIf is IF
type NodeIf struct {
	Index int
	Cond  NodeExp
	Then  []Node
	Else  []Node
}

// Name returns the name of the node
func (n *NodeIf) Name() string {
	return "IF"
}

// Show displays the node
func (n *NodeIf) Show(printf FuncPrintf) {
	printf("[")
	printf(n.Name())
	printf(" <")
	printf(n.Cond.String())
	printf("> THEN ")
	for _, t := range n.Then {
		printf("<")
		t.Show(printf)
		printf(">")
	}
	printf(" ELSE ")
	for _, t := range n.Else {
		printf("<")
		t.Show(printf)
		printf(">")
	}
	printf("]")
}

// Build generates code
func (n *NodeIf) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// %s Index=%d %s THEN ", n.Name(), n.Index, n.Cond.String())
	for _, t := range n.Then {
		t.Show(outputf)
	}
	outputf(" ELSE ")
	for _, t := range n.Else {
		t.Show(outputf)
	}
	outputf("\n")

	var hasElse bool
	for _, t := range n.Else {
		if _, empty := t.(*NodeEmpty); !empty {
			hasElse = true // found non-empty node under ELSE
			break
		}
	}

	outputf("if 0==(%s) {\n", n.Cond.Exp(options))
	if hasElse {
		outputf("  goto if_else_%d\n", n.Index)
	} else {
		outputf("  goto if_exit_%d\n", n.Index)
	}
	outputf("}\n")

	for _, t := range n.Then {
		t.Build(options, outputf)
	}

	if hasElse {
		outputf("goto if_exit_%d\n", n.Index)
		outputf("if_else_%d:\n", n.Index)
		for _, t := range n.Else {
			t.Build(options, outputf)
		}
	}

	outputf("if_exit_%d:\n", n.Index)
}

// FindUsedVars finds used vars
func (n *NodeIf) FindUsedVars(options *BuildOptions) {
	n.Cond.FindUsedVars(options)
	for _, t := range n.Then {
		t.FindUsedVars(options)
	}
	for _, t := range n.Else {
		t.FindUsedVars(options)
	}
}

// NodeInput handles input
type NodeInput struct {
	Variables    []NodeExp
	PromptString NodeExp
	AddQuestion  bool
}

// Name returns the name of the node
func (n *NodeInput) Name() string {
	return "INPUT"
}

// Show displays the node
func (n *NodeInput) Show(printf FuncPrintf) {
	printf("[%s vars=%d prompt=%s question=%v]", n.Name(), len(n.Variables), n.PromptString, n.AddQuestion)
}

// Build generates code
func (n *NodeInput) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	promptStr := "``"
	if prompt, hasPrompt := n.PromptString.(*NodeExpString); hasPrompt {
		promptStr = prompt.Exp(options)
	}
	var questionMark string
	if n.AddQuestion {
		questionMark = "? "
	}

	outputf("{\n")
	outputf("  fields := baslib.Input(%s,`%s`,%d)\n", promptStr, questionMark, len(n.Variables))

	for i, v := range n.Variables {
		str := v.Exp(options)
		var code string
		t := v.Type(options.TypeTable)
		switch t {
		case TypeString:
			code = fmt.Sprintf("%s = fields[%d]", str, i)
		case TypeInteger:
			code = fmt.Sprintf("%s = baslib.InputParseInteger(fields[%d])", str, i)
		case TypeFloat:
			code = fmt.Sprintf("%s = baslib.InputParseFloat(fields[%d])", str, i)
		}
		used := VarOrArrayIsUsed(options, v)
		if used {
			outputf("  %s // INPUT var %d\n", code, i)
		} else {
			outputf("  // %s // unnused INPUT var %d %s suppressed\n", code, i, v.String())
		}
	}

	outputf("}\n")
}

// FindUsedVars finds used vars
func (n *NodeInput) FindUsedVars(options *BuildOptions) {
	// INPUT does not actually use var
}

// NodeList lists lines
type NodeList struct {
}

// Name returns the name of the node
func (n *NodeList) Name() string {
	return "LIST"
}

// Show displays the node
func (n *NodeList) Show(printf FuncPrintf) {
	printf("[" + n.Name() + "]")
}

// Build generates code
func (n *NodeList) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// %s currently not supported by compiler\n", n.Name())
}

// FindUsedVars finds used vars
func (n *NodeList) FindUsedVars(options *BuildOptions) {
	// do nothing
}

// NodeRead is read
type NodeRead struct {
	Variables []NodeExp
}

// Name returns the name of the node
func (n *NodeRead) Name() string {
	return "READ"
}

// Show displays the node
func (n *NodeRead) Show(printf FuncPrintf) {
	printf("[%s", n.Name())
	for _, v := range n.Variables {
		printf(" ")
		printf(v.String())
	}
	printf("]")
}

func VarOrArrayIsUsed(options *BuildOptions, e NodeExp) bool {
	switch ee := e.(type) {
	case *NodeExpIdent:
		return options.VarIsUsed(ee.Value, ee.Type(options.TypeTable))
	case *NodeExpArray:
		return ArrayIsUsed(options.Arrays, ee.Name)
	}
	log.Printf("VarIsUsed: unexpected '%s' non-var non-array: %v", e.String(), e)
	return false
}

// Build generates code
func (n *NodeRead) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	for _, e := range n.Variables {
		v := e.String()      // cosmetic error reporting
		vv := e.Exp(options) // go code
		t := e.Type(options.TypeTable)
		var code string
		switch t {
		case TypeString:
			code = fmt.Sprintf(`%s = baslib.ReadDataString(readData, "%s")`, vv, vv)
		case TypeInteger:
			code = fmt.Sprintf(`%s = baslib.ReadDataInteger(readData, "%s")`, vv, vv)
		case TypeFloat:
			code = fmt.Sprintf(`%s = baslib.ReadDataFloat(readData, "%s")`, vv, vv)
		default:
			msg := fmt.Sprintf("NodeRead.Build: unsupported var %s type: %d", v, t)
			log.Printf(msg)
			code = fmt.Sprintf("println(%s)\n", msg)
		}

		used := VarOrArrayIsUsed(options, e)

		if used {
			outputf(code+" // READ %s\n", v)
		} else {
			outputf("// %s // READ suppressed: '%s' not used\n", code, v)
		}
	}
}

// FindUsedVars finds used vars
func (n *NodeRead) FindUsedVars(options *BuildOptions) {
	// assign value to var does not use it
}

// NodeRestore is restore
type NodeRestore struct {
	Line string
}

// Name returns the name of the node
func (n *NodeRestore) Name() string {
	return "RESTORE"
}

// Show displays the node
func (n *NodeRestore) Show(printf FuncPrintf) {
	printf("[" + n.Name() + " " + n.Line + "]")
}

// Build generates code
func (n *NodeRestore) Build(options *BuildOptions, outputf FuncPrintf) {
	if n.Line == "" {
		outputf("baslib.Restore(readData,\"\",0)\n")
		return
	}

	offset, found := options.RestoreTable[n.Line]
	if !found {
		log.Printf("NodeRestore.Build: line not found: %s", n.Line)
		outputf("baslib.Restore(readData,\"%s\",%d) // RESTORE line NOT FOUND\n", n.Line, offset)
		outputf("panic(`RESTORE line not found: %s`)\n", n.Line)
		return
	}
	outputf("baslib.Restore(readData,\"%s\",%d)\n", n.Line, offset)
}

// FindUsedVars finds used vars
func (n *NodeRestore) FindUsedVars(options *BuildOptions) {
	// do nothing
}

// NodeRem is rem
type NodeRem struct {
	Value string
}

// Name returns the name of the node
func (n *NodeRem) Name() string {
	return "REM"
}

// Show displays the node
func (n *NodeRem) Show(printf FuncPrintf) {
	printf("[" + n.Name() + "]")
}

// Build generates code
func (n *NodeRem) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// REM: '%s'\n", n.Value)
}

// FindUsedVars finds used vars
func (n *NodeRem) FindUsedVars(options *BuildOptions) {
	// do nothing
}

// NodeReturn is return
type NodeReturn struct {
	Line string
}

// Name returns the name of the node
func (n *NodeReturn) Name() string {
	return "RETURN"
}

// Show displays the node
func (n *NodeReturn) Show(printf FuncPrintf) {
	printf("[" + n.Name() + " " + n.Line + "]")
}

func BuildReturnGlobal(options *BuildOptions, outputf FuncPrintf) {
	outputf(`panic("global return label reached without jump")`)
	outputf("\n")

	outputf("return_global:\n")
	outputf("if len(gosubStack) < 1 {\n")
	outputf(`   panic("RETURN without GOSUB")`)
	outputf("\n}\n")

	outputf("{\n") // block open

	// pop
	outputf("last := len(gosubStack) - 1\n")
	outputf("index := gosubStack[last] // save top\n")
	outputf("gosubStack = gosubStack[:last] // pop\n")

	// actual jump
	outputf("switch index {\n") // switch open
	for i := 0; i < options.CountGosub; i++ {
		outputf("case %d: goto gosub_return_%d\n", i, i)
	}
	outputf("}\n") // switch close

	outputf("}\n") // block close

	outputf(`panic("global return switch fallthru")`)
	outputf("\n")
}

// Build generates code
func (n *NodeReturn) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	if n.Line == "" {
		// return
		outputf("goto return_global // %s\n", n.Name())
	} else {
		// return line
		outputf("goto line%s // %s %s\n", n.Line, n.Name(), n.Line)
	}
}

// FindUsedVars finds used vars
func (n *NodeReturn) FindUsedVars(options *BuildOptions) {
	// do nothing
}

// NodeSwap is swap
type NodeSwap struct {
	Var1 NodeExp
	Var2 NodeExp
}

// Name returns the name of the node
func (n *NodeSwap) Name() string {
	return "SWAP"
}

// Show displays the node
func (n *NodeSwap) Show(printf FuncPrintf) {
	printf("[%s %s,%s]", n.Name(), n.Var1.String(), n.Var2.String)
}

// Build generates code
func (n *NodeSwap) Build(options *BuildOptions, outputf FuncPrintf) {
	v1 := n.Var1.Exp(options)
	v2 := n.Var2.Exp(options)
	outputf("%s,%s = %s,%s // SWAP\n", v1, v2, v2, v1)
}

// FindUsedVars finds used vars
func (n *NodeSwap) FindUsedVars(options *BuildOptions) {
	n.Var1.FindUsedVars(options)
	n.Var2.FindUsedVars(options)
}

// NodeWend is wend
type NodeWend struct {
	While *NodeWhile
}

// Name returns the name of the node
func (n *NodeWend) Name() string {
	return "WEND"
}

// Show displays the node
func (n *NodeWend) Show(printf FuncPrintf) {
	printf("[%s index=%d]", n.Name(), n.While.Index)
}

// Build generates code
func (n *NodeWend) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	outputf("goto while_loop_%d\n", n.While.Index)
	outputf("wend_exit_%d:\n", n.While.Index)
}

// FindUsedVars finds used vars
func (n *NodeWend) FindUsedVars(options *BuildOptions) {
	// do nothing
}

// NodeWhile is while
type NodeWhile struct {
	Index int // WHILE and WEND are linked thru the same index
	Cond  NodeExp
}

// Name returns the name of the node
func (n *NodeWhile) Name() string {
	return "WHILE"
}

// Show displays the node
func (n *NodeWhile) Show(printf FuncPrintf) {
	printf("[" + n.Name())
	printf(" " + n.Cond.String())
	printf(" Index=%d", n.Index)
	printf("]")
}

// Build generates code
func (n *NodeWhile) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	outputf("while_loop_%d:\n", n.Index)
	outputf("if 0==(%s) {\n", n.Cond.Exp(options))
	outputf("  goto wend_exit_%d\n", n.Index)
	outputf("}\n")
}

// FindUsedVars finds used vars
func (n *NodeWhile) FindUsedVars(options *BuildOptions) {
	n.Cond.FindUsedVars(options)
}

// NodeGoproc is _GOPROC
type NodeGoproc struct {
	ProcName  *NodeExpString
	Arguments []NodeExp
}

// Name returns the name of the node
func (n *NodeGoproc) Name() string {
	return "_GOPROC"
}

// Show displays the node
func (n *NodeGoproc) Show(printf FuncPrintf) {
	printf("[" + n.Name())
	printf(" " + n.ProcName.String())
	printf("]")
}

// Build generates code
func (n *NodeGoproc) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	outputf(n.ProcName.String())
	outputf("(")
	argc := len(n.Arguments)
	if argc > 0 {
		outputf(n.Arguments[0].Exp(options))
		for i := 1; i < argc; i++ {
			outputf(",")
			outputf(n.Arguments[i].Exp(options))
		}
	}
	outputf(") /* <-- _GOPROC */\n")
}

// FindUsedVars finds used vars
func (n *NodeGoproc) FindUsedVars(options *BuildOptions) {
	for _, a := range n.Arguments {
		a.FindUsedVars(options)
	}
}

// NodeGodecl is _GODECL
type NodeGodecl struct {
	Value *NodeExpString
}

// Name returns the name of the node
func (n *NodeGodecl) Name() string {
	return "_GODECL"
}

// Show displays the node
func (n *NodeGodecl) Show(printf FuncPrintf) {
	printf("[" + n.Name())
	printf(" " + n.Value.String())
	printf("]")
}

// Build generates code
func (n *NodeGodecl) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")
}

// FindUsedVars finds used vars
func (n *NodeGodecl) FindUsedVars(options *BuildOptions) {
}

// NodeGoimport is _GOIMPORT
type NodeGoimport struct {
	Value *NodeExpString
}

// Name returns the name of the node
func (n *NodeGoimport) Name() string {
	return "_GOIMPORT"
}

// Show displays the node
func (n *NodeGoimport) Show(printf FuncPrintf) {
	printf("[" + n.Name())
	printf(" " + n.Value.String())
	printf("]")
}

// Build generates code
func (n *NodeGoimport) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")
}

// FindUsedVars finds used vars
func (n *NodeGoimport) FindUsedVars(options *BuildOptions) {
}

// NodeRandomize is randomize
type NodeRandomize struct {
	Seed NodeExp
}

// Name returns the name of the node
func (n *NodeRandomize) Name() string {
	return "RANDOMIZE"
}

// Show displays the node
func (n *NodeRandomize) Show(printf FuncPrintf) {
	printf("[" + n.Name())
	if n.Seed != NodeExp(nil) {
		printf(" " + n.Seed.String())
	}
	printf("]")
}

// Build generates code
func (n *NodeRandomize) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	if n.Seed == NodeExp(nil) {
		outputf("baslib.RandomizeAuto()\n")
		return
	}

	outputf("baslib.Randomize(" + forceFloat(options, n.Seed) + ")\n")
}

// FindUsedVars finds used vars
func (n *NodeRandomize) FindUsedVars(options *BuildOptions) {
	if n.Seed != NodeExp(nil) {
		n.Seed.FindUsedVars(options)
	}
}

const (
	OpenRandom = iota
	OpenInput  = iota
	OpenOutput = iota
	OpenAppend = iota
)

// NodeOpen is open
type NodeOpen struct {
	File   NodeExp
	Number NodeExp
	Mode   int
}

// Name returns the name of the node
func (n *NodeOpen) Name() string {
	return "OPEN"
}

// Show displays the node
func (n *NodeOpen) Show(printf FuncPrintf) {
	printf("[" + n.Name() + " ")
	printf(n.File.String())
	printf("]")
}

// Build generates code
func (n *NodeOpen) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

}

// FindUsedVars finds used vars
func (n *NodeOpen) FindUsedVars(options *BuildOptions) {
	n.File.FindUsedVars(options)
	n.Number.FindUsedVars(options)
}

// NodeScreen is screen
type NodeScreen struct {
	Mode int
}

// Name returns the name of the node
func (n *NodeScreen) Name() string {
	return "SCREEN"
}

// Show displays the node
func (n *NodeScreen) Show(printf FuncPrintf) {
	printf("[" + n.Name() + " " + strconv.Itoa(n.Mode) + "]")
}

// Build generates code
func (n *NodeScreen) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	outputf("baslib.Screen(%d)\n", n.Mode)
}

// FindUsedVars finds used vars
func (n *NodeScreen) FindUsedVars(options *BuildOptions) {
}

// NodeCls is cls
type NodeCls struct{}

// Name returns the name of the node
func (n *NodeCls) Name() string {
	return "CLS"
}

// Show displays the node
func (n *NodeCls) Show(printf FuncPrintf) {
	printf("[" + n.Name() + "]")
}

// Build generates code
func (n *NodeCls) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	outputf("baslib.Cls()\n")
}

// FindUsedVars finds used vars
func (n *NodeCls) FindUsedVars(options *BuildOptions) {
}

// NodeLocate is locate
type NodeLocate struct {
	Row    NodeExp
	Col    NodeExp
	Cursor NodeExp
}

// Name returns the name of the node
func (n *NodeLocate) Name() string {
	return "LOCATE"
}

// Show displays the node
func (n *NodeLocate) Show(printf FuncPrintf) {
	printf("[")
	printf(n.Name())
	if n.Row != NodeExp(nil) {
		printf(" ")
		printf(n.Row.String())
	}
	if n.Col != NodeExp(nil) {
		printf(" ")
		printf(n.Col.String())
	}
	if n.Cursor != NodeExp(nil) {
		printf(" ")
		printf(n.Cursor.String())
	}
	printf("]")
}

// Build generates code
func (n *NodeLocate) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	row := "-1"
	col := "-1"

	if n.Row != NodeExp(nil) {
		row = forceInt(options, n.Row)
	}
	if n.Col != NodeExp(nil) {
		col = forceInt(options, n.Col)
	}
	if n.Cursor != NodeExp(nil) {
		cursor := "(0 != " + forceInt(options, n.Cursor) + ")"
		outputf("baslib.LocateCursor(%s,%s,%s)\n", row, col, cursor)
		return
	}

	outputf("baslib.Locate(%s,%s)\n", row, col)
}

// FindUsedVars finds used vars
func (n *NodeLocate) FindUsedVars(options *BuildOptions) {
	if n.Row != NodeExp(nil) {
		n.Row.FindUsedVars(options)
	}
	if n.Col != NodeExp(nil) {
		n.Col.FindUsedVars(options)
	}
	if n.Cursor != NodeExp(nil) {
		n.Cursor.FindUsedVars(options)
	}
}

// NodeWidth is width
type NodeWidth struct {
	Width NodeExp
}

// Name returns the name of the node
func (n *NodeWidth) Name() string {
	return "WIDTH"
}

// Show displays the node
func (n *NodeWidth) Show(printf FuncPrintf) {
	printf("[")
	printf(n.Name())
	printf(" ")
	printf(n.Width.String())
	printf("]")
}

// Build generates code
func (n *NodeWidth) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	outputf("baslib.Width(")
	outputf(forceInt(options, n.Width))
	outputf(")\n")
}

// FindUsedVars finds used vars
func (n *NodeWidth) FindUsedVars(options *BuildOptions) {
	n.Width.FindUsedVars(options)
}

// NodeViewPrint is view print
type NodeViewPrint struct {
	Top    NodeExp
	Bottom NodeExp
}

// Name returns the name of the node
func (n *NodeViewPrint) Name() string {
	return "VIEW PRINT"
}

// Show displays the node
func (n *NodeViewPrint) Show(printf FuncPrintf) {
	printf("[")
	printf(n.Name())
	if n.Top != NodeExp(nil) {
		printf(" ")
		printf(n.Top.String())
	}
	if n.Bottom != NodeExp(nil) {
		printf(" ")
		printf(n.Bottom.String())
	}
	printf("]")
}

// Build generates code
func (n *NodeViewPrint) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	if n.Top == NodeExp(nil) || n.Bottom == NodeExp(nil) {
		outputf("baslib.ViewPrintReset()\n")
		return
	}

	top := forceInt(options, n.Top)
	bottom := forceInt(options, n.Bottom)

	outputf("baslib.ViewPrint(%s,%s)\n", top, bottom)
}

// FindUsedVars finds used vars
func (n *NodeViewPrint) FindUsedVars(options *BuildOptions) {
	if n.Top != NodeExp(nil) {
		n.Top.FindUsedVars(options)
	}
	if n.Bottom != NodeExp(nil) {
		n.Bottom.FindUsedVars(options)
	}
}

// NodeColor is color
type NodeColor struct {
	Foreground NodeExp
	Background NodeExp
}

// Name returns the name of the node
func (n *NodeColor) Name() string {
	return "COLOR"
}

// Show displays the node
func (n *NodeColor) Show(printf FuncPrintf) {
	printf("[")
	printf(n.Name())
	if n.Foreground != NodeExp(nil) {
		printf(" ")
		printf(n.Foreground.String())
	}
	if n.Background != NodeExp(nil) {
		printf(" ")
		printf(n.Background.String())
	}
	printf("]")
}

// Build generates code
func (n *NodeColor) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	fg := "-1" // do not change color
	bg := "-1" // do not change color

	if n.Foreground != NodeExp(nil) {
		fg = forceInt(options, n.Foreground)
	}
	if n.Background != NodeExp(nil) {
		bg = forceInt(options, n.Background)
	}

	outputf("baslib.Color(%s,%s)\n", fg, bg)
}

// FindUsedVars finds used vars
func (n *NodeColor) FindUsedVars(options *BuildOptions) {
	if n.Foreground != NodeExp(nil) {
		n.Foreground.FindUsedVars(options)
	}
	if n.Background != NodeExp(nil) {
		n.Background.FindUsedVars(options)
	}
}

// NodeFiles is files
type NodeFiles struct {
	Pattern NodeExp
}

// Name returns the name of the node
func (n *NodeFiles) Name() string {
	return "FILES"
}

// Show displays the node
func (n *NodeFiles) Show(printf FuncPrintf) {
	printf("[")
	printf(n.Name())
	printf(" <")
	printf(n.Pattern.String())
	printf(">]")
}

// Build generates code
func (n *NodeFiles) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	outputf("baslib.Files(%s)\n", n.Pattern.Exp(options))
}

// FindUsedVars finds used vars
func (n *NodeFiles) FindUsedVars(options *BuildOptions) {
	n.Pattern.FindUsedVars(options)
}

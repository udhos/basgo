package node

import (
	"fmt"
	"log"
	//"bufio"
	"strconv"
	"strings"
)

type LineNumber struct {
	Used    bool // GOTO 10, GOSUB 10 etc
	Defined bool // 10 print
}

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

type BuildOptions struct {
	Headers     map[string]struct{}
	Vars        map[string]struct{}
	LineNumbers map[string]LineNumber // numbers used by GOTO, GOSUB etc
}

// RenameVar renames a.b$ => str_a_b
func RenameVar(name string) string {
	name = strings.Replace(name, ".", "_", -1)
	last := len(name) - 1
	if last < 0 {
		return "sng_" + name
	}
	switch name[last] {
	case '$':
		return "str_" + name[:last]
	case '%':
		return "int_" + name[:last]
	case '!':
		return "sng_" + name[:last]
	case '#':
		return "dbl_" + name[:last]
	}
	return "sng_" + name
}

// VarType finds var type: a$ => string
func VarType(name string) int {
	last := len(name) - 1
	if last < 1 {
		return TypeFloat
	}
	switch name[last] {
	case '$':
		return TypeString
	case '%':
		return TypeInteger
	}
	return TypeFloat
}

// Node is element for syntax tree
type Node interface {
	Show(printf FuncPrintf)
	Name() string
	Build(options *BuildOptions, outputf FuncPrintf)
	FindUsedVars(vars map[string]struct{})
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
func (n *LineNumbered) FindUsedVars(vars map[string]struct{}) {
	for _, n := range n.Nodes {
		n.FindUsedVars(vars)
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
func (n *LineImmediate) FindUsedVars(vars map[string]struct{}) {
	// do nothing
}

// NodeEmpty is empty
type NodeEmpty struct {
}

// Name returns the name of the node
func (n *NodeEmpty) Name() string {
	return "EMPTY-NODE"
}

// Show displays the node
func (n *NodeEmpty) Show(printf FuncPrintf) {
	printf("[" + n.Name() + "]")
}

// Build generates code
func (n *NodeEmpty) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// empty node ignored\n")
}

// FindUsedVars finds used vars
func (n *NodeEmpty) FindUsedVars(vars map[string]struct{}) {
	// do nothing
}

// NodeAssign is assignment
type NodeAssign struct {
	Left  string
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
	printf(n.Left)
	printf("=")
	printf(n.Right.String())
	printf("]")
}

// Build generates code
func (n *NodeAssign) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	if _, found := options.Vars[n.Left]; found {
		outputf("%s = %s\n", RenameVar(n.Left), n.Right.Exp(options))
	} else {
		outputf("// %s = %s // suppressed \n", RenameVar(n.Left), n.Right.Exp(options))
	}
}

// FindUsedVars finds used vars
func (n *NodeAssign) FindUsedVars(vars map[string]struct{}) {
	n.Right.FindUsedVars(vars)
}

// NodePrint is print
type NodePrint struct {
	Newline     bool
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
		outputf("fmt.Print(%s)\n", e.Exp(options))
	}

	if n.Newline {
		outputf("fmt.Println()\n")
	}

	options.Headers["fmt"] = struct{}{} // used package
}

// FindUsedVars finds used vars
func (n *NodePrint) FindUsedVars(vars map[string]struct{}) {
	for _, e := range n.Expressions {
		e.FindUsedVars(vars)
	}
}

// NodeEnd is end
type NodeEnd struct {
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
	outputf("os.Exit(0) // %s\n", n.Name())
	options.Headers["os"] = struct{}{} // used package
}

// FindUsedVars finds used vars
func (n *NodeEnd) FindUsedVars(vars map[string]struct{}) {
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
func (n *NodeGoto) FindUsedVars(vars map[string]struct{}) {
	// do nothing
}

// NodeIf is IF
type NodeIf struct {
	Cond NodeExp
	Then []Node
	Else []Node
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
	outputf("// %s %s THEN ", n.Name(), n.Cond.String())
	for _, t := range n.Then {
		t.Show(outputf)
	}
	outputf(" ELSE ")
	for _, t := range n.Else {
		t.Show(outputf)
	}
	outputf("\n")

	outputf("if 0!=(%s) {\n", n.Cond.Exp(options))

	for _, t := range n.Then {
		t.Build(options, outputf)
	}

	var hasElse bool

	for _, t := range n.Else {
		if _, empty := t.(*NodeEmpty); !empty {
			hasElse = true // found non-empty node under ELSE
			break
		}
	}

	if hasElse {
		outputf("} else {\n")

		for _, t := range n.Else {
			t.Build(options, outputf)
		}
	}

	outputf("}\n")
}

// FindUsedVars finds used vars
func (n *NodeIf) FindUsedVars(vars map[string]struct{}) {
	n.Cond.FindUsedVars(vars)
	for _, t := range n.Then {
		t.FindUsedVars(vars)
	}
	for _, t := range n.Else {
		t.FindUsedVars(vars)
	}
}

// NodeInput handles input
type NodeInput struct {
	Variable string
}

// Name returns the name of the node
func (n *NodeInput) Name() string {
	return "INPUT"
}

// Show displays the node
func (n *NodeInput) Show(printf FuncPrintf) {
	printf("[%s %s]", n.Name(), n.Variable)
}

const InputString = "inputString()"
const InputInteger = "inputInteger()"
const InputFloat = "inputFloat()"

// Build generates code
func (n *NodeInput) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	var code string

	t := VarType(n.Variable) // a$ => string
	switch t {
	case TypeString:
		code = InputString
	case TypeInteger:
		code = InputInteger
	case TypeFloat:
		code = InputFloat
	default:
		msg := fmt.Sprintf("NodeInput.Build: unknown var '%s' type: %d", n.Variable, t)
		log.Printf(msg)
		outputf("panic(%s) // INPUT bad var type\n", msg)
		return
	}

	v := RenameVar(n.Variable) // a$ => str_a

	if _, varUsed := options.Vars[n.Variable]; varUsed {
		outputf("%s = %s\n", v, code)
		return
	}

	outputf("%s // unnused INPUT variable %s/%s suppressed\n", code, n.Variable, v)
}

// FindUsedVars finds used vars
func (n *NodeInput) FindUsedVars(vars map[string]struct{}) {
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
func (n *NodeList) FindUsedVars(vars map[string]struct{}) {
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
	outputf("// REM: '%s'", n.Value)
}

// FindUsedVars finds used vars
func (n *NodeRem) FindUsedVars(vars map[string]struct{}) {
	// do nothing
}

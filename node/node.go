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
	Rnd         bool                  // using lib RND
	Input       bool                  // using lib INPUT
	Left        bool                  // using lib LEFT
}

func (o *BuildOptions) VarSetUsed(name string) {
	o.Vars[strings.ToLower(name)] = struct{}{}
}

func (o *BuildOptions) VarIsUsed(name string) bool {
	_, used := o.Vars[strings.ToLower(name)]
	return used
}

func VarMatch(s1, s2 string) bool {
	return strings.ToLower(s1) == strings.ToLower(s2)
}

// RenameVar renames a.B$ => str_a_b
func RenameVar(name string) string {
	name = strings.ToLower(name)
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
func (n *NodeEmpty) FindUsedVars(options *BuildOptions) {
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

	ti := VarType(n.Left)
	te := n.Right.Type()

	v := RenameVar(n.Left)
	e := n.Right.Exp(options)

	code := assignCode(options, "=", v, e, ti, te)

	if options.VarIsUsed(n.Left) {
		outputf(code + "\n")
	} else {
		outputf("// %s // suppressed: '%s' not used\n", code, n.Left)
	}
}

// FindUsedVars finds used vars
func (n *NodeAssign) FindUsedVars(options *BuildOptions) {
	n.Right.FindUsedVars(options)
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
		outputf("fmt.Print(%s)\n", e.Exp(options))
	}

	if n.Tab {
		outputf(`fmt.Print("        ") // PRINT tab due to ending comma\n`)
	}

	if n.Newline {
		outputf("fmt.Println() // PRINT newline not suppressed\n")
	}

	options.Headers["fmt"] = struct{}{} // used package
}

// FindUsedVars finds used vars
func (n *NodePrint) FindUsedVars(options *BuildOptions) {
	for _, e := range n.Expressions {
		e.FindUsedVars(options)
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
func (n *NodeEnd) FindUsedVars(options *BuildOptions) {
	// do nothing
}

// NodeFor is for
type NodeFor struct {
	Index    int // FOR and NEXT are linked thru the same index
	Variable string
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
	printf(" " + n.Variable)
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

	v := RenameVar(n.Variable)
	first := n.First.Exp(options)
	typeV := VarType(n.Variable)
	typeFirst := n.First.Type()
	code := assignCode(options, "=", v, first, typeV, typeFirst)
	outputf("%s // FOR %d initialization\n", code, n.Index)

	last := n.Last.Exp(options)
	typeLast := n.Last.Type()
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
	options.VarSetUsed(n.Variable)
	n.First.FindUsedVars(options)
	n.Last.FindUsedVars(options)
	n.Step.FindUsedVars(options)
}

// NodeNext is next
type NodeNext struct {
	Variables []string
	Fors      []*NodeFor
}

// Name returns the name of the node
func (n *NodeNext) Name() string {
	return "NEXT"
}

// Show displays the node
func (n *NodeNext) Show(printf FuncPrintf) {
	printf("[%s vars=%q fors_size=%d]", n.Name(), n.Variables, len(n.Fors))
}

// Build generates code
func (n *NodeNext) Build(options *BuildOptions, outputf FuncPrintf) {
	outputf("// ")
	n.Show(outputf)
	outputf("\n")

	for _, f := range n.Fors {

		v := RenameVar(f.Variable)
		step := f.Step.Exp(options)
		typeV := VarType(f.Variable)
		typeStep := f.Step.Type()
		code := assignCode(options, "+=", v, step, typeV, typeStep)
		outputf("%s // FOR %d step\n", code, f.Index)

		outputf("goto for_loop_%d\n", f.Index)
		outputf("for_exit_%d:\n", f.Index)
	}
}

// FindUsedVars finds used vars
func (n *NodeNext) FindUsedVars(options *BuildOptions) {
	for _, v := range n.Variables {
		options.VarSetUsed(v)
	}
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

	if options.VarIsUsed(n.Variable) {
		outputf("%s = %s\n", v, code)
		return
	}

	outputf("%s // unnused INPUT variable %s/%s suppressed\n", code, n.Variable, v)
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

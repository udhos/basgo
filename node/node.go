package node

import (
	"log"
	//"fmt"
	//"bufio"
	"strconv"
	"strings"
)

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
	Headers map[string]struct{}
}

// Node is element for syntax tree
type Node interface {
	Show(printf FuncPrintf)
	Name() string
	Build(options *BuildOptions, outputf FuncPrintf)
}

// LineNumbered is empty
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
	for _, n := range n.Nodes {
		n.Build(options, outputf)
	}
}

// LineImmediate is empty
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

// NodePrint is print
type NodePrint struct {
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
	outputf("fmt.Println()\n")

	options.Headers["fmt"] = struct{}{} // used package
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

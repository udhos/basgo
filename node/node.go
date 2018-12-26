package node

import (
//"log"
//"fmt"
//"bufio"
)

type FuncPrintf func(format string, v ...interface{}) (int, error)

// Node is element for syntax tree
type Node interface {
	Show(printf FuncPrintf)
	Name() string
	Build(outputf FuncPrintf)
}

// LineNumbered is empty
type LineNumbered struct {
	LineNumber string
	Nodes      []Node
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
func (n *LineNumbered) Build(outputf FuncPrintf) {
	outputf("// line %s\n", n.LineNumber)
	for _, n := range n.Nodes {
		n.Build(outputf)
	}
}

// LineImmediate is empty
type LineImmediate struct {
	Nodes []Node
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
func (n *LineImmediate) Build(outputf FuncPrintf) {
	outputf("// unnumbered line ignored\n")
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
func (n *NodeEmpty) Build(outputf FuncPrintf) {
	outputf("// empty node ignored\n")
}

// NodePrint is print
type NodePrint struct {
}

// Name returns the name of the node
func (n *NodePrint) Name() string {
	return "PRINT"
}

// Show displays the node
func (n *NodePrint) Show(printf FuncPrintf) {
	printf("[" + n.Name() + "]")
}

// Build generates code
func (n *NodePrint) Build(outputf FuncPrintf) {
	outputf("fmt.Println() // %s\n", n.Name())
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
func (n *NodeEnd) Build(outputf FuncPrintf) {
	outputf("os.Exit(0) // %s\n", n.Name())
}

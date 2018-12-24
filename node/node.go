package node

import (
//"log"
//"fmt"
//"bufio"

//"github.com/udhos/basgo/basgo"
)

type funcPrintf func(format string, v ...interface{}) (int, error)

// Node is element for syntax tree
type Node interface {
	Show(printf funcPrintf)
	//Run(b *Basgo, printf funcPrintf, r *bufio.Reader)
}

// LineNumbered is empty
type LineNumbered struct {
	LineNumber string
	Nodes      []Node
}

// Show displays the node
func (n *LineNumbered) Show(printf funcPrintf) {
	printf("line[%s]: ", n.LineNumber)
	for _, n := range n.Nodes {
		n.Show(printf)
	}
	printf("\n")
}

// LineImmediate is empty
type LineImmediate struct {
	Nodes []Node
}

// Show displays the node
func (n *LineImmediate) Show(printf funcPrintf) {
	printf("immediate: ")
	for _, n := range n.Nodes {
		n.Show(printf)
	}
	printf("\n")
}

// NodeEmpty is empty
type NodeEmpty struct {
}

// Show displays the node
func (n *NodeEmpty) Show(printf funcPrintf) {
	printf("[empty]")
}

// NodePrint is print
type NodePrint struct {
}

// Show displays the node
func (n *NodePrint) Show(printf funcPrintf) {
	printf("[PRINT]")
}

// NodeEnd is end
type NodeEnd struct {
}

// Show displays the node
func (n *NodeEnd) Show(printf funcPrintf) {
	printf("[END]")
}

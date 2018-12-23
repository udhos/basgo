package basparser

import (
	//"log"
	"fmt"
)

// Node is element for syntax tree
type Node interface {
	Run()
}

// LineNumbered is empty
type LineNumbered struct {
	LineNumber int
	Nodes      []Node
}

// Run executes the node
func (n *LineNumbered) Run() {
	fmt.Printf("LineNumbered.Run: %d\n", n.LineNumber)
	for _, n := range n.Nodes {
		n.Run()
	}
}

// LineImmediate is empty
type LineImmediate struct {
	Nodes []Node
}

// Run executes the node
func (n *LineImmediate) Run() {
	fmt.Printf("LineImmediate.Run\n")
	for _, n := range n.Nodes {
		n.Run()
	}
}

// NodeEmpty is empty
type NodeEmpty struct {
}

// Run executes the node
func (n *NodeEmpty) Run() {
	fmt.Printf("NodeEmpty.Run\n")
}

// NodePrint is print
type NodePrint struct {
}

// Run executes the node
func (n *NodePrint) Run() {
	fmt.Printf("NodePrint.Run\n")
}

// NodeEnd is end
type NodeEnd struct {
}

// Run executes the node
func (n *NodeEnd) Run() {
	fmt.Printf("NodeEnd.Run\n")
}

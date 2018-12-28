package node

import (
//"log"
//"fmt"
//"bufio"
)

// NodeExp is interface for expressions
type NodeExp interface {
	Exp() string
}

// NodeExpNumber holds value
type NodeExpNumber struct{ Value string }

// Exp returns value
func (e *NodeExpNumber) Exp() string {
	return e.Value
}

// NodeExpFloat holds value
type NodeExpFloat struct{ Value string }

// Exp returns value
func (e *NodeExpFloat) Exp() string {
	return e.Value
}

// NodeExpString holds value
type NodeExpString struct{ Value string }

// Exp returns value
func (e *NodeExpString) Exp() string {
	return e.Value
}

// NodeExpIdentifier holds value
type NodeExpIdentifier struct{ Value string }

// Exp returns value
func (e *NodeExpIdentifier) Exp() string {
	return e.Value
}

// NodeExpPlus holds value
type NodeExpPlus struct {
	Left  NodeExp
	Right NodeExp
}

// Exp returns value
func (e *NodeExpPlus) Exp() string {
	return e.Left.Exp() + "+" + e.Right.Exp()
}

// NodeExpMinus holds value
type NodeExpMinus struct {
	Left  NodeExp
	Right NodeExp
}

// Exp returns value
func (e *NodeExpMinus) Exp() string {
	return e.Left.Exp() + "-" + e.Right.Exp()
}

// NodeExpMult holds value
type NodeExpMult struct {
	Left  NodeExp
	Right NodeExp
}

// Exp returns value
func (e *NodeExpMult) Exp() string {
	return e.Left.Exp() + "*" + e.Right.Exp()
}

// NodeExpDiv holds value
type NodeExpDiv struct {
	Left  NodeExp
	Right NodeExp
}

// Exp returns value
func (e *NodeExpDiv) Exp() string {
	return e.Left.Exp() + "/" + e.Right.Exp()
}

// NodeExpDivInt holds value
type NodeExpDivInt struct {
	Left  NodeExp
	Right NodeExp
}

// Exp returns value
func (e *NodeExpDivInt) Exp() string {
	return e.Left.Exp() + "\\" + e.Right.Exp()
}

// NodeExpUnaryPlus holds value
type NodeExpUnaryPlus struct{ Value NodeExp }

// Exp returns value
func (e *NodeExpUnaryPlus) Exp() string {
	return "+" + e.Value.Exp()
}

// NodeExpUnaryMinus holds value
type NodeExpUnaryMinus struct{ Value NodeExp }

// Exp returns value
func (e *NodeExpUnaryMinus) Exp() string {
	return "-" + e.Value.Exp()
}

// NodeExpGroup holds value
type NodeExpGroup struct{ Value NodeExp }

// Exp returns value
func (e *NodeExpGroup) Exp() string {
	return "(" + e.Value.Exp() + ")"
}

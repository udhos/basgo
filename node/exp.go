package node

import (
//"log"
//"fmt"
//"bufio"
)

// NodeExp is interface for expressions
type NodeExp interface {
	String() string                   // Literal cosmetic display
	Exp(options *BuildOptions) string // For code generation in Go
}

// NodeExpNumber holds value
type NodeExpNumber struct{ Value string }

// String returns value
func (e *NodeExpNumber) String() string {
	return e.Value
}

// Exp returns value
func (e *NodeExpNumber) Exp(options *BuildOptions) string {
	return expFloat(e.Value)
}

func expFloat(v string) string {
	return "float32(" + v + ")"
}

// NodeExpFloat holds value
type NodeExpFloat struct{ Value string }

// String returns value
func (e *NodeExpFloat) String() string {
	return e.Value
}

// Exp returns value
func (e *NodeExpFloat) Exp(options *BuildOptions) string {
	return expFloat(e.Value)
}

// NodeExpString holds value
type NodeExpString struct{ Value string }

// String returns value
func (e *NodeExpString) String() string {
	return e.Value
}

// Exp returns value
func (e *NodeExpString) Exp(options *BuildOptions) string {
	return e.Value
}

// NodeExpIdentifier holds value
type NodeExpIdentifier struct{ Value string }

// String returns value
func (e *NodeExpIdentifier) String() string {
	return e.Value
}

// Exp returns value
func (e *NodeExpIdentifier) Exp(options *BuildOptions) string {
	return e.Value
}

// NodeExpPlus holds value
type NodeExpPlus struct {
	Left  NodeExp
	Right NodeExp
}

// String returns value
func (e *NodeExpPlus) String() string {
	return e.Left.String() + "+" + e.Right.String()
}

// Exp returns value
func (e *NodeExpPlus) Exp(options *BuildOptions) string {
	return e.Left.Exp(options) + "+" + e.Right.Exp(options)
}

// NodeExpMinus holds value
type NodeExpMinus struct {
	Left  NodeExp
	Right NodeExp
}

// String returns value
func (e *NodeExpMinus) String() string {
	return e.Left.String() + "-" + e.Right.String()
}

// Exp returns value
func (e *NodeExpMinus) Exp(options *BuildOptions) string {
	return e.Left.Exp(options) + "-" + e.Right.Exp(options)
}

// NodeExpMult holds value
type NodeExpMult struct {
	Left  NodeExp
	Right NodeExp
}

// String returns value
func (e *NodeExpMult) String() string {
	return e.Left.String() + "*" + e.Right.String()
}

// Exp returns value
func (e *NodeExpMult) Exp(options *BuildOptions) string {
	return e.Left.Exp(options) + "*" + e.Right.Exp(options)
}

// NodeExpDiv holds value
type NodeExpDiv struct {
	Left  NodeExp
	Right NodeExp
}

// String returns value
func (e *NodeExpDiv) String() string {
	return e.Left.String() + "/" + e.Right.String()
}

// Exp returns value
func (e *NodeExpDiv) Exp(options *BuildOptions) string {
	return e.Left.Exp(options) + "/" + e.Right.Exp(options)
}

// NodeExpDivInt holds value
type NodeExpDivInt struct {
	Left  NodeExp
	Right NodeExp
}

// String returns value
func (e *NodeExpDivInt) String() string {
	return e.Left.String() + "\\" + e.Right.String()
}

// Exp returns value
func (e *NodeExpDivInt) Exp(options *BuildOptions) string {
	options.Headers["math"] = struct{}{}
	return trunc(round(e.Left.Exp(options)) + "/" + round(e.Right.Exp(options)))
}

func trunc(s string) string {
	return "math.Trunc(" + s + ")"
}

func round(s string) string {
	return "math.Round(float64(" + s + "))"
}

// NodeExpUnaryPlus holds value
type NodeExpUnaryPlus struct{ Value NodeExp }

// String returns value
func (e *NodeExpUnaryPlus) String() string {
	return "+" + e.Value.String()
}

// Exp returns value
func (e *NodeExpUnaryPlus) Exp(options *BuildOptions) string {
	return "+" + e.Value.Exp(options)
}

// NodeExpUnaryMinus holds value
type NodeExpUnaryMinus struct{ Value NodeExp }

// String returns value
func (e *NodeExpUnaryMinus) String() string {
	return "-" + e.Value.String()
}

// Exp returns value
func (e *NodeExpUnaryMinus) Exp(options *BuildOptions) string {
	return "-" + e.Value.Exp(options)
}

// NodeExpGroup holds value
type NodeExpGroup struct{ Value NodeExp }

// String returns value
func (e *NodeExpGroup) String() string {
	return "(" + e.Value.String() + ")"
}

// Exp returns value
func (e *NodeExpGroup) Exp(options *BuildOptions) string {
	return "(" + e.Value.Exp(options) + ")"
}

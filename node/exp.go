package node

import (
//"log"
//"fmt"
//"bufio"
//"strconv"
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
	return toFloat(e.Value)
}

func toFloat(v string) string {
	return "float64(" + v + ")"
}

// NodeExpFloat holds value
type NodeExpFloat struct{ Value string }

// String returns value
func (e *NodeExpFloat) String() string {
	return e.Value
}

// Exp returns value
func (e *NodeExpFloat) Exp(options *BuildOptions) string {
	/*
		value, err := strconv.ParseFloat(e.Value, 64)
		if err != nil {
			return fmt.Sprintf("NodeExpFloat.Exp: '%s' error: %v", e.Value, err)
		}
		return fmt.Sprintf("%v", value)
	*/
	if e.Value == "." {
		return "0.0"
	}
	return e.Value
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

// NodeExpMod holds value
type NodeExpMod struct {
	Left  NodeExp
	Right NodeExp
}

// String returns value
func (e *NodeExpMod) String() string {
	return e.Left.String() + " MOD " + e.Right.String()
}

// Exp returns value
func (e *NodeExpMod) Exp(options *BuildOptions) string {
	return toInt(round(options, e.Left.Exp(options))) + `%%` + toInt(round(options, e.Right.Exp(options)))
}

func toInt(s string) string {
	return "int(" + s + ")"
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
	return trunc(options, round(options, e.Left.Exp(options))+"/"+round(options, e.Right.Exp(options)))
}

// NodeExpPow holds value
type NodeExpPow struct {
	Left  NodeExp
	Right NodeExp
}

// String returns value
func (e *NodeExpPow) String() string {
	return e.Left.String() + "^" + e.Right.String()
}

// Exp returns value
func (e *NodeExpPow) Exp(options *BuildOptions) string {
	options.Headers["math"] = struct{}{}
	return "math.Pow(" + e.Left.Exp(options) + "," + e.Right.Exp(options) + ")"
}

func trunc(options *BuildOptions, s string) string {
	options.Headers["math"] = struct{}{}
	return "math.Trunc(" + s + ")"
}

func round(options *BuildOptions, s string) string {
	options.Headers["math"] = struct{}{}
	return "math.Round(" + toFloat(s) + ")"
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

package node

import (
	//"log"
	"fmt"
	//"bufio"
	//"strconv"
)

// Types
const (
	TypeUnknown = iota
	TypeString  = iota
	TypeFloat   = iota
	TypeInteger = iota
)

// NodeExp is interface for expressions
type NodeExp interface {
	String() string                   // Literal cosmetic display
	Exp(options *BuildOptions) string // For code generation in Go
	Type() int
	FindUsedVars(vars map[string]struct{})
}

// NodeExpNumber holds value
type NodeExpNumber struct{ Value string }

// Type returns type
func (e *NodeExpNumber) Type() int {
	return TypeInteger
}

// String returns value
func (e *NodeExpNumber) String() string {
	return e.Value
}

// Exp returns value
func (e *NodeExpNumber) Exp(options *BuildOptions) string {
	return toFloat(e.Value)
}

// FindUsedVars finds used vars
func (e *NodeExpNumber) FindUsedVars(vars map[string]struct{}) {
	// do nothing
}

func toFloat(v string) string {
	return "float64(" + v + ")"
}

// NodeExpFloat holds value
type NodeExpFloat struct{ Value float64 }

// Type returns type
func (e *NodeExpFloat) Type() int {
	return TypeFloat
}

// String returns value
func (e *NodeExpFloat) String() string {
	return fmt.Sprintf("%v", e.Value)
}

// Exp returns value
func (e *NodeExpFloat) Exp(options *BuildOptions) string {
	return fmt.Sprintf("%v", e.Value)
}

// FindUsedVars finds used vars
func (e *NodeExpFloat) FindUsedVars(vars map[string]struct{}) {
	// do nothing
}

// NodeExpString holds value
type NodeExpString struct{ Value string }

// Type returns type
func (e *NodeExpString) Type() int {
	return TypeString
}

// String returns value
func (e *NodeExpString) String() string {
	return e.Value
}

// Exp returns value
func (e *NodeExpString) Exp(options *BuildOptions) string {
	return e.Value
}

// FindUsedVars finds used vars
func (e *NodeExpString) FindUsedVars(vars map[string]struct{}) {
	// do nothing
}

// NodeExpIdentifier holds value
type NodeExpIdentifier struct{ Value string }

// Type returns type
func (e *NodeExpIdentifier) Type() int {
	last := len(e.Value) - 1
	if last < 1 {
		return TypeFloat
	}
	switch e.Value[last] {
	case '$':
		return TypeString
	case '%':
		return TypeInteger
	}
	return TypeFloat
}

// String returns value
func (e *NodeExpIdentifier) String() string {
	return e.Value
}

// Exp returns value
func (e *NodeExpIdentifier) Exp(options *BuildOptions) string {
	return RenameVar(e.Value)
}

// FindUsedVars finds used vars
func (e *NodeExpIdentifier) FindUsedVars(vars map[string]struct{}) {
	vars[e.Value] = struct{}{}
}

// NodeExpPlus holds value
type NodeExpPlus struct {
	Left  NodeExp
	Right NodeExp
}

// Type returns type
func (e *NodeExpPlus) Type() int {
	return combineType(e.Left.Type(), e.Right.Type())
}

func combineType(t1, t2 int) int {
	if t1 == TypeString && t2 == TypeString {
		return TypeString
	}
	if t1 == TypeInteger && t2 == TypeInteger {
		return TypeInteger
	}
	if t1 == TypeFloat && t2 == TypeFloat {
		return TypeFloat
	}
	if t1 == TypeInteger && t2 == TypeFloat {
		return TypeFloat
	}
	if t1 == TypeFloat && t2 == TypeInteger {
		return TypeFloat
	}
	return TypeUnknown
}

// String returns value
func (e *NodeExpPlus) String() string {
	return e.Left.String() + "+" + e.Right.String()
}

// Exp returns value
func (e *NodeExpPlus) Exp(options *BuildOptions) string {
	return e.Left.Exp(options) + "+" + e.Right.Exp(options)
}

// FindUsedVars finds used vars
func (e *NodeExpPlus) FindUsedVars(vars map[string]struct{}) {
	e.Left.FindUsedVars(vars)
	e.Right.FindUsedVars(vars)
}

// NodeExpMinus holds value
type NodeExpMinus struct {
	Left  NodeExp
	Right NodeExp
}

// Type returns type
func (e *NodeExpMinus) Type() int {
	return combineType(e.Left.Type(), e.Right.Type())
}

// String returns value
func (e *NodeExpMinus) String() string {
	return e.Left.String() + "-" + e.Right.String()
}

// Exp returns value
func (e *NodeExpMinus) Exp(options *BuildOptions) string {
	return e.Left.Exp(options) + "-" + e.Right.Exp(options)
}

// FindUsedVars finds used vars
func (e *NodeExpMinus) FindUsedVars(vars map[string]struct{}) {
	e.Left.FindUsedVars(vars)
	e.Right.FindUsedVars(vars)
}

// NodeExpMod holds value
type NodeExpMod struct {
	Left  NodeExp
	Right NodeExp
}

// Type returns type
func (e *NodeExpMod) Type() int {
	return TypeInteger
}

// String returns value
func (e *NodeExpMod) String() string {
	return e.Left.String() + " MOD " + e.Right.String()
}

// Exp returns value
func (e *NodeExpMod) Exp(options *BuildOptions) string {
	return toInt(round(options, e.Left.Exp(options))) + `%%` + toInt(round(options, e.Right.Exp(options)))
}

// FindUsedVars finds used vars
func (e *NodeExpMod) FindUsedVars(vars map[string]struct{}) {
	e.Left.FindUsedVars(vars)
	e.Right.FindUsedVars(vars)
}

func toInt(s string) string {
	return "int(" + s + ")"
}

// NodeExpMult holds value
type NodeExpMult struct {
	Left  NodeExp
	Right NodeExp
}

// Type returns type
func (e *NodeExpMult) Type() int {
	return combineType(e.Left.Type(), e.Right.Type())
}

// String returns value
func (e *NodeExpMult) String() string {
	return e.Left.String() + "*" + e.Right.String()
}

// Exp returns value
func (e *NodeExpMult) Exp(options *BuildOptions) string {
	return e.Left.Exp(options) + "*" + e.Right.Exp(options)
}

// FindUsedVars finds used vars
func (e *NodeExpMult) FindUsedVars(vars map[string]struct{}) {
	e.Left.FindUsedVars(vars)
	e.Right.FindUsedVars(vars)
}

// NodeExpDiv holds value
type NodeExpDiv struct {
	Left  NodeExp
	Right NodeExp
}

// Type returns type
func (e *NodeExpDiv) Type() int {
	return TypeFloat
}

// String returns value
func (e *NodeExpDiv) String() string {
	return e.Left.String() + "/" + e.Right.String()
}

// Exp returns value
func (e *NodeExpDiv) Exp(options *BuildOptions) string {
	return e.Left.Exp(options) + "/" + e.Right.Exp(options)
}

// FindUsedVars finds used vars
func (e *NodeExpDiv) FindUsedVars(vars map[string]struct{}) {
	e.Left.FindUsedVars(vars)
	e.Right.FindUsedVars(vars)
}

// NodeExpDivInt holds value
type NodeExpDivInt struct {
	Left  NodeExp
	Right NodeExp
}

// Type returns type
func (e *NodeExpDivInt) Type() int {
	return combineType(e.Left.Type(), e.Right.Type())
}

// String returns value
func (e *NodeExpDivInt) String() string {
	return e.Left.String() + "\\" + e.Right.String()
}

// Exp returns value
func (e *NodeExpDivInt) Exp(options *BuildOptions) string {
	return trunc(options, round(options, e.Left.Exp(options))+"/"+round(options, e.Right.Exp(options)))
}

// FindUsedVars finds used vars
func (e *NodeExpDivInt) FindUsedVars(vars map[string]struct{}) {
	e.Left.FindUsedVars(vars)
	e.Right.FindUsedVars(vars)
}

// NodeExpPow holds value
type NodeExpPow struct {
	Left  NodeExp
	Right NodeExp
}

// Type returns type
func (e *NodeExpPow) Type() int {
	return combineType(e.Left.Type(), e.Right.Type())
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

// FindUsedVars finds used vars
func (e *NodeExpPow) FindUsedVars(vars map[string]struct{}) {
	e.Left.FindUsedVars(vars)
	e.Right.FindUsedVars(vars)
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

// Type returns type
func (e *NodeExpUnaryPlus) Type() int {
	return e.Value.Type()
}

// String returns value
func (e *NodeExpUnaryPlus) String() string {
	return "+" + e.Value.String()
}

// Exp returns value
func (e *NodeExpUnaryPlus) Exp(options *BuildOptions) string {
	return "+" + e.Value.Exp(options)
}

// FindUsedVars finds used vars
func (e *NodeExpUnaryPlus) FindUsedVars(vars map[string]struct{}) {
	e.Value.FindUsedVars(vars)
}

// NodeExpUnaryMinus holds value
type NodeExpUnaryMinus struct{ Value NodeExp }

// Type returns type
func (e *NodeExpUnaryMinus) Type() int {
	return e.Value.Type()
}

// String returns value
func (e *NodeExpUnaryMinus) String() string {
	return "-" + e.Value.String()
}

// Exp returns value
func (e *NodeExpUnaryMinus) Exp(options *BuildOptions) string {
	return "-" + e.Value.Exp(options)
}

// FindUsedVars finds used vars
func (e *NodeExpUnaryMinus) FindUsedVars(vars map[string]struct{}) {
	e.Value.FindUsedVars(vars)
}

// NodeExpGroup holds value
type NodeExpGroup struct{ Value NodeExp }

// Type returns type
func (e *NodeExpGroup) Type() int {
	return e.Value.Type()
}

// String returns value
func (e *NodeExpGroup) String() string {
	return "(" + e.Value.String() + ")"
}

// Exp returns value
func (e *NodeExpGroup) Exp(options *BuildOptions) string {
	return "(" + e.Value.Exp(options) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpGroup) FindUsedVars(vars map[string]struct{}) {
	e.Value.FindUsedVars(vars)
}

// NodeExpLen holds value
type NodeExpLen struct {
	Value NodeExp
}

// Type returns type
func (e *NodeExpLen) Type() int {
	return TypeInteger
}

// String returns value
func (e *NodeExpLen) String() string {
	return "LEN(" + e.Value.String() + ")"
}

// Exp returns value
func (e *NodeExpLen) Exp(options *BuildOptions) string {
	if e.Value.Type() == TypeString {
		return "len(" + e.Value.Exp(options) + ")"
	}
	return "8 /* LEN(non-string) */"
}

// FindUsedVars finds used vars
func (e *NodeExpLen) FindUsedVars(vars map[string]struct{}) {
	if e.Value.Type() == TypeString {
		e.Value.FindUsedVars(vars)
	}
}

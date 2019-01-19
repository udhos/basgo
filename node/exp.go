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

// TypeNumeric reports whether type is numeric.
func TypeNumeric(t int) bool {
	return t == TypeFloat || t == TypeInteger
}

// TypeCompare reports whether types are comparable.
func TypeCompare(t1, t2 int) bool {
	return (t1 == TypeString && t2 == TypeString) || (TypeNumeric(t1) && TypeNumeric(t2))
}

// NodeExp is interface for expressions
type NodeExp interface {
	String() string                   // Literal cosmetic display
	Exp(options *BuildOptions) string // For code generation in Go
	Type() int
	FindUsedVars(options *BuildOptions)
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
	return e.Value
}

// FindUsedVars finds used vars
func (e *NodeExpNumber) FindUsedVars(options *BuildOptions) {
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
func (e *NodeExpFloat) FindUsedVars(options *BuildOptions) {
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
func (e *NodeExpString) FindUsedVars(options *BuildOptions) {
	// do nothing
}

// NodeExpIdentifier holds value
type NodeExpIdentifier struct{ Value string }

// Type returns type
func (e *NodeExpIdentifier) Type() int {
	return VarType(e.Value)
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
func (e *NodeExpIdentifier) FindUsedVars(options *BuildOptions) {
	options.VarSetUsed(e.Value)
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
	return "(" + e.Left.String() + ") + (" + e.Right.String() + ")"
}

// Exp returns value
func (e *NodeExpPlus) Exp(options *BuildOptions) string {
	return "(" + e.Left.Exp(options) + ")+(" + e.Right.Exp(options) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpPlus) FindUsedVars(options *BuildOptions) {
	e.Left.FindUsedVars(options)
	e.Right.FindUsedVars(options)
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
	return "(" + e.Left.String() + ") - (" + e.Right.String() + ")"
}

// Exp returns value
func (e *NodeExpMinus) Exp(options *BuildOptions) string {
	return "(" + e.Left.Exp(options) + ")-(" + e.Right.Exp(options) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpMinus) FindUsedVars(options *BuildOptions) {
	e.Left.FindUsedVars(options)
	e.Right.FindUsedVars(options)
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
	return "(" + e.Left.String() + ") MOD (" + e.Right.String() + ")"
}

// Exp returns value
func (e *NodeExpMod) Exp(options *BuildOptions) string {
	return "(" + forceInt(options, e.Left) + `)%%(` + forceInt(options, e.Right) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpMod) FindUsedVars(options *BuildOptions) {
	e.Left.FindUsedVars(options)
	e.Right.FindUsedVars(options)
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
	return "(" + e.Left.String() + ") * (" + e.Right.String() + ")"
}

// Exp returns value
func (e *NodeExpMult) Exp(options *BuildOptions) string {
	return "(" + e.Left.Exp(options) + ")*(" + e.Right.Exp(options) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpMult) FindUsedVars(options *BuildOptions) {
	e.Left.FindUsedVars(options)
	e.Right.FindUsedVars(options)
}

// NodeExpDiv holds value
type NodeExpDiv struct {
	Left  NodeExp
	Right NodeExp
}

// Type returns type
func (e *NodeExpDiv) Type() int {
	return TypeFloat // remember: 5 / 2 = float
}

// String returns value
func (e *NodeExpDiv) String() string {
	return "(" + e.Left.String() + ") / (" + e.Right.String() + ")"
}

// Exp returns value
func (e *NodeExpDiv) Exp(options *BuildOptions) string {
	return "(" + forceFloat(options, e.Left) + ")/(" + forceFloat(options, e.Right) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpDiv) FindUsedVars(options *BuildOptions) {
	e.Left.FindUsedVars(options)
	e.Right.FindUsedVars(options)
}

// NodeExpDivInt holds value
type NodeExpDivInt struct {
	Left  NodeExp
	Right NodeExp
}

// Type returns type
func (e *NodeExpDivInt) Type() int {
	return TypeInteger
}

// String returns value
func (e *NodeExpDivInt) String() string {
	return "(" + e.Left.String() + ") \\ (" + e.Right.String() + ")"
}

// Exp returns value
func (e *NodeExpDivInt) Exp(options *BuildOptions) string {
	return "(" + forceInt(options, e.Left) + `)/(` + forceInt(options, e.Right) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpDivInt) FindUsedVars(options *BuildOptions) {
	e.Left.FindUsedVars(options)
	e.Right.FindUsedVars(options)
}

// NodeExpPow holds value
type NodeExpPow struct {
	Left  NodeExp
	Right NodeExp
}

// Type returns type
func (e *NodeExpPow) Type() int {
	return TypeFloat
}

// String returns value
func (e *NodeExpPow) String() string {
	return "(" + e.Left.String() + " ^ " + e.Right.String() + ")"
}

// Exp returns value
func (e *NodeExpPow) Exp(options *BuildOptions) string {
	options.Headers["math"] = struct{}{}
	return "math.Pow(" + forceFloat(options, e.Left) + "," + forceFloat(options, e.Right) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpPow) FindUsedVars(options *BuildOptions) {
	e.Left.FindUsedVars(options)
	e.Right.FindUsedVars(options)
}

// NodeExpUnaryPlus holds value
type NodeExpUnaryPlus struct{ Value NodeExp }

// Type returns type
func (e *NodeExpUnaryPlus) Type() int {
	return e.Value.Type()
}

// String returns value
func (e *NodeExpUnaryPlus) String() string {
	return "+(" + e.Value.String() + ")"
}

// Exp returns value
func (e *NodeExpUnaryPlus) Exp(options *BuildOptions) string {
	return "+(" + e.Value.Exp(options) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpUnaryPlus) FindUsedVars(options *BuildOptions) {
	e.Value.FindUsedVars(options)
}

// NodeExpUnaryMinus holds value
type NodeExpUnaryMinus struct{ Value NodeExp }

// Type returns type
func (e *NodeExpUnaryMinus) Type() int {
	return e.Value.Type()
}

// String returns value
func (e *NodeExpUnaryMinus) String() string {
	return "-(" + e.Value.String() + ")"
}

// Exp returns value
func (e *NodeExpUnaryMinus) Exp(options *BuildOptions) string {
	return "-(" + e.Value.Exp(options) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpUnaryMinus) FindUsedVars(options *BuildOptions) {
	e.Value.FindUsedVars(options)
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
func (e *NodeExpGroup) FindUsedVars(options *BuildOptions) {
	e.Value.FindUsedVars(options)
}

// NodeExpInt holds value
type NodeExpInt struct {
	Value NodeExp
}

// Type returns type
func (e *NodeExpInt) Type() int {
	return TypeInteger
}

// String returns value
func (e *NodeExpInt) String() string {
	return "INT(" + e.Value.String() + ")"
}

// Exp returns value
func (e *NodeExpInt) Exp(options *BuildOptions) string {
	return "(" + floorInt(options, e.Value) + ")"
}

func floorInt(options *BuildOptions, e NodeExp) string {
	s := e.Exp(options)
	if e.Type() != TypeInteger {
		options.Headers["math"] = struct{}{}
		return toInt("math.Floor(" + s + ") /* <- floorInt(non-int) */")
	}
	return s
}

// FindUsedVars finds used vars
func (e *NodeExpInt) FindUsedVars(options *BuildOptions) {
	e.Value.FindUsedVars(options)
}

// NodeExpLeft holds value
type NodeExpLeft struct {
	Value NodeExp
	Size  NodeExp
}

// Type returns type
func (e *NodeExpLeft) Type() int {
	return TypeString
}

// String returns value
func (e *NodeExpLeft) String() string {
	return "LEFT(" + e.Value.String() + "," + e.Size.String() + ")"
}

// Exp returns value
func (e *NodeExpLeft) Exp(options *BuildOptions) string {
	options.Left = true
	return "stringLeft(" + e.Value.Exp(options) + "," + forceInt(options, e.Size) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpLeft) FindUsedVars(options *BuildOptions) {
	e.Value.FindUsedVars(options)
	e.Size.FindUsedVars(options)
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
	return "8 /* <- LEN(non-string) */"
}

// FindUsedVars finds used vars
func (e *NodeExpLen) FindUsedVars(options *BuildOptions) {
	if e.Value.Type() == TypeString {
		e.Value.FindUsedVars(options)
	}
}

// NodeExpRnd holds value
type NodeExpRnd struct {
	Value NodeExp
}

// Type returns type
func (e *NodeExpRnd) Type() int {
	return TypeFloat
}

// String returns value
func (e *NodeExpRnd) String() string {
	return "RND"
}

// Exp returns value
func (e *NodeExpRnd) Exp(options *BuildOptions) string {
	options.Headers["time"] = struct{}{}
	options.Headers["math/rand"] = struct{}{}
	options.Rnd = true

	return "randomFloat64(" + e.Value.Exp(options) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpRnd) FindUsedVars(options *BuildOptions) {
	// do nothing
}

// NodeExpNot holds value
type NodeExpNot struct {
	Value NodeExp
}

// Type returns type
func (e *NodeExpNot) Type() int {
	return TypeInteger
}

// String returns value
func (e *NodeExpNot) String() string {
	return "NOT(" + e.Value.String() + ")"
}

// Exp returns value
func (e *NodeExpNot) Exp(options *BuildOptions) string {
	return "^(" + forceInt(options, e.Value) + ")"
}

func forceInt(options *BuildOptions, e NodeExp) string {
	s := e.Exp(options)
	if e.Type() != TypeInteger {
		options.Headers["math"] = struct{}{}
		return toInt("math.Round(" + s + ") /* <- forceInt(non-int) */")
	}
	return s
}

func forceFloat(options *BuildOptions, e NodeExp) string {
	s := e.Exp(options)
	if e.Type() != TypeFloat {
		return toFloat(s) + " /* <- forceFloat(non-float) */"
	}
	return s
}

// FindUsedVars finds used vars
func (e *NodeExpNot) FindUsedVars(options *BuildOptions) {
	e.Value.FindUsedVars(options)
}

// NodeExpAnd holds value
type NodeExpAnd struct {
	Left  NodeExp
	Right NodeExp
}

// Type returns type
func (e *NodeExpAnd) Type() int {
	return TypeInteger
}

// String returns value
func (e *NodeExpAnd) String() string {
	return "(" + e.Left.String() + ") AND (" + e.Right.String() + ")"
}

// Exp returns value
func (e *NodeExpAnd) Exp(options *BuildOptions) string {
	return "(" + forceInt(options, e.Left) + ")&(" + forceInt(options, e.Right) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpAnd) FindUsedVars(options *BuildOptions) {
	e.Left.FindUsedVars(options)
	e.Right.FindUsedVars(options)
}

// NodeExpEqv holds value
type NodeExpEqv struct {
	Left  NodeExp
	Right NodeExp
}

// Type returns type
func (e *NodeExpEqv) Type() int {
	return TypeInteger
}

// String returns value
func (e *NodeExpEqv) String() string {
	return "(" + e.Left.String() + ") EQV (" + e.Right.String() + ")"
}

// Exp returns value
func (e *NodeExpEqv) Exp(options *BuildOptions) string {
	return "^((" + forceInt(options, e.Left) + ")^(" + forceInt(options, e.Right) + "))"
}

// FindUsedVars finds used vars
func (e *NodeExpEqv) FindUsedVars(options *BuildOptions) {
	e.Left.FindUsedVars(options)
	e.Right.FindUsedVars(options)
}

// NodeExpImp holds value
type NodeExpImp struct {
	Left  NodeExp
	Right NodeExp
}

// Type returns type
func (e *NodeExpImp) Type() int {
	return TypeInteger
}

// String returns value
func (e *NodeExpImp) String() string {
	return "(" + e.Left.String() + ") IMP (" + e.Right.String() + ")"
}

// Exp returns value
func (e *NodeExpImp) Exp(options *BuildOptions) string {
	return "((^(" + forceInt(options, e.Left) + "))|(" + forceInt(options, e.Right) + "))"
}

// FindUsedVars finds used vars
func (e *NodeExpImp) FindUsedVars(options *BuildOptions) {
	e.Left.FindUsedVars(options)
	e.Right.FindUsedVars(options)
}

// NodeExpOr holds value
type NodeExpOr struct {
	Left  NodeExp
	Right NodeExp
}

// Type returns type
func (e *NodeExpOr) Type() int {
	return TypeInteger
}

// String returns value
func (e *NodeExpOr) String() string {
	return "(" + e.Left.String() + ") OR (" + e.Right.String() + ")"
}

// Exp returns value
func (e *NodeExpOr) Exp(options *BuildOptions) string {
	return "(" + forceInt(options, e.Left) + ")|(" + forceInt(options, e.Right) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpOr) FindUsedVars(options *BuildOptions) {
	e.Left.FindUsedVars(options)
	e.Right.FindUsedVars(options)
}

// NodeExpXor holds value
type NodeExpXor struct {
	Left  NodeExp
	Right NodeExp
}

// Type returns type
func (e *NodeExpXor) Type() int {
	return TypeInteger
}

// String returns value
func (e *NodeExpXor) String() string {
	return "(" + e.Left.String() + ") XOR (" + e.Right.String() + ")"
}

// Exp returns value
func (e *NodeExpXor) Exp(options *BuildOptions) string {
	return "(" + forceInt(options, e.Left) + ")^(" + forceInt(options, e.Right) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpXor) FindUsedVars(options *BuildOptions) {
	e.Left.FindUsedVars(options)
	e.Right.FindUsedVars(options)
}

type NodeExpBinary interface {
	Values() (left, right NodeExp)
}

// NodeExpEqual holds value
type NodeExpEqual struct {
	Left  NodeExp
	Right NodeExp
}

// Type returns type
func (e *NodeExpEqual) Type() int {
	return TypeInteger
}

// Values returns children from binary exp
func (e *NodeExpEqual) Values() (NodeExp, NodeExp) {
	return e.Left, e.Right
}

// String returns value
func (e *NodeExpEqual) String() string {
	return "(" + e.Left.String() + ") = (" + e.Right.String() + ")"
}

// Exp returns value
func (e *NodeExpEqual) Exp(options *BuildOptions) string {
	return compareOp(e, options, "==")
}

func compareOp(e NodeExpBinary, options *BuildOptions, golangOp string) string {
	left, right := e.Values()
	return boolToInt("(" + left.Exp(options) + ")" + golangOp + "(" + right.Exp(options) + ")")
}

func boolToInt(s string) string {
	return fmt.Sprintf("boolToInt(%s)", s)
}

// FindUsedVars finds used vars
func (e *NodeExpEqual) FindUsedVars(options *BuildOptions) {
	e.Left.FindUsedVars(options)
	e.Right.FindUsedVars(options)
}

// NodeExpUnequal holds value
type NodeExpUnequal struct {
	Left  NodeExp
	Right NodeExp
}

// Type returns type
func (e *NodeExpUnequal) Type() int {
	return TypeInteger
}

// Values returns children from binary exp
func (e *NodeExpUnequal) Values() (NodeExp, NodeExp) {
	return e.Left, e.Right
}

// String returns value
func (e *NodeExpUnequal) String() string {
	return "(" + e.Left.String() + ") <> (" + e.Right.String() + ")"
}

// Exp returns value
func (e *NodeExpUnequal) Exp(options *BuildOptions) string {
	return compareOp(e, options, "!=")
}

// FindUsedVars finds used vars
func (e *NodeExpUnequal) FindUsedVars(options *BuildOptions) {
	e.Left.FindUsedVars(options)
	e.Right.FindUsedVars(options)
}

// NodeExpGT holds value
type NodeExpGT struct {
	Left  NodeExp
	Right NodeExp
}

// Type returns type
func (e *NodeExpGT) Type() int {
	return TypeInteger
}

// Values returns children from binary exp
func (e *NodeExpGT) Values() (NodeExp, NodeExp) {
	return e.Left, e.Right
}

// String returns value
func (e *NodeExpGT) String() string {
	return "(" + e.Left.String() + ") > (" + e.Right.String() + ")"
}

// Exp returns value
func (e *NodeExpGT) Exp(options *BuildOptions) string {
	return compareOp(e, options, ">")
}

// FindUsedVars finds used vars
func (e *NodeExpGT) FindUsedVars(options *BuildOptions) {
	e.Left.FindUsedVars(options)
	e.Right.FindUsedVars(options)
}

// NodeExpLT holds value
type NodeExpLT struct {
	Left  NodeExp
	Right NodeExp
}

// Type returns type
func (e *NodeExpLT) Type() int {
	return TypeInteger
}

// Values returns children from binary exp
func (e *NodeExpLT) Values() (NodeExp, NodeExp) {
	return e.Left, e.Right
}

// String returns value
func (e *NodeExpLT) String() string {
	return "(" + e.Left.String() + ") < (" + e.Right.String() + ")"
}

// Exp returns value
func (e *NodeExpLT) Exp(options *BuildOptions) string {
	return compareOp(e, options, "<")
}

// FindUsedVars finds used vars
func (e *NodeExpLT) FindUsedVars(options *BuildOptions) {
	e.Left.FindUsedVars(options)
	e.Right.FindUsedVars(options)
}

// NodeExpGE holds value
type NodeExpGE struct {
	Left  NodeExp
	Right NodeExp
}

// Type returns type
func (e *NodeExpGE) Type() int {
	return TypeInteger
}

// Values returns children from binary exp
func (e *NodeExpGE) Values() (NodeExp, NodeExp) {
	return e.Left, e.Right
}

// String returns value
func (e *NodeExpGE) String() string {
	return "(" + e.Left.String() + ") >= (" + e.Right.String() + ")"
}

// Exp returns value
func (e *NodeExpGE) Exp(options *BuildOptions) string {
	return compareOp(e, options, ">=")
}

// FindUsedVars finds used vars
func (e *NodeExpGE) FindUsedVars(options *BuildOptions) {
	e.Left.FindUsedVars(options)
	e.Right.FindUsedVars(options)
}

// NodeExpLE holds value
type NodeExpLE struct {
	Left  NodeExp
	Right NodeExp
}

// Type returns type
func (e *NodeExpLE) Type() int {
	return TypeInteger
}

// Values returns children from binary exp
func (e *NodeExpLE) Values() (NodeExp, NodeExp) {
	return e.Left, e.Right
}

// String returns value
func (e *NodeExpLE) String() string {
	return "(" + e.Left.String() + ") <= (" + e.Right.String() + ")"
}

// Exp returns value
func (e *NodeExpLE) Exp(options *BuildOptions) string {
	return compareOp(e, options, "<=")
}

// FindUsedVars finds used vars
func (e *NodeExpLE) FindUsedVars(options *BuildOptions) {
	e.Left.FindUsedVars(options)
	e.Right.FindUsedVars(options)
}

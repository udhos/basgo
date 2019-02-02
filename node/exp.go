package node

import (
	"fmt"
	//"log"
	//"bufio"
	//"strconv"
	"strings"
)

// Types
const (
	TypeUnknown = iota
	TypeString  = iota
	TypeFloat   = iota
	TypeInteger = iota
)

func TypeLabel(t int) string {
	switch t {
	case TypeString:
		return "STRING"
	case TypeFloat:
		return "FLOAT"
	case TypeInteger:
		return "INTEGER"
	}
	return "UNKNOWN"
}

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

func NewNodeExpString(s string) *NodeExpString {
	q := `"`
	return &NodeExpString{Value: strings.TrimSuffix(strings.TrimPrefix(s, q), q)}
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
	return "`" + e.Value + "`"
}

// FindUsedVars finds used vars
func (e *NodeExpString) FindUsedVars(options *BuildOptions) {
	// do nothing
}

// NodeExpIdentifier holds value
type NodeExpIdentifier struct {
	Value string
}

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

// NodeExpArray holds value
type NodeExpArray struct {
	Name    string
	Indices []NodeExp
}

// Type returns type
func (e *NodeExpArray) Type() int {
	return VarType(e.Name)
}

// String returns value
func (e *NodeExpArray) String() string {
	str := e.Name + "(" + e.Indices[0].String()
	for i := 1; i < len(e.Indices); i++ {
		str += "," + e.Indices[i].String()
	}
	str += ")"
	return str
}

// Exp returns value
func (e *NodeExpArray) Exp(options *BuildOptions) string {
	str := RenameArray(e.Name)
	for _, i := range e.Indices {
		str += "[" + forceInt(options, i) + "]"
	}
	return str
}

// FindUsedVars finds used vars
func (e *NodeExpArray) FindUsedVars(options *BuildOptions) {
	//options.ArraySetUsed(e.Name, len(e.Indices))
	for _, i := range e.Indices {
		i.FindUsedVars(options)
	}
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

// Promotes Integer to Float if needed
func combineNumeric(options *BuildOptions, e1, e2 NodeExp) (string, string) {
	t1 := e1.Type()
	t2 := e2.Type()
	if t1 == TypeInteger && t2 == TypeFloat {
		return forceFloat(options, e1), e2.Exp(options)
	}
	if t1 == TypeFloat && t2 == TypeInteger {
		return e1.Exp(options), forceFloat(options, e2)
	}
	return e1.Exp(options), e2.Exp(options)
}

// String returns value
func (e *NodeExpPlus) String() string {
	return "(" + e.Left.String() + ") + (" + e.Right.String() + ")"
}

// Exp returns value
func (e *NodeExpPlus) Exp(options *BuildOptions) string {
	left, right := combineNumeric(options, e.Left, e.Right)
	return "(" + left + ")+(" + right + ")"
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
	left, right := combineNumeric(options, e.Left, e.Right)
	return "(" + left + ")-(" + right + ")"
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
	left, right := combineNumeric(options, e.Left, e.Right)
	return "(" + left + ")*(" + right + ")"
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

// NodeExpInkey holds value
type NodeExpInkey struct{}

// Type returns type
func (e *NodeExpInkey) Type() int {
	return TypeString
}

// String returns value
func (e *NodeExpInkey) String() string {
	return "INKEY$"
}

// Exp returns value
func (e *NodeExpInkey) Exp(options *BuildOptions) string {
	return "inputInkey()"
}

// FindUsedVars finds used vars
func (e *NodeExpInkey) FindUsedVars(options *BuildOptions) {
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
	return "LEFT$(" + e.Value.String() + "," + e.Size.String() + ")"
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

// NodeExpLeft holds value
type NodeExpRight struct {
	Value NodeExp
	Size  NodeExp
}

// Type returns type
func (e *NodeExpRight) Type() int {
	return TypeString
}

// String returns value
func (e *NodeExpRight) String() string {
	return "RIGHT$(" + e.Value.String() + "," + e.Size.String() + ")"
}

// Exp returns value
func (e *NodeExpRight) Exp(options *BuildOptions) string {
	return "stringRight(" + e.Value.Exp(options) + "," + forceInt(options, e.Size) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpRight) FindUsedVars(options *BuildOptions) {
	e.Value.FindUsedVars(options)
	e.Size.FindUsedVars(options)
}

// NodeExpMid holds value
type NodeExpMid struct {
	Value NodeExp
	Begin NodeExp
	Size  NodeExp
}

// Type returns type
func (e *NodeExpMid) Type() int {
	return TypeString
}

// String returns value
func (e *NodeExpMid) String() string {
	if e.Size == NodeExp(nil) {
		return "MID$(" + e.Value.String() + "," + e.Begin.String() + ")"
	}
	return "MID$(" + e.Value.String() + "," + e.Begin.String() + "," + e.Size.String() + ")"
}

// Exp returns value
func (e *NodeExpMid) Exp(options *BuildOptions) string {
	options.Mid = true
	if e.Size == NodeExp(nil) {
		return "stringMid(" + e.Value.Exp(options) + "," + forceInt(options, e.Begin) + ")"
	}
	return "stringMidSize(" + e.Value.Exp(options) + "," + forceInt(options, e.Begin) + "," + forceInt(options, e.Size) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpMid) FindUsedVars(options *BuildOptions) {
	e.Value.FindUsedVars(options)
	e.Begin.FindUsedVars(options)
	if e.Size != NodeExp(nil) {
		e.Size.FindUsedVars(options)
	}
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

	return "randomFloat64(" + forceFloat(options, e.Value) + ")"
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

// NodeExpBinary helper interface
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
	strLeft, strRight := combineNumeric(options, left, right)
	return boolToInt("(" + strLeft + ")" + golangOp + "(" + strRight + ")")
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

// NodeExpStr holds value
type NodeExpStr struct {
	Value NodeExp
}

// Type returns type
func (e *NodeExpStr) Type() int {
	return TypeString
}

// String returns value
func (e *NodeExpStr) String() string {
	return "STR$(" + e.Value.String() + ")"
}

// Exp returns value
func (e *NodeExpStr) Exp(options *BuildOptions) string {
	options.Headers["strconv"] = struct{}{}

	v := e.Value.Exp(options)

	if e.Value.Type() == TypeInteger {
		return "strconv.Itoa(" + v + ")"
	} else {
		return "strconv.FormatFloat(" + v + ", 'f', -1, 64)"
	}
}

// FindUsedVars finds used vars
func (e *NodeExpStr) FindUsedVars(options *BuildOptions) {
	e.Value.FindUsedVars(options)
}

// NodeExpVal holds value
type NodeExpVal struct {
	Value NodeExp
}

// Type returns type
func (e *NodeExpVal) Type() int {
	return TypeFloat
}

// String returns value
func (e *NodeExpVal) String() string {
	return "VAL(" + e.Value.String() + ")"
}

// Exp returns value
func (e *NodeExpVal) Exp(options *BuildOptions) string {
	return "stringToFloat(" + e.Value.Exp(options) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpVal) FindUsedVars(options *BuildOptions) {
	e.Value.FindUsedVars(options)
}

// NodeExpTab holds value
type NodeExpTab struct {
	Value NodeExp
}

// Type returns type
func (e *NodeExpTab) Type() int {
	return TypeString
}

// String returns value
func (e *NodeExpTab) String() string {
	return "TAB(" + e.Value.String() + ")"
}

// Exp returns value
func (e *NodeExpTab) Exp(options *BuildOptions) string {
	return `stringRepeat(" ",` + forceInt(options, e.Value) + "-1)"
}

// FindUsedVars finds used vars
func (e *NodeExpTab) FindUsedVars(options *BuildOptions) {
	e.Value.FindUsedVars(options)
}

// NodeExpSpc holds value
type NodeExpSpc struct {
	Value NodeExp
}

// Type returns type
func (e *NodeExpSpc) Type() int {
	return TypeString
}

// String returns value
func (e *NodeExpSpc) String() string {
	return "SPC(" + e.Value.String() + ")"
}

// Exp returns value
func (e *NodeExpSpc) Exp(options *BuildOptions) string {
	return `stringRepeat(" ",` + forceInt(options, e.Value) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpSpc) FindUsedVars(options *BuildOptions) {
	e.Value.FindUsedVars(options)
}

// NodeExpSpace holds value
type NodeExpSpace struct {
	Value NodeExp
}

// Type returns type
func (e *NodeExpSpace) Type() int {
	return TypeString
}

// String returns value
func (e *NodeExpSpace) String() string {
	return "SPACE$(" + e.Value.String() + ")"
}

// Exp returns value
func (e *NodeExpSpace) Exp(options *BuildOptions) string {
	return `stringRepeat(" ",` + forceInt(options, e.Value) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpSpace) FindUsedVars(options *BuildOptions) {
	e.Value.FindUsedVars(options)
}

// NodeExpFuncString holds value
type NodeExpFuncString struct {
	Value NodeExp
	Char  NodeExp
}

// Type returns type
func (e *NodeExpFuncString) Type() int {
	return TypeString
}

// String returns value
func (e *NodeExpFuncString) String() string {
	return "STRING$(" + e.Value.String() + "," + e.Char.String() + ")"
}

// Exp returns value
func (e *NodeExpFuncString) Exp(options *BuildOptions) string {
	if TypeNumeric(e.Char.Type()) {
		str := "string([]byte{byte(" + forceInt(options, e.Char) + ")})"
		return "stringRepeat(" + str + "," + forceInt(options, e.Value) + ")"
	}

	options.Left = true

	left := "stringLeft(" + e.Char.Exp(options) + ",1)"

	return `stringRepeat(` + left + `,` + forceInt(options, e.Value) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpFuncString) FindUsedVars(options *BuildOptions) {
	e.Value.FindUsedVars(options)
}

// NodeExpAsc holds value
type NodeExpAsc struct {
	Value NodeExp
}

// Type returns type
func (e *NodeExpAsc) Type() int {
	return TypeInteger
}

// String returns value
func (e *NodeExpAsc) String() string {
	return "ASC(" + e.Value.String() + ")"
}

// Exp returns value
func (e *NodeExpAsc) Exp(options *BuildOptions) string {
	return "firstByte(" + e.Value.Exp(options) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpAsc) FindUsedVars(options *BuildOptions) {
	e.Value.FindUsedVars(options)
}

// NodeExpChr holds value
type NodeExpChr struct {
	Value NodeExp
}

// Type returns type
func (e *NodeExpChr) Type() int {
	return TypeString
}

// String returns value
func (e *NodeExpChr) String() string {
	return "CHR$(" + e.Value.String() + ")"
}

// Exp returns value
func (e *NodeExpChr) Exp(options *BuildOptions) string {
	return "string(" + forceInt(options, e.Value) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpChr) FindUsedVars(options *BuildOptions) {
	e.Value.FindUsedVars(options)
}

// NodeExpDate holds value
type NodeExpDate struct {
}

// Type returns type
func (e *NodeExpDate) Type() int {
	return TypeString
}

// String returns value
func (e *NodeExpDate) String() string {
	return "DATE$"
}

// Exp returns value
func (e *NodeExpDate) Exp(options *BuildOptions) string {
	return "timeDate()"
}

// FindUsedVars finds used vars
func (e *NodeExpDate) FindUsedVars(options *BuildOptions) {
}

// NodeExpTime holds value
type NodeExpTime struct {
}

// Type returns type
func (e *NodeExpTime) Type() int {
	return TypeString
}

// String returns value
func (e *NodeExpTime) String() string {
	return "TIME$"
}

// Exp returns value
func (e *NodeExpTime) Exp(options *BuildOptions) string {
	return "timeTime()"
}

// FindUsedVars finds used vars
func (e *NodeExpTime) FindUsedVars(options *BuildOptions) {
}

// NodeExpTimer holds value
type NodeExpTimer struct {
}

// Type returns type
func (e *NodeExpTimer) Type() int {
	return TypeFloat
}

// String returns value
func (e *NodeExpTimer) String() string {
	return "TIMER"
}

// Exp returns value
func (e *NodeExpTimer) Exp(options *BuildOptions) string {
	return "timeTimer()"
}

// FindUsedVars finds used vars
func (e *NodeExpTimer) FindUsedVars(options *BuildOptions) {
}

// NodeExpAbs holds value
type NodeExpAbs struct {
	Value NodeExp
}

// Type returns type
func (e *NodeExpAbs) Type() int {
	return TypeFloat
}

// String returns value
func (e *NodeExpAbs) String() string {
	return "ABS(" + e.Value.String() + ")"
}

// Exp returns value
func (e *NodeExpAbs) Exp(options *BuildOptions) string {
	return "math.Abs(" + forceFloat(options, e.Value) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpAbs) FindUsedVars(options *BuildOptions) {
	e.Value.FindUsedVars(options)
}

// NodeExpSgn holds value
type NodeExpSgn struct {
	Value NodeExp
}

// Type returns type
func (e *NodeExpSgn) Type() int {
	return TypeInteger
}

// String returns value
func (e *NodeExpSgn) String() string {
	return "SGN(" + e.Value.String() + ")"
}

// Exp returns value
func (e *NodeExpSgn) Exp(options *BuildOptions) string {
	return "baslib.Sgn(" + forceFloat(options, e.Value) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpSgn) FindUsedVars(options *BuildOptions) {
	e.Value.FindUsedVars(options)
}

// NodeExpCos holds value
type NodeExpCos struct {
	Value NodeExp
}

// Type returns type
func (e *NodeExpCos) Type() int {
	return TypeFloat
}

// String returns value
func (e *NodeExpCos) String() string {
	return "COS"
}

// Exp returns value
func (e *NodeExpCos) Exp(options *BuildOptions) string {
	return "math.Cos(" + forceFloat(options, e.Value) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpCos) FindUsedVars(options *BuildOptions) {
	e.Value.FindUsedVars(options)
}

// NodeExpSin holds value
type NodeExpSin struct {
	Value NodeExp
}

// Type returns type
func (e *NodeExpSin) Type() int {
	return TypeFloat
}

// String returns value
func (e *NodeExpSin) String() string {
	return "SIN"
}

// Exp returns value
func (e *NodeExpSin) Exp(options *BuildOptions) string {
	return "math.Sin(" + forceFloat(options, e.Value) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpSin) FindUsedVars(options *BuildOptions) {
	e.Value.FindUsedVars(options)
}

// NodeExpSqr holds value
type NodeExpSqr struct {
	Value NodeExp
}

// Type returns type
func (e *NodeExpSqr) Type() int {
	return TypeFloat
}

// String returns value
func (e *NodeExpSqr) String() string {
	return "SQR"
}

// Exp returns value
func (e *NodeExpSqr) Exp(options *BuildOptions) string {
	return "math.Sqrt(" + forceFloat(options, e.Value) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpSqr) FindUsedVars(options *BuildOptions) {
	e.Value.FindUsedVars(options)
}

// NodeExpTan holds value
type NodeExpTan struct {
	Value NodeExp
}

// Type returns type
func (e *NodeExpTan) Type() int {
	return TypeFloat
}

// String returns value
func (e *NodeExpTan) String() string {
	return "TAN"
}

// Exp returns value
func (e *NodeExpTan) Exp(options *BuildOptions) string {
	return "math.Tan(" + forceFloat(options, e.Value) + ")"
}

// FindUsedVars finds used vars
func (e *NodeExpTan) FindUsedVars(options *BuildOptions) {
	e.Value.FindUsedVars(options)
}

// NodeExpFuncCall holds value
type NodeExpFuncCall struct {
	Name       string
	Parameters []NodeExp
}

// Type returns type
func (e *NodeExpFuncCall) Type() int {
	return VarType(e.Name)
}

// String returns value
func (e *NodeExpFuncCall) String() string {
	return "FN:" + e.Name + "()"
}

// Exp returns value
func (e *NodeExpFuncCall) Exp(options *BuildOptions) string {
	name := RenameFunc(e.Name)
	call := name + "("
	if len(e.Parameters) > 0 {
		call += e.Parameters[0].Exp(options)
		for i := 1; i < len(e.Parameters); i++ {
			call += "," + e.Parameters[i].Exp(options)
		}
	}
	call += ") /* <-- call DEF FN func */ "
	return call
}

// FindUsedVars finds used vars
func (e *NodeExpFuncCall) FindUsedVars(options *BuildOptions) {
	for _, p := range e.Parameters {
		p.FindUsedVars(options)
	}
}

// NodeExpGofunc holds value
type NodeExpGofunc struct {
	Name      *NodeExpString
	Arguments []NodeExp
}

// Type returns type
func (e *NodeExpGofunc) Type() int {
	return VarType(e.Name.Value)
}

// String returns value
func (e *NodeExpGofunc) String() string {
	return "GOFUNC(" + e.Name.Value + ")"
}

func RemoveSigil(s string) string {
	last := len(s) - 1
	if last < 0 {
		return s
	}
	switch s[last] {
	case '!', '$', '%', '#':
		return s[:last]
	}
	return s
}

// Exp returns value
func (e *NodeExpGofunc) Exp(options *BuildOptions) string {

	call := RemoveSigil(e.Name.Value) + "("
	if len(e.Arguments) > 0 {
		call += e.Arguments[0].Exp(options)
		for i := 1; i < len(e.Arguments); i++ {
			call += "," + e.Arguments[i].Exp(options)
		}
	}
	call += ") /* <-- GOFUNC */ "

	return call
}

// FindUsedVars finds used vars
func (e *NodeExpGofunc) FindUsedVars(options *BuildOptions) {
	for _, a := range e.Arguments {
		a.FindUsedVars(options)
	}
}

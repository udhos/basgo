package basparser

// misc.go extracted from header/footer of parser.y

import (
	"fmt"
	"log"
	"strings"

	"github.com/udhos/basgo/node"
)

type ParserResult struct {
	Root           []node.Node
	LineNumbers    map[string]node.LineNumber // used by GOTO GOSUB etc
	LibReadData    bool
	LibGosubReturn bool
	LibMath        bool
	Baslib         bool
	ForStack       []*node.NodeFor
	WhileStack     []*node.NodeWhile
	CountFor       int
	CountNext      int
	ArrayTable     map[string]node.ArraySymbol
	CountGosub     int
	CountReturn    int
	CountWhile     int
	CountWend      int
	CountIf        int
	FuncTable      map[string]node.FuncSymbol
	Imports        map[string]struct{}
	Declarations   []string
	RestoreTable   map[string]int
	DataOffset     int
	TypeTable      []int
}

// parser auxiliary variables
var (
	Result = newResult()

	nodeListStack [][]node.Node    // support nested node lists (1)
	expListStack  [][]node.NodeExp // support nested exp lists (2)
	lineList      []node.Node
	constList     []node.NodeExp
	varList       []node.NodeExp
	numberList    []string
	identList     []string
	lastLineNum   string // basic line number for parser error reporting
	rangeList     [][]string

	// (1) stmt IF-THEN can nest list of stmt: THEN CLS:IF:CLS
	// (2) exp can nest list of exp: array(exp,exp,exp)

	nodeExpNull = &node.NodeExpNull{}
)

func newResult() ParserResult {
	r := ParserResult{
		LineNumbers:  map[string]node.LineNumber{},
		ArrayTable:   map[string]node.ArraySymbol{},
		FuncTable:    map[string]node.FuncSymbol{},
		Imports:      map[string]struct{}{},
		RestoreTable: map[string]int{},
	}
	r.TypeTable = make([]int, 26, 26)
	defineType(&r, 0, 25, node.TypeFloat) // DEFSNG A-Z
	return r
}

func defineType(r *ParserResult, first, last, t int) {
	log.Printf("defineType: range %c-%c as %s", byte('a'+first), byte('a'+last), node.TypeLabel(t))
	for i := first; i <= last; i++ {
		r.TypeTable[i] = t
	}
}

func defineTypeRange(r *ParserResult, list [][]string, t int) {
	for _, p := range list {
		first := int(p[0][0] - 'a')
		last := int(p[1][0] - 'a')
		defineType(&Result, first, last, t)
	}
}

func Reset() {
	Result = newResult()

	nodeListStack = [][]node.Node{}
	expListStack = [][]node.NodeExp{}
}

func isSymbol(ident, symbol string) bool {
	return strings.ToLower(ident) == strings.ToLower(symbol)
}

func unsupportedEmpty(keyword string) *node.NodeEmpty {
	log.Printf("ignoring unsupported keyword %s", keyword)
	return &node.NodeEmpty{}
}

func createEndNode(result *ParserResult, msg string) *node.NodeEnd {
	result.Baslib = true
	return &node.NodeEnd{Message: msg}
}

func unsupportedEnd(result *ParserResult, keyword string) *node.NodeEnd {
	log.Printf("unsupported keyword %s will halt the program", keyword)
	msg := fmt.Sprintf("stopping on unsupported keyword %s", keyword)
	result.Imports["log"] = struct{}{} // NodeEnd.Message uses log
	return createEndNode(result, msg)
}

func captureRawLine(label string, list []node.Node, rawLine string) {
	last := len(list) - 1
	if last < 0 {
		log.Printf("captureRawLine: %s last line index=%d < 0", label, last)
		return
	}

	switch n := list[last].(type) {
	case *node.LineNumbered:
		n.RawLine = rawLine
		list[last] = n
		//log.Printf("captureRawLine: %s numbered index=%d raw=[%s]", label, last, n.RawLine)
	case *node.LineImmediate:
		n.RawLine = rawLine
		list[last] = n
		//log.Printf("captureRawLine: %s immediate index=%d raw=[%s]", label, last, n.RawLine)
	default:
		log.Printf("captureRawLine: %s non-line node: %v", label, list[last])
	}
}

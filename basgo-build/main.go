package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sort"
	//"strings"

	"github.com/udhos/basgo/basgo"
	"github.com/udhos/basgo/basparser"
	"github.com/udhos/basgo/node"
)

const (
	basgoLabel = "basgo-build"
)

func main() {
	log.Printf("%s version %s runtime %s GOMAXPROC=%d", basgoLabel, basgo.Version, runtime.Version(), runtime.GOMAXPROCS(0))

	status, errors := compile(os.Stdin, fmt.Printf)
	if status != 0 || errors != 0 {
		log.Printf("%s: status=%d errors=%d\n", basgoLabel, status, errors)
		os.Exit(1)
	}
}

func compile(input io.Reader, outputf node.FuncPrintf) (int, int) {

	log.Printf("%s: parsing", basgoLabel)

	result, status, errors := parse(input, outputf)
	lineNumbersTab := result.LineNumbers
	nodes := result.Root

	if status != 0 || errors != 0 {
		return status, errors
	}

	if result.CountFor != result.CountNext {
		log.Printf("%s: FOR count=%d NEXT count=%d", basgoLabel, result.CountFor, result.CountNext)
		return status, 1000
	}

	if result.CountWhile != result.CountWend {
		log.Printf("%s: WHILE count=%d WEND count=%d", basgoLabel, result.CountWhile, result.CountWend)
		return status, 2000
	}

	log.Printf("%s: FIXME WRITEME replace duplicate lines", basgoLabel)

	log.Printf("%s: checking lines used/defined", basgoLabel)

	var undefLines int

	for n, ln := range lineNumbersTab {
		//log.Printf("%s: line %s used=%v defined=%v", basgoLabel, n, ln.Used, ln.Defined)
		if ln.Used && !ln.Defined {
			undefLines++
			log.Printf("%s: line %s used but not defined", basgoLabel, n)
		}
	}

	if undefLines != 0 {
		return 0, 3000 + undefLines
	}

	log.Printf("%s: sorting lines", basgoLabel)

	sort.Sort(node.ByLineNumber(result.Root))

	header := `
package main

`
	mainOpen := `
func main() {
`
	mainClose := `
}
`

	options := node.BuildOptions{
		Headers:     map[string]struct{}{},
		Vars:        map[string]struct{}{},
		Arrays:      result.ArrayTable,
		LineNumbers: lineNumbersTab,
		CountGosub:  result.CountGosub,
		CountReturn: result.CountReturn,
	}

	if result.LibReadData {
		result.Baslib = true
	}

	if result.Baslib {
		options.Headers["github.com/udhos/basgo/baslib"] = struct{}{}
	}

	if result.LibMath {
		options.Headers["math"] = struct{}{}
	}

	log.Printf("%s: scanning used vars", basgoLabel)

	for _, n := range nodes {
		n.FindUsedVars(&options)
	}

	log.Printf("%s: issuing code", basgoLabel)

	buf := bytes.Buffer{}

	funcGen := func(format string, v ...interface{}) (int, error) {
		s := fmt.Sprintf(format, v...)
		return buf.WriteString(s)
	}

	for _, n := range nodes {
		n.Build(&options, funcGen)
	}

	outputf(header)

	// Copy headers from options into result imports
	for h, v := range options.Headers {
		result.Imports[h] = v
	}

	writeImport(result.Imports, outputf)

	if result.LibReadData {
		outputf("var readData = []interface{}{\n")
		for _, d := range options.ReadData {
			outputf("%s,\n", d)
		}
		outputf("}\n")
	}

	outputf(mainOpen)

	if result.LibGosubReturn {
		outputf("gosubStack := []int{} // used by GOSUB/RETURN lib\n")
	}

	writeVar(options.Vars, outputf)
	writeArrays(&options, outputf)
	declareFuncs(&options, result.FuncTable, outputf)

	outputf(buf.String())

	outputf(mainClose)

	// Declarations from _GODECL
	if len(result.Declarations) > 0 {
		outputf("\n// _GODECL declarations - begin\n")
		for _, d := range result.Declarations {
			outputf(d)
			outputf("\n")
		}
		outputf("// _GODECL declarations - end\n\n")
	}

	return status, errors
}

func writeImport(imp map[string]struct{}, outputf node.FuncPrintf) {
	if len(imp) < 1 {
		return
	}

	outputf("import (\n")
	for h := range imp {
		outputf("\"%s\"\n", h)
	}
	outputf(")\n")
}

func declareFuncs(options *node.BuildOptions, funcTable map[string]node.FuncSymbol, outputf node.FuncPrintf) {
	if len(funcTable) < 1 {
		return
	}

	outputf("var (\n")
	for n, symb := range funcTable {
		f := node.RenameFunc(n)
		funcType := node.FuncBuildType(options, n, symb.Func.Variables)
		outputf("  %s %s // DEF FN [%s] used=%v\n", f, funcType, n, symb.Used)
	}
	outputf(")\n")
}

func writeArrays(options *node.BuildOptions, outputf node.FuncPrintf) {
	if len(options.Arrays) < 1 {
		return
	}

	outputf("var (\n")
	for v, arr := range options.Arrays {
		a := node.RenameArray(v)
		arrayType := arr.ArrayType(v)
		outputf("  %s %s // array [%s]\n", a, arrayType, v)
	}
	outputf(")\n")
}

func writeVar(vars map[string]struct{}, outputf node.FuncPrintf) {
	if len(vars) < 1 {
		return
	}

	outputf("var (\n")
	for v := range vars {
		vv := node.RenameVar(v)
		t := node.VarType(v)
		tt := node.TypeName(v, t)
		outputf("  %s %s // var [%s]\n", vv, tt, v)
	}
	outputf(")\n")
}

func parse(input io.Reader, outputf node.FuncPrintf) (basparser.ParserResult, int, int) {
	debug := false
	byteInput := bufio.NewReader(input)
	lex := basparser.NewInputLex(byteInput, debug)
	basparser.Reset()
	status := basparser.InputParse(lex)

	return basparser.Result, status, lex.Errors()
}

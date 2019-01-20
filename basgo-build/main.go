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
	libInput := result.LibInput
	lineNumbersTab := result.LineNumbers
	nodes := result.Root

	if status != 0 || errors != 0 {
		return status, errors
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
		return 0, 1000 + undefLines
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
		LineNumbers: lineNumbersTab,
		Input:       libInput,
	}

	if options.Input {
		inputHeaders(options.Headers)
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
	writeImport(options.Headers, outputf)

	if options.Input {
		outputf("var stdin = bufio.NewReader(os.Stdin) // stdin used by INPUT lib\n")
	}

	if options.Rnd {
		outputf("var rnd *rand.Rand // used by RND lib\n")
		outputf("var rndLast float64 // used by RND lib\n")
	}

	outputf(mainOpen)

	if options.Rnd {
		outputf("rnd = rand.New(rand.NewSource(time.Now().UnixNano())) // used by RND lib\n")
		outputf("rndLast = rnd.Float64() // used by RND lib\n")
	}

	if result.CountNextAnon > 0 {
		// create variable for tracking last FOR for NEXT-novar
		outputf("var for_last int // last active FOR\n")
	}
	if result.CountNextIdent > 0 {
		// create variables for tracking last FOR for NEXT-var
		for _, f := range result.ForTable {
			v := node.RenameVar(f.Variable)
			outputf("var for_%s int // last FOR for NEXT %s\n", v, f.Variable)
		}
	}

	writeVar(options.Vars, outputf)
	outputf(buf.String())

	outputf(mainClose)

	lib(outputf, options.Input, options.Rnd, options.Left)

	return status, errors
}

func inputHeaders(h map[string]struct{}) {
	h["bufio"] = struct{}{}
	h["log"] = struct{}{}
	h["os"] = struct{}{}
	h["strconv"] = struct{}{}
	h["strings"] = struct{}{}
}

func lib(outputf node.FuncPrintf, input, rnd, left bool) {

	funcBoolToInt := `
func boolToInt(v bool) int {
	if v {
		return -1
	}
	return 0
}
`

	funcInputFmt := `
func %s string {

	buf, isPrefix, err := stdin.ReadLine()
	if err != nil {
		log.Printf("input error: %%%%v", err)
	}
	if isPrefix {
		log.Printf("input too big has been truncated")
	}

	return string(buf)
}

func %s int {
        str := strings.TrimSpace(inputString())
	v, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("input: integer '%%%%s' error: %%%%v", str, err)
	}
	return v 
}

func %s float64 {
        str := strings.TrimSpace(inputString())
	v, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Printf("input: float '%%%%s' error: %%%%v", str, err)
	}
	return v 
}
`
	outputf(funcBoolToInt)

	if input {
		funcInput := fmt.Sprintf(funcInputFmt, node.InputString, node.InputInteger, node.InputFloat)
		outputf(funcInput)
	}

	if rnd {
		funcRnd := `
func randomFloat64(v float64) float64 {
	if v > 0 {
		rndLast = rnd.Float64() // generate new number
	}
	return rndLast
}
`
		outputf(funcRnd)
	}

	if left {
		funcLeft := `
func stringLeft(s string, size int) string {
	if size < 1 {
		return ""
	}
	if size > len(s) {
		return s
	}
	return s[:size]
}
`
		outputf(funcLeft)
	}
}

func writeImport(headers map[string]struct{}, outputf node.FuncPrintf) {
	if len(headers) > 0 {
		outputf("import (\n")
		for h := range headers {
			outputf("\"%s\"\n", h)
		}
		outputf(")\n")
	}
}

func writeVar(vars map[string]struct{}, outputf node.FuncPrintf) {
	if len(vars) > 0 {
		outputf("var (\n")
		for v := range vars {
			vv := node.RenameVar(v)
			t := node.VarType(v)
			var tt string
			switch t {
			case node.TypeString:
				tt = "string"
			case node.TypeInteger:
				tt = "int"
			case node.TypeFloat:
				tt = "float64"
			default:
				log.Printf("writeVar: unknown var %s/%s type: %d", v, vv, t)
				tt = "TYPE_UNKNOWN_writeVar"
			}
			outputf("  %s %s // [%s]\n", vv, tt, v)
		}
		outputf(")\n")
	}
}

func parse(input io.Reader, outputf node.FuncPrintf) (basparser.ParserResult, int, int) {
	debug := false
	byteInput := bufio.NewReader(input)
	lex := basparser.NewInputLex(byteInput, debug)
	basparser.Reset()
	status := basparser.InputParse(lex)

	return basparser.Result, status, lex.Errors()
}

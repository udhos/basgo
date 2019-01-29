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
		Input:       libInput,
		CountGosub:  result.CountGosub,
		CountReturn: result.CountReturn,
	}

	if options.Input {
		inputHeaders(options.Headers)
	}

	if result.LibVal {
		options.Headers["log"] = struct{}{}
		options.Headers["strings"] = struct{}{}
		options.Headers["strconv"] = struct{}{}
	}

	if result.LibRepeat {
		options.Headers["log"] = struct{}{}
		options.Headers["strings"] = struct{}{}
	}

	if result.LibAsc {
		options.Headers["log"] = struct{}{}
	}

	if result.LibTime {
		options.Headers["time"] = struct{}{}
		options.Headers["fmt"] = struct{}{}
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

	if result.LibReadData {
		options.Headers["log"] = struct{}{}
	}

	writeImport(options.Headers, outputf)

	if result.LibReadData {
		outputf("var dataPos int // READ-DATA cursor\n")
		outputf("var data = []interface{}{\n")
		for _, d := range options.Data {
			outputf("%s,\n", d)
		}
		outputf("}\n")
	}

	if options.Input {
		outputf("var stdin = bufio.NewReader(os.Stdin) // stdin used by INPUT lib\n")
	}

	if options.Rnd {
		outputf("var rnd *rand.Rand // used by RND lib\n")
		outputf("var rndLast float64 // used by RND lib\n")
	}

	outputf(mainOpen)

	if result.LibGosubReturn {
		outputf("gosubStack := []int{} // used by GOSUB/RETURN lib\n")
	}

	if options.Rnd {
		outputf("rnd = rand.New(rand.NewSource(time.Now().UnixNano())) // used by RND lib\n")
		outputf("rndLast = rnd.Float64() // used by RND lib\n")
	}

	writeVar(options.Vars, outputf)
	writeArrays(&options, outputf)
	outputf(buf.String())

	outputf(mainClose)

	lib(outputf, options.Input, options.Rnd, options.Left, result.LibReadData, options.Mid, result.LibVal, result.LibRight, result.LibRepeat, result.LibAsc, result.LibBool, result.LibTime)

	return status, errors
}

func inputHeaders(h map[string]struct{}) {
	h["bufio"] = struct{}{}
	h["log"] = struct{}{}
	h["os"] = struct{}{}
	h["strconv"] = struct{}{}
	h["strings"] = struct{}{}
}

func lib(outputf node.FuncPrintf, input, rnd, left, libReadData, mid, val, right, repeat, asc, libBool, libTime bool) {

	if libTime {

		funcTime := `
func timeDate() string {
	y, m, d := time.Now().Date()
	return fmt.Sprintf("%%02d-%%02d-%%04d", m, d, y)
}
func timeTime() string {
	h, m, s := time.Now().Clock()
	return fmt.Sprintf("%%02d:%%02d:%%02d", h, m, s)
}
func timeTimer() float64 {
        now := time.Now()
        y, m, d := now.Date()
        midnight := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
        elapsed := now.Sub(midnight)
        return elapsed.Seconds()
}
`
		outputf(funcTime)
	}

	if libBool {

		funcBoolToInt := `
func boolToInt(v bool) int {
	if v {
		return -1
	}
	return 0
}
`
		outputf(funcBoolToInt)
	}

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
	if size >= len(s) {
		return s
	}
	return s[:size]
}
`
		outputf(funcLeft)
	}

	if right {
		funcRight := `
func stringRight(s string, size int) string {
	if size < 1 {
		return ""
	}
	if size >= len(s) {
		return s
	}
	return s[len(s)-size:]
}
`
		outputf(funcRight)
	}

	if repeat {
		funcRepeat := `
func stringRepeat(s string, count int) string {
	if count < 0 {
		log.Printf("repeat string negative count")
		count = 0
	}
	return strings.Repeat(s, count)
}
`
		outputf(funcRepeat)
	}

	if val {
		funcVal := `
func stringToFloat(s string) float64 {
        s = strings.TrimSpace(s)
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		m := "value for number: '" + s + "' error: " + err.Error()
		log.Printf(m)
	}
	return v 
}
`
		outputf(funcVal)
	}

	if asc {
		funcAsc := `
func firstByte(s string) int {
	if len(s) < 1 {
		log.Printf("ASC: firstByte: empty string")
		return 0
	}
	return int(s[0])
}
`
		outputf(funcAsc)
	}

	if mid {
		funcMid := `
func stringMid(s string, begin, size int) string {
	if size < 1 {
		return ""
	}
	begin--
	if begin >= len(s) {
		return ""
	}
	if begin < 0 {
		begin = 0
	} 
	avail := len(s) - begin
	if size > avail {
		size = avail
	}
	return s[begin:begin+size]
}
`
		outputf(funcMid)
	}

	if libReadData {
		funcData := `
func readDataString(name string) string {
	if dataPos >= len(data) {
		log.Fatalf("readDataString overflow error: var=%%s pos=%%d\n", name, dataPos)
	}
	d := data[dataPos]
	dataPos++
	v, t := d.(string)
	if t {
		return v
	}
	log.Fatalf("readDataString type error: var=%%s pos=%%d\n", name, dataPos)
	return v
}
func readDataInteger(name string) int {
	if dataPos >= len(data) {
		log.Fatalf("readDataInteger overflow error: var=%%s pos=%%d\n", name, dataPos)
	}
	d := data[dataPos]
	dataPos++
	v, t := d.(int)
	if t {
		return v
	}
	log.Fatalf("readDataInteger type error: var=%%s pos=%%d\n", name, dataPos)
	return v
}
func readDataFloat(name string) float64 {
	if dataPos >= len(data) {
		log.Fatalf("readDataFloat overflow error: var=%%s pos=%%d\n", name, dataPos)
	}
	d := data[dataPos]
	dataPos++
	v, t := d.(float64)
	if t {
		return v
	}
	v1, t1 := d.(int)
	if t1 {
		return float64(v1)
	}
	log.Fatalf("readDataFloat type error: var=%%s pos=%%d\n", name, dataPos)
	return v
}
`
		outputf(funcData)
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

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
	"strings"

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

	log.Printf("%s: parsing\n", basgoLabel)

	lineNumbersFound, nodes, status, errors := parse(input, outputf)

	if status != 0 || errors != 0 {
		return status, errors
	}

	log.Printf("%s: FIXME WRITEME replace duplicate lines\n", basgoLabel)

	log.Printf("%s: FIXME WRITEME check that used line numbers in lineNumbersFound exist as nodes\n", basgoLabel)

	log.Printf("%s: sorting lines\n", basgoLabel)

	sort.Sort(node.ByLineNumber(nodes))

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
		LineNumbers: lineNumbersFound,
	}

	log.Printf("%s: scanning used vars\n", basgoLabel)

	for _, n := range nodes {
		n.FindUsedVars(options.Vars)
	}

	log.Printf("%s: issuing code\n", basgoLabel)

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
	outputf(mainOpen)
	writeVar(options.Vars, outputf)
	outputf(buf.String())

	outputf(mainClose)

	funcBoolToInt := `
func boolToInt(v bool) int {
	if v {
		return -1
	}
	return 0
}
`
	outputf(funcBoolToInt)

	return status, errors
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
			var t string
			switch {
			case strings.HasSuffix(v, "$"):
				t = "string"
			case strings.HasSuffix(v, "%"):
				t = "int"
			default:
				t = "float64"
			}
			outputf("  %s %s // [%s]\n", node.RenameVar(v), t, v)
		}
		outputf(")\n")
	}
}

func parse(input io.Reader, outputf node.FuncPrintf) (map[string]struct{}, []node.Node, int, int) {
	debug := false
	byteInput := bufio.NewReader(input)
	lex := basparser.NewInputLex(byteInput, debug)
	status := basparser.InputParse(lex)

	return basparser.LineNumbers, basparser.Root, status, lex.Errors()
}

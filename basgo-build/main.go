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

	compile(os.Stdin, fmt.Printf)
}

func compile(input io.Reader, outputf node.FuncPrintf) {

	log.Printf("%s: parsing\n", basgoLabel)

	nodes, errors := parse(input, outputf)

	if errors != 0 {
		log.Printf("%s: syntax errors: %d\n", basgoLabel, errors)
		os.Exit(1)
	}

	log.Printf("%s: sorting lines\n", basgoLabel)

	sort.Sort(node.ByLineNumber(nodes))

	log.Printf("%s: issuing code\n", basgoLabel)

	header := `
package main
`
	mainOpen := `
func main() {
`
	mainClose := `
}
`

	log.Printf("%s: issuing code FIXME WRITEME generate runtime\n", basgoLabel)
	log.Printf("%s: issuing code FIXME WRITEME sort lines\n", basgoLabel)
	log.Printf("%s: issuing code FIXME WRITEME replace duplicate lines\n", basgoLabel)

	buf := bytes.Buffer{}
	options := node.BuildOptions{
		Headers: map[string]struct{}{},
		Vars:    map[string]struct{}{},
	}

	options.Headers["os"] = struct{}{}

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

	outputf("// below we use all vars to prevent Go compiler error\n")
	outputf("os.Exit(0)\n")
	for v := range options.Vars {
		outputf("println(%s) // [%s]\n", node.RenameVar(v), v)
	}

	outputf(mainClose)
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

func parse(input io.Reader, outputf node.FuncPrintf) ([]node.Node, int) {
	debug := false
	byteInput := bufio.NewReader(input)
	lex := basparser.NewInputLex(byteInput, debug)
	status := basparser.InputParse(lex)

	log.Printf("%s: parse status=%d", basgoLabel, status)

	return basparser.Root, lex.Errors()
}

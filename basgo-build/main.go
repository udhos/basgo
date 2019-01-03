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

	nodes, status := parse(input, outputf)

	if status != 0 {
		log.Printf("%s: syntax error\n", basgoLabel)
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
	options := node.BuildOptions{Headers: map[string]struct{}{}}

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
	outputf(buf.String())
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

func parse(input io.Reader, outputf node.FuncPrintf) ([]node.Node, int) {
	debug := false
	byteInput := bufio.NewReader(input)
	lex := basparser.NewInputLex(byteInput, debug)
	status := basparser.InputParse(lex)

	return basparser.Root, status
}

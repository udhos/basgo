package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

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

	log.Printf("%s: issuing code\n", basgoLabel)

	header := `package main

import (
        "fmt"
        "os"
)

`

	mainOpen := `func main() {
`

	mainClose := `}
`

	outputf(header)
	outputf(mainOpen)

	log.Printf("%s: issuing code FIXME WRITEME generate runtime\n", basgoLabel)
	log.Printf("%s: issuing code FIXME WRITEME sort lines\n", basgoLabel)
	log.Printf("%s: issuing code FIXME WRITEME replace duplicate lines\n", basgoLabel)

	for _, n := range nodes {
		n.Build(outputf)
	}

	outputf(mainClose)

}

func parse(input io.Reader, outputf node.FuncPrintf) ([]node.Node, int) {
	debug := false
	byteInput := bufio.NewReader(input)
	lex := basparser.NewInputLex(byteInput, debug)
	status := basparser.InputParse(lex)

	return basparser.Root, status
}

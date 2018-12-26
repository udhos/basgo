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

	debug := false
	byteInput := bufio.NewReader(input)
	lex := basparser.NewInputLex(byteInput, debug)
	status := basparser.InputParse(lex)

	if status != 0 {
		log.Printf("%s: syntax error\n", basgoLabel)
		os.Exit(1)
	}

	log.Printf("%s: building\n", basgoLabel)

	header := `package main

import (
        "fmt"
        "os"
)

func main() {
`

	footer := `
}
`

	outputf(header)

	for _, n := range basparser.Root {
		n.Build(outputf)
	}

	outputf(footer)

}

package basgo

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/udhos/basgo/basparser"
	"github.com/udhos/basgo/node"
)

// Version reports basgo version
const Version = "0.0"

// Basgo holds a full environment
type Basgo struct {
	In  io.Reader
	Out io.Writer
}

// New creates a new basgo environment
func New() *Basgo {
	return &Basgo{In: os.Stdin, Out: os.Stdout}
}

// REPL is read-evaluate-print-loop
func (b *Basgo) REPL() (errCount int) {
	r := bufio.NewReader(b.In)
	w := bufio.NewWriter(b.Out)

	printf := func(format string, v ...interface{}) (int, error) {
		s := fmt.Sprintf(format, v...)
		_, err := w.Write([]byte(s))
		if err != nil {
			log.Printf("REPL printf: %v", err)
		}
		return 0, nil
	}

	errCount = b.execReader(printf, r, w.Flush)
	return
}

type funcPrintf node.FuncPrintf // shortcut

func (b *Basgo) execReader(printf funcPrintf, r *bufio.Reader, flush func() error) (errorCount int) {
	for {
		s, errRead := r.ReadString('\n')
		if errRead != nil {
			log.Printf("REPL input: %v", errRead)
			break
		}
		line := strings.TrimSpace(s)
		if line == "" {
			continue
		}
		errLine := b.execLine(printf, r, line)
		if errLine != nil {
			errorCount++ // only count execLine() errors
		}
		flush()
	}
	return
}

func (b *Basgo) execLine(printf funcPrintf, r *bufio.Reader, line string) error {

	debug := false
	input := bufio.NewReader(strings.NewReader(line))
	lex := basparser.NewInputLex(input, debug)
	status := basparser.InputParse(lex)

	if status != 0 {
		err := fmt.Errorf("execLine: syntax error")
		printf("%v\n", err)
		return err
	}

	for _, n := range basparser.Root {
		scanStatements(printf, b, n, "")
	}

	printf("execLine: FIXME WRITEME insert numbered lines, execute immediate lines\n")

	return nil
}

func scanStatements(printf funcPrintf, b *Basgo, n node.Node, lineNum string) {
	switch nn := n.(type) {
	case *node.LineNumbered:
		for _, c := range nn.Nodes {
			scanStatements(printf, b, c, nn.LineNumber)
		}
	case *node.LineImmediate:
		for _, c := range nn.Nodes {
			scanStatements(printf, b, c, "")
		}
	default:
		printf("line [%s] statement [%s]\n", lineNum, nn.Name())
	}
}

func (b *Basgo) printf(format string, v ...interface{}) (int, error) {
	b.write(fmt.Sprintf(format, v...))
	return 0, nil
}

func (b *Basgo) write(s string) {
	_, err := b.Out.Write([]byte(s))
	if err != nil {
		log.Printf("write: %v", err)
	}
}

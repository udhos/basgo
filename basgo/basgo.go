package basgo

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/udhos/basgo/basparser"
)

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
func (b *Basgo) REPL() {
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

	b.execReader(printf, r, w.Flush)
}

type funcPrintf func(format string, v ...interface{}) (int, error)

/*
type hasReadString interface {
	ReadString(delim byte) (string, error)
}
*/

func (b *Basgo) execReader(printf funcPrintf, r *bufio.Reader, flush func() error) {
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
		b.execLine(printf, r, line)
		flush()
	}
}

func (b *Basgo) execLine(printf funcPrintf, r *bufio.Reader, line string) {

	debug := false
	input := bufio.NewReader(strings.NewReader(line))
	lex := basparser.NewInputLex(input, debug)
	status := basparser.InputParse(lex)

	if status != 0 {
		printf("execLine: syntax error\n")
		return
	}

	printf("execLine: begin\n")

	//basparser.Root.Run(b, printf, r)

	printf("execLine: end\n")
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

/*

// ExecuteLine executes a single line
func (b *Basgo) ExecuteLine(line string) {
	r := bufio.NewReader(strings.NewReader(line))
	b.execLine(b.printf, r, line)
}

// ExecuteString executes a multi-line string
func (b *Basgo) ExecuteString(s string) {
	r := bufio.NewReader(strings.NewReader(s))
	w := bufio.NewWriter(b.Out)
	b.execReader(b.printf, r, w.Flush)
}

*/

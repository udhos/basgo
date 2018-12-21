package basgo

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/udhos/basgo/baslex"
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

	printf := func(format string, v ...interface{}) {
		s := fmt.Sprintf(format, v...)
		_, err := w.Write([]byte(s))
		if err != nil {
			log.Printf("REPL printf: %v", err)
		}
	}

	b.execReader(printf, r, w.Flush)
}

type funcPrintf func(format string, v ...interface{})

type hasReadString interface {
	ReadString(delim byte) (string, error)
}

func (b *Basgo) execReader(printf funcPrintf, r hasReadString, flush func() error) {
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
		b.execLine(printf, line)
		flush()
	}
}

func (b *Basgo) execLine(printf funcPrintf, line string) {
	//printf("execLine: [%s]\n", line)
	lex := baslex.NewStr(line)
	for lex.HasToken() {
		tok := lex.Next()
		printf("execLine: [%s] token: %s [%s]\n", line, tok.Type(), tok.Value)
	}
}

func (b *Basgo) printf(format string, v ...interface{}) {
	b.write(fmt.Sprintf(format, v...))
}

func (b *Basgo) write(s string) {
	_, err := b.Out.Write([]byte(s))
	if err != nil {
		log.Printf("write: %v", err)
	}
}

// ExecuteLine executes a single line
func (b *Basgo) ExecuteLine(line string) {
	b.execLine(b.printf, line)
}

// ExecuteString executes a multi-line string
func (b *Basgo) ExecuteString(s string) {
	r := bufio.NewReader(strings.NewReader(s))
	w := bufio.NewWriter(b.Out)
	b.execReader(b.printf, r, w.Flush)
}

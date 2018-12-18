package basgo

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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

	for {
		s, errRead := r.ReadString('\n')
		if errRead != nil {
			log.Printf("REPL input: %v", errRead)
			break
		}
		line := strings.TrimSpace(s)
		log.Printf("REPL: [%s]", line)
		if line == "" {
			continue
		}
		b.executeLine(printf, line)
		w.Flush()
	}
}

type funcPrintf func(format string, v ...interface{})

func (b *Basgo) executeLine(printf funcPrintf, line string) {
	printf("executeLine: [%s]\n", line)
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

// ExecuteString executes a string in environment
func (b *Basgo) ExecuteString(s string) {
	b.printf("ExecuteString: FIXME WRITEME")
}

// ExecuteCommandList executes LIST command
func (b *Basgo) ExecuteCommandList() {
	b.printf("ExecuteCommandList: FIXME WRITEME")
}

// ExecuteCommandRun executes RUN command
func (b *Basgo) ExecuteCommandRun() {
	b.printf("ExecuteCommandRun: FIXME WRITEME")
}

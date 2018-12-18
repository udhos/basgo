package basgo

import (
	"fmt"
	"io"
	"log"
	"os"
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
func REPL() {
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

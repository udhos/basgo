package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
)

const (
	basgoVersion = "0.0"
	basgoLabel   = "basgo-run"
)

func main() {
	log.Printf("%s version %s runtime %s GOMAXPROC=%d", basgoLabel, basgoVersion, runtime.Version(), runtime.GOMAXPROCS(0))
}

// Basgo holds a full environment
type Basgo struct {
	In  io.Reader
	Out io.Writer
}

// NewBasgo creates a new basgo environment
func NewBasgo() *Basgo {
	return &Basgo{In: os.Stdin, Out: os.Stdout}
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

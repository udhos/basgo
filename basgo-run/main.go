package main

import (
	//"fmt"
	//"io"
	"log"
	//"os"
	"runtime"

	"github.com/udhos/basgo/basgo"
)

const (
	basgoLabel = "basgo-run"
)

func main() {
	log.Printf("%s version %s runtime %s GOMAXPROC=%d", basgoLabel, basgo.Version, runtime.Version(), runtime.GOMAXPROCS(0))

	b := basgo.New()

	errCount := b.REPL()

	log.Printf("%s: syntax errors: %d", basgoLabel, errCount)
}

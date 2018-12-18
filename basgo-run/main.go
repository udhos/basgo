package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github/udhos/basgo/basgo"
)

const (
	basgoVersion = "0.0"
	basgoLabel   = "basgo-run"
)

func main() {
	log.Printf("%s version %s runtime %s GOMAXPROC=%d", basgoLabel, basgoVersion, runtime.Version(), runtime.GOMAXPROCS(0))

	b := basgo.NewBasgo()

	b.REPL()
}

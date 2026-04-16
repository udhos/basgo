package main

import (
	//"fmt"
	//"io"
	"log"
	//"os"
	//"runtime"

	"github.com/udhos/basgo/basgo"
)

const (
	basgoLabel = "basgo-run"
)

func main() {
	basgo.ShowVersion(basgoLabel)

	b := basgo.New()

	errCount := b.REPL()

	log.Printf("%s: syntax errors: %d", basgoLabel, errCount)
}

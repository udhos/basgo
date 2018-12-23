package main

import (
	"bufio"
	//"fmt"
	"log"
	"os"

	"github.com/udhos/basgo/basparser"
)

func main() {
	me := os.Args[0]
	log.Printf("%s: reading from stdin...", me)

	input := bufio.NewReader(os.Stdin)
	lex := basparser.NewInputLex(input)
	status := basparser.InputParse(lex)

	log.Printf("%s: reading from stdin...done status=%d", me, status)
}

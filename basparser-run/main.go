package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/udhos/basgo/basparser"
)

func main() {
	me := os.Args[0]
	log.Printf("%s: reading from stdin...", me)

	debug := true
	input := bufio.NewReader(os.Stdin)
	lex := basparser.NewInputLex(input, debug)
	status := basparser.InputParse(lex)

	log.Printf("%s: reading from stdin...done status=%d", me, status)

	log.Printf("%s: syntax tree lines=%d:", me, len(basparser.Root))

	for i, n := range basparser.Root {
		fmt.Printf("%s: input line %d: ", me, i)
		n.Show(fmt.Printf)
	}
}

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
	basparser.Reset()
	status := basparser.InputParse(lex)
	nodes := basparser.Result.Root

	log.Printf("%s: reading from stdin...done", me)

	log.Printf("%s: status=%d errors=%d", me, status, lex.Errors())

	log.Printf("%s: syntax tree lines=%d:", me, len(nodes))

	for i, n := range nodes {
		fmt.Printf("%s: input line %d: ", me, i)
		n.Show(fmt.Printf)
	}
}

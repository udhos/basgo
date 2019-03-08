package main

import (
	//"bufio"
	"fmt"
	"log"
	"os"

	"github.com/udhos/basgo/basparser"
)

func main() {
	me := os.Args[0]
	log.Printf("%s: reading from stdin...", me)

	/*
		debug := true
		input := bufio.NewReader(os.Stdin)
		log.Printf("%s: input buffer size: %d", me, input.Size())
		lex := basparser.NewInputLex(input, debug)
		basparser.Reset()
		status := basparser.InputParse(lex)
		nodes := basparser.Result.Root
	*/

	result, status, errors := basparser.Run(me, os.Stdin)

	nodes := result.Root

	log.Printf("%s: reading from stdin...done", me)

	log.Printf("%s: status=%d errors=%d", me, status, errors)

	log.Printf("%s: FOR count=%d NEXT count=%d", me, basparser.Result.CountFor, basparser.Result.CountNext)
	log.Printf("%s: WHILE count=%d WEND count=%d", me, basparser.Result.CountWhile, basparser.Result.CountWend)

	log.Printf("%s: syntax tree lines=%d:", me, len(nodes))

	for i, n := range nodes {
		fmt.Printf("%s: input line %d: ", me, i+1)
		n.Show(fmt.Printf)
	}
}

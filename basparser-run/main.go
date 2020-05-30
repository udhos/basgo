package main

import (
	//"bufio"
	"fmt"
	"log"
	"os"

	"github.com/udhos/basgo/basgo"
	"github.com/udhos/basgo/basparser"
)

func main() {
	me := os.Args[0]
	basgo.ShowVersion(me)

	log.Printf("%s: reading BASIC code from stdin...", me)

	result, status, errors := basparser.Run(me, os.Stdin)

	nodes := result.Root

	log.Printf("%s: reading BASIC code from stdin...done", me)

	log.Printf("%s: status=%d errors=%d", me, status, errors)

	log.Printf("%s: FOR count=%d NEXT count=%d", me, result.CountFor, result.CountNext)
	log.Printf("%s: WHILE count=%d WEND count=%d", me, result.CountWhile, result.CountWend)

	log.Printf("%s: syntax tree lines=%d:", me, len(nodes))

	for i, n := range nodes {
		fmt.Printf("%s: input line %d: ", me, i+1)
		n.Show(fmt.Printf)
	}
}

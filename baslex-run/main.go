package main

import (
	//"bufio"
	"fmt"
	"log"
	"os"

	"github.com/udhos/basgo/baslex"
)

func main() {
	me := os.Args[0]

	input := baslex.NewInput(me, os.Stdin)

	log.Printf("%s: reading BASIC code from stdin...", me)

	lex := baslex.New(input, false)
	for lex.HasToken() {
		t := lex.Next()
		s := baslex.TokenString("", t, lex.Offset())
		fmt.Println(s)
	}

	log.Printf("%s: reading BASIC code from stdin...done", me)

	log.Printf("%s: stopped at line=%d column=%d offset=%d", me, lex.Line(), lex.Column(), lex.Offset())
}

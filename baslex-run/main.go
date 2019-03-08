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

	log.Printf("%s: reading input from stdin...", me)
	/*
		input := bufio.NewReader(os.Stdin)
		log.Printf("%s: input buffer size: %d", me, input.Size())
	*/
	input := baslex.NewInput(me, os.Stdin)

	lex := baslex.New(input, false)
	for lex.HasToken() {
		t := lex.Next()
		s := baslex.TokenString("", t, lex.Offset())
		fmt.Println(s)
	}

	log.Printf("%s: reading input from stdin...done", me)

	log.Printf("%s: stopped at line=%d column=%d offset=%d", me, lex.Line(), lex.Column(), lex.Offset())
}

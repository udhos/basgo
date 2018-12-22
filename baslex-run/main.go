package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/udhos/basgo/baslex"
)

func main() {
	me := os.Args[0]
	log.Printf("%s: reading input from stdin...", me)
	lex := baslex.New(bufio.NewReader(os.Stdin))
	for lex.HasToken() {
		tok := lex.Next()
		fmt.Printf("lineCount=%03d id=%02d %-s [%-s]\n", tok.LineCount, tok.ID, tok.Type(), tok.Value)
	}
	log.Printf("%s: reading input from stdin...done", me)
}

package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/udhos/basgo/basparser"
)

func main() {
	fi := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("input: ")

		if line, ok := readline(fi); ok {
			lex := basparser.NewInputLex(line)
			basparser.InputParse(lex)
		} else {
			break
		}
	}
}

func readline(fi *bufio.Reader) (string, bool) {
	s, err := fi.ReadString('\n')
	if err != nil {
		return "", false
	}
	return s, true
}

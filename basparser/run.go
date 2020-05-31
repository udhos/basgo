package basparser

import (
	//"bufio"
	"io"
	"log"
	"os"
	"strconv"

	//"github.com/udhos/basgo/node"
	"github.com/udhos/basgo/baslex"
)

func Run(label string, input io.Reader) (ParserResult, int, int) {
	d := os.Getenv("DEBUG")
	debug := d != ""
	if debug {
		level, err := strconv.Atoi(d)
		if err == nil {
			InputDebug = level
		}
	}
	log.Printf("%s: DEBUG=[%s] debug=%v level=%d", label, d, debug, InputDebug)

	inputBuf := baslex.NewInput(label, input)
	lex := NewInputLex(inputBuf, debug)
	Reset()
	status := InputParse(lex)

	return Result, status, lex.Errors()
}

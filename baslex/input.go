package baslex

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
)

func NewInput(label string, reader io.Reader) *bufio.Reader {
	var size int

	sz := os.Getenv("INPUTSZ")
	if sz != "" {
		v, err := strconv.Atoi(sz)
		if err == nil {
			size = v
		}
	}

	log.Printf("%s: INPUTSZ=[%s] size=%d", label, sz, size)

	var input *bufio.Reader

	if size > 0 {
		input = bufio.NewReaderSize(reader, size)
	} else {
		input = bufio.NewReader(reader)
	}

	log.Printf("%s: input buffer size: %d", label, input.Size())

	return input
}

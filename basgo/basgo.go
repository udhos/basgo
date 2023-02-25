package basgo

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/udhos/basgo/basparser"
	"github.com/udhos/basgo/node"
)

// Version reports basgo version
const Version = "0.12.0"

const (
	//DefaultBaslibModule = "github.com/udhos/baslib/baslib@master"
	//DefaultBaslibModule = "github.com/udhos/baslib/baslib"
	DefaultBaslibModule = "github.com/udhos/baslib@v0.11.0"
	DefaultBaslibImport = "github.com/udhos/baslib/baslib"
)

func ShowVersion(label string) {
	log.Printf("%s version %s runtime %s GOMAXPROC=%d", label, Version, runtime.Version(), runtime.GOMAXPROCS(0))
}

// Basgo holds a full environment
type Basgo struct {
	In  io.Reader
	Out io.Writer

	lines []lineEntry
}

type lineEntry struct {
	number   int
	raw      string
	commands []command
}

// New creates a new basgo environment
func New() *Basgo {
	return &Basgo{In: os.Stdin, Out: os.Stdout}
}

// REPL is read-evaluate-print-loop
func (b *Basgo) REPL() (errCount int) {
	r := bufio.NewReader(b.In)
	w := bufio.NewWriter(b.Out)

	printf := func(format string, v ...interface{}) (int, error) {
		s := fmt.Sprintf(format, v...)
		_, err := w.Write([]byte(s))
		if err != nil {
			log.Printf("REPL printf: %v", err)
		}
		return 0, nil
	}

	errCount = b.execReader(printf, r, w.Flush)
	return
}

type funcPrintf node.FuncPrintf // shortcut

func (b *Basgo) execReader(printf funcPrintf, r *bufio.Reader, flush func() error) (errorCount int) {
	for {
		s, errRead := r.ReadString('\n')
		if errRead != nil {
			log.Printf("REPL input: %v", errRead)
			break
		}
		line := strings.TrimSpace(s)
		if line == "" {
			continue
		}
		errLine := b.execLine(printf, r, line)
		if errLine != nil {
			errorCount++ // only count execLine() errors
		}
		flush()
	}
	return
}

func (b *Basgo) execLine(printf funcPrintf, r *bufio.Reader, line string) error {

	debug := false
	input := bufio.NewReader(strings.NewReader(line))
	lex := basparser.NewInputLex(input, debug)
	status := basparser.InputParse(lex)
	nodes := basparser.Result.Root

	if status != 0 {
		err := fmt.Errorf("execLine: parser error status: %d", status)
		printf("%v\n", err)
		return err
	}

	if errors := lex.Errors(); errors != 0 {
		err := fmt.Errorf("execLine: syntax errors: %d", errors)
		printf("%v\n", err)
		return err
	}

	for _, n := range nodes {
		b.scanSingleLine(printf, n, line)
	}

	return nil
}

func (b *Basgo) scanSingleLine(printf funcPrintf, n node.Node, rawLine string) {
	switch nn := n.(type) {
	case *node.LineNumbered:
		b.installLine(printf, nn.Nodes, rawLine, nn.LineNumber)
	case *node.LineImmediate:
		// execute line
		for _, stmt := range nn.Nodes {
			if b.execStatement(printf, stmt) {
				break // some commands like END might stop execution
			}
		}
	default:
		printf("unexpected non-line: %v\n", nn)
	}
}

func (b *Basgo) findLine(num int) (int, bool) {
	for i, line := range b.lines {
		if line.number > num {
			return i, false
		}
		if line.number == num {
			return i, true
		}
	}
	return -1, false
}

func (b *Basgo) installLine(printf funcPrintf, statements []node.Node, rawLine, lineNum string) {

	num, errAtoi := strconv.Atoi(lineNum)
	if errAtoi != nil {
		printf("line number error: %v\n", errAtoi)
		return
	}

	newLine := lineEntry{number: num, raw: rawLine}
	for _, n := range statements {
		cmd, errCmd := commandNew(n)
		if errCmd != nil {
			printf("line %d: command error: %v\n", num, errCmd)
			return
		}
		newLine.commands = append(newLine.commands, cmd)
	}

	index, found := b.findLine(num)

	// insert or delete line?

	if len(newLine.commands) == 1 {
		// single command
		if _, empty := newLine.commands[0].(*commandEmpty); empty {
			// single empty command: means delete line
			if found {
				b.lines = append(b.lines[:index], b.lines[index+1:]...)
			}
			return
		}
	}

	// replace?

	if found {
		b.lines[index] = newLine // found: replace existing index
		return
	}

	// !found = insert or append?

	if index >= 0 {
		// insert

		// insert https://github.com/golang/go/wiki/SliceTricks
		b.lines = append(b.lines, newLine) // just grow, element will be lost
		copy(b.lines[index+1:], b.lines[index:])
		b.lines[index] = newLine

		return
	}

	// append

	b.lines = append(b.lines, newLine)
}

func (b *Basgo) execStatement(printf funcPrintf, stmt node.Node) (stop bool) {

	cmd, errCmd := commandNew(stmt)
	if errCmd != nil {
		printf("command error: %v\n", errCmd)
		return
	}

	stop = cmd.exec(b, printf)

	return
}

func (b *Basgo) printf(format string, v ...interface{}) (int, error) {
	b.write(fmt.Sprintf(format, v...))
	return 0, nil
}

func (b *Basgo) write(s string) {
	_, err := b.Out.Write([]byte(s))
	if err != nil {
		log.Printf("write: %v", err)
	}
}

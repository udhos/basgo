package baslib

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/udhos/basgo/baslib/file"
)

type fileInfo struct {
	file   *os.File
	reader *bufio.Reader
	number int
	eof    bool
}

var fileTable = map[int]fileInfo{}

func Files(pattern string) {
	files, errFiles := filepath.Glob(pattern)
	if errFiles != nil {
		alert("FILES %s: error: %v", pattern, errFiles)
	}
	for _, f := range files {
		Println(f)
	}
}

func Eof(number int) int {
	return BoolToInt(hitEof(number))
}

func hitEof(number int) bool {
	i, found := fileTable[number]
	if !found {
		alert("EOF %d: file not open", number)
		return true
	}
	if i.eof {
		return true
	}
	if i.reader == nil {
		return true
	}
	return false
}

func isOpen(number int) bool {
	_, found := fileTable[number]
	return found
}

func Open(name string, number, mode int) {

	if isOpen(number) {
		alert("OPEN %d: file already open", number)
		return
	}

	var f *os.File
	var errOpen error

	switch mode {
	case file.OpenInput:
		f, errOpen = os.Open(name)
	case file.OpenOutput:
		f, errOpen = os.Create(name)
	default:
		alert("OPEN unsupported mode: %d", mode)
		return
	}

	if errOpen != nil {
		alert("OPEN error: %v", errOpen)
		return
	}

	i := fileInfo{
		file:   f,
		number: number,
	}

	switch mode {
	case file.OpenInput:
		i.reader = bufio.NewReader(f)
	}

	fileTable[number] = i
}

func Close(number int) {
	i, found := fileTable[number]
	if !found {
		alert("CLOSE %d: file not open", number)
		return
	}
	fileClose(i)
}

func fileClose(i fileInfo) {
	errClose := i.file.Close()
	if errClose != nil {
		alert("CLOSE %d: %v", i.number, errClose)
	}
	delete(fileTable, i.number)
}

func CloseAll() {
	for _, i := range fileTable {
		fileClose(i)
	}
}

func getReader(number int) *bufio.Reader {
	if hitEof(number) {
		return nil
	}
	i, _ := fileTable[number]
	return i.reader
}

func FileInputString(number int) string {
	return fileInputString(number)
}

func FileInputInteger(number int) int {
	s := fileInputString(number)
	if s == "" {
		return 0
	}
	return InputParseInteger(s)
}

func FileInputFloat(number int) float64 {
	s := fileInputString(number)
	if s == "" {
		return 0
	}
	return InputParseFloat(s)
}

func setEof(number int) {
	i, found := fileTable[number]
	if !found {
		alert("EOF on non-open file: %d", number)
		return
	}
	if i.eof {
		return // noop
	}
	i.eof = true
	fileTable[number] = i
}

func fileInputString(number int) string {
	reader := getReader(number)
	if reader == nil {
		return ""
	}
	buf, err := reader.ReadBytes('\n')
	switch err {
	case nil:
	case io.EOF:
		setEof(number)
	default:
		alert("INPUT# %d error: %v", number, err)
	}

	buf = bytes.TrimRight(buf, "\n")
	buf = bytes.TrimRight(buf, "\r")

	return string(buf)
}

func FileInputCount(count, number int) string {
	if count < 1 {
		alert("INPUT$ #%d bad length: %d", number, count)
		return ""
	}

	reader := getReader(number)
	if reader == nil {
		return ""
	}

	buf := make([]byte, count)

	n, err := reader.Read(buf)
	switch err {
	case nil:
	case io.EOF:
		setEof(number)
	default:
		alert("INPUT$ #%d error: %v", number, err)
	}

	if n != count {
		alert("INPUT$ #%d found=%d < request=%d", number, n, count)
	}

	return string(buf[:n])

}

func FilePrint(number int, value string) {
	i, found := fileTable[number]
	if !found {
		alert("PRINT# %d: file not open", number)
		return
	}
	_, err := i.file.WriteString(value)
	if err != nil {
		alert("PRINT# %d error: %v", number, err)
	}
}

func FilePrintInt(number, value int) {
	FilePrint(number, itoa(value))
}

func FilePrintFloat(number int, value float64) {
	FilePrint(number, ftoa(value))
}

func FileNewline(number int) {
	i, found := fileTable[number]
	if !found {
		alert("PRINT# %d: file not open", number)
		return
	}
	_, err := i.file.WriteString("\n")
	if err != nil {
		alert("PRINT# %d error: %v", number, err)
	}
}

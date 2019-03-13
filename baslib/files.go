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
	writer *bufio.Writer
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

func Lof(number int) int {
	i, found := fileTable[number]
	if !found {
		alert("LOF %d: file not open", number)
		return 0
	}
	info, err := i.file.Stat()
	if err != nil {
		alert("LOF %d: %v", number, err)
	}
	return int(info.Size())
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
		alert("EOF %d: file not open for input", number)
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
	case file.OpenAppend:
		f, errOpen = os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	default:
		alert("OPEN %d: unsupported mode: %d", number, mode)
		return
	}

	if errOpen != nil {
		alert("OPEN %d: %v", number, errOpen)
		return
	}

	i := fileInfo{
		file:   f,
		number: number,
	}

	switch mode {
	case file.OpenInput:
		i.reader = bufio.NewReader(f)
	case file.OpenOutput, file.OpenAppend:
		i.writer = bufio.NewWriter(f)
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
	if i.writer != nil {
		if errFlush := i.writer.Flush(); errFlush != nil {
			alert("CLOSE %d: flush: %v", i.number, errFlush)
		}
	}
	if errClose := i.file.Close(); errClose != nil {
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
	if i.writer == nil {
		alert("PRINT# %d: file not open for output", number)
		return
	}
	_, err := i.writer.WriteString(value)
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
	FilePrint(number, "\n")
}

func Kill(pattern string) {
	files, errFiles := filepath.Glob(pattern)
	if errFiles != nil {
		alert("KILL %s: %v", pattern, errFiles)
	}
	for _, f := range files {
		if errRem := os.Remove(f); errRem != nil {
			alert("KILL '%s': %s: %v", pattern, f, errRem)
		}
	}
}

func Name(from, to string) {
	if errRename := os.Rename(from, to); errRename != nil {
		alert("NAME '%s' AS '%s': %v", from, to, errRename)
	}
}

func Chdir(dir string) {
	if errChdir := os.Chdir(dir); errChdir != nil {
		alert("CHDIR '%s': %v", dir, errChdir)
	}
}

func Mkdir(dir string) {
	if errMkdir := os.Mkdir(dir, 0750); errMkdir != nil {
		alert("MKDIR '%s': %v", dir, errMkdir)
	}
}

func Rmdir(dir string) {
	if errRmdir := os.Remove(dir); errRmdir != nil {
		alert("RMDIR '%s': %v", dir, errRmdir)
	}
}

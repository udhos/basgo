package baslib

import (
	"os"
	"path/filepath"
)

const (
	OpenRandom = iota
	OpenInput  = iota
	OpenOutput = iota
	OpenAppend = iota
)

type fileInfo struct {
	file   *os.File
	number int
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
	case OpenInput:
		f, errOpen = os.Open(name)
	case OpenOutput:
		f, errOpen = os.Create(name)
	default:
		alert("OPEN unsupported mode: %d", mode)
		return
	}

	if errOpen != nil {
		alert("OPEN error: %v", errOpen)
		return
	}

	fileTable[number] = fileInfo{
		file:   f,
		number: number,
	}
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

package codepage

import (
	"bytes"
	"io"
	"log"
	"strconv"
	"strings"
)

var (
	tableCodepage    [256]int
	tableCodepageRev map[int]int
)

func init() {
	loadCodepage437()
}

func loadCodepage437() {
	loadCodepage("437", cp437) // Load from baslib/codepage/cp437.go

	// Overwrite some code points
	//
	// https://en.wikipedia.org/wiki/Code_page_437

	loadOne(1, 0x263a)
	loadOne(2, 0x263b)
	loadOne(3, 0x2665)
	loadOne(4, 0x2666)
	loadOne(5, 0x2663)
	loadOne(6, 0x2660)
	loadOne(7, 0x2022)
	loadOne(8, 0x25d8)
	loadOne(9, 0x25cb)
	loadOne(10, 0x25d9)
	loadOne(11, 0x2642)
	loadOne(12, 0x2640)
	loadOne(13, 0x266a)
	loadOne(14, 0x266b)
	loadOne(15, 0x263c)
	loadOne(16, 0x25ba)
	loadOne(17, 0x25c4)
	loadOne(18, 0x2195)
	loadOne(19, 0x203c)
	loadOne(20, 0x00b6)
	loadOne(21, 0x00a7)
	loadOne(22, 0x25ac)
	loadOne(23, 0x21a8)
	loadOne(24, 0x2191)
	loadOne(25, 0x2193)
	loadOne(26, 0x2192)
	loadOne(27, 0x2190)
	loadOne(28, 0x221f)
	loadOne(29, 0x2194)
	loadOne(30, 0x25b2)
	loadOne(31, 0x25bc)
}

func loadCodepage(label, s string) {
	log.Printf("loading codepage %s", label)

	buf := bytes.NewBufferString(s)

	tableCodepageRev = map[int]int{}

LOOP:
	for {
		line, errRead := buf.ReadString('\n')
		t := strings.TrimSpace(line)
		if t != "" && !strings.HasPrefix(t, "#") {
			f := strings.Fields(line)
			if len(f) > 1 {
				b, errParseByte := strconv.ParseInt(f[0], 0, 64)
				if errParseByte != nil {
					log.Printf("loading code page %s: %v: %s", label, errParseByte, line)
				}
				u, errParseUnicode := strconv.ParseInt(f[1], 0, 64)
				if errParseUnicode != nil {
					log.Printf("loading code page %s: %v: %s", label, errParseUnicode, line)
				}
				if errParseByte == nil && errParseUnicode == nil {
					loadOne(int(b), int(u))
				}
			} else {
				log.Printf("loading code page %s: line fields=%d underflow: [%s]", label, len(f), line)
			}
		}
		switch errRead {
		case nil:
		case io.EOF:
			break LOOP
		default:
			log.Printf("loading code page %s: %v", label, errRead)
		}
	}

	log.Printf("loading codepage %s: found %d symbols", label, len(tableCodepageRev))
}

func loadOne(b, u int) {
	uu := ByteToUnicode(b)
	if uu == u {
		log.Printf("codepage.loadOne(): dup byte=%d unicode=%x", b, u)
	}
	tableCodepage[b] = u
	tableCodepageRev[u] = b
}

func ByteToUnicode(b int) int {
	if b >= 0 && b < len(tableCodepage) {
		return tableCodepage[b]
	}
	return b
}

func UnicodeToByte(u int) int {
	if b, found := tableCodepageRev[u]; found {
		return b
	}
	return u
}

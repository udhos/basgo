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
	loadCodepage("437", cp437)
}

func loadCodepage(label, s string) {
	log.Printf("loading codepage %s", label)

	buf := bytes.NewBufferString(s)

	tableCodepageRev = map[int]int{}

LOOP:
	for {
		line, errRead := buf.ReadString('\n')
		if line != "" {
			t := strings.TrimSpace(line)
			if !strings.HasPrefix(t, "#") {
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
					log.Printf("loading code page %s: line fields=%d underflow: %s", label, len(f), line)
				}
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

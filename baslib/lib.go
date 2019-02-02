package baslib

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	stdin   = bufio.NewReader(os.Stdin)                       // INPUT
	rnd     = rand.New(rand.NewSource(time.Now().UnixNano())) // RND
	rndLast = rnd.Float64()                                   // RND
)

func BoolToInt(v bool) int {
	if v {
		return -1
	}
	return 0
}

func Sgn(v float64) int {
	switch {
	case v < 0:
		return -1
	case v > 0:
		return 1
	}
	return 0
}

func Date() string {
	y, m, d := time.Now().Date()
	return fmt.Sprintf("%02d-%02d-%04d", m, d, y)
}

func Time() string {
	h, m, s := time.Now().Clock()
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func Timer() float64 {
	now := time.Now()
	y, m, d := now.Date()
	midnight := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	elapsed := now.Sub(midnight)
	return elapsed.Seconds()
}

func Inkey() string {
	b, err := stdin.ReadByte()
	if err != nil {
		log.Printf("input byte error: %v", err)
		return ""
	}
	return string([]byte{b})
}

func inputString() string {

	buf, isPrefix, err := stdin.ReadLine()
	if err != nil {
		log.Printf("input error: %v", err)
	}
	if isPrefix {
		log.Printf("input too big has been truncated")
	}

	return string(buf)
}

func Input(prompt, question string, count int) []string {
	for {
		if prompt != "" {
			fmt.Print(prompt)
		}
		if question != "" {
			fmt.Print(question)
		}
		buf := inputString()
		fields := strings.Split(buf, ",")
		if n := len(fields); n != count {
			log.Printf("input: found %d of %d required comma-separated fields, please retry.", n, count)
			continue
		}
		return fields
	}
}

func InputParseInteger(str string) int {
	v, err := strconv.Atoi(strings.TrimSpace(str))
	if err != nil {
		log.Printf("input: integer '%s' error: %v", str, err)
	}
	return v
}

func InputParseFloat(str string) float64 {
	v, err := strconv.ParseFloat(strings.TrimSpace(str), 64)
	if err != nil {
		log.Printf("input: float '%s' error: %v", str, err)
	}
	return v
}

func Left(s string, size int) string {
	if size < 1 {
		return ""
	}
	if size >= len(s) {
		return s
	}
	return s[:size]
}

func MidSize(s string, begin, size int) string {
	if size < 1 {
		return ""
	}
	begin--
	if begin >= len(s) {
		return ""
	}
	if begin < 0 {
		begin = 0
	}
	avail := len(s) - begin
	if size > avail {
		size = avail
	}
	return s[begin : begin+size]
}

func Mid(s string, begin int) string {
	begin--
	if begin >= len(s) {
		return ""
	}
	if begin < 0 {
		begin = 0
	}
	return s[begin:]
}

func String(s string, count int) string {
	if count < 0 {
		log.Printf("repeat string negative count")
		count = 0
	}
	return strings.Repeat(Left(s, 1), count)
}

func Right(s string, size int) string {
	if size < 1 {
		return ""
	}
	if size >= len(s) {
		return s
	}
	return s[len(s)-size:]
}

func Rnd(v float64) float64 {
	if v > 0 {
		rndLast = rnd.Float64() // generate new number
	}
	return rndLast
}

package baslib

import (
	//"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/udhos/inkey/inkey"
)

const baslibVersion = "0.6.0"

type keyInput interface {
	Inkey() (byte, bool)
	Read(buf []byte) (int, error)
	ReadBytes(delim byte) (line []byte, err error)
}

var (
	stdin       = newInkey()                                      // INPUT
	rnd         = rand.New(rand.NewSource(time.Now().UnixNano())) // RND
	rndLast     = rnd.Float64()                                   // RND
	readDataPos int                                               // READ-DATA cursor
	screenPos   = 1                                               // PRINT COLUMN
	screenRow   = 1                                               // PRINT ROW
	showAlert   = true
)

func Begin() {
	envAlert := "BASLIB_ALERT_OFF"
	envAlertValue := os.Getenv(envAlert)
	if envAlertValue != "" {
		showAlert = false
	}
	if showAlert {
		log.Printf("baslib: version %s runtime %s GOMAXPROC=%d", baslibVersion, runtime.Version(), runtime.GOMAXPROCS(0))
		log.Printf("baslib: %s=%s showAlert=%v", envAlert, envAlertValue, showAlert)
		log.Printf("baslib: env var %s is empty, set it to non-empty to disable alerts", envAlert)
	}
}

func newInkey() keyInput {
	alert("newInkey(): will consume os.Stdin")
	return inkey.New(os.Stdin)
}

func alert(format string, v ...interface{}) {
	if !showAlert {
		return
	}
	s := "BASLIB ALERT: " + fmt.Sprintf(format, v...)
	if screenMode0() {
		Println(s)
		return
	}
	log.Print(s)
}

func fatal(format string, v ...interface{}) {
	s := "BASLIB FATAL: " + fmt.Sprintf(format, v...)
	if screenMode0() {
		Println(s)
	} else {
		log.Print(s)
	}
	os.Exit(1)
}

func Asc(s string) int {
	if len(s) < 1 {
		alert("asc empty string")
		return 0
	}
	return int(s[0])
}

func Val(s string) float64 {
	s = strings.TrimSpace(s)
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		alert("numeric value from: '%s' error: %v", s, err)
	}
	return v
}

func Int(v float64) int {
	return int(math.Floor(v))
}

func Fix(v float64) int {
	return int(math.Trunc(v))
}

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
	b, found := stdin.Inkey()
	if !found {
		return ""
	}
	return string(b)
}

func inputString() string {

	screenCursor(true) // enable

	buf, err := stdin.ReadBytes('\n')
	if err != nil && err != io.EOF {
		alert("input error: %v", err)
	}

	screenCursor(false) // disable

	buf = bytes.TrimRight(buf, "\n")
	buf = bytes.TrimRight(buf, "\r")

	// remove backspaces
	buf2 := make([]byte, 0, len(buf))
	for _, b := range buf {
		switch b {
		case 8, 127:
			if len(buf2) > 0 {
				buf2 = buf2[:len(buf2)-1]
			}
		default:
			buf2 = append(buf2, b)
		}
	}

	s := string(buf2)

	return s
}

func InputCount(count int) string {

	//log.Printf("InputCount(%d)", count)

	if count < 1 {
		return ""
	}

	buf := make([]byte, count)

	n, err := stdin.Read(buf)
	if err != nil {
		alert("InputCount: error: %v", err)
	}

	//log.Printf("InputCount(%d): %d", count, n)

	return string(buf[:n])
}

func Input(prompt, question string, count int) []string {
	for {
		if prompt != "" {
			Print(prompt)
		}
		if question != "" {
			Print(question)
		}
		buf := inputString()
		fields := strings.Split(buf, ",")
		if n := len(fields); n != count {
			alert("input: found %d of %d required comma-separated fields, please retry.", n, count)
			continue
		}
		return fields
	}
}

func InputParseInteger(str string) int {
	v, err := strconv.Atoi(strings.TrimSpace(str))
	if err != nil {
		alert("input: integer '%s' error: %v", str, err)
	}
	return v
}

func InputParseFloat(str string) float64 {
	v, err := strconv.ParseFloat(strings.TrimSpace(str), 64)
	if err != nil {
		alert("input: float '%s' error: %v", str, err)
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
		alert("string repeat negative count")
	}
	if count < 1 {
		return ""
	}
	if count == 1 {
		return s
	}
	if len(s) > 1 {
		s = s[:1]
	}
	return strings.Repeat(s, count)
}

func ReadDataString(data []interface{}, name string) string {
	if readDataPos >= len(data) {
		fatal("ReadDataString overflow error: var=%s pos=%d", name, readDataPos)
		return ""
	}
	d := data[readDataPos]
	readDataPos++
	v, t := d.(string)
	if t {
		return v
	}
	fatal("ReadDataString type error: var=%s pos=%d", name, readDataPos)
	return v
}

func ReadDataInteger(data []interface{}, name string) int {
	if readDataPos >= len(data) {
		fatal("ReadDataInteger overflow error: var=%s pos=%d", name, readDataPos)
		return 0
	}
	d := data[readDataPos]
	readDataPos++
	v, t := d.(int)
	if t {
		return v
	}
	fatal("ReadDataInteger type error: var=%s pos=%d", name, readDataPos)
	return v
}

func ReadDataFloat(data []interface{}, name string) float64 {
	if readDataPos >= len(data) {
		fatal("ReadDataFloat overflow error: var=%s pos=%d", name, readDataPos)
		return 0
	}
	d := data[readDataPos]
	readDataPos++
	v, t := d.(float64)
	if t {
		return v
	}
	v1, t1 := d.(int)
	if t1 {
		return float64(v1)
	}
	fatal("ReadDataFloat type error: var=%s pos=%d\n", name, readDataPos)
	return v
}

func Restore(data []interface{}, line string, pos int) {
	if pos < 0 {
		// warn only, actual fault hit in READ
		fatal("Restore underflow error: line=%s pos=%d", line, pos)
	}
	if pos >= len(data) {
		// warn only, actual fault hit in READ
		fatal("Restore overflow error: line=%s pos=%d", line, pos)
	}
	readDataPos = pos
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

func RandomizeAuto() {
	rnd.Seed(time.Now().UnixNano())
}

func Randomize(seed float64) {
	rnd.Seed(int64(seed))
}

func Rnd(v float64) float64 {
	if v < 0 {
		Randomize(v)
	}
	if v > 0 {
		rndLast = rnd.Float64() // generate new number
	}
	return rndLast
}

func StrInt(v int) string {
	return " " + itoa(v)
}

func itoa(v int) string {
	return strconv.Itoa(v)
}

func StrFloat(v float64) string {
	return " " + ftoa(v)
}

func ftoa(v float64) string {
	//return strconv.FormatFloat(v, 'f', -1, 64)
	return fmt.Sprint(v)
}

func Instr(begin int, str string, sub string) int {
	if begin > len(str) {
		return 0
	}
	begin--
	if begin < 0 {
		begin = 0
	}
	i := strings.Index(str[begin:], sub)
	if i < 0 {
		return 0
	}
	return i + begin + 1
}

func Pos() int {
	return screenPos
}

func Csrlin() int {
	return screenRow
}

func PrintInt(i int) {
	Print(" ")
	Print(itoa(i))
	Print(" ")
}

func PrintFloat(f float64) {
	Print(" ")
	Print(ftoa(f))
	Print(" ")
}

func Print(s string) {
	for _, b := range s {
		switch b {
		case 13: // CR
			cr()
		default:
			printItem(b)
		}
	}
	if screenMode0() {
		screenShow()
	}
}

func printItem(b rune) {

	if screenMode0() {
		// running in terminal

		screenPut(b)
		screenPos++
		if screenPos > screenWidth {
			screenCR()
		}
		return
	}

	// running in console

	fmt.Print(string(b))
	screenPos++
}

func screenCR() {
	screenPos = 1
	screenRow++
	if screenRow > screenHeight {
		screenScroll()
		screenRow = screenHeight
	}
}

func cr() {

	if screenMode0() {
		// terminal
		screenCR()
		return
	}

	// console
	printItem('\n')
	screenPos = 1
}

func Println(s string) {
	Print(s)
	cr()
	if screenMode0() {
		screenShow()
	}
}

func Tab(col int) string {
	if col < 1 {
		col = 1
	}
	if col == screenPos {
		return ""
	}
	if col < screenPos {
		return string(13) + String(" ", col-1)
	}
	return String(" ", col-screenPos)
}

func Hex(v int) string {
	return fmt.Sprintf("%X", v)
}

func Oct(v int) string {
	return strconv.FormatInt(int64(v), 8)
}

func LenStr(v string) int {
	return len(v)
}

func LenInt(v int) int {
	return int(unsafe.Sizeof(v))
}

func LenFloat(v float64) int {
	return 8
}

func Chr(v int) string {
	return string(v)
}

func Peek(addr int) int {
	return 0 // PEEK unsupported
}

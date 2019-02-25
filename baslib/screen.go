package baslib

import (
	"fmt"
	"io"
	//"log"
	"unicode/utf8"

	"github.com/gdamore/tcell"
	"github.com/udhos/inkey/inkey"
)

var (
	scr              screen
	screenWidth      = 80
	screenHeight     = 25
	screenViewTop    = 1
	screenCursorShow bool
)

func End() {
	//alert("baslib.End()")
	scr.close()
}

func screenMode0() bool {
	return scr.s != nil
}

func Screen(mode int) {
	if mode != 0 {
		alert("SCREEN %d: only screen 0 is supported", mode)
		return
	}

	if screenMode0() {
		alert("Screen: text mode 0 is running already")
		return
	}

	scr.start()

	stdin = inkey.New(&scr) // replace inkey(os.Stdin) with inkey(tcell)
}

func Cls() {
	if screenMode0() {
		if screenViewTop == 1 && screenHeight > 24 {
			scr.s.Clear() // clear terminal
		} else {
			cls() // clear view print window
		}
	}
	screenPos = 1
	screenRow = 1
}

func cls() {
	lastRow := screenLastRow()

	for row := screenViewTop; row <= lastRow; row++ {
		for col := 0; col < screenWidth; col++ {
			scr.s.SetContent(col, row-1, ' ', nil, 0)
		}
	}

	screenShow()
}

func Locate(row, col int) {
	if row > 0 && row <= screenHeight {
		screenRow = row
	}
	if col > 0 && col <= screenWidth {
		screenPos = col
	}
}

func LocateCursor(row, col int, cursor bool) {
	Locate(row, col)
	alert("LOCATE FIXME: handle cursor enable=%v", cursor)
}

func Width(w int) {
	if w < 1 || w > 1000 {
		alert("WIDTH value out-of-range: %d", w)
		return
	}
	screenWidth = w
}

func ViewPrint(top, bottom int) {
	if top < 1 || top > 1000 {
		alert("VIEW PRINT top line out-of-range: %d", top)
		return
	}
	if bottom < 1 || bottom > 1000 {
		alert("VIEW PRINT bottom line out-of-range: %d", bottom)
		return
	}
	if top > bottom {
		alert("VIEW PRINT top line must not be greater than bottom line")
		return
	}
	screenViewTop = top
	screenHeight = bottom - top + 1
}

func ViewPrintReset() {
	screenViewTop = 1
	screenHeight = 25
}

func screenLastRow() int {
	return screenViewTop + screenHeight - 1
}

type screen struct {
	s    tcell.Screen
	keys chan tcell.EventKey
}

func locateAlert(s string) {
			x := screenPos
			 y:= screenRow
			Locate(15,30)
			alert(itoa(x) + " " + itoa(y))
			Locate(y,x)
}

func (s *screen) Read(buf []byte) (int, error) {
	for {
		if screenCursorShow {
			s.s.ShowCursor(screenPos-1, screenRow-1)

			locateAlert(itoa(screenRow) + " " + itoa(screenPos))
		} else {
			s.s.HideCursor()
		}
		key, ok := <-s.keys
		if !ok {
			return 0, io.EOF
		}
		kType := key.Key()
		switch kType {
		case tcell.KeyBackspace:
			r := key.Rune()
			need := utf8.RuneLen(r)
			avail := len(buf)
			if need > avail {
				return 0, fmt.Errorf("screen.Read: rune short buffer: need=%d avail=%d", need, avail)
			}
			size := utf8.EncodeRune(buf, r)
			locateAlert("backspace")
			if screenCursorShow {
				Locate(screenRow, screenPos-1)
				screenPut(' ')
			}
			return size, nil
		case tcell.KeyEnter:
			need := 1
			avail := len(buf)
			if need > avail {
				return 0, fmt.Errorf("screen.Read: enter short buffer: need=%d avail=%d", need, avail)
			}
			buf[0] = '\n'
			if screenCursorShow {
				screenCR()
			}
			return 1, nil
		case tcell.KeyRune:
			r := key.Rune()
			need := utf8.RuneLen(r)
			avail := len(buf)
			if need > avail {
				return 0, fmt.Errorf("screen.Read: rune short buffer: need=%d avail=%d", need, avail)
			}
			size := utf8.EncodeRune(buf, r)
			if screenCursorShow {
				printItem(r)
			}
			return size, nil
		}
	}
}

func (s *screen) start() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)

	sNew, errScreen := tcell.NewScreen()
	if errScreen != nil {
		alert("tcell create screen: %v", errScreen)
		return
	}

	if errInit := sNew.Init(); errInit != nil {
		alert("tcell init screen: %v", errInit)
		sNew.Fini()
		return
	}

	s.s = sNew

	s.s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))

	s.s.Clear()

	scr.keys = make(chan tcell.EventKey)

	go screenEvents()
}

func (s *screen) close() {
	if s.s != nil {
		s.s.Fini()
	}
}

func screenEvents() {
	for {
		ev := scr.s.PollEvent()
		switch ev := ev.(type) {
		case nil: // PollEvent() will return nil if the Screen is finalized
			close(scr.keys)
			return
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyCtrlL:
				scr.s.Sync()
			}
			scr.keys <- *ev
		case *tcell.EventResize:
			scr.s.Sync()
		}
	}
}

func screenPut(r rune) {
	scr.s.SetContent(screenPos-1, screenViewTop+screenRow-2, r, nil, 0)
}

func screenScroll() {
	lastRow := screenLastRow() - 1

	// shift rows up
	for row := screenViewTop - 1; row < lastRow; row++ {
		for col := 0; col < screenWidth; col++ {
			mainc, combc, style, _ := scr.s.GetContent(col, row+1)
			scr.s.SetContent(col, row, mainc, combc, style)
		}
	}

	// clear last line
	for col := 0; col < screenWidth; col++ {
		scr.s.SetContent(col, lastRow, ' ', nil, 0)
	}
}

func screenShow() {
	scr.s.Show()
}

func screenCursor(enable bool) {
	screenCursorShow = enable
	if screenMode0() {
		if screenCursorShow {
			scr.s.ShowCursor(screenPos-1, screenRow-1)
		} else {
			scr.s.HideCursor()
		}
	}
}

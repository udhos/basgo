package baslib

import (
	"fmt"
	"io"
	//"log"
	"unicode/utf8"

	"github.com/gdamore/tcell"
	"github.com/udhos/inkey/inkey"

	"github.com/udhos/basgo/baslib/codepage"
)

var (
	scr                   screen
	screenWidth           = 80
	screenHeight          = 25
	screenViewTop         = 1
	screenCursorShow      bool
	screenColorForeground = tcell.ColorWhite
	screenColorBackground = tcell.ColorBlack
	screenStyle           tcell.Style
)

func End() {
	//alert("baslib.End()")
	CloseAll()
	scr.close()
}

func screenMode0() bool {
	return scr.s != nil
}

func Color(fg, bg int) {
	if fg >= 0 {
		fg = colorTerm(fg)
		screenColorForeground = tcell.Color(fg)
	}

	if bg >= 0 {
		bg = colorTerm(bg)
		screenColorBackground = tcell.Color(bg)
	}

	if screenModeGraphics() {
		graphicsColorUpload()
	}

	screenStyle = tcell.StyleDefault.Foreground(screenColorForeground).Background(screenColorBackground)
}

func ScreenFunc(row, col int, colorFlag bool) int {
	if !screenMode0() {
		return 0
	}

	mainc, _, style, _ := scr.s.GetContent(col-1, row-1)

	if colorFlag {
		fg, bg, _ := style.Decompose() // fg,bg,attr

		f := colorBas(int(fg))
		b := colorBas(int(bg))

		//locateAlert(16, 20, "SCREEN(): fg="+itoa(f)+" bg="+itoa(b))

		return b*16 + f
	}

	b := codepage.UnicodeToByte(int(mainc))

	return b
}

func Screen(mode int) {

	if mode == 0 {

		if screenMode0() {
			alert("Screen: text mode 0 is running already")
			return
		}

		graphicsStop()

		// start SCREEN 0
		scr.start()
		stdin = inkey.New(&scr) // replace inkey(os.Stdin) with inkey(tcell)

		return
	}

	if screenMode0() {
		scr.close() // stop SCREEN 0
	}

	graphicsStart(mode)
}

func Cls() {
	switch {
	case screenMode0():
		if screenViewTop == 1 && screenHeight > 24 {
			scr.s.Fill(' ', screenStyle) // clear terminal
		} else {
			cls() // clear view print window
		}
		screenShow()
	case screenModeGraphics():
		graphicsCls()
	}
	screenPos = 1
	screenRow = 1
}

func cls() {
	lastRow := screenLastRow()

	for row := screenViewTop - 1; row < lastRow; row++ {
		for col := 0; col < screenWidth; col++ {
			scr.s.SetContent(col, row, ' ', nil, screenStyle)
		}
	}
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
	s       tcell.Screen
	keys    chan tcell.EventKey
	bufSize int
}

func locateAlert(row, col int, s string) {
	x := screenPos
	y := screenRow
	Locate(row, col)
	alert(s)
	Locate(y, x)
	screenShow()
}

func (s *screen) Read(buf []byte) (int, error) {
	for {
		if screenCursorShow {
			s.s.ShowCursor(screenPos-1, screenRow-1)

			//locateAlert(15, 20, itoa(screenRow)+" "+itoa(screenPos))
		} else {
			s.s.HideCursor()
		}
		screenShow()
		key, ok := <-s.keys
		if !ok {
			return 0, io.EOF
		}
		f := func() {
			//locateAlert(20, 10, itoa(s.bufSize)+" "+itoa(screenRow)+","+itoa(screenPos)+"           ")
		}
		kType := key.Key()
		switch kType {
		case tcell.KeyBackspace, tcell.KeyDEL:
			r := key.Rune()
			need := utf8.RuneLen(r)
			avail := len(buf)
			if need > avail {
				return 0, fmt.Errorf("screen.Read: rune short buffer: need=%d avail=%d", need, avail)
			}
			size := utf8.EncodeRune(buf, r)
			//locateAlert(16, 20, "backspace="+itoa(int(r))+"        ")
			if screenCursorShow {
				if s.bufSize > 0 {
					if screenPos == 1 {
						if screenRow > 1 {
							// wrap line up
							screenRow--
							screenPos = 80
						}
					} else {
						screenPos-- // backspace
					}
					Locate(screenRow, screenPos)
					screenPut(' ')
					s.bufSize--
				}
			}
			f()
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
				s.bufSize = 0
			}
			f()
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
				s.bufSize++
			}
			f()
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

	Color(7, 0)

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
	u := codepage.ByteToUnicode(int(r))
	scr.s.SetContent(screenPos-1, screenViewTop+screenRow-2, rune(u), nil, screenStyle)
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
		scr.s.SetContent(col, lastRow, ' ', nil, screenStyle)
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
		screenShow()
	}
}

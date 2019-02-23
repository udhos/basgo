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
	scr          screen
	screenWidth  = 80
	screenHeight = 25
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
		scr.s.Clear()
	}
	screenPos = 1
	screenRow = 1
}

type screen struct {
	s    tcell.Screen
	keys chan tcell.EventKey
}

func (s *screen) Read(buf []byte) (int, error) {
	for {
		key, ok := <-s.keys
		if !ok {
			return 0, io.EOF
		}
		kType := key.Key()
		switch kType {
		case tcell.KeyEnter:
			//Println(fmt.Sprintf("[enter]"))
			need := 1
			avail := len(buf)
			if need > avail {
				return 0, fmt.Errorf("screen.Read: enter short buffer: need=%d avail=%d", need, avail)
			}
			buf[0] = '\n'
			return 1, nil
		case tcell.KeyRune:
			r := key.Rune()
			//Println(fmt.Sprintf("[%v]", r))
			need := utf8.RuneLen(r)
			avail := len(buf)
			if need > avail {
				return 0, fmt.Errorf("screen.Read: rune short buffer: need=%d avail=%d", need, avail)
			}
			size := utf8.EncodeRune(buf, r)
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
	scr.s.SetContent(screenPos-1, screenRow-1, r, nil, 0)
}

func screenScroll() {
	// shift rows up
	for row := 0; row < screenHeight; row++ {
		for col := 0; col < screenWidth; col++ {
			mainc, combc, style, _ := scr.s.GetContent(col, row+1)
			scr.s.SetContent(col, row, mainc, combc, style)
		}
	}

	// clear last line
	for col := 0; col < screenWidth; col++ {
		scr.s.SetContent(col, screenHeight-1, ' ', nil, 0)
	}
}

func screenShow() {
	scr.s.Show()
}

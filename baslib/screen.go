package baslib

import (
	//"fmt"
	"io"
	"log"
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
	log.Printf("baslib.End()")
	scr.close()
}

func screenMode0() bool {
	return scr.s != nil
}

func Screen(mode int) {
	if mode != 0 {
		log.Printf("SCREEN %d: only screen 0 is supported", mode)
		return
	}

	if screenMode0() {
		log.Printf("Screen: text mode 0 is running already")
		return
	}

	scr.start()

	stdin = inkey.New(&scr) // replace inkey(os.Stdin) with inkey(tcell)
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
			if 1 > (cap(buf) - len(buf)) {
				return 0, io.ErrShortBuffer
			}
			buf = append(buf, '\n')
			return 1, nil
		case tcell.KeyRune:
			r := key.Rune()
			need := utf8.RuneLen(r)
			if need > (cap(buf) - len(buf)) {
				return 0, io.ErrShortBuffer
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
		log.Printf("tcell create screen: %v", errScreen)
		return
	}

	if errInit := sNew.Init(); errInit != nil {
		log.Printf("tcell init screen: %v", errInit)
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

	//log.Printf("tcell screen initialized")
}

func (s *screen) close() {
	if s.s != nil {
		s.s.Fini()
	}
	//log.Printf("tcell screen finalized")
}

func screenEvents() {
	for {
		ev := scr.s.PollEvent()
		//log.Printf("screenEvents: %v", ev)
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
			//log.Printf("tcell screen resized")
		default:
			//log.Printf("tcell unhandled event")
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
			mainc, combc, style, _ := scr.s.GetContent(col, row)
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

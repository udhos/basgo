package baslib

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell"
)

var (
	scr screen
)

func End() {
	log.Printf("baslib.End()")
	scr.close()
}

func Screen(mode int) {
	if mode != 0 {
		log.Printf("SCREEN %d: only screen 0 is supported", mode)
		return
	}

	if scr.s != nil {
		log.Printf("Screen: text mode 0 is running already")
		return
	}

	scr.start()

	stdin = &scr
}

type screen struct {
	s    tcell.Screen
	keys chan tcell.EventKey
}

func (s *screen) Inkey() (byte, bool) {
	log.Printf("screen.Inkey: FIXME WRITEME")
	return 0, false
}

func (s *screen) Read(buf []byte) (int, error) {
	return 0, fmt.Errorf("screen.Read: FIXME WRITEME")
}

func (s *screen) ReadBytes(delim byte) (line []byte, err error) {
	return nil, fmt.Errorf("screen.Read: FIXME WRITEME")
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
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))

	s.s.Clear()

	scr.keys = make(chan tcell.EventKey)

	go screenEvents()

	stdin = s

	log.Printf("tcell screen initialized")
}

func (s *screen) close() {
	if s.s != nil {
		s.s.Fini()
	}
	log.Printf("tcell screen finalized")
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
			log.Printf("tcell screen resized")
		default:
			log.Printf("tcell unhandled event")
		}
	}
}

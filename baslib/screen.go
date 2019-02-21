package baslib

import (
	"log"

	"github.com/gdamore/tcell"
)

var (
	screen tcell.Screen
)

func End() {
	log.Printf("baslib.End()")
	if screen != nil {
		log.Printf("tcell screen finalized")
		screen.Fini()
	}
}

func Screen(mode int) {
	if mode != 0 {
		log.Printf("SCREEN %d: only screen 0 is supported", mode)
		return
	}

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)

	var errScreen error

	screen, errScreen = tcell.NewScreen()
	if errScreen != nil {
		log.Printf("tcell create screen: %v", errScreen)
		return
	}
	if errScreen = screen.Init(); errScreen != nil {
		log.Printf("tcell init screen: %v", errScreen)
		return
	}

	screen.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))

	screen.Clear()

	go screenEvents()

	log.Printf("tcell screen initialized")
}

func screenEvents() {
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyCtrlL:
				screen.Sync()
			}
		case *tcell.EventResize:
			screen.Sync()
			log.Printf("tcell screen resize")
		}
	}
}

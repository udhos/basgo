package baslib

import (
	"log"
)

func Screen(mode int) {
	if mode != 0 {
		log.Printf("SCREEN %d: only screen 0 is supported", mode)
		return
	}
}

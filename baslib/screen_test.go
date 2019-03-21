package baslib

import (
	"testing"
)

func TestScreen(t *testing.T) {
	//Screen(99) // Screen(0) would disrupt the terminal
	End()
}

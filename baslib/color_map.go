package baslib

import (
	//"os"
	//"strings"
	"github.com/gdamore/tcell"
)

var colorBasToTerm = [16]int{
	int(tcell.NewRGBColor(0, 0, 0)),          // 0 black
	int(tcell.NewRGBColor(0, 0, 0xa8)),       // 1 blue
	int(tcell.NewRGBColor(0, 0xa8, 0)),       // 2 green
	int(tcell.NewRGBColor(0, 0xa8, 0xa8)),    // 3 cyan
	int(tcell.NewRGBColor(0xa8, 0, 0)),       // 4 red
	int(tcell.NewRGBColor(0xa8, 0, 0xa8)),    // 5 magenta
	int(tcell.NewRGBColor(0xa8, 0x54, 0)),    // 6 brown (dark yellow)
	int(tcell.NewRGBColor(0xa8, 0xa8, 0xa8)), // 7 gray (dark white)
	int(tcell.NewRGBColor(0x54, 0x54, 0x54)), // 8 dark gray (bright black)
	int(tcell.NewRGBColor(0x54, 0x54, 0xfc)), // 9 bright blue
	int(tcell.NewRGBColor(0x54, 0xfc, 0x54)), // 10 bright green
	int(tcell.NewRGBColor(0x54, 0xfc, 0xfc)), // 11 bright cyan
	int(tcell.NewRGBColor(0xfc, 0x54, 0x54)), // 12 bright red
	int(tcell.NewRGBColor(0xfc, 0x54, 0xfc)), // 13 bright magenta (pink)
	int(tcell.NewRGBColor(0xfc, 0xfc, 0x54)), // 14 yellow
	int(tcell.NewRGBColor(0xfc, 0xfc, 0xfc)), // 15 white
}

var colorTermToBas map[int]int

func init() {
	colorTermToBas = map[int]int{}
	for i, v := range colorBasToTerm {
		colorTermToBas[v] = i
	}
}

func colorTerm(c int) int {
	if c >= 0 && c < len(colorBasToTerm) {
		return colorBasToTerm[c]
	}
	return c
}

func colorBas(c int) int {
	if b, found := colorTermToBas[c]; found {
		return b
	}
	return c
}

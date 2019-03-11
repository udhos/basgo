package baslib

import (
	//"os"
	//"strings"
	"github.com/gdamore/tcell"
)

var colorBasToTerm = [16]int{
	int(tcell.ColorBlack),       // 0
	int(tcell.ColorDarkBlue),    // 1
	int(tcell.ColorDarkGreen),   // 2
	int(tcell.ColorDarkCyan),    // 3
	int(tcell.ColorRed),         // 4
	int(tcell.ColorDarkMagenta), // 5
	int(tcell.ColorSaddleBrown), // 6
	int(tcell.ColorLightGray),   // 7
	int(tcell.ColorGray),        // 8
	int(tcell.ColorBlue),        // 9
	int(tcell.ColorLightGreen),  // 10
	int(tcell.ColorLightCyan),   // 11
	int(tcell.ColorTomato),      // 12
	int(tcell.ColorPink),        // 13
	int(tcell.ColorYellow),      // 14
	int(tcell.ColorWhite),       // 15
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

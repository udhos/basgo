package main

import (
	"github.com/faiface/mainthread"
	"github.com/udhos/baslib/baslib"
	//"github.com/go-gl/glfw/v3.2/glfw"
	"math"
	"os"
)

func main() {
	mainthread.Run(run)
}

func run() {
	mainthread.Call(func() {
		baslib.G = baslib.InitWin(640, 480)
	})
	var (
		sng_i float64 // var [i]
	)
	// line 100
	// empty node ignored: [EMPTY-NODE]
	// [CLS]
	baslib.Cls()
	// [SCREEN 9]
	baslib.Screen(9)
	// line 200
	// [COLOR 7 5]
	baslib.Color(7, 5)
	// REM: 'rem next CLS will clear bg color to 5'
	// line 205
	// [LINE <0> <0> <639> <479> <> <>]
	baslib.Line(0, 0, 639, 479, -1, -1)
	// line 206
	// [LINE <50> <0> <0> <50> <> <>]
	baslib.Line(50, 0, 0, 50, -1, -1)
	// line 210
	// REM: 'rem _goproc("sleep1")'
	// line 215
	// [PRINT <INPUT$(1)> NEWLINE]
	baslib.Print(baslib.InputCount(1))
	baslib.Println(``) // PRINT newline not suppressed
	// line 220
	// [CLS]
	baslib.Cls()
	// line 230
	// [COLOR 2]
	baslib.Color(2, -1)
	// [FOR i = 50 TO 300 STEP 1 Index=0]
	sng_i = float64(50) // FOR 0 initialization
for_loop_0:
	if (1) >= 0 { // FOR step non-negative?
		if sng_i > float64(300) {
			goto for_exit_0
		}
	} else {
		if sng_i < float64(300) {
			goto for_exit_0
		}
	}
	// [LINE <100> <50> <319> <i> <> <>]
	baslib.Line(100, 50, 319, int(math.Round(sng_i) /* <- forceInt(non-int) */), -1, -1)
	// [NEXT vars_size=0 fors_size=1]
	sng_i += float64(1) // FOR 0 step
	goto for_loop_0
for_exit_0:
	// line 240
	// [COLOR 4]
	baslib.Color(4, -1)
	// line 250
	// [LINE <10> <100> <40> <130> <1> <> box]
	baslib.LineBox(10, 100, 40, 130, 1, -1, false)
	// line 260
	// [LINE <15> <105> <35> <125> <> <> box fill]
	baslib.LineBox(15, 105, 35, 125, -1, -1, true)
	// line 270
	// [LINE <80> <130> <50> <100> <1> <> box]
	baslib.LineBox(80, 130, 50, 100, 1, -1, false)
	// line 280
	// [LINE <75> <125> <55> <105> <> <> box fill]
	baslib.LineBox(75, 125, 55, 105, -1, -1, true)
	// line 290
	// [LINE <40> <140> <10> <170> <1> <> box]
	baslib.LineBox(40, 140, 10, 170, 1, -1, false)
	// line 300
	// [LINE <15> <165> <35> <145> <> <> box fill]
	baslib.LineBox(15, 165, 35, 145, -1, -1, true)
	// line 310
	// [LINE <50> <170> <80> <140> <1> <> box]
	baslib.LineBox(50, 170, 80, 140, 1, -1, false)
	// line 320
	// [LINE <55> <165> <75> <145> <> <> box fill]
	baslib.LineBox(55, 165, 75, 145, -1, -1, true)
	// line 900
	// REM: 'rem _goimport("time")'
	// line 910
	// REM: 'rem _godecl("func sleep1() { time.Sleep(1*time.Second) }")'
	// line 920
	// REM: 'rem _godecl("func sleep3() { time.Sleep(3*time.Second) }")'
	// line 930
	// REM: 'rem _goproc("sleep3")'
	// line 935
	// [PRINT <INPUT$(1)> NEWLINE]
	baslib.Print(baslib.InputCount(1))
	baslib.Println(``) // PRINT newline not suppressed
	// line 940
	baslib.End()
	os.Exit(0) // END
	// unnumbered line ignored: ''
	baslib.End()

}

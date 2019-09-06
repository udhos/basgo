package baslib

import (
	"fmt"
	"io"
	"log"
	//"runtime"
	"unicode/utf8"

	"github.com/faiface/mainthread"
	"github.com/gdamore/tcell"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/udhos/inkey/inkey"
)

var G *glfw.Window

type graph struct {
	mode        int
	window      *glfw.Window
	program     uint32
	width       int
	height      int
	geom        []float32
	u_color     int32
	keys        chan int
	bufferPoint uint32
}

var graphics graph

func geomBuf(size int) {
	grow := size - len(graphics.geom)
	if grow < 1 {
		return
	}
	graphics.geom = append(graphics.geom, make([]float32, grow)...)
}

func screenModeGraphics() bool {
	return graphics.mode != 0
}

const swap = false

func swapOne() {
	if swap {
		graphics.window.SwapBuffers()
	}
}

func swapTwo() {
	if swap {
		graphics.window.SwapBuffers()
	} else {
		gl.Flush()
	}
}

func graphicsCls() {
	mainthread.Call(func() {
		swapOne()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		swapTwo()
	})
}

func graphicsStop() {
	graphics.mode = 0
	log.Printf("baslib graphicsStop()")
	//runtime.UnlockOSThread()
}

func graphicsStart(mode int) {
	log.Printf("baslib graphicsStart(%d)", mode)

	//runtime.LockOSThread()

	graphics.mode = mode

	graphics.width = 640
	graphics.height = 480

	log.Printf("graphicsStart(%d): %d x %d", mode, graphics.width, graphics.height)

	//graphics.window = initWin(graphics.width, graphics.height)
	graphics.window = G
	if graphics.window == nil {
		log.Printf("graphicsStart(%d): failed", mode)
		return
	}

	mainthread.Call(func() {
		if err := gl.Init(); err != nil {
			log.Printf("OpenGL init: %v", err)
			return
		}
		version := gl.GoStr(gl.GetString(gl.VERSION))
		log.Println("OpenGL version", version)
	})

	mainthread.Call(func() {
		prog, errProg := initProg()
		if errProg != nil {
			log.Printf("OpenGL program: %v", errProg)
			return
		}
		log.Printf("OpenGL program: %d", prog)
		graphics.program = prog
	})

	mainthread.Call(func() {
		graphics.u_color = getUniformLocation("u_color")
	})

	geomBuf(18)

	mainthread.Call(func() {
		gl.UseProgram(graphics.program)
	})

	graphics.mode = mode

	mainthread.Call(func() {
		gl.ClearDepthf(1)         // default
		gl.Disable(gl.DEPTH_TEST) // disable depth testing
		gl.Disable(gl.CULL_FACE)  // disable face culling

		gl.CreateBuffers(1, &graphics.bufferPoint) // buffer for point
	})

	graphics.keys = make(chan int, 10)

	mainthread.Call(func() {
		graphics.window.SetKeyCallback(keyCallback)
	})

	graphicsColorUpload()

	stdin = inkey.New(&graphics) // replace inkey(os.Stdin) with inkey(graph)

	log.Printf("baslib graphicsStart(%d) done", mode)
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Release {
		return // ignore key release
	}

	shift := mods&glfw.ModShift != 0
	capslock := mods&0x0010 != 0 // GLFW 3.3
	upper := shift || capslock

	k := int(key)

	if !upper && key >= glfw.KeyA && key <= glfw.KeyZ {
		k += 32 // lower case
	}

	log.Printf("keyCallback: key: '%c' '%c' %d shift=%v capslock=%v", byte(key), byte(k), k, shift, capslock)

	graphics.keys <- k
}

func remapKey(k int) int {
	switch k {
	case 257:
		log.Printf("remapKey: ENTER")
		return '\n' // enter
	case 258:
		log.Printf("remapKey: TAB")
		return 9 // tab
	case 259:
		log.Printf("remapKey: BACKSPACE")
		return 8 // backspace
	}
	return k
}

func (g *graph) Read(buf []byte) (int, error) {
	for {
		select {
		case k, ok := <-g.keys:
			if !ok {
				log.Printf("graph.Read: EOF")
				return 0, io.EOF
			}
			log.Printf("graph.Read: key: %d", k)

			kk := remapKey(k)

			switch {
			case kk != k:
				need := 1
				avail := len(buf)
				if need > avail {
					return 0, fmt.Errorf("graph.Read: remap short buffer: need=%d avail=%d", need, avail)
				}
				buf[0] = byte(kk)
				return 1, nil
			case k < 256:
				need := 1
				avail := len(buf)
				if need > avail {
					return 0, fmt.Errorf("graph.Read: short buffer: need=%d avail=%d", need, avail)
				}
				buf[0] = byte(k)
				return 1, nil
			default:
				r := rune(k)
				need := utf8.RuneLen(r)
				avail := len(buf)
				if need > avail {
					return 0, fmt.Errorf("graph.Read: rune short buffer: need=%d avail=%d", need, avail)
				}
				size := utf8.EncodeRune(buf, r)
				return size, nil
			}
		default:
			mainthread.Call(func() {
				glfw.PollEvents()
			})
		}

	}
}

func getUniformLocation(name string) int32 {
	u := gl.GetUniformLocation(graphics.program, gl.Str(name+"\x00"))
	if u < 0 {
		log.Printf("getUniformLocation: uniform location not found for: %s", name)
	}
	return u
}

// upload foreground color
func graphicsColorFg(fg tcell.Color) {
	r, g, b := fg.RGB()
	rr, gg, bb := rgbFloat(r, g, b)
	mainthread.Call(func() {
		gl.Uniform4f(graphics.u_color, rr, gg, bb, 1)
	})
}

func graphicsColorUpload() {
	graphicsColorFg(screenColorForeground)

	// upload background color
	r, g, b := screenColorBackground.RGB()
	rr, gg, bb := rgbFloat(r, g, b)
	mainthread.Call(func() {
		gl.ClearColor(rr, gg, bb, 1)
	})
}

func rgbFloat(r, g, b int32) (float32, float32, float32) {
	return float32(r) / 255, float32(g) / 255, float32(b) / 255
}

func draw(mode, vao uint32, window *glfw.Window, count int32) {
	mainthread.Call(func() {
		swapOne()

		gl.BindVertexArray(vao)
		gl.DrawArrays(mode, 0, count)

		swapTwo()
	})
}

func drawPoint() {
	mainthread.Call(func() {
		swapOne()

		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

		gl.BindBuffer(gl.ARRAY_BUFFER, graphics.bufferPoint)
		gl.BufferData(gl.ARRAY_BUFFER, 4*3*1, gl.Ptr(graphics.geom), gl.DYNAMIC_DRAW)
		gl.DrawArrays(gl.POINTS, 0, 1)

		swapTwo()
	})
}

// 0..639 to -1..1
func pix2Clip(x, w int) float32 {
	if w == 1 {
		return 0 // ugh
	}

	w--
	x *= 2                       // 0..w-1 -> 0..2(w-1)
	x -= w                       // -(w-1)..w-1
	c := float32(x) / float32(w) // -1..1

	return c
}

func pixelToClip(x, y int) (float32, float32) {
	a, b := pix2Clip(x, graphics.width), -pix2Clip(y, graphics.height)
	//log.Printf("pixelToClip: %d x %d => %f x %f", x, y, a, b)
	return a, b
}

const pointVao = false

func PSet(x, y, color int) {
	a, b := pixelToClip(x, y)

	graphics.geom[0] = a
	graphics.geom[1] = b
	graphics.geom[2] = 0 // clear

	if color >= 0 {
		// draw with specified color
		graphicsColorFg(tcell.Color(colorTerm(color)))
	}

	if pointVao {
		vao := makeVao(graphics.geom, 3)
		vaoIndices := int32(1)
		draw(gl.POINTS, vao, graphics.window, vaoIndices)
	} else {
		drawPoint()
	}

	if color >= 0 {
		// restore color
		graphicsColorFg(screenColorForeground)
	}
}

func PReset(x, y, color int) {
	a, b := pixelToClip(x, y)

	graphics.geom[0] = a
	graphics.geom[1] = b
	graphics.geom[2] = 0 // clear

	if color >= 0 {
		// draw with specified color
		graphicsColorFg(tcell.Color(colorTerm(color)))
	} else {
		// draw with background color
		graphicsColorFg(screenColorBackground)
	}

	if pointVao {
		vao := makeVao(graphics.geom, 3)
		vaoIndices := int32(1)
		draw(gl.POINTS, vao, graphics.window, vaoIndices)
	} else {
		drawPoint()
	}

	// restore color
	graphicsColorFg(screenColorForeground)
}

func Line(x1, y1, x2, y2, color, style int) {

	a1, b1 := pixelToClip(x1, y1)
	a2, b2 := pixelToClip(x2, y2)

	graphics.geom[0] = a1
	graphics.geom[1] = b1
	graphics.geom[2] = 0 // clear
	graphics.geom[3] = a2
	graphics.geom[4] = b2
	graphics.geom[5] = 0 // clear

	vao := makeVao(graphics.geom, 6)
	vaoIndices := int32(2)

	if color >= 0 {
		// draw with specified color
		graphicsColorFg(tcell.Color(colorTerm(color)))
	}

	draw(gl.LINES, vao, graphics.window, vaoIndices)

	if color >= 0 {
		// restore color
		graphicsColorFg(screenColorForeground)
	}
}

func LineBox(x1, y1, x2, y2, color, style int, fill bool) {

	a1, b1 := pixelToClip(x1, y1)
	a2, b2 := pixelToClip(x2, y2)

	var mode uint32
	var vao uint32
	var vaoIndices int32

	if fill {
		graphics.geom[0] = a1
		graphics.geom[1] = b1
		graphics.geom[2] = 0 // clear
		graphics.geom[3] = a1
		graphics.geom[4] = b2
		graphics.geom[5] = 0 // clear
		graphics.geom[6] = a2
		graphics.geom[7] = b2
		graphics.geom[8] = 0 // clear

		graphics.geom[9] = a2
		graphics.geom[10] = b2
		graphics.geom[11] = 0 // clear
		graphics.geom[12] = a2
		graphics.geom[13] = b1
		graphics.geom[14] = 0 // clear
		graphics.geom[15] = a1
		graphics.geom[16] = b1
		graphics.geom[17] = 0 // clear

		mode = gl.TRIANGLES
		vao = makeVao(graphics.geom, 18)
		vaoIndices = 6
	} else {
		graphics.geom[0] = a1
		graphics.geom[1] = b1
		graphics.geom[2] = 0 // clear
		graphics.geom[3] = a1
		graphics.geom[4] = b2
		graphics.geom[5] = 0 // clear
		graphics.geom[6] = a2
		graphics.geom[7] = b2
		graphics.geom[8] = 0 // clear
		graphics.geom[9] = a2
		graphics.geom[10] = b1
		graphics.geom[11] = 0 // clear

		mode = gl.LINE_LOOP
		vao = makeVao(graphics.geom, 12)
		vaoIndices = 4
	}

	if color >= 0 {
		// draw with specified color
		graphicsColorFg(tcell.Color(colorTerm(color)))
	}

	draw(mode, vao, graphics.window, vaoIndices)

	if color >= 0 {
		// restore color
		graphicsColorFg(screenColorForeground)
	}
}

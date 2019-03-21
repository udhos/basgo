package baslib

import (
	"fmt"
	"io"
	"log"
	"runtime"
	"unicode/utf8"

	"github.com/gdamore/tcell"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/udhos/inkey/inkey"
)

type graph struct {
	mode    int
	window  *glfw.Window
	program uint32
	width   int
	height  int
	geom    []float32
	u_color int32
	keys    chan rune
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

func graphicsCls() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func graphicsStop() {
	graphics.mode = 0
	log.Printf("baslib graphicsStop()")
	runtime.UnlockOSThread()
}

func graphicsStart(mode int) {
	log.Printf("baslib graphicsStart(%d)", mode)

	runtime.LockOSThread()

	graphics.mode = mode

	graphics.width = 640
	graphics.height = 480

	log.Printf("graphicsStart(%d): %d x %d", mode, graphics.width, graphics.height)

	graphics.window = initWin(graphics.width, graphics.height)
	if graphics.window == nil {
		log.Printf("graphicsStart(%d): failed", mode)
		return
	}

	if err := gl.Init(); err != nil {
		log.Printf("OpenGL init: %v", err)
		return
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	prog, errProg := initProg()
	if errProg != nil {
		log.Printf("OpenGL program: %v", errProg)
		return
	}

	log.Printf("OpenGL program: %d", prog)

	graphics.program = prog

	graphics.u_color = getUniformLocation("u_color")

	geomBuf(18)

	gl.UseProgram(graphics.program)

	graphics.mode = mode

	gl.ClearDepthf(1)         // default
	gl.Disable(gl.DEPTH_TEST) // disable depth testing
	gl.Disable(gl.CULL_FACE)  // disable face culling

	graphics.keys = make(chan rune, 10)

	graphics.window.SetCharCallback(charCallback)

	graphicsColorUpload()

	stdin = inkey.New(&graphics) // replace inkey(os.Stdin) with inkey(graph)

	//drawTriangle()

	log.Printf("baslib graphicsStart(%d) done", mode)
}

func charCallback(w *glfw.Window, r rune) {
	log.Printf("charCallback: key: %d", r)
	graphics.keys <- r
}

func (g *graph) Read(buf []byte) (int, error) {
LOOP:
	for {
		select {
		case r, ok := <-g.keys:
			if !ok {
				log.Printf("graph.Read: EOF")
				return 0, io.EOF
			}
			log.Printf("graph.Read: key: %d", r)

			switch r {
			case 13: // discard CR
				continue LOOP
			case 10: // LF
				need := 1
				avail := len(buf)
				if need > avail {
					return 0, fmt.Errorf("graph.Read: enter short buffer: need=%d avail=%d", need, avail)
				}
				buf[0] = '\n'
				return 1, nil
			default:
				need := utf8.RuneLen(r)
				avail := len(buf)
				if need > avail {
					return 0, fmt.Errorf("graph.Read: rune short buffer: need=%d avail=%d", need, avail)
				}
				size := utf8.EncodeRune(buf, r)
				return size, nil
			}
		default:
			glfw.PollEvents()
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
	gl.Uniform4f(graphics.u_color, rr, gg, bb, 1)
}

func graphicsColorUpload() {
	graphicsColorFg(screenColorForeground)

	// upload background color
	r, g, b := screenColorBackground.RGB()
	rr, gg, bb := rgbFloat(r, g, b)
	gl.ClearColor(rr, gg, bb, 1)
}

func rgbFloat(r, g, b int32) (float32, float32, float32) {
	return float32(r) / 255, float32(g) / 255, float32(b) / 255
}

func drawTriangle() {

	graphics.geom[0] = 0
	graphics.geom[1] = .5
	graphics.geom[3] = -.5
	graphics.geom[4] = -.5
	graphics.geom[6] = .5
	graphics.geom[7] = -.5

	vao := makeVao(graphics.geom, 9)
	vaoIndices := int32(3)

	draw(gl.TRIANGLES, vao, graphics.window, vaoIndices)
}

func draw(mode, vao uint32, window *glfw.Window, count int32) {
	gl.BindVertexArray(vao)
	gl.DrawArrays(mode, 0, count)

	window.SwapBuffers()
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
	return pix2Clip(x, graphics.width), -pix2Clip(y, graphics.height)
}

func Line(x1, y1, x2, y2, color, style int) {

	a1, b1 := pixelToClip(x1, y1)
	a2, b2 := pixelToClip(x2, y2)

	graphics.geom[0] = a1
	graphics.geom[1] = b1
	graphics.geom[3] = a2
	graphics.geom[4] = b2

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
		graphics.geom[3] = a1
		graphics.geom[4] = b2
		graphics.geom[6] = a2
		graphics.geom[7] = b2

		graphics.geom[9] = a2
		graphics.geom[10] = b2
		graphics.geom[12] = a2
		graphics.geom[13] = b1
		graphics.geom[15] = a1
		graphics.geom[16] = b1

		mode = gl.TRIANGLES
		vao = makeVao(graphics.geom, 18)
		vaoIndices = 6
	} else {
		graphics.geom[0] = a1
		graphics.geom[1] = b1
		graphics.geom[3] = a1
		graphics.geom[4] = b2
		graphics.geom[6] = a2
		graphics.geom[7] = b2
		graphics.geom[9] = a2
		graphics.geom[10] = b1

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

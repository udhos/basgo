package baslib

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type graph struct {
	mode    int
	window  *glfw.Window
	program uint32
	width   int
	height  int
	geom    []float32
	u_color int32
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

	//drawTriangle()

	graphics.mode = mode

	gl.ClearColor(0, 0, 0, 0) // clear color
	gl.ClearDepthf(1)         // default
	gl.Disable(gl.DEPTH_TEST) // disable depth testing
	gl.Disable(gl.CULL_FACE)  // disable face culling

	graphicsColorUpload()
}

func getUniformLocation(name string) int32 {
	u := gl.GetUniformLocation(graphics.program, gl.Str(name+"\x00"))
	if u < 0 {
		log.Printf("getUniformLocation: uniform location not found for: %s", name)
	}
	return u
}

func graphicsColorUpload() {
	r, g, b := screenColorForeground.RGB()
	graphicsColor(r, g, b)
}

func graphicsColor(r, g, b int32) {
	rr := float32(r) / 255
	gg := float32(g) / 255
	bb := float32(b) / 255
	gl.Uniform4f(graphics.u_color, rr, gg, bb, 1)
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

// 1..640 to -1..1
func pix2Clip(x, w int) float32 {
	if w == 1 {
		return 0 // ugh
	}

	w--
	x--                          // 1..w -> 0..w-1
	x *= 2                       // 0..2(w-1)
	x -= w                       // -(w-1)..w-1
	c := float32(x) / float32(w) // -1..1

	return c
}

func pixelToClip(x, y int) (float32, float32) {
	return pix2Clip(x, graphics.width), pix2Clip(y, graphics.height)
}

func Line(x1, y1, x2, y2, color, style int) {

	a1, b1 := pixelToClip(x1, y1)
	a2, b2 := pixelToClip(x2, y2)

	graphics.geom[0] = a1
	graphics.geom[1] = -b1 // invert y1
	graphics.geom[3] = a2
	graphics.geom[4] = -b2 // invert y2

	vao := makeVao(graphics.geom, 6)
	vaoIndices := int32(2)

	draw(gl.LINES, vao, graphics.window, vaoIndices)
}

func LineBox(x1, y1, x2, y2, color, style int, fill bool) {
	minx := min(x1, x2)
	miny := min(y1, y2)
	maxx := max(x1, x2)
	maxy := max(y1, y2)

	a1, b1 := pixelToClip(minx, miny)
	a2, b2 := pixelToClip(maxx, maxy)

	b1 = -b1 // invert y
	b2 = -b2 // invert y

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

	draw(mode, vao, graphics.window, vaoIndices)
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

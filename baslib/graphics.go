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
}

var graphics graph

func graphicsStop() {
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

	//gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(graphics.program)

	//drawTriangle()
}

func drawTriangle() {
	triangle := []float32{
		0, 0.5, 0, // top
		-0.5, -0.5, 0, // left
		0.5, -0.5, 0, // right
	}

	vao := makeVao(triangle)
	vaoIndices := int32(len(triangle) / 3)
	log.Printf("triangle vao: %d", vao)

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

func Line(x1, y1, x2, y2 int) {

	a1, b1 := pixelToClip(x1, y1)
	a2, b2 := pixelToClip(x2, y2)

	b1 = -b1 // invert y1
	b2 = -b2 // invert y2

	data := []float32{
		a1, b1, 0,
		a2, b2, 0,
	}

	vao := makeVao(data)
	vaoIndices := int32(len(data) / 3)

	draw(gl.LINES, vao, graphics.window, vaoIndices)
}

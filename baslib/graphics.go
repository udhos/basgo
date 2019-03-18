package baslib

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type graph struct {
	mode   int
	window *glfw.Window
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

	width := 640
	height := 480

	log.Printf("graphicsStart(%d): %d x %d", mode, width, height)

	graphics.window = initWin(width, height)
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

	triangle := []float32{
		0, 0.5, 0, // top
		-0.5, -0.5, 0, // left
		0.5, -0.5, 0, // right
	}

	vao := makeVao(triangle)
	count := int32(len(triangle) / 3)

	log.Printf("triangle vao: %d", vao)

	draw(vao, graphics.window, prog, count)
}

func draw(vao uint32, window *glfw.Window, program uint32, count int32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, count)

	//glfw.PollEvents()
	window.SwapBuffers()
}

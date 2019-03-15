package baslib

import (
	"log"
	"runtime"

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
}

func initWin(width, height int) *glfw.Window {
	if err := glfw.Init(); err != nil {
		log.Printf("%v", err)
		return nil
	}

	major := 3
	minor := 3

	log.Printf("requesting window for OpenGL %d.%d", major, minor)

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, major)
	glfw.WindowHint(glfw.ContextVersionMinor, minor)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "basgo", nil, nil)
	if err != nil {
		log.Printf("%v", err)
		return nil
	}
	window.MakeContextCurrent()

	return window
}

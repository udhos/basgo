package baslib

import (
	"fmt"
	"log"
	"strings"

	"github.com/faiface/mainthread"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	vertexShaderSource = `
    #version 330
    in vec3 vp;
    void main() {
        gl_Position = vec4(vp, 1.0);
    }
` + "\x00"

	fragmentShaderSource = `
    #version 330
    uniform vec4 u_color;
    out vec4 frag_color;
    void main() {
        frag_color = u_color;
    }
` + "\x00"
)

func InitWin(width, height int) *glfw.Window {
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
	if swap {
		glfw.WindowHint(glfw.DoubleBuffer, glfw.True)
	} else {
		glfw.WindowHint(glfw.DoubleBuffer, glfw.False) // no SwapBuffers
	}

	window, err := glfw.CreateWindow(width, height, "basgo", nil, nil)
	if err != nil {
		log.Printf("%v", err)
		return nil
	}
	window.MakeContextCurrent()

	return window
}

func initProg() (uint32, error) {

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	return prog, nil
}

// https://github.com/go-gl/examples/blob/master/gl41core-cube/cube.go
func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32, size int) uint32 {
	var vao uint32

	mainthread.Call(func() {
		var vbo uint32
		gl.GenBuffers(1, &vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, 4*size, gl.Ptr(points), gl.STATIC_DRAW)

		gl.GenVertexArrays(1, &vao)
		gl.BindVertexArray(vao)
		gl.EnableVertexAttribArray(0)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	})

	return vao
}

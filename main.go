package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"runtime"
	"strings"
)

const (
	vertexShader = `
		#version 410
		in vec3 pos;
		void main() {
			gl_Position = vec4(pos, 1.0);
		}
	`
	fragmentShader = `
		#version 410
		out vec4 color;
		void main() {
			color = vec4(1, 1, 1, 1.0);
		}
	`
)

var triangle = []float32{
	0, 0.5, 0,
	-0.5, -0.5, 0,
	0.5, -0.5, 0,
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func initGlfw() *glfw.Window {
	check(glfw.Init())
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	win, err := glfw.CreateWindow(800, 600, "comanche", nil, nil)
	check(err)
	win.MakeContextCurrent()
	return win
}

func attachShader(program uint32, kind uint32, src string) {
	s := gl.CreateShader(kind)
	strs, free := gl.Strs(src + "\x00")
	gl.ShaderSource(s, 1, strs, nil)
	free()
	gl.CompileShader(s)
	gl.AttachShader(program, s)

	var status int32
	gl.GetShaderiv(s, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var length int32
		gl.GetShaderiv(s, gl.INFO_LOG_LENGTH, &length)
		log := strings.Repeat("\x00", int(length+1))
		gl.GetShaderInfoLog(s, length, nil, gl.Str(log))
		panic(fmt.Sprintf("shader compilation failed: %v\n%v", log, src))
	}
}

func initGl() (uint32, uint32) {
	check(gl.Init())
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println(version)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(triangle), gl.Ptr(triangle), gl.STATIC_DRAW)
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	p := gl.CreateProgram()
	attachShader(p, gl.VERTEX_SHADER, vertexShader)
	attachShader(p, gl.FRAGMENT_SHADER, fragmentShader)
	gl.LinkProgram(p)
	return vao, p
}

func render(vao, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(3))
}

func main() {
	runtime.LockOSThread()

	win := initGlfw()
	defer glfw.Terminate()
	vao, program := initGl()
	for !win.ShouldClose() {
		render(vao, program)
		win.SwapBuffers()
		glfw.PollEvents()
	}
}

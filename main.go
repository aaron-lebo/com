package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"runtime"
	"strings"
)

const (
	vertexShader = `
		#version 410
		uniform mat4 mvp;
		in vec3 pos;
		void main() {
			gl_Position = mvp * vec4(pos, 1.0);
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
var indices = []uint16{
	0, 1, 2,
}

type Shader struct {
	program uint32
	mvp     int32
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

func attachShader(program, kind uint32, src string) {
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

func initGl() *Shader {
	check(gl.Init())
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println(version)

	var vao, vbo, ibo uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(triangle), gl.Ptr(triangle), gl.STATIC_DRAW)
	gl.GenBuffers(1, &ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 2*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	p := gl.CreateProgram()
	attachShader(p, gl.VERTEX_SHADER, vertexShader)
	attachShader(p, gl.FRAGMENT_SHADER, fragmentShader)
	gl.LinkProgram(p)
	return &Shader{p, gl.GetUniformLocation(p, gl.Str("mvp\x00"))}
}

func render(shader *Shader) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(shader.program)
	mvp := mgl32.Ident4()
	gl.UniformMatrix4fv(shader.mvp, 1, false, &mvp[0])
	gl.EnableVertexAttribArray(0)
	gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_SHORT, nil)
	gl.DisableVertexAttribArray(0)
}

func main() {
	runtime.LockOSThread()

	win := initGlfw()
	defer glfw.Terminate()
	program := initGl()
	for !win.ShouldClose() {
		render(program)
		win.SwapBuffers()
		glfw.PollEvents()
	}
}

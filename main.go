package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"math"
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
	a = 0.5
	b = -a
)

var (
	block = []float32{
		b, b, a,
		a, b, a,
		a, a, a,
		b, a, a,
		b, b, b,
		b, a, b,
		a, a, b,
		a, b, b,
	}
	chunk          [24 * 4096]float32
	chunkIndices   = make([]uint16, 0, 36*4096)
	keys           [512]bool
	mouseX, mouseY float64
	position       = mgl32.Vec3{0, 0, 10}
	direction      = mgl32.Vec3{0, 0, -1}
	up             = mgl32.Vec3{0, 1, 0}
	pitch          = 0.0
	yaw            = -90.0
)

func addBlock(x, y, z int) {
	var blockIndices = []uint16{
		0, 1, 2, 2, 3, 0, // +z
		4, 5, 6, 6, 7, 4, // -z
		3, 2, 6, 6, 5, 3, // +y
		0, 4, 7, 7, 1, 0, // -y
		1, 7, 6, 6, 2, 1, // +x
		0, 3, 5, 5, 4, 0, // -x
	}
	for _, idx := range blockIndices {
		chunkIndices = append(chunkIndices, idx)
	}
}

type ShaderProgram struct {
	program uint32
	mvp     int32
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func keyCallback(win *glfw.Window, key glfw.Key, scancode int, act glfw.Action, mod glfw.ModifierKey) {
	if act == glfw.Press {
		keys[key] = true
	} else if act == glfw.Release {
		keys[key] = false
	}
}

func mouseCallback(win *glfw.Window, x, y float64) {
	const sensitivity = 0.05
	yaw += sensitivity * (x - mouseX)
	pitch += sensitivity * (mouseY - y)
	if pitch < -89.0 {
		pitch = -89.0
	} else if pitch > 89.0 {
		pitch = 89.0
	}
	mouseX = x
	mouseY = y
}

func sizeCallback(win *glfw.Window, w, h int) {
	gl.Viewport(0, 0, int32(w), int32(h))
	mouseX = float64(w / 2)
	mouseY = float64(h / 2)
	win.SetCursorPos(mouseX, mouseY)
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
	win.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	win.SetKeyCallback(keyCallback)
	win.SetCursorPosCallback(mouseCallback)
	win.SetSizeCallback(sizeCallback)
	win.SetCursorPos(mouseX, mouseY)
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

func initGl() *ShaderProgram {
	check(gl.Init())
	fmt.Println(gl.GoStr(gl.GetString(gl.VERSION)))

	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)

	var vao, vbo, ibo uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(block), gl.Ptr(block), gl.STATIC_DRAW)
	gl.GenBuffers(1, &ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 2*len(chunkIndices), gl.Ptr(chunkIndices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	p := gl.CreateProgram()
	attachShader(p, gl.VERTEX_SHADER, vertexShader)
	attachShader(p, gl.FRAGMENT_SHADER, fragmentShader)
	gl.LinkProgram(p)
	return &ShaderProgram{p, gl.GetUniformLocation(p, gl.Str("mvp\x00"))}
}

func radians(angle float64) float64 {
	return angle * math.Pi / 180
}

func update() {
	direction = mgl32.Vec3{
		float32(math.Cos(radians(yaw)) * math.Cos(radians(pitch))),
		float32(math.Sin(radians(pitch))),
		float32(math.Sin(radians(yaw)) * math.Cos(radians(pitch))),
	}.Normalize()

	const speed = 0.5
	pos := direction.Mul(speed)
	if keys[glfw.KeyW] {
		position = position.Add(pos)
	}
	if keys[glfw.KeyA] {
		position = position.Sub(pos.Cross(up).Normalize())
	}
	if keys[glfw.KeyS] {
		position = position.Sub(pos)
	}
	if keys[glfw.KeyD] {
		position = position.Add(pos.Cross(up).Normalize())
	}
}

func render(program *ShaderProgram) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program.program)
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), 4/3, 0.1, 1000.0)
	view := mgl32.LookAtV(position, position.Add(direction), up)
	model := mgl32.Ident4()
	mvp := projection.Mul4(view).Mul4(model)
	gl.UniformMatrix4fv(program.mvp, 1, false, &mvp[0])
	gl.EnableVertexAttribArray(0)
	gl.DrawElements(gl.TRIANGLES, int32(len(chunkIndices)), gl.UNSIGNED_SHORT, nil)
	gl.DisableVertexAttribArray(0)
}

func main() {
	runtime.LockOSThread()

	addBlock(0, 0, 0)
	win := initGlfw()
	defer glfw.Terminate()
	program := initGl()
	for !win.ShouldClose() {
		update()
		render(program)
		win.SwapBuffers()
		glfw.PollEvents()
	}
}

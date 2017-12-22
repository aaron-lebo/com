package main

import (
	"comanche/block"
	_ "comanche/text"
	. "comanche/util"
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"runtime"
)

var (
	keys           [512]bool
	mouseX, mouseY float64
	position       = mgl32.Vec3{0, 0, 10}
	direction      = mgl32.Vec3{0, 0, -1}
	up             = mgl32.Vec3{0, 1, 0}
	pitch          = 0.0
	yaw            = -90.0
)

func init() {
	runtime.LockOSThread()
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

func initGl() *glfw.Window {
	Check(glfw.Init())
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	win, err := glfw.CreateWindow(800, 600, "comanche", nil, nil)
	Check(err)
	win.MakeContextCurrent()
	win.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	win.SetKeyCallback(keyCallback)
	win.SetCursorPosCallback(mouseCallback)
	win.SetSizeCallback(sizeCallback)
	win.SetCursorPos(mouseX, mouseY)

	Check(gl.Init())
	fmt.Println(gl.GoStr(gl.GetString(gl.VERSION)))

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	return win
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

func main() {
	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {
			for z := 0; z < 16; z++ {
				block.Add(float32(x), float32(y), float32(z))
			}
		}
	}

	win := initGl()
	defer glfw.Terminate()
	block.Init()
	//text.Init()
	for !win.ShouldClose() {
		update()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		projection := mgl32.Perspective(mgl32.DegToRad(45.0), 4/3, 0.1, 1000.0)
		view := mgl32.LookAtV(position, position.Add(direction), up)
		model := mgl32.Ident4()
		block.Render(projection.Mul4(view).Mul4(model))

		//text.Render("test", 10, 10)

		win.SwapBuffers()
		glfw.PollEvents()
	}
}

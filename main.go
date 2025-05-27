package main

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/lunararch/helios/pkg/graphics/shader"
	"runtime"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.True)

	window, err := glfw.CreateWindow(640, 480, "Helios", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	glfw.SwapInterval(1)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	width, height := window.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.ClearColor(0.2, 0.3, 0.8, 1.0)

	shaderProgram, err := shader.New("assets/shaders/basic.vert", "assets/shaders/basic.frag")
	if err != nil {
		panic(err)
	}
	defer shaderProgram.Delete()

	shaderProgram.Use()

	shaderProgram.SetVec4("color", mgl32.Vec4{1.0, 0.5, 0.2, 1.0})

	projection := mgl32.Ortho(0, 640, 480, 0, -1, 1)
	shaderProgram.SetMat4("projection", projection)

	view := mgl32.Ident4()
	shaderProgram.SetMat4("view", view)

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		// Close on Escape key press
		if key == glfw.KeyEscape && action == glfw.Press {
			w.SetShouldClose(true)
		}
	})

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		shaderProgram.Use()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

package main

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/lunararch/helios/pkg/graphics/shader"
	"github.com/lunararch/helios/pkg/graphics/texture"
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

	vertices := []float32{
		320.0, 100.0, 0.0, 0.5, 0.0,
		220.0, 300.0, 0.0, 0.0, 1.0,
		420.0, 300.0, 0.0, 1.0, 1.0,
	}

	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	textureImg, err := texture.LoadFromFile("assets/textures/knight.png")
	if err != nil {
		panic(err)
	}
	defer textureImg.Delete()

	shaderProgram.Use()
	shaderProgram.SetInt("texture1", 0)
	shaderProgram.SetVec4("color", mgl32.Vec4{1.0, 1.0, 1.0, 1.0})

	projection := mgl32.Ortho(0, 640, 480, 0, -1, 1)
	shaderProgram.SetMat4("projection", projection)

	view := mgl32.Ident4()
	shaderProgram.SetMat4("view", view)

	model := mgl32.Ident4()
	shaderProgram.SetMat4("model", model)

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		// Close on Escape key press
		if key == glfw.KeyEscape && action == glfw.Press {
			w.SetShouldClose(true)
		}
	})

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		shaderProgram.Use()

		textureImg.Bind(0)

		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		window.SwapBuffers()
		glfw.PollEvents()
	}

	gl.DeleteVertexArrays(1, &VAO)
	gl.DeleteBuffers(1, &VBO)
}

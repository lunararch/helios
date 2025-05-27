package main

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/lunararch/helios/pkg/graphics/camera"
	"github.com/lunararch/helios/pkg/graphics/shader"
	"github.com/lunararch/helios/pkg/graphics/sprite"
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

	width, height = window.GetSize()
	gameCamera := camera.New(float32(width), float32(height))

	shaderProgram, err := shader.New("assets/shaders/basic.vert", "assets/shaders/basic.frag")
	if err != nil {
		panic(err)
	}
	defer shaderProgram.Delete()

	spriteRenderer := sprite.NewRenderer(shaderProgram)
	defer spriteRenderer.Delete()

	textureImg, err := texture.LoadFromFile("assets/textures/knight.png")
	if err != nil {
		panic(err)
	}
	defer textureImg.Delete()

	knightSprite := sprite.NewSprite(
		textureImg,
		mgl32.Vec3{200.0, 100.0, 0.0},
		mgl32.Vec2{float32(textureImg.Width), float32(textureImg.Height)},
	)

	projection := mgl32.Ortho(0, 640, 480, 0, -1, 1)
	shaderProgram.Use()
	shaderProgram.SetInt("texture1", 0)
	shaderProgram.SetMat4("projection", projection)
	//shaderProgram.SetMat4("view", mgl32.Ident4())

	gameCamera.Position = mgl32.Vec2{float32(width) / 2, float32(height) / 2}

	cameraSpeed := float32(200.0)
	lastFrameTime := glfw.GetTime()

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		// Close on Escape key press
		if key == glfw.KeyEscape && action == glfw.Press {
			w.SetShouldClose(true)
		}
	})

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		gameCamera.Size = mgl32.Vec2{float32(width), float32(height)}
		projection := mgl32.Ortho(0, float32(width), float32(height), 0, -1, 1)
		shaderProgram.Use()
		shaderProgram.SetMat4("projection", projection)
	})

	for !window.ShouldClose() {
		currentTime := glfw.GetTime()
		deltaTime := float32(currentTime - lastFrameTime)
		lastFrameTime = currentTime

		if window.GetKey(glfw.KeyS) == glfw.Press {
			gameCamera.Position[1] -= cameraSpeed * deltaTime
		}
		if window.GetKey(glfw.KeyW) == glfw.Press {
			gameCamera.Position[1] += cameraSpeed * deltaTime
		}
		if window.GetKey(glfw.KeyD) == glfw.Press {
			gameCamera.Position[0] -= cameraSpeed * deltaTime
		}
		if window.GetKey(glfw.KeyA) == glfw.Press {
			gameCamera.Position[0] += cameraSpeed * deltaTime
		}

		if window.GetKey(glfw.KeyQ) == glfw.Press {
			gameCamera.Zoom *= (1.0 - deltaTime)
			if gameCamera.Zoom < 0.1 {
				gameCamera.Zoom = 0.1
			}
		}
		if window.GetKey(glfw.KeyE) == glfw.Press {
			gameCamera.Zoom *= (1.0 + deltaTime)
			if gameCamera.Zoom > 10.0 {
				gameCamera.Zoom = 10.0
			}
		}

		gameCamera.ClampToBounds()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		shaderProgram.Use()
		shaderProgram.SetMat4("view", gameCamera.GetViewMatrix())

		spriteRenderer.DrawSprite(knightSprite)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

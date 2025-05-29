package main

import (
	"runtime"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/lunararch/helios/pkg/engine"
	"github.com/lunararch/helios/pkg/graphics/camera"
	"github.com/lunararch/helios/pkg/graphics/shader"
	"github.com/lunararch/helios/pkg/graphics/sprite"
	"github.com/lunararch/helios/pkg/graphics/texture"
	"github.com/lunararch/helios/pkg/input"
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

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.ClearColor(0.2, 0.3, 0.8, 1.0)

	inputManager := input.NewInputManager(window)
	inputMapping := input.NewInputMapping()

	inputMapping.MapKey("move_up", glfw.KeyW)
	inputMapping.MapKey("move_up", glfw.KeyUp)
	inputMapping.MapKey("move_down", glfw.KeyS)
	inputMapping.MapKey("move_down", glfw.KeyDown)
	inputMapping.MapKey("move_left", glfw.KeyA)
	inputMapping.MapKey("move_left", glfw.KeyLeft)
	inputMapping.MapKey("move_right", glfw.KeyD)
	inputMapping.MapKey("move_right", glfw.KeyRight)
	inputMapping.MapKey("zoom_in", glfw.KeyE)
	inputMapping.MapKey("zoom_out", glfw.KeyQ)
	inputMapping.MapKey("quit", glfw.KeyEscape)
	inputMapping.MapKey("pause", glfw.KeyP)
	inputMapping.MapKey("slow_motion", glfw.KeyLeftShift)
	inputMapping.MapKey("fast_forward", glfw.KeyTab)
	inputMapping.MapKey("reset_time", glfw.KeyR)

	inputManager.AddInputCallback(func(event input.InputEvent) {
		switch e := event.(type) {
		case input.KeyPressEvent:
			if e.Key == glfw.KeyEscape {
				window.SetShouldClose(true)
			}
		case input.MousePressEvent:
			// Handle mouse clicks
		case input.MouseScrollEvent:
			// Handle scroll events
		}
	})

	gameCamera := camera.New(float32(width), float32(height))

	batchShader, err := shader.New("assets/shaders/batch.vert", "assets/shaders/batch.frag")
	if err != nil {
		panic(err)
	}
	defer batchShader.Delete()

	projection := mgl32.Ortho(0, 640, 480, 0, -1, 1)
	batchShader.Use()
	batchShader.SetInt("texture1", 0)
	batchShader.SetMat4("projection", projection)

	spriteBatch := sprite.NewSpriteBatch(batchShader)
	defer spriteBatch.Delete()

	knightTexture, err := texture.LoadFromFile("assets/textures/knight.png")
	if err != nil {
		panic(err)
	}
	defer knightTexture.Delete()

	hornetTexture, err := texture.LoadFromFile("assets/textures/hornet.png")
	if err != nil {
		panic(err)
	}
	defer hornetTexture.Delete()

	knightSprite := sprite.NewSprite(
		knightTexture,
		mgl32.Vec3{100.0, 100.0, 0.0},
		mgl32.Vec2{float32(knightTexture.Width), float32(knightTexture.Height)},
	)

	hornetSprite := sprite.NewSprite(
		hornetTexture,
		mgl32.Vec3{300.0, 200.0, 0.0},
		mgl32.Vec2{float32(hornetTexture.Width), float32(hornetTexture.Height)},
	)

	gameCamera.Position = mgl32.Vec2{float32(width) / 2, float32(height) / 2}
	gameCamera.SetBounds(0, 0, float32(width), float32(height))
	cameraSpeed := float32(200.0)

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		gameCamera.Size = mgl32.Vec2{float32(width), float32(height)}
		projection := mgl32.Ortho(0, float32(width), float32(height), 0, -1, 1)
		batchShader.Use()
		batchShader.SetMat4("projection", projection)
	})

	gameLoop := engine.NewGameLoop(window)

	gameLoop.UseFixedTimestep(true)
	gameLoop.SetTargetFPS(60)

	rotationTimer := engine.NewRepeatingTimer(2.0)
	rotationTimer.SetOnComplete(func() {
		// Reverse the rotation direction
		// You could implement this by storing rotation direction in a variable
	})
	rotationTimer.Start()

	printTimer := engine.NewTimer(5.0)
	printTimer.SetOnComplete(func() {
		timeManager := gameLoop.TimeManager()
		println("=== Performance Stats ===")
		println("FPS:", timeManager.FPS())
		println("Avg Delta:", timeManager.AvgDeltaTime())
		println("Min Delta:", timeManager.MinDeltaTime())
		println("Max Delta:", timeManager.MaxDeltaTime())
		println("Frame Count:", timeManager.FrameCount())
		println("Time Scale:", timeManager.TimeScale())
		println("========================")
		printTimer.Restart()
	})
	printTimer.Start()

	gameLoop.SetUpdateFunc(func(deltaTime float32) {
		inputManager.Update()

		if inputMapping.IsActionPressed("pause", inputManager) {
			gameLoop.TogglePause()
		}

		if inputMapping.IsActionPressed("reset_time", inputManager) {
			gameLoop.TimeManager().ResetPerformanceStats()
		}

		if inputMapping.IsActionHeld("slow_motion", inputManager) {
			gameLoop.SetTimeScale(0.3)
		} else if inputMapping.IsActionHeld("fast_forward", inputManager) {
			gameLoop.SetTimeScale(2.0)
		} else {
			gameLoop.SetTimeScale(1.0)
		}

		rotationTimer.Update(deltaTime)
		printTimer.Update(deltaTime)

		if inputMapping.IsActionHeld("move_up", inputManager) {
			gameCamera.Position[1] -= cameraSpeed * deltaTime
		}
		if inputMapping.IsActionHeld("move_down", inputManager) {
			gameCamera.Position[1] += cameraSpeed * deltaTime
		}
		if inputMapping.IsActionHeld("move_left", inputManager) {
			gameCamera.Position[0] -= cameraSpeed * deltaTime
		}
		if inputMapping.IsActionHeld("move_right", inputManager) {
			gameCamera.Position[0] += cameraSpeed * deltaTime
		}

		if inputMapping.IsActionHeld("zoom_out", inputManager) {
			gameCamera.Zoom *= (1.0 - deltaTime)
			if gameCamera.Zoom < 0.1 {
				gameCamera.Zoom = 0.1
			}
		}
		if inputMapping.IsActionHeld("zoom_in", inputManager) {
			gameCamera.Zoom *= (1.0 + deltaTime)
			if gameCamera.Zoom > 10.0 {
				gameCamera.Zoom = 10.0
			}
		}

		// Example of using mouse input
		if inputManager.IsMouseButtonPressed(input.MouseButtonLeft) {
			mousePos := inputManager.GetMousePosition()
			// Do something with mouse click at mousePos
			_ = mousePos
		}

		// Handle scroll wheel for zoom
		scrollDelta := inputManager.GetScrollDelta()
		if scrollDelta.Y() != 0 {
			zoomFactor := 1.0 + scrollDelta.Y()*0.1
			gameCamera.Zoom *= zoomFactor
			if gameCamera.Zoom < 0.1 {
				gameCamera.Zoom = 0.1
			}
			if gameCamera.Zoom > 10.0 {
				gameCamera.Zoom = 10.0
			}
		}

		// Animate sprite rotation (affected by time scale)
		knightSprite.Rotation += deltaTime * 1.0

		gameCamera.Update(deltaTime)
		gameCamera.ClampToBounds()
	})

	gameLoop.SetRenderFunc(func(alpha float32) {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		batchShader.Use()
		batchShader.SetMat4("view", gameCamera.GetViewMatrix())

		spriteBatch.Begin()
		spriteBatch.Draw(knightSprite)
		spriteBatch.Draw(hornetSprite)
		spriteBatch.End()
	})

	gameLoop.Start()
}

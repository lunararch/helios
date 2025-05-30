package main

import (
	"runtime"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/lunararch/helios/pkg/engine"
	"github.com/lunararch/helios/pkg/graphics/camera"
	"github.com/lunararch/helios/pkg/input"
	"github.com/lunararch/helios/pkg/scene"
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

	// Setup input system
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
	inputMapping.MapKey("menu", glfw.KeyM)

	// Setup camera
	gameCamera := camera.New(float32(width), float32(height))
	gameCamera.Position = mgl32.Vec2{float32(width) / 2, float32(height) / 2}
	gameCamera.SetBounds(0, 0, float32(width), float32(height))

	// Setup scene management
	sceneManager := scene.NewSceneManager()
	defer sceneManager.Cleanup()

	// Create and register scenes
	menuScene := scene.NewMenuScene(gameCamera)
	gameplayScene := scene.NewGameplayScene(gameCamera)

	sceneManager.RegisterScene(menuScene)
	sceneManager.RegisterScene(gameplayScene)

	// Start with menu scene
	sceneManager.SwitchToScene("menu")

	// Setup input callbacks for scene switching
	inputManager.AddInputCallback(func(event input.InputEvent) {
		switch e := event.(type) {
		case input.KeyPressEvent:
			switch e.Key {
			case glfw.KeyEscape:
				window.SetShouldClose(true)
			case glfw.KeyEnter:
				if sceneManager.GetCurrentScene().GetName() == "menu" {
					sceneManager.SwitchToScene("gameplay")
				}
			case glfw.KeyM:
				currentScene := sceneManager.GetCurrentScene().GetName()
				if currentScene == "gameplay" {
					sceneManager.PushScene("menu")
				} else if currentScene == "menu" {
					sceneManager.PopScene()
				}
			}
		}
	})

	// Setup window resize callback
	window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		gameCamera.Size = mgl32.Vec2{float32(width), float32(height)}
	})

	// Setup game loop
	gameLoop := engine.NewGameLoop(window)
	gameLoop.UseFixedTimestep(true)
	gameLoop.SetTargetFPS(60)

	gameLoop.SetUpdateFunc(func(deltaTime float32) {
		inputManager.SetDeltaTime(deltaTime)
		inputManager.Update()

		// Global input handling
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

		// Update scene manager
		sceneManager.Update(deltaTime)

		// Let current scene handle input
		sceneManager.HandleInput(inputManager, inputMapping)
	})

	gameLoop.SetRenderFunc(func(alpha float32) {
		// Scene manager handles rendering
		sceneManager.Render(alpha)
	})

	println("Controls:")
	println("Enter - Switch from menu to gameplay")
	println("M - Toggle between menu and gameplay")
	println("WASD/Arrow Keys - Move camera (in gameplay)")
	println("E/Q - Zoom in/out (in gameplay)")
	println("P - Pause")
	println("Escape - Quit")

	gameLoop.Start()
}

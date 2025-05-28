package engine

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"time"
)

const (
	DefaultFixedDeltaTime = 1.0 / 60.0 // 60 FPS
	MaxDeltaTime          = 0.25       // Maximum delta time to prevent large jumps in time
)

type UpdateFunc func(deltaTime float32)

type RenderFunc func(alpha float32)

type GameLoop struct {
	window         *glfw.Window
	fixedTimestep  bool
	targetFPS      int
	fixedDeltaTime float64
	running        bool
	accumulator    float64
	updateFunc     UpdateFunc
	renderFunc     RenderFunc
	lastFrameTime  float64
	currentFPS     float32
	fpsUpdateTime  float64
	frameCount     int
}

func NewGameLoop(window *glfw.Window) *GameLoop {
	return &GameLoop{
		window:         window,
		fixedTimestep:  true,
		targetFPS:      60,
		fixedDeltaTime: DefaultFixedDeltaTime,
		running:        false,
	}
}

func (gl *GameLoop) SetUpdateFunc(update UpdateFunc) {
	gl.updateFunc = update
}

func (gl *GameLoop) SetRenderFunc(render RenderFunc) {
	gl.renderFunc = render
}

func (gl *GameLoop) UseFixedTimestep(fixed bool) {
	gl.fixedTimestep = fixed
}

func (gl *GameLoop) SetTargetFPS(fps int) {
	gl.targetFPS = fps
	if fps > 0 {
		gl.fixedDeltaTime = 1.0 / float64(fps)
	}
}

func (gl *GameLoop) GetCurrentFPS() float32 {
	return gl.currentFPS
}

func (gl *GameLoop) Start() {
	if gl.updateFunc == nil || gl.renderFunc == nil {
		panic("Update and Render functions must be set before starting the game loop")
	}

	gl.running = true
	gl.lastFrameTime = glfw.GetTime()
	gl.accumulator = 0
	gl.frameCount = 0
	gl.fpsUpdateTime = gl.lastFrameTime

	if gl.fixedTimestep {
		gl.fixedTimeStepLoop()
	} else {
		gl.variableTimestepLoop()
	}
}

func (gl *GameLoop) fixedTimeStepLoop() {
	for gl.running && !gl.window.ShouldClose() {
		currentTime := glfw.GetTime()
		deltaTime := currentTime - gl.lastFrameTime
		gl.lastFrameTime = currentTime

		if deltaTime > MaxDeltaTime {
			deltaTime = MaxDeltaTime
		}

		gl.accumulator += deltaTime

		for gl.accumulator >= gl.fixedDeltaTime {
			gl.updateFunc(float32(gl.fixedDeltaTime))
			gl.accumulator -= gl.fixedDeltaTime
		}

		alpha := float32(gl.accumulator / gl.fixedDeltaTime)

		gl.renderFunc(alpha)

		gl.updateFPS(currentTime)

		gl.window.SwapBuffers()
		glfw.PollEvents()

		if gl.targetFPS > 0 {
			gl.limitFPS()
		}
	}
}

func (gl *GameLoop) variableTimestepLoop() {
	for gl.running && !gl.window.ShouldClose() {
		currentTime := glfw.GetTime()
		deltaTime := currentTime - gl.lastFrameTime
		gl.lastFrameTime = currentTime

		if deltaTime > MaxDeltaTime {
			deltaTime = MaxDeltaTime
		}

		gl.updateFunc(float32(deltaTime))

		gl.renderFunc(1.0)

		gl.updateFPS(currentTime)

		gl.window.SwapBuffers()
		glfw.PollEvents()

		if gl.targetFPS > 0 {
			gl.limitFPS()
		}
	}
}

func (gl *GameLoop) updateFPS(currentTime float64) {
	gl.frameCount++

	if currentTime-gl.fpsUpdateTime >= 1.0 {
		gl.currentFPS = float32(gl.frameCount) / float32(currentTime-gl.fpsUpdateTime)
		gl.fpsUpdateTime = currentTime
		gl.frameCount = 0
	}
}

func (gl *GameLoop) limitFPS() {
	if gl.targetFPS <= 0 {
		return
	}

	frameDuration := 1.0 / float64(gl.targetFPS)

	elapsed := glfw.GetTime() - gl.lastFrameTime

	if elapsed < frameDuration {
		sleepTime := time.Duration((frameDuration - elapsed) * 1000000000)
		time.Sleep(sleepTime)
	}
}

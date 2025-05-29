package engine

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	DefaultFixedDeltaTime = 1.0 / 60.0 // 60 FPS
	MaxDeltaTime          = 0.25       // Maximum delta time to prevent large jumps in time
)

type UpdateFunc func(deltaTime float32)
type RenderFunc func(alpha float32)

type GameLoop struct {
	window        *glfw.Window
	timeManager   *TimeManager
	fixedTimestep bool
	running       bool
	accumulator   float64
	updateFunc    UpdateFunc
	renderFunc    RenderFunc
}

func NewGameLoop(window *glfw.Window) *GameLoop {
	return &GameLoop{
		window:        window,
		timeManager:   NewTimeManager(),
		fixedTimestep: true,
		running:       false,
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
	gl.timeManager.SetTargetFPS(fps)
}

func (gl *GameLoop) TimeManager() *TimeManager {
	return gl.timeManager
}

func (gl *GameLoop) GetCurrentFPS() float32 {
	return gl.timeManager.FPS()
}

func (gl *GameLoop) PauseGame() {
	gl.timeManager.Pause()
}

func (gl *GameLoop) ResumeGame() {
	gl.timeManager.Resume()
}

func (gl *GameLoop) TogglePause() {
	gl.timeManager.TogglePause()
}

func (gl *GameLoop) SetTimeScale(scale float64) {
	gl.timeManager.SetTimeScale(scale)
}

func (gl *GameLoop) Start() {
	if gl.updateFunc == nil || gl.renderFunc == nil {
		panic("Update and Render functions must be set before starting the game loop")
	}

	gl.running = true
	gl.accumulator = 0

	if gl.fixedTimestep {
		gl.fixedTimeStepLoop()
	} else {
		gl.variableTimestepLoop()
	}
}

func (gl *GameLoop) Stop() {
	gl.running = false
}

func (gl *GameLoop) fixedTimeStepLoop() {
	for gl.running && !gl.window.ShouldClose() {
		gl.timeManager.Update()

		deltaTime := gl.timeManager.UnscaledDeltaTime()

		if deltaTime > MaxDeltaTime {
			deltaTime = MaxDeltaTime
		}

		gl.accumulator += float64(deltaTime)

		for gl.accumulator >= DefaultFixedDeltaTime {
			gl.updateFunc(gl.timeManager.DeltaTime())
			gl.accumulator -= DefaultFixedDeltaTime
		}

		alpha := float32(gl.accumulator / DefaultFixedDeltaTime)
		gl.renderFunc(alpha)

		gl.window.SwapBuffers()
		glfw.PollEvents()

		gl.timeManager.SleepForFrameLimit()
	}
}

func (gl *GameLoop) variableTimestepLoop() {
	for gl.running && !gl.window.ShouldClose() {
		gl.timeManager.Update()

		deltaTime := gl.timeManager.DeltaTime()

		if deltaTime > MaxDeltaTime {
			deltaTime = MaxDeltaTime
		}

		gl.updateFunc(deltaTime)
		gl.renderFunc(1.0)

		gl.window.SwapBuffers()
		glfw.PollEvents()

		gl.timeManager.SleepForFrameLimit()
	}
}

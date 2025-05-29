package engine

import (
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type TimeManager struct {
	startTime     float64
	currentTime   float64
	lastFrameTime float64
	deltaTime     float64
	unscaledDelta float64
	totalTime     float64
	frameCount    uint64

	timeScale      float64
	isPaused       bool
	pauseStartTime float64
	totalPauseTime float64

	fps           float32
	fpsUpdateTime float64
	fpsFrameCount int
	targetFPS     int

	frameTimeTarget float64
	lastSleepTime   time.Duration

	minDeltaTime float64
	maxDeltaTime float64
	avgDeltaTime float64
	deltaHistory []float64
	historyIndex int
	historySize  int
}

func NewTimeManager() *TimeManager {
	historySize := 60
	return &TimeManager{
		startTime:       glfw.GetTime(),
		timeScale:       1.0,
		isPaused:        false,
		targetFPS:       60,
		frameTimeTarget: 1.0 / 60.0,
		minDeltaTime:    1.0,
		maxDeltaTime:    0.0,
		deltaHistory:    make([]float64, historySize),
		historySize:     historySize,
	}
}

func (tm *TimeManager) Update() {
	tm.currentTime = glfw.GetTime()

	rawDelta := tm.currentTime - tm.lastFrameTime
	tm.lastFrameTime = tm.currentTime

	tm.unscaledDelta = rawDelta

	if tm.isPaused {
		tm.deltaTime = 0.0
	} else {
		tm.deltaTime = rawDelta * tm.timeScale
		tm.totalTime += tm.deltaTime
	}

	tm.frameCount++

	tm.updatePerformanceStats(rawDelta)

	tm.updateFPS()
}

func (tm *TimeManager) updatePerformanceStats(delta float64) {
	if delta < tm.minDeltaTime {
		tm.minDeltaTime = delta
	}
	if delta > tm.maxDeltaTime {
		tm.maxDeltaTime = delta
	}

	tm.deltaHistory[tm.historyIndex] = delta
	tm.historyIndex = (tm.historyIndex + 1) % tm.historySize

	sum := 0.0
	count := 0
	for _, d := range tm.deltaHistory {
		if d > 0 {
			sum += d
			count++
		}
	}
	if count > 0 {
		tm.avgDeltaTime = sum / float64(count)
	}
}

func (tm *TimeManager) updateFPS() {
	tm.fpsFrameCount++

	if tm.currentTime-tm.fpsUpdateTime >= 1.0 {
		tm.fps = float32(tm.fpsFrameCount) / float32(tm.currentTime-tm.fpsUpdateTime)
		tm.fpsUpdateTime = tm.currentTime
		tm.fpsFrameCount = 0
	}
}

func (tm *TimeManager) Pause() {
	if !tm.isPaused {
		tm.isPaused = true
		tm.pauseStartTime = tm.currentTime
	}
}

func (tm *TimeManager) Resume() {
	if tm.isPaused {
		tm.isPaused = false
		pauseDuration := tm.currentTime - tm.pauseStartTime
		tm.totalPauseTime += pauseDuration
	}
}

func (tm *TimeManager) TogglePause() {
	if tm.isPaused {
		tm.Resume()
	} else {
		tm.Pause()
	}
}

func (tm *TimeManager) SetTimeScale(scale float64) {
	if scale < 0 {
		scale = 0
	}
	tm.timeScale = scale
}

func (tm *TimeManager) SetTargetFPS(fps int) {
	tm.targetFPS = fps
	if fps > 0 {
		tm.frameTimeTarget = 1.0 / float64(fps)
	}
}

func (tm *TimeManager) SleepForFrameLimit() {
	if tm.targetFPS <= 0 {
		return
	}

	elapsed := glfw.GetTime() - tm.lastFrameTime

	if elapsed < tm.frameTimeTarget {
		sleepTime := time.Duration((tm.frameTimeTarget - elapsed) * 1000000000)
		time.Sleep(sleepTime)
		tm.lastSleepTime = sleepTime
	} else {
		tm.lastSleepTime = 0
	}
}

func (tm *TimeManager) DeltaTime() float32         { return float32(tm.deltaTime) }
func (tm *TimeManager) UnscaledDeltaTime() float32 { return float32(tm.unscaledDelta) }
func (tm *TimeManager) TotalTime() float32         { return float32(tm.totalTime) }
func (tm *TimeManager) CurrentTime() float32       { return float32(tm.currentTime) }
func (tm *TimeManager) FrameCount() uint64         { return tm.frameCount }
func (tm *TimeManager) FPS() float32               { return tm.fps }
func (tm *TimeManager) TimeScale() float64         { return tm.timeScale }
func (tm *TimeManager) IsPaused() bool             { return tm.isPaused }

func (tm *TimeManager) MinDeltaTime() float32        { return float32(tm.minDeltaTime) }
func (tm *TimeManager) MaxDeltaTime() float32        { return float32(tm.maxDeltaTime) }
func (tm *TimeManager) AvgDeltaTime() float32        { return float32(tm.avgDeltaTime) }
func (tm *TimeManager) LastSleepTime() time.Duration { return tm.lastSleepTime }

func (tm *TimeManager) GetTimeSinceStart() float32 {
	return float32(tm.currentTime - tm.startTime - tm.totalPauseTime)
}

func (tm *TimeManager) ResetPerformanceStats() {
	tm.minDeltaTime = 1.0
	tm.maxDeltaTime = 0.0
	tm.avgDeltaTime = 0.0
	for i := range tm.deltaHistory {
		tm.deltaHistory[i] = 0.0
	}
	tm.historyIndex = 0
}

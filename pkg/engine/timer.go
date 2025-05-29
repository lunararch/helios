package engine

type Timer struct {
	duration    float32
	remaining   float32
	isRunning   bool
	isRepeating bool
	onComplete  func()
}

func NewTimer(duration float32) *Timer {
	return &Timer{
		duration:  duration,
		remaining: duration,
		isRunning: false,
	}
}

func NewRepeatingTimer(duration float32) *Timer {
	return &Timer{
		duration:    duration,
		remaining:   duration,
		isRunning:   false,
		isRepeating: true,
	}
}

func (t *Timer) Start() {
	t.isRunning = true
}

func (t *Timer) Stop() {
	t.isRunning = false
}

func (t *Timer) Reset() {
	t.remaining = t.duration
}

func (t *Timer) Restart() {
	t.Reset()
	t.Start()
}

func (t *Timer) Update(deltaTime float32) {
	if !t.isRunning {
		return
	}

	t.remaining -= deltaTime

	if t.remaining <= 0 {
		if t.onComplete != nil {
			t.onComplete()
		}

		if t.isRepeating {
			t.remaining = t.duration
		} else {
			t.isRunning = false
			t.remaining = 0
		}
	}
}

func (t *Timer) SetOnComplete(callback func()) {
	t.onComplete = callback
}

func (t *Timer) IsRunning() bool {
	return t.isRunning
}

func (t *Timer) IsComplete() bool {
	return !t.isRunning && t.remaining <= 0
}

func (t *Timer) GetRemaining() float32 {
	return t.remaining
}

func (t *Timer) GetProgress() float32 {
	if t.duration <= 0 {
		return 1.0
	}
	return 1.0 - (t.remaining / t.duration)
}

func (t *Timer) GetDuration() float32 {
	return t.duration
}

func (t *Timer) SetDuration(duration float32) {
	t.duration = duration
	if t.remaining > duration {
		t.remaining = duration
	}
}

type Stopwatch struct {
	elapsed   float32
	isRunning bool
}

func NewStopwatch() *Stopwatch {
	return &Stopwatch{
		elapsed:   0,
		isRunning: false,
	}
}

func (s *Stopwatch) Start() {
	s.isRunning = true
}

func (s *Stopwatch) Stop() {
	s.isRunning = false
}

func (s *Stopwatch) Reset() {
	s.elapsed = 0
}

func (s *Stopwatch) Restart() {
	s.Reset()
	s.Start()
}

func (s *Stopwatch) Update(deltaTime float32) {
	if s.isRunning {
		s.elapsed += deltaTime
	}
}

func (s *Stopwatch) GetElapsed() float32 {
	return s.elapsed
}

func (s *Stopwatch) IsRunning() bool {
	return s.isRunning
}

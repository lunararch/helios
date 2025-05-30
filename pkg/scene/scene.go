package scene

import (
	"github.com/lunararch/helios/pkg/graphics/camera"
	"github.com/lunararch/helios/pkg/input"
)

type Scene interface {
	Load() error
	Unload() error

	Enter(prevScene Scene) error
	Exit(nextScene Scene) error

	Update(deltaTime float32) error
	Render(alpha float32) error

	HandleInput(inputManager *input.InputManager, inputMapping *input.InputMapping) error

	GetName() string

	Pause()
	Resume()
	IsPaused() bool

	IsLoaded() bool
}

type BaseScene struct {
	name   string
	loaded bool
	paused bool
	camera *camera.Camera
}

func NewBaseScene(name string, camera *camera.Camera) *BaseScene {
	return &BaseScene{
		name:   name,
		loaded: false,
		paused: false,
		camera: camera,
	}
}

func (s *BaseScene) GetName() string {
	return s.name
}

func (s *BaseScene) IsLoaded() bool {
	return s.loaded
}

func (s *BaseScene) SetLoaded(loaded bool) {
	s.loaded = loaded
}

func (s *BaseScene) Pause() {
	s.paused = true
}

func (s *BaseScene) Resume() {
	s.paused = false
}

func (s *BaseScene) IsPaused() bool {
	return s.paused
}

func (s *BaseScene) GetCamera() *camera.Camera {
	return s.camera
}

func (s *BaseScene) Load() error {
	s.loaded = true
	return nil
}

func (s *BaseScene) Unload() error {
	s.loaded = false
	return nil
}

func (s *BaseScene) Enter(prevScene Scene) error {
	s.paused = false
	return nil
}

func (s *BaseScene) Exit(nextScene Scene) error {
	return nil
}

func (s *BaseScene) Update(deltaTime float32) error {
	if s.paused {
		return nil
	}

	if s.camera != nil {
		s.camera.Update(deltaTime)
	}

	return nil
}

func (s *BaseScene) Render(alpha float32) error {
	return nil
}

func (s *BaseScene) HandleInput(inputManager *input.InputManager, inputMapping *input.InputMapping) error {
	return nil
}

package scene

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/lunararch/helios/pkg/graphics/camera"
	"github.com/lunararch/helios/pkg/input"
)

type MenuScene struct {
	*BaseScene
}

func NewMenuScene(camera *camera.Camera) *MenuScene {
	return &MenuScene{
		BaseScene: NewBaseScene("menu", camera),
	}
}

func (s *MenuScene) Load() error {
	if err := s.BaseScene.Load(); err != nil {
		return err
	}

	println("Menu scene loaded")
	return nil
}

func (s *MenuScene) Unload() error {
	println("Menu scene unloaded")
	return s.BaseScene.Unload()
}

func (s *MenuScene) Enter(prevScene Scene) error {
	if err := s.BaseScene.Enter(prevScene); err != nil {
		return err
	}

	println("Entered menu scene")
	return nil
}

func (s *MenuScene) Exit(nextScene Scene) error {
	println("Exiting menu scene")
	return s.BaseScene.Exit(nextScene)
}

func (s *MenuScene) Render(alpha float32) error {
	gl.ClearColor(0.1, 0.1, 0.2, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Here you would render menu UI elements

	return nil
}

func (s *MenuScene) HandleInput(inputManager *input.InputManager, inputMapping *input.InputMapping) error {
	if inputManager.IsKeyPressed(glfw.KeyEnter) {
		println("Enter pressed - should switch to gameplay")
	}

	return nil
}

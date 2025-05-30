package scene

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/lunararch/helios/pkg/engine"
	"github.com/lunararch/helios/pkg/graphics/camera"
	"github.com/lunararch/helios/pkg/graphics/shader"
	"github.com/lunararch/helios/pkg/graphics/sprite"
	"github.com/lunararch/helios/pkg/graphics/texture"
	"github.com/lunararch/helios/pkg/input"
)

type GameplayScene struct {
	*BaseScene

	batchShader *shader.Shader
	spriteBatch *sprite.SpriteBatch

	knightSprite *sprite.Sprite
	hornetSprite *sprite.Sprite

	knightTexture *texture.Texture
	hornetTexture *texture.Texture

	rotationTimer *engine.Timer
	printTimer    *engine.Timer

	cameraSpeed float32
}

func NewGameplayScene(camera *camera.Camera) *GameplayScene {
	return &GameplayScene{
		BaseScene:   NewBaseScene("gameplay", camera),
		cameraSpeed: 200.0,
	}
}

func (s *GameplayScene) Load() error {
	if err := s.BaseScene.Load(); err != nil {
		return err
	}

	var err error
	s.batchShader, err = shader.New("assets/shaders/batch.vert", "assets/shaders/batch.frag")
	if err != nil {
		return err
	}

	projection := mgl32.Ortho(0, s.camera.Size.X(), s.camera.Size.Y(), 0, -1, 1)
	s.batchShader.Use()
	s.batchShader.SetInt("texture1", 0)
	s.batchShader.SetMat4("projection", projection)

	s.spriteBatch = sprite.NewSpriteBatch(s.batchShader)

	s.knightTexture, err = texture.LoadFromFile("assets/textures/knight.png")
	if err != nil {
		return err
	}

	s.hornetTexture, err = texture.LoadFromFile("assets/textures/hornet.png")
	if err != nil {
		return err
	}

	s.knightSprite = sprite.NewSprite(
		s.knightTexture,
		mgl32.Vec3{100.0, 100.0, 0.0},
		mgl32.Vec2{float32(s.knightTexture.Width), float32(s.knightTexture.Height)},
	)

	s.hornetSprite = sprite.NewSprite(
		s.hornetTexture,
		mgl32.Vec3{300.0, 200.0, 0.0},
		mgl32.Vec2{float32(s.hornetTexture.Width), float32(s.hornetTexture.Height)},
	)

	s.rotationTimer = engine.NewRepeatingTimer(2.0)
	s.rotationTimer.SetOnComplete(func() {
		// Could implement rotation direction reversal here
	})
	s.rotationTimer.Start()

	s.printTimer = engine.NewTimer(5.0)
	s.printTimer.SetOnComplete(func() {
		println("=== Gameplay Scene Stats ===")
		println("Knight position:", s.knightSprite.Position.X(), s.knightSprite.Position.Y())
		println("Camera position:", s.camera.Position.X(), s.camera.Position.Y())
		println("Camera zoom:", s.camera.Zoom)
		println("============================")
		s.printTimer.Restart()
	})
	s.printTimer.Start()

	return nil
}

func (s *GameplayScene) Unload() error {
	if s.batchShader != nil {
		s.batchShader.Delete()
	}
	if s.spriteBatch != nil {
		s.spriteBatch.Delete()
	}
	if s.knightTexture != nil {
		s.knightTexture.Delete()
	}
	if s.hornetTexture != nil {
		s.hornetTexture.Delete()
	}

	return s.BaseScene.Unload()
}

func (s *GameplayScene) Update(deltaTime float32) error {
	if err := s.BaseScene.Update(deltaTime); err != nil {
		return err
	}

	if s.paused {
		return nil
	}

	s.rotationTimer.Update(deltaTime)
	s.printTimer.Update(deltaTime)

	s.knightSprite.Rotation += deltaTime * 1.0

	return nil
}

func (s *GameplayScene) Render(alpha float32) error {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	s.batchShader.Use()
	s.batchShader.SetMat4("view", s.camera.GetViewMatrix())

	s.spriteBatch.Begin()
	s.spriteBatch.Draw(s.knightSprite)
	s.spriteBatch.Draw(s.hornetSprite)
	s.spriteBatch.End()

	return nil
}

func (s *GameplayScene) HandleInput(inputManager *input.InputManager, inputMapping *input.InputMapping) error {
	if inputMapping.IsActionHeld("move_up", inputManager) {
		s.camera.Position[1] -= s.cameraSpeed * inputManager.GetDeltaTime()
	}
	if inputMapping.IsActionHeld("move_down", inputManager) {
		s.camera.Position[1] += s.cameraSpeed * inputManager.GetDeltaTime()
	}
	if inputMapping.IsActionHeld("move_left", inputManager) {
		s.camera.Position[0] -= s.cameraSpeed * inputManager.GetDeltaTime()
	}
	if inputMapping.IsActionHeld("move_right", inputManager) {
		s.camera.Position[0] += s.cameraSpeed * inputManager.GetDeltaTime()
	}

	deltaTime := inputManager.GetDeltaTime()
	if inputMapping.IsActionHeld("zoom_out", inputManager) {
		s.camera.Zoom *= (1.0 - deltaTime)
		if s.camera.Zoom < 0.1 {
			s.camera.Zoom = 0.1
		}
	}
	if inputMapping.IsActionHeld("zoom_in", inputManager) {
		s.camera.Zoom *= (1.0 + deltaTime)
		if s.camera.Zoom > 10.0 {
			s.camera.Zoom = 10.0
		}
	}

	scrollDelta := inputManager.GetScrollDelta()
	if scrollDelta.Y() != 0 {
		zoomFactor := 1.0 + scrollDelta.Y()*0.1
		s.camera.Zoom *= zoomFactor
		if s.camera.Zoom < 0.1 {
			s.camera.Zoom = 0.1
		}
		if s.camera.Zoom > 10.0 {
			s.camera.Zoom = 10.0
		}
	}

	s.camera.ClampToBounds()

	return nil
}

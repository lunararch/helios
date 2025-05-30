package scene

import (
	"math"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/lunararch/helios/pkg/engine"
	"github.com/lunararch/helios/pkg/entity"
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

	world *entity.World

	knightTexture *texture.Texture
	hornetTexture *texture.Texture

	knightEntity *entity.Entity
	hornetEntity *entity.Entity

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

	s.world = entity.NewWorld()

	s.knightEntity = s.world.CreateEntity("Knight")
	s.knightEntity.GetTransform().SetPosition2D(100.0, 100.0)

	knightSprite := entity.NewSpriteComponent(s.knightTexture, s.spriteBatch)
	s.knightEntity.AddComponent(knightSprite)

	s.hornetEntity = s.world.CreateEntity("Hornet")
	s.hornetEntity.GetTransform().SetPosition2D(300.0, 200.0)

	hornetSprite := entity.NewSpriteComponent(s.hornetTexture, s.spriteBatch)
	s.hornetEntity.AddComponent(hornetSprite)

	weaponEntity := s.world.CreateEntity("Knight's Weapon")
	weaponEntity.SetParent(s.knightEntity)
	weaponEntity.GetTransform().SetPosition2D(50.0, 0.0) // Relative to knight
	weaponEntity.GetTransform().SetUniformScale(0.5)

	weaponSprite := entity.NewSpriteComponent(s.hornetTexture, s.spriteBatch)
	weaponSprite.SetColor(mgl32.Vec4{1.0, 0.5, 0.5, 1.0}) // Reddish tint
	weaponEntity.AddComponent(weaponSprite)

	s.rotationTimer = engine.NewRepeatingTimer(2.0)
	s.rotationTimer.SetOnComplete(func() {
		s.knightEntity.GetTransform().Rotate(0.5)
	})
	s.rotationTimer.Start()

	s.printTimer = engine.NewTimer(5.0)
	s.printTimer.SetOnComplete(func() {
		println("=== Gameplay Scene Stats ===")
		println("Entity count:", s.world.GetEntityCount())
		println("Knight position:", s.knightEntity.GetTransform().Position.X(), s.knightEntity.GetTransform().Position.Y())
		println("Knight world position:", s.knightEntity.GetWorldPosition().X(), s.knightEntity.GetWorldPosition().Y())
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
	if s.world != nil {
		s.world.Cleanup()
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

	s.world.Update(deltaTime)

	if s.hornetEntity != nil {
		time := s.printTimer.GetProgress() * 6.28 // 2 * PI
		transform := s.hornetEntity.GetTransform()
		baseX := float32(300.0)
		baseY := float32(200.0)
		radius := float32(50.0)

		newX := baseX + radius*float32(math.Cos(float64(time)))
		newY := baseY + radius*float32(math.Sin(float64(time)))

		transform.SetPosition2D(newX, newY)
	}

	return nil
}

func (s *GameplayScene) Render(alpha float32) error {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	s.batchShader.Use()
	s.batchShader.SetMat4("view", s.camera.GetViewMatrix())

	s.spriteBatch.Begin()

	s.world.Render(alpha)

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

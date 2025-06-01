package scene

import (
	"math"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/lunararch/helios/pkg/engine"
	"github.com/lunararch/helios/pkg/entity"
	"github.com/lunararch/helios/pkg/graphics/animation"
	"github.com/lunararch/helios/pkg/graphics/camera"
	"github.com/lunararch/helios/pkg/graphics/shader"
	"github.com/lunararch/helios/pkg/graphics/sprite"
	"github.com/lunararch/helios/pkg/graphics/texture"
	"github.com/lunararch/helios/pkg/input"
)

type AnimatedGameplayScene struct {
	*BaseScene

	batchShader *shader.Shader
	spriteBatch *sprite.SpriteBatch

	world *entity.World

	knightTexture  *texture.Texture
	hornetTexture  *texture.Texture
	characterSheet *texture.Texture

	knightEntity   *entity.Entity
	hornetEntity   *entity.Entity
	animatedEntity *entity.Entity

	characterSpriteSheet *animation.SpriteSheet
	idleAnimation        *animation.AnimationClip
	walkAnimation        *animation.AnimationClip
	jumpAnimation        *animation.AnimationClip

	rotationTimer *engine.Timer
	printTimer    *engine.Timer

	cameraSpeed float32
}

func NewAnimatedGameplayScene(camera *camera.Camera) *AnimatedGameplayScene {
	return &AnimatedGameplayScene{
		BaseScene:   NewBaseScene("animated_gameplay", camera),
		cameraSpeed: 200.0,
	}
}

func (s *AnimatedGameplayScene) Load() error {
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

	// Create a simple character sprite sheet (using existing textures as example)
	// In a real project, you'd have an actual sprite sheet
	s.characterSheet = s.knightTexture                                          // Using knight texture as example
	s.characterSpriteSheet = animation.NewSpriteSheet(s.characterSheet, 32, 32) // Assuming 32x32 frames

	s.createAnimations()

	s.world = entity.NewWorld()

	s.knightEntity = s.world.CreateEntity("Knight")
	s.knightEntity.GetTransform().SetPosition2D(100.0, 100.0)

	knightSprite := entity.NewSpriteComponent(s.knightTexture, s.spriteBatch)
	s.knightEntity.AddComponent(knightSprite)

	s.hornetEntity = s.world.CreateEntity("Hornet")
	s.hornetEntity.GetTransform().SetPosition2D(300.0, 200.0)

	hornetSprite := entity.NewSpriteComponent(s.hornetTexture, s.spriteBatch)
	s.hornetEntity.AddComponent(hornetSprite)

	s.animatedEntity = s.world.CreateEntity("Animated Character")
	s.animatedEntity.GetTransform().SetPosition2D(200.0, 150.0)

	animatedSprite := entity.NewSpriteComponent(s.characterSheet, s.spriteBatch)
	s.animatedEntity.AddComponent(animatedSprite)

	animationComp := entity.NewAnimationComponent(animatedSprite)
	s.setupAnimationStateMachine(animationComp)
	s.animatedEntity.AddComponent(animationComp)

	s.rotationTimer = engine.NewRepeatingTimer(2.0)
	s.rotationTimer.SetOnComplete(func() {
		s.knightEntity.GetTransform().Rotate(0.5)
	})
	s.rotationTimer.Start()

	s.printTimer = engine.NewTimer(5.0)
	s.printTimer.SetOnComplete(func() {
		println("=== Animated Gameplay Scene Stats ===")
		println("Entity count:", s.world.GetEntityCount())

		if animComp, ok := s.animatedEntity.GetComponent(entity.ComponentTypeAnimation); ok {
			animationComponent := animComp.(*entity.AnimationComponent)
			println("Current animation state:", animationComponent.GetCurrentStateName())
			println("Animation playing:", animationComponent.IsPlaying())
		}

		println("========================================")
		s.printTimer.Restart()
	})
	s.printTimer.Start()

	return nil
}

func (s *AnimatedGameplayScene) createAnimations() {
	// Create sample animations
	// Note: These are simplified examples using the same texture
	// In a real project, you'd have actual sprite sheet frames

	builder := animation.NewAnimationBuilder(s.characterSpriteSheet, "idle", true)
	s.idleAnimation = builder.
		AddFrame(0, 0.5). // Frame 0 for 0.5 seconds
		AddFrame(1, 0.5). // Frame 1 for 0.5 seconds
		Build()

	builder = animation.NewAnimationBuilder(s.characterSpriteSheet, "walk", true)
	s.walkAnimation = builder.
		AddFrameRange(2, 5, 0.2). // Frames 2-5, each for 0.2 seconds
		Build()

	builder = animation.NewAnimationBuilder(s.characterSpriteSheet, "jump", false)
	s.jumpAnimation = builder.
		AddFrame(6, 0.1). // Jump start
		AddFrame(7, 0.3). // In air
		AddFrame(8, 0.1). // Landing
		Build()
}

func (s *AnimatedGameplayScene) setupAnimationStateMachine(animComp *entity.AnimationComponent) {
	stateMachine := animComp.GetStateMachine()

	// Create states
	idleState := animation.NewAnimationState("idle", s.idleAnimation)
	walkState := animation.NewAnimationState("walk", s.walkAnimation)
	jumpState := animation.NewAnimationState("jump", s.jumpAnimation)

	// Add transitions
	idleState.AddTransition("walk", "walk", func(sm *animation.AnimationStateMachine) bool {
		return sm.GetBool("isWalking")
	})

	idleState.AddTransition("jump", "jump", func(sm *animation.AnimationStateMachine) bool {
		return sm.GetBool("isJumping")
	})

	walkState.AddTransition("idle", "idle", func(sm *animation.AnimationStateMachine) bool {
		return !sm.GetBool("isWalking")
	})

	walkState.AddTransition("jump", "jump", func(sm *animation.AnimationStateMachine) bool {
		return sm.GetBool("isJumping")
	})

	jumpState.AddTimedTransition("complete", "idle", jumpState.Clip.TotalTime, func(sm *animation.AnimationStateMachine) bool {
		return true
	})

	// Add states to state machine
	stateMachine.AddState(idleState)
	stateMachine.AddState(walkState)
	stateMachine.AddState(jumpState)

	stateMachine.SetState("idle")
}

func (s *AnimatedGameplayScene) Update(deltaTime float32) error {
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

	if s.animatedEntity != nil {
		if animComp, ok := s.animatedEntity.GetComponent(entity.ComponentTypeAnimation); ok {
			animationComponent := animComp.(*entity.AnimationComponent)
			stateMachine := animationComponent.GetStateMachine()

			time := s.printTimer.GetProgress()
			if time < 2.0 {
				stateMachine.SetParameter("isWalking", false)
				stateMachine.SetParameter("isJumping", false)
			} else if time < 4.0 {
				stateMachine.SetParameter("isWalking", true)
				stateMachine.SetParameter("isJumping", false)
			} else {
				stateMachine.SetParameter("isWalking", false)
				stateMachine.SetParameter("isJumping", true)
			}
		}
	}

	return nil
}

func (s *AnimatedGameplayScene) Render(alpha float32) error {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	s.batchShader.Use()
	s.batchShader.SetMat4("view", s.camera.GetViewMatrix())

	s.spriteBatch.Begin()
	s.world.Render(alpha)
	s.spriteBatch.End()

	return nil
}

func (s *AnimatedGameplayScene) HandleInput(inputManager *input.InputManager, inputMapping *input.InputMapping) error {
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

	if s.animatedEntity != nil {
		if animComp, ok := s.animatedEntity.GetComponent(entity.ComponentTypeAnimation); ok {
			animationComponent := animComp.(*entity.AnimationComponent)
			stateMachine := animationComponent.GetStateMachine()

			if inputManager.IsKeyPressed(glfw.Key1) {
				stateMachine.SetParameter("isWalking", false)
				stateMachine.SetParameter("isJumping", false)
			}
			if inputManager.IsKeyPressed(glfw.Key2) {
				stateMachine.SetParameter("isWalking", true)
				stateMachine.SetParameter("isJumping", false)
			}
			if inputManager.IsKeyPressed(glfw.Key3) {
				stateMachine.SetParameter("isWalking", false)
				stateMachine.SetParameter("isJumping", true)
			}
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

func (s *AnimatedGameplayScene) Unload() error {
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
	if s.characterSheet != nil {
		s.characterSheet.Delete()
	}
	if s.world != nil {
		s.world.Cleanup()
	}

	return s.BaseScene.Unload()
}

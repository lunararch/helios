package entity

import (
	"github.com/lunararch/helios/pkg/graphics/animation"
	"github.com/lunararch/helios/pkg/graphics/texture"
)

type AnimationComponent struct {
	*BaseComponent
	stateMachine  *animation.AnimationStateMachine
	spriteComp    *SpriteComponent
	currentFrame  *animation.Frame
	defaultRegion *texture.TextureRegion
}

func NewAnimationComponent(spriteComponent *SpriteComponent) *AnimationComponent {
	comp := &AnimationComponent{
		BaseComponent: NewBaseComponent(ComponentTypeAnimation),
		stateMachine:  animation.NewAnimationStateMachine(),
		spriteComp:    spriteComponent,
	}

	if spriteComponent != nil && spriteComponent.GetSprite() != nil {
		sprite := spriteComponent.GetSprite()
		comp.defaultRegion = sprite.Region
	}

	return comp
}

func (ac *AnimationComponent) Initialize() error {
	if err := ac.BaseComponent.Initialize(); err != nil {
		return err
	}

	ac.updateCurrentFrame()
	return nil
}

func (ac *AnimationComponent) Update(deltaTime float32) {
	if !ac.active || ac.stateMachine == nil {
		return
	}

	ac.stateMachine.Update(deltaTime)
	ac.updateCurrentFrame()
}

func (ac *AnimationComponent) updateCurrentFrame() {
	if ac.spriteComp == nil || ac.spriteComp.GetSprite() == nil {
		return
	}

	sprite := ac.spriteComp.GetSprite()

	if frame, err := ac.stateMachine.GetCurrentFrame(); err == nil {
		ac.currentFrame = frame
		sprite.Region = frame.TextureRegion

		if frame.TextureRegion != nil && frame.TextureRegion.Texture != nil {
			sprite.Texture = frame.TextureRegion.Texture
		}
	} else {
		sprite.Region = ac.defaultRegion
	}
}

func (ac *AnimationComponent) GetStateMachine() *animation.AnimationStateMachine {
	return ac.stateMachine
}

func (ac *AnimationComponent) AddState(state *animation.AnimationState) {
	ac.stateMachine.AddState(state)
}

func (ac *AnimationComponent) SetState(stateName string) error {
	return ac.stateMachine.SetState(stateName)
}

func (ac *AnimationComponent) SetTrigger(triggerName string) {
	ac.stateMachine.SetTrigger(triggerName)
}

func (ac *AnimationComponent) SetParameter(name string, value interface{}) {
	ac.stateMachine.SetParameter(name, value)
}

func (ac *AnimationComponent) Play() {
	ac.stateMachine.Play()
}

func (ac *AnimationComponent) Pause() {
	ac.stateMachine.Pause()
}

func (ac *AnimationComponent) Stop() {
	ac.stateMachine.Stop()
}

func (ac *AnimationComponent) IsPlaying() bool {
	return ac.stateMachine.IsPlaying()
}

func (ac *AnimationComponent) GetCurrentStateName() string {
	return ac.stateMachine.GetCurrentStateName()
}

func (ac *AnimationComponent) GetCurrentFrame() *animation.Frame {
	return ac.currentFrame
}

func (ac *AnimationComponent) SetSpriteComponent(spriteComp *SpriteComponent) {
	ac.spriteComp = spriteComp

	if spriteComp != nil && spriteComp.GetSprite() != nil {
		sprite := spriteComp.GetSprite()
		ac.defaultRegion = sprite.Region
	}
}

func (ac *AnimationComponent) Cleanup() {
	ac.stateMachine = nil
	ac.spriteComp = nil
	ac.currentFrame = nil
	ac.defaultRegion = nil
	ac.BaseComponent.Cleanup()
}

package entity

import (
	"github.com/go-gl/mathgl/mgl32"
)

type ScriptComponent struct {
	*BaseComponent
	script Script
}

type Script interface {
	Start(entity *Entity)
	Update(entity *Entity, deltaTime float32)
	Stop(entity *Entity)
}

func NewScriptComponent(script Script) *ScriptComponent {
	return &ScriptComponent{
		BaseComponent: NewBaseComponent(ComponentTypeScript),
		script:        script,
	}
}

func (sc *ScriptComponent) Initialize() error {
	if err := sc.BaseComponent.Initialize(); err != nil {
		return err
	}

	if sc.script != nil && sc.entity != nil {
		sc.script.Start(sc.entity)
	}

	return nil
}

func (sc *ScriptComponent) Update(deltaTime float32) {
	if !sc.active || sc.script == nil || sc.entity == nil {
		return
	}

	sc.script.Update(sc.entity, deltaTime)
}

func (sc *ScriptComponent) Cleanup() {
	if sc.script != nil && sc.entity != nil {
		sc.script.Stop(sc.entity)
	}
	sc.BaseComponent.Cleanup()
}

func (sc *ScriptComponent) GetScript() Script {
	return sc.script
}

func (sc *ScriptComponent) SetScript(script Script) {
	if sc.script != nil && sc.entity != nil {
		sc.script.Stop(sc.entity)
	}

	sc.script = script

	if sc.initialized && sc.script != nil && sc.entity != nil {
		sc.script.Start(sc.entity)
	}
}

type MovementScript struct {
	speed     float32
	direction mgl32.Vec2
}

func NewMovementScript(speed float32, direction mgl32.Vec2) *MovementScript {
	return &MovementScript{
		speed:     speed,
		direction: direction.Normalize(),
	}
}

func (ms *MovementScript) Start(entity *Entity) {
	// Called when script starts
}

func (ms *MovementScript) Update(entity *Entity, deltaTime float32) {
	offset := ms.direction.Mul(ms.speed * deltaTime)
	entity.GetTransform().Translate2D(offset.X(), offset.Y())
}

func (ms *MovementScript) Stop(entity *Entity) {
	// Called when script stops
}

func (ms *MovementScript) SetSpeed(speed float32) {
	ms.speed = speed
}

func (ms *MovementScript) SetDirection(direction mgl32.Vec2) {
	ms.direction = direction.Normalize()
}

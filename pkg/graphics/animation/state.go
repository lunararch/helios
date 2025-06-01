package animation

import (
	"fmt"
)

type AnimationState struct {
	Name     string
	Clip     *AnimationClip
	Speed    float32
	Triggers map[string]*AnimationTransition
}

func NewAnimationState(name string, clip *AnimationClip) *AnimationState {
	return &AnimationState{
		Name:     name,
		Clip:     clip,
		Speed:    1.0,
		Triggers: make(map[string]*AnimationTransition),
	}
}

func (as *AnimationState) AddTransition(triggerName string, targetState string, condition TransitionCondition) {
	transition := &AnimationTransition{
		FromState:   as.Name,
		ToState:     targetState,
		Trigger:     triggerName,
		Condition:   condition,
		HasExitTime: false,
		ExitTime:    0,
	}
	as.Triggers[triggerName] = transition
}

func (as *AnimationState) AddTimedTransition(triggerName string, targetState string, exitTime float32, condition TransitionCondition) {
	transition := &AnimationTransition{
		FromState:   as.Name,
		ToState:     targetState,
		Trigger:     triggerName,
		Condition:   condition,
		HasExitTime: true,
		ExitTime:    exitTime,
	}
	as.Triggers[triggerName] = transition
}

func (as *AnimationState) SetSpeed(speed float32) {
	as.Speed = speed
}

type TransitionCondition func(stateMachine *AnimationStateMachine) bool

type AnimationTransition struct {
	FromState   string
	ToState     string
	Trigger     string
	Condition   TransitionCondition
	HasExitTime bool
	ExitTime    float32
}

func (at *AnimationTransition) CanTransition(stateMachine *AnimationStateMachine, currentTime float32) bool {
	if at.HasExitTime && currentTime < at.ExitTime {
		return false
	}

	if at.Condition != nil {
		return at.Condition(stateMachine)
	}

	return true
}

type AnimationStateMachine struct {
	States       map[string]*AnimationState
	CurrentState *AnimationState
	CurrentTime  float32
	Triggers     map[string]bool
	Parameters   map[string]interface{}
	Playing      bool
}

func NewAnimationStateMachine() *AnimationStateMachine {
	return &AnimationStateMachine{
		States:     make(map[string]*AnimationState),
		Triggers:   make(map[string]bool),
		Parameters: make(map[string]interface{}),
		Playing:    true,
	}
}

func (asm *AnimationStateMachine) AddState(state *AnimationState) {
	asm.States[state.Name] = state

	if asm.CurrentState == nil {
		asm.CurrentState = state
	}
}

func (asm *AnimationStateMachine) SetState(stateName string) error {
	state, exists := asm.States[stateName]
	if !exists {
		return fmt.Errorf("animation state '%s' not found", stateName)
	}

	asm.CurrentState = state
	asm.CurrentTime = 0
	return nil
}

func (asm *AnimationStateMachine) SetTrigger(triggerName string) {
	asm.Triggers[triggerName] = true
}

func (asm *AnimationStateMachine) ResetTrigger(triggerName string) {
	asm.Triggers[triggerName] = false
}

func (asm *AnimationStateMachine) SetParameter(name string, value interface{}) {
	asm.Parameters[name] = value
}

func (asm *AnimationStateMachine) GetParameter(name string) (interface{}, bool) {
	value, exists := asm.Parameters[name]
	return value, exists
}

func (asm *AnimationStateMachine) GetBool(name string) bool {
	if value, exists := asm.Parameters[name]; exists {
		if boolValue, ok := value.(bool); ok {
			return boolValue
		}
	}
	return false
}

func (asm *AnimationStateMachine) GetFloat(name string) float32 {
	if value, exists := asm.Parameters[name]; exists {
		if floatValue, ok := value.(float32); ok {
			return floatValue
		}
	}
	return 0.0
}

func (asm *AnimationStateMachine) GetInt(name string) int {
	if value, exists := asm.Parameters[name]; exists {
		if intValue, ok := value.(int); ok {
			return intValue
		}
	}
	return 0
}

func (asm *AnimationStateMachine) Update(deltaTime float32) {
	if !asm.Playing || asm.CurrentState == nil {
		return
	}

	asm.CurrentTime += deltaTime * asm.CurrentState.Speed

	for triggerName, isSet := range asm.Triggers {
		if !isSet {
			continue
		}

		if transition, exists := asm.CurrentState.Triggers[triggerName]; exists {
			if transition.CanTransition(asm, asm.CurrentTime) {
				if targetState, exists := asm.States[transition.ToState]; exists {
					asm.CurrentState = targetState
					asm.CurrentTime = 0
					asm.ResetTrigger(triggerName) // Reset trigger after use
					break
				}
			}
		}
	}

	if asm.CurrentState.Clip != nil && !asm.CurrentState.Clip.Loop {
		if asm.CurrentTime >= asm.CurrentState.Clip.TotalTime {
			asm.CurrentTime = asm.CurrentState.Clip.TotalTime
		}
	}
}

func (asm *AnimationStateMachine) GetCurrentFrame() (*Frame, error) {
	if asm.CurrentState == nil || asm.CurrentState.Clip == nil {
		return nil, fmt.Errorf("no current state or clip")
	}

	return asm.CurrentState.Clip.GetFrameAt(asm.CurrentTime)
}

func (asm *AnimationStateMachine) GetCurrentClip() *AnimationClip {
	if asm.CurrentState == nil {
		return nil
	}
	return asm.CurrentState.Clip
}

func (asm *AnimationStateMachine) GetCurrentStateName() string {
	if asm.CurrentState == nil {
		return ""
	}
	return asm.CurrentState.Name
}

func (asm *AnimationStateMachine) IsPlaying() bool {
	return asm.Playing
}

func (asm *AnimationStateMachine) Play() {
	asm.Playing = true
}

func (asm *AnimationStateMachine) Pause() {
	asm.Playing = false
}

func (asm *AnimationStateMachine) Stop() {
	asm.Playing = false
	asm.CurrentTime = 0
}

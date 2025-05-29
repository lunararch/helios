package input

import "github.com/go-gl/glfw/v3.3/glfw"

type Action string

type InputMapping struct {
	keyMappings   map[Action][]glfw.Key
	mouseMappings map[Action][]MouseButton
}

func NewInputMapping() *InputMapping {
	return &InputMapping{
		keyMappings:   make(map[Action][]glfw.Key),
		mouseMappings: make(map[Action][]MouseButton),
	}
}

func (im *InputMapping) MapKey(action Action, key glfw.Key) {
	if _, exists := im.keyMappings[action]; !exists {
		im.keyMappings[action] = make([]glfw.Key, 0)
	}
	im.keyMappings[action] = append(im.keyMappings[action], key)
}

func (im *InputMapping) MapMouseButton(action Action, button MouseButton) {
	if _, exists := im.mouseMappings[action]; !exists {
		im.mouseMappings[action] = make([]MouseButton, 0)
	}
	im.mouseMappings[action] = append(im.mouseMappings[action], button)
}

func (im *InputMapping) IsActionPressed(action Action, inputMgr *InputManager) bool {
	if keys, exists := im.keyMappings[action]; exists {
		for _, key := range keys {
			if inputMgr.IsKeyPressed(key) {
				return true
			}
		}
	}

	if buttons, exists := im.mouseMappings[action]; exists {
		for _, button := range buttons {
			if inputMgr.IsMouseButtonPressed(button) {
				return true
			}
		}
	}

	return false
}

func (im *InputMapping) IsActionHeld(action Action, inputMgr *InputManager) bool {
	if keys, exists := im.keyMappings[action]; exists {
		for _, key := range keys {
			if inputMgr.IsKeyHeld(key) {
				return true
			}
		}
	}

	if buttons, exists := im.mouseMappings[action]; exists {
		for _, button := range buttons {
			if inputMgr.IsMouseButtonHeld(button) {
				return true
			}
		}
	}

	return false
}

func (im *InputMapping) IsActionReleased(action Action, inputMgr *InputManager) bool {
	if keys, exists := im.keyMappings[action]; exists {
		for _, key := range keys {
			if inputMgr.IsKeyReleased(key) {
				return true
			}
		}
	}

	if buttons, exists := im.mouseMappings[action]; exists {
		for _, button := range buttons {
			if inputMgr.IsMouseButtonReleased(button) {
				return true
			}
		}
	}

	return false
}

func (im *InputMapping) ClearAction(action Action) {
	delete(im.keyMappings, action)
	delete(im.mouseMappings, action)
}

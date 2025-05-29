package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type KeyState int

const (
	KeyStateReleased KeyState = iota
	KeyStatePressed
	KeyStateHeld
)

type MouseButton int

const (
	MouseButtonLeft MouseButton = iota
	MouseButtonRight
	MouseButtonMiddle
)

type InputManager struct {
	window          *glfw.Window
	keyStates       map[glfw.Key]KeyState
	prevKeyStates   map[glfw.Key]bool
	mouseStates     map[MouseButton]KeyState
	prevMouseStates map[MouseButton]bool
	mousePosition   mgl32.Vec2
	prevMousePos    mgl32.Vec2
	mouseDelta      mgl32.Vec2
	scrollDelta     mgl32.Vec2
	inputCallbacks  []InputCallback
}

type InputCallback func(event InputEvent)

type InputEvent interface {
	GetType() InputEventType
}

type InputEventType int

const (
	EventTypeKeyPress InputEventType = iota
	EventTypeKeyRelease
	EventTypeMousePress
	EventTypeMouseRelease
	EventTypeMouseMove
	EventTypeMouseScroll
)

type KeyPressEvent struct {
	Key glfw.Key
}

func (e KeyPressEvent) GetType() InputEventType { return EventTypeKeyPress }

type KeyReleaseEvent struct {
	Key glfw.Key
}

func (e KeyReleaseEvent) GetType() InputEventType { return EventTypeKeyRelease }

type MousePressEvent struct {
	Button   MouseButton
	Position mgl32.Vec2
}

func (e MousePressEvent) GetType() InputEventType { return EventTypeMousePress }

type MouseReleaseEvent struct {
	Button   MouseButton
	Position mgl32.Vec2
}

func (e MouseReleaseEvent) GetType() InputEventType { return EventTypeMouseRelease }

type MouseMoveEvent struct {
	Position mgl32.Vec2
	Delta    mgl32.Vec2
}

func (e MouseMoveEvent) GetType() InputEventType { return EventTypeMouseMove }

type MouseScrollEvent struct {
	Delta mgl32.Vec2
}

func (e MouseScrollEvent) GetType() InputEventType { return EventTypeMouseScroll }

func NewInputManager(window *glfw.Window) *InputManager {
	im := &InputManager{
		window:          window,
		keyStates:       make(map[glfw.Key]KeyState),
		prevKeyStates:   make(map[glfw.Key]bool),
		mouseStates:     make(map[MouseButton]KeyState),
		prevMouseStates: make(map[MouseButton]bool),
		inputCallbacks:  make([]InputCallback, 0),
	}

	im.setupCallbacks()
	return im
}

func (im *InputManager) setupCallbacks() {
	im.window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		switch action {
		case glfw.Press:
			for _, callback := range im.inputCallbacks {
				callback(KeyPressEvent{Key: key})
			}
		case glfw.Release:
			for _, callback := range im.inputCallbacks {
				callback(KeyReleaseEvent{Key: key})
			}
		}
	})

	im.window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		mouseBtn := MouseButton(button)
		pos := im.GetMousePosition()

		switch action {
		case glfw.Press:
			for _, callback := range im.inputCallbacks {
				callback(MousePressEvent{Button: mouseBtn, Position: pos})
			}
		case glfw.Release:
			for _, callback := range im.inputCallbacks {
				callback(MouseReleaseEvent{Button: mouseBtn, Position: pos})
			}
		}
	})

	im.window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		newPos := mgl32.Vec2{float32(xpos), float32(ypos)}
		delta := newPos.Sub(im.mousePosition)
		im.mouseDelta = delta

		for _, callback := range im.inputCallbacks {
			callback(MouseMoveEvent{Position: newPos, Delta: delta})
		}
	})

	im.window.SetScrollCallback(func(w *glfw.Window, xoffset, yoffset float64) {
		delta := mgl32.Vec2{float32(xoffset), float32(yoffset)}
		im.scrollDelta = delta

		for _, callback := range im.inputCallbacks {
			callback(MouseScrollEvent{Delta: delta})
		}
	})
}

func (im *InputManager) Update() {
	for key := range im.keyStates {
		currentPressed := im.window.GetKey(key) == glfw.Press
		wasPressed := im.prevKeyStates[key]

		if currentPressed && !wasPressed {
			im.keyStates[key] = KeyStatePressed
		} else if currentPressed && wasPressed {
			im.keyStates[key] = KeyStateHeld
		} else {
			im.keyStates[key] = KeyStateReleased
		}

		im.prevKeyStates[key] = currentPressed
	}

	// Update mouse button states
	for button := range im.mouseStates {
		currentPressed := im.window.GetMouseButton(glfw.MouseButton(button)) == glfw.Press
		wasPressed := im.prevMouseStates[button]

		if currentPressed && !wasPressed {
			im.mouseStates[button] = KeyStatePressed
		} else if currentPressed && wasPressed {
			im.mouseStates[button] = KeyStateHeld
		} else {
			im.mouseStates[button] = KeyStateReleased
		}

		im.prevMouseStates[button] = currentPressed
	}

	x, y := im.window.GetCursorPos()
	im.prevMousePos = im.mousePosition
	im.mousePosition = mgl32.Vec2{float32(x), float32(y)}

	im.scrollDelta = mgl32.Vec2{0, 0}
}

func (im *InputManager) IsKeyPressed(key glfw.Key) bool {
	if state, exists := im.keyStates[key]; exists {
		return state == KeyStatePressed
	}
	im.keyStates[key] = KeyStateReleased
	im.prevKeyStates[key] = false
	return false
}

func (im *InputManager) IsKeyHeld(key glfw.Key) bool {
	if state, exists := im.keyStates[key]; exists {
		return state == KeyStateHeld || state == KeyStatePressed
	}
	im.keyStates[key] = KeyStateReleased
	im.prevKeyStates[key] = false
	return false
}

func (im *InputManager) IsKeyReleased(key glfw.Key) bool {
	if state, exists := im.keyStates[key]; exists {
		return state == KeyStateReleased
	}
	im.keyStates[key] = KeyStateReleased
	im.prevKeyStates[key] = false
	return false
}

func (im *InputManager) IsMouseButtonPressed(button MouseButton) bool {
	if state, exists := im.mouseStates[button]; exists {
		return state == KeyStatePressed
	}
	im.mouseStates[button] = KeyStateReleased
	im.prevMouseStates[button] = false
	return false
}

func (im *InputManager) IsMouseButtonHeld(button MouseButton) bool {
	if state, exists := im.mouseStates[button]; exists {
		return state == KeyStateHeld || state == KeyStatePressed
	}
	im.mouseStates[button] = KeyStateReleased
	im.prevMouseStates[button] = false
	return false
}

func (im *InputManager) IsMouseButtonReleased(button MouseButton) bool {
	if state, exists := im.mouseStates[button]; exists {
		return state == KeyStateReleased
	}
	im.mouseStates[button] = KeyStateReleased
	im.prevMouseStates[button] = false
	return false
}

func (im *InputManager) GetMousePosition() mgl32.Vec2 {
	return im.mousePosition
}

func (im *InputManager) GetMouseDelta() mgl32.Vec2 {
	return im.mouseDelta
}

func (im *InputManager) GetScrollDelta() mgl32.Vec2 {
	return im.scrollDelta
}

func (im *InputManager) AddInputCallback(callback InputCallback) {
	im.inputCallbacks = append(im.inputCallbacks, callback)
}

func (im *InputManager) IsShiftHeld() bool {
	return im.IsKeyHeld(glfw.KeyLeftShift) || im.IsKeyHeld(glfw.KeyRightShift)
}

func (im *InputManager) IsCtrlHeld() bool {
	return im.IsKeyHeld(glfw.KeyLeftControl) || im.IsKeyHeld(glfw.KeyRightControl)
}

func (im *InputManager) IsAltHeld() bool {
	return im.IsKeyHeld(glfw.KeyLeftAlt) || im.IsKeyHeld(glfw.KeyRightAlt)
}

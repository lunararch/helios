package scene

import (
	"fmt"

	"github.com/lunararch/helios/pkg/input"
)

type SceneManager struct {
	scenes        map[string]Scene
	currentScene  Scene
	nextScene     Scene
	sceneStack    []Scene
	transitioning bool
}

func NewSceneManager() *SceneManager {
	return &SceneManager{
		scenes:        make(map[string]Scene),
		sceneStack:    make([]Scene, 0),
		transitioning: false,
	}
}

func (sm *SceneManager) RegisterScene(scene Scene) error {
	name := scene.GetName()
	if _, exists := sm.scenes[name]; exists {
		return fmt.Errorf("scene '%s' already registered", name)
	}

	sm.scenes[name] = scene
	return nil
}

func (sm *SceneManager) UnregisterScene(name string) error {
	scene, exists := sm.scenes[name]
	if !exists {
		return fmt.Errorf("scene '%s' not found", name)
	}

	// Unload scene if it was loaded
	if scene.IsLoaded() {
		scene.Unload()
	}

	delete(sm.scenes, name)
	return nil
}

func (sm *SceneManager) SwitchToScene(sceneName string) error {
	scene, exists := sm.scenes[sceneName]
	if !exists {
		return fmt.Errorf("scene '%s' not found", sceneName)
	}

	sm.nextScene = scene
	sm.transitioning = true
	return nil
}

func (sm *SceneManager) PushScene(sceneName string) error {
	scene, exists := sm.scenes[sceneName]
	if !exists {
		return fmt.Errorf("scene '%s' not found", sceneName)
	}

	if sm.currentScene != nil {
		sm.sceneStack = append(sm.sceneStack, sm.currentScene)
		sm.currentScene.Pause()
	}

	sm.nextScene = scene
	sm.transitioning = true
	return nil
}

func (sm *SceneManager) PopScene() error {
	if len(sm.sceneStack) == 0 {
		return fmt.Errorf("no scenes in stack to pop")
	}

	prevScene := sm.sceneStack[len(sm.sceneStack)-1]
	sm.sceneStack = sm.sceneStack[:len(sm.sceneStack)-1]

	sm.nextScene = prevScene
	sm.transitioning = true
	return nil
}

func (sm *SceneManager) Update(deltaTime float32) error {
	if sm.transitioning {
		if err := sm.performTransition(); err != nil {
			return err
		}
		sm.transitioning = false
	}

	if sm.currentScene != nil {
		return sm.currentScene.Update(deltaTime)
	}

	return nil
}

func (sm *SceneManager) Render(alpha float32) error {
	if sm.currentScene != nil {
		return sm.currentScene.Render(alpha)
	}
	return nil
}

func (sm *SceneManager) HandleInput(inputManager *input.InputManager, inputMapping *input.InputMapping) error {
	if sm.currentScene != nil {
		return sm.currentScene.HandleInput(inputManager, inputMapping)
	}
	return nil
}

func (sm *SceneManager) GetCurrentScene() Scene {
	return sm.currentScene
}

func (sm *SceneManager) GetSceneCount() int {
	return len(sm.scenes)
}

func (sm *SceneManager) IsTransitioning() bool {
	return sm.transitioning
}

func (sm *SceneManager) performTransition() error {
	if sm.nextScene == nil {
		return fmt.Errorf("no next scene set for transition")
	}

	prevScene := sm.currentScene

	if sm.currentScene != nil {
		if err := sm.currentScene.Exit(sm.nextScene); err != nil {
			return fmt.Errorf("failed to exit scene '%s': %w", sm.currentScene.GetName(), err)
		}
	}

	sm.currentScene = sm.nextScene
	sm.nextScene = nil

	if !sm.currentScene.IsLoaded() {
		if err := sm.currentScene.Load(); err != nil {
			return fmt.Errorf("failed to load scene '%s': %w", sm.currentScene.GetName(), err)
		}
	}

	if err := sm.currentScene.Enter(prevScene); err != nil {
		return fmt.Errorf("failed to enter scene '%s': %w", sm.currentScene.GetName(), err)
	}

	sm.currentScene.Resume()

	return nil
}

func (sm *SceneManager) Cleanup() error {
	if sm.currentScene != nil {
		sm.currentScene.Exit(nil)
	}

	for _, scene := range sm.scenes {
		if scene.IsLoaded() {
			if err := scene.Unload(); err != nil {
				return fmt.Errorf("failed to unload scene '%s': %w", scene.GetName(), err)
			}
		}
	}

	sm.scenes = make(map[string]Scene)
	sm.currentScene = nil
	sm.nextScene = nil
	sm.sceneStack = sm.sceneStack[:0]
	sm.transitioning = false

	return nil
}

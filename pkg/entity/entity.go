package entity

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

type EntityID uint64

type Entity struct {
	ID         EntityID
	name       string
	active     bool
	components map[ComponentType]Component
	transform  *Transform
	parent     *Entity
	children   []*Entity
	world      *World
	destroying bool // Add this flag to prevent circular destruction
}

func NewEntity(id EntityID, name string) *Entity {
	entity := &Entity{
		ID:         id,
		name:       name,
		active:     true,
		components: make(map[ComponentType]Component),
		children:   make([]*Entity, 0),
		destroying: false,
	}

	transform := NewTransform()
	entity.transform = transform
	entity.AddComponent(transform)

	return entity
}

func (e *Entity) GetName() string {
	return e.name
}

func (e *Entity) SetName(name string) {
	e.name = name
}

func (e *Entity) IsActive() bool {
	return e.active
}

func (e *Entity) SetActive(active bool) {
	e.active = active
}

func (e *Entity) GetTransform() *Transform {
	return e.transform
}

func (e *Entity) AddComponent(component Component) error {
	componentType := component.GetType()

	if _, exists := e.components[componentType]; exists {
		return fmt.Errorf("entity already has component of type %v", componentType)
	}

	e.components[componentType] = component
	component.SetEntity(e)

	if !component.IsInitialized() {
		component.Initialize()
	}

	return nil
}

func (e *Entity) RemoveComponent(componentType ComponentType) error {
	if componentType == ComponentTypeTransform {
		return fmt.Errorf("cannot remove transform component")
	}

	component, exists := e.components[componentType]
	if !exists {
		return fmt.Errorf("entity does not have component of type %v", componentType)
	}

	component.Cleanup()
	delete(e.components, componentType)

	return nil
}

func (e *Entity) GetComponent(componentType ComponentType) (Component, bool) {
	component, exists := e.components[componentType]
	return component, exists
}

func (e *Entity) HasComponent(componentType ComponentType) bool {
	_, exists := e.components[componentType]
	return exists
}

func (e *Entity) GetComponents() map[ComponentType]Component {
	return e.components
}

func (e *Entity) SetParent(parent *Entity) {
	if e.parent != nil {
		e.parent.RemoveChild(e)
	}

	e.parent = parent

	if parent != nil {
		parent.AddChild(e)
	}
}

func (e *Entity) GetParent() *Entity {
	return e.parent
}

func (e *Entity) AddChild(child *Entity) {
	for _, existingChild := range e.children {
		if existingChild.ID == child.ID {
			return
		}
	}

	e.children = append(e.children, child)
	child.parent = e
}

func (e *Entity) RemoveChild(child *Entity) {
	for i, existingChild := range e.children {
		if existingChild.ID == child.ID {
			e.children = append(e.children[:i], e.children[i+1:]...)
			child.parent = nil
			break
		}
	}
}

func (e *Entity) GetChildren() []*Entity {
	return e.children
}

func (e *Entity) GetWorldPosition() mgl32.Vec3 {
	if e.parent == nil {
		return e.transform.Position
	}

	parentWorld := e.parent.GetWorldPosition()
	return parentWorld.Add(e.transform.Position)
}

func (e *Entity) GetWorldRotation() float32 {
	if e.parent == nil {
		return e.transform.Rotation
	}

	parentRotation := e.parent.GetWorldRotation()
	return parentRotation + e.transform.Rotation
}

func (e *Entity) GetWorldScale() mgl32.Vec2 {
	if e.parent == nil {
		return e.transform.Scale
	}

	parentScale := e.parent.GetWorldScale()
	return mgl32.Vec2{
		parentScale.X() * e.transform.Scale.X(),
		parentScale.Y() * e.transform.Scale.Y(),
	}
}

func (e *Entity) Update(deltaTime float32) {
	if !e.active {
		return
	}

	for _, component := range e.components {
		if component.IsActive() {
			component.Update(deltaTime)
		}
	}

	for _, child := range e.children {
		child.Update(deltaTime)
	}
}

func (e *Entity) Render(alpha float32) {
	if !e.active {
		return
	}

	for _, component := range e.components {
		if renderable, ok := component.(RenderableComponent); ok && component.IsActive() {
			renderable.Render(alpha)
		}
	}

	for _, child := range e.children {
		child.Render(alpha)
	}
}

func (e *Entity) Destroy() {
	if e.destroying {
		return
	}
	e.destroying = true

	if e.parent != nil {
		e.parent.RemoveChild(e)
	}

	for _, child := range e.children {
		child.Destroy()
	}

	for _, component := range e.components {
		component.Cleanup()
	}

	e.components = nil
	e.children = nil
	e.parent = nil

	// Note: Don't call world.DestroyEntity here to avoid circular reference
	// The world will handle removing the entity from its maps
	e.world = nil
}

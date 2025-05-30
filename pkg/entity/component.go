package entity

type ComponentType int

const (
	ComponentTypeTransform ComponentType = iota
	ComponentTypeSprite
	ComponentTypeRigidbody
	ComponentTypeCollider
	ComponentTypeScript
	ComponentTypeAnimation
	ComponentTypeAudio
)

type Component interface {
	GetType() ComponentType
	IsActive() bool
	SetActive(active bool)
	IsInitialized() bool
	Initialize() error
	Update(deltaTime float32)
	Cleanup()
	SetEntity(entity *Entity)
	GetEntity() *Entity
}

type RenderableComponent interface {
	Component
	Render(alpha float32)
}

type BaseComponent struct {
	componentType ComponentType
	active        bool
	initialized   bool
	entity        *Entity
}

func NewBaseComponent(componentType ComponentType) *BaseComponent {
	return &BaseComponent{
		componentType: componentType,
		active:        true,
		initialized:   false,
	}
}

func (c *BaseComponent) GetType() ComponentType {
	return c.componentType
}

func (c *BaseComponent) IsActive() bool {
	return c.active
}

func (c *BaseComponent) SetActive(active bool) {
	c.active = active
}

func (c *BaseComponent) IsInitialized() bool {
	return c.initialized
}

func (c *BaseComponent) Initialize() error {
	c.initialized = true
	return nil
}

func (c *BaseComponent) Update(deltaTime float32) {
	// Default implementation does nothing
}

func (c *BaseComponent) Cleanup() {
	c.initialized = false
}

func (c *BaseComponent) SetEntity(entity *Entity) {
	c.entity = entity
}

func (c *BaseComponent) GetEntity() *Entity {
	return c.entity
}

package entity

import (
	"fmt"
	"sort"
)

type World struct {
	entities     map[EntityID]*Entity
	nextID       EntityID
	rootEntities []*Entity
}

func NewWorld() *World {
	return &World{
		entities:     make(map[EntityID]*Entity),
		nextID:       1,
		rootEntities: make([]*Entity, 0),
	}
}

func (w *World) CreateEntity(name string) *Entity {
	entity := NewEntity(w.nextID, name)
	w.nextID++

	w.entities[entity.ID] = entity
	entity.world = w

	w.rootEntities = append(w.rootEntities, entity)

	return entity
}

func (w *World) CreateEntityWithID(id EntityID, name string) (*Entity, error) {
	if _, exists := w.entities[id]; exists {
		return nil, fmt.Errorf("entity with ID %d already exists", id)
	}

	entity := NewEntity(id, name)
	w.entities[entity.ID] = entity
	entity.world = w

	if id >= w.nextID {
		w.nextID = id + 1
	}

	w.rootEntities = append(w.rootEntities, entity)

	return entity, nil
}

func (w *World) GetEntity(id EntityID) (*Entity, bool) {
	entity, exists := w.entities[id]
	return entity, exists
}

func (w *World) FindEntity(name string) *Entity {
	for _, entity := range w.entities {
		if entity.GetName() == name {
			return entity
		}
	}
	return nil
}

func (w *World) FindEntitiesWithName(name string) []*Entity {
	var found []*Entity
	for _, entity := range w.entities {
		if entity.GetName() == name {
			found = append(found, entity)
		}
	}
	return found
}

func (w *World) DestroyEntity(id EntityID) error {
	entity, exists := w.entities[id]
	if !exists {
		return fmt.Errorf("entity with ID %d not found", id)
	}

	for i, rootEntity := range w.rootEntities {
		if rootEntity.ID == id {
			w.rootEntities = append(w.rootEntities[:i], w.rootEntities[i+1:]...)
			break
		}
	}

	entity.Destroy()

	delete(w.entities, id)

	return nil
}

func (w *World) GetEntitiesWithComponent(componentType ComponentType) []*Entity {
	var entities []*Entity
	for _, entity := range w.entities {
		if entity.HasComponent(componentType) {
			entities = append(entities, entity)
		}
	}
	return entities
}

func (w *World) GetActiveEntities() []*Entity {
	var entities []*Entity
	for _, entity := range w.entities {
		if entity.IsActive() {
			entities = append(entities, entity)
		}
	}
	return entities
}

func (w *World) GetRootEntities() []*Entity {
	return w.rootEntities
}

func (w *World) GetEntityCount() int {
	return len(w.entities)
}

func (w *World) Update(deltaTime float32) {
	for _, entity := range w.rootEntities {
		entity.Update(deltaTime)
	}
}

func (w *World) Render(alpha float32) {
	var renderableEntities []*Entity
	for _, entity := range w.entities {
		if entity.IsActive() {
			for _, component := range entity.GetComponents() {
				if _, ok := component.(RenderableComponent); ok {
					renderableEntities = append(renderableEntities, entity)
					break
				}
			}
		}
	}

	sort.Slice(renderableEntities, func(i, j int) bool {
		entityA := renderableEntities[i]
		entityB := renderableEntities[j]

		spriteA, hasA := entityA.GetComponent(ComponentTypeSprite)
		spriteB, hasB := entityB.GetComponent(ComponentTypeSprite)

		if hasA && hasB {
			layerA := spriteA.(*SpriteComponent).GetLayer()
			layerB := spriteB.(*SpriteComponent).GetLayer()
			return layerA < layerB
		}

		return false
	})

	for _, entity := range renderableEntities {
		entity.Render(alpha)
	}
}

func (w *World) Clear() {
	entityIDs := make([]EntityID, 0, len(w.entities))
	for id := range w.entities {
		entityIDs = append(entityIDs, id)
	}

	for _, id := range entityIDs {
		if entity, exists := w.entities[id]; exists {
			entity.Destroy()
			delete(w.entities, id)
		}
	}

	w.rootEntities = make([]*Entity, 0)
	w.nextID = 1
}

func (w *World) Cleanup() {
	w.Clear()
}

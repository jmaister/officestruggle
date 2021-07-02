package ecs

type Entity struct {
	id         int
	components []*Component
}

type Component interface {
	ComponentType() string
}

type System interface {
	SystemType() string
}

type Engine struct {
	currentId  int
	components []*Component
	entities   []*Entity
}

/**
 * Engine
 */

func NewEngine() *Engine {
	return &Engine{
		currentId: 1,
	}
}

func (engine *Engine) RegisterComponent(component *Component) {
	engine.components = append(engine.components, component)
}

func (engine *Engine) NewEntity() *Entity {
	newEntity := &Entity{
		id: engine.currentId,
	}
	engine.currentId = engine.currentId + 1
	engine.entities = append(engine.entities, newEntity)
	return newEntity
}

/**
 * Entity
 */

func (entity *Entity) AddComponent(component Component) {
	entity.components = append(entity.components, &component)
}

func (entity *Entity) HasComponent(componentType string) bool {
	for _, cmp := range entity.components {
		if (*cmp).ComponentType() == componentType {
			return true
		}
	}
	return false
}

func (entity *Entity) GetComponent(componentType string) *Component {
	for _, cmp := range entity.components {
		if (*cmp).ComponentType() == componentType {
			return cmp
		}
	}
	return nil
}

func (entity *Entity) GetComponents(componentType string) []*Component {
	var found []*Component

	for _, cmp := range entity.components {
		if (*cmp).ComponentType() == componentType {
			found = append(found, cmp)
		}
	}
	return found
}

/**
 * Component
 */

/**
 * System
 */

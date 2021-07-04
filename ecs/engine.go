package ecs

import (
	"strconv"
)

type Entity struct {
	id         int
	components map[string]Component
}

type Component interface {
	ComponentType() string
}

type System interface {
	SystemType() string
}

type Engine struct {
	currentId int
	entities  []*Entity
}

/**
 * Engine
 */

func NewEngine() *Engine {
	return &Engine{
		currentId: 1,
	}
}

func (engine *Engine) NewEntity() *Entity {
	newEntity := &Entity{
		id:         engine.currentId,
		components: make(map[string]Component),
	}
	engine.currentId = engine.currentId + 1
	engine.entities = append(engine.entities, newEntity)
	return newEntity
}

func (engine Engine) GetEntities(types []string) []*Entity {
	var found []*Entity
	for _, entity := range engine.entities {
		if entity.HasComponents(types) {
			found = append(found, entity)
		}
	}
	return found
}

/**
 * Entity
 */

func (entity *Entity) String() string {
	var str = "Entity " + strconv.Itoa(entity.id) + "["
	for _, c := range entity.components {
		str += c.ComponentType() + ","
	}
	str += "]"
	return str
}

func (entity *Entity) AddComponent(componentType string, component Component) {
	entity.components[componentType] = component
}

func (entity *Entity) RemoveComponent(componentType string) Component {
	cmp, ok := entity.components[componentType]
	if ok {
		delete(entity.components, componentType)
		return cmp
	}
	return nil
}

func (entity *Entity) HasComponent(componentType string) bool {
	if _, ok := entity.components[componentType]; ok {
		return true
	} else {
		return false
	}
}

func (entity *Entity) HasComponents(componentTypes []string) bool {
	// Check to see if the entity has the given components
	containsAll := true
	if entity != nil {
		for i := 0; i < len(componentTypes); i++ {
			if !entity.HasComponent(componentTypes[i]) {
				containsAll = false
				break
			}
		}
	} else {
		return false
	}
	return containsAll
}

func (entity *Entity) GetComponent(componentType string) Component {
	if cmp, ok := entity.components[componentType]; ok {
		return cmp
	} else {
		return nil
	}
}

/**
 * Component
 */

/**
 * System
 */

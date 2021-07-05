package ecs

import (
	"strconv"
)

type Entity struct {
	id         int
	components map[string]Component
	engine     *Engine
}

type Component interface {
	ComponentType() string
}

type OnAddComponent interface {
	OnAdd(engine *Engine, entity *Entity)
}

type OnRemoveComponent interface {
	OnRemove(engine *Engine, entity *Entity)
}

type System interface {
	SystemType() string
}

type Engine struct {
	currentId int
	entities  []*Entity
	PosCache  PositionCache
}

/**
 * Engine
 */

func NewEngine() *Engine {
	return &Engine{
		currentId: 1,
		PosCache: PositionCache{
			Entities: make(map[string]*Entity),
		},
	}
}

func (engine *Engine) NewEntity() *Entity {
	newEntity := &Entity{
		id:         engine.currentId,
		components: make(map[string]Component),
		engine:     engine,
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

func (engine Engine) GetAllEntities() []*Entity {
	return engine.entities
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

	// Call event
	cmp, ok := component.(OnAddComponent)
	if ok {
		cmp.OnAdd(entity.engine, entity)
	}
}

func (entity *Entity) RemoveComponent(componentType string) Component {
	component, ok := entity.components[componentType]
	if ok {
		delete(entity.components, componentType)

		// Call event
		cmp, ok := component.(OnRemoveComponent)
		if ok {
			cmp.OnRemove(entity.engine, entity)
		}
		return component
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

/**
 * Position cache
 */
type PositionCache struct {
	Entities map[string]*Entity
}

func (c *PositionCache) Add(key string, value *Entity) {
	c.Entities[key] = value
}

func (c *PositionCache) Delete(key string) {
	delete(c.Entities, key)
}

func (c *PositionCache) Get(key string) (*Entity, bool) {
	e, ok := c.Entities[key]
	return e, ok
}

package ecs

import (
	"fmt"
	"strconv"
	"strings"
)

type Entity struct {
	Id         int
	components map[string]Component
	Engine     *Engine
}

type EntityList []*Entity

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
	Entities  EntityList
	PosCache  PositionCache
}

/**
 * Engine
 */

func NewEngine() *Engine {
	return &Engine{
		currentId: 1,
		PosCache: PositionCache{
			Entities: map[string]EntitySet{},
		},
	}
}

func (engine *Engine) NewEntity() *Entity {
	newEntity := &Entity{
		Id:         engine.currentId,
		components: make(map[string]Component),
		Engine:     engine,
	}
	engine.currentId = engine.currentId + 1
	engine.Entities = append(engine.Entities, newEntity)
	return newEntity
}

func (engine *Engine) DestroyEntity(entity *Entity) {
	// Remove all components to trigger possible actions
	for k := range entity.components {
		entity.RemoveComponent(k)
	}
	engine.Entities.RemoveEntity(entity)

	entity.components = nil
	entity.Engine = nil
}

/**
 * Entity
 */

func (entity *Entity) String() string {
	var str = "Entity " + strconv.Itoa(entity.Id) + "["
	for _, c := range entity.components {
		str += c.ComponentType() + ","
	}
	str += "]"
	return str
}

func (entity *Entity) AddComponent(componentType string, component Component) {
	entity.components[componentType] = component

	// Call event if possible
	cmp, ok := component.(OnAddComponent)
	if ok {
		cmp.OnAdd(entity.Engine, entity)
	}
}

func (entity *Entity) RemoveComponent(componentType string) Component {
	component, ok := entity.components[componentType]
	if ok {
		delete(entity.components, componentType)

		// Call event if possible
		cmp, ok := component.(OnRemoveComponent)
		if ok {
			cmp.OnRemove(entity.Engine, entity)
		}
		return component
	}
	return nil
}

func (entity *Entity) ReplaceComponent(componentType string, newComponent Component) {
	entity.RemoveComponent(componentType)
	entity.AddComponent(componentType, newComponent)
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
 * EntityList
 */

func (entityList *EntityList) GetEntities(types []string) EntityList {
	var found EntityList
	for _, entity := range *entityList {
		if entity.HasComponents(types) {
			found = append(found, entity)
		}
	}
	return found
}

// Only one Entity expected, nil if not
func (entityList *EntityList) GetEntity(types []string) *Entity {
	found := entityList.GetEntities(types)
	if len(found) == 1 {
		return found[0]
	}
	fmt.Println("Warning, more than one entity found for types: " + strings.Join(types, ","))
	return nil
}

func (entityList *EntityList) RemoveEntity(entity *Entity) {
	old := *entityList
	for i, e := range old {
		if e.Id == entity.Id {
			old = append(old[:i], old[i+1:]...)
			break
		}
	}
	entityList = &old
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

type EntitySet map[*Entity]bool

type PositionCache struct {
	Entities map[string]EntitySet
}

func (c *PositionCache) Add(key string, value *Entity) {
	_, ok := c.Entities[key]
	if !ok {
		c.Entities[key] = make(EntitySet)
	}
	c.Entities[key][value] = true
}

func (c *PositionCache) Delete(key string, value *Entity) {
	set, ok := c.Entities[key]
	if ok {
		delete(set, value)
		if len(set) == 0 {
			delete(c.Entities, key)
		}
	}
}

func (c *PositionCache) GetByCoord(x int, y int) (EntityList, bool) {
	key := strconv.Itoa(x) + "," + strconv.Itoa(y)
	return c.Get(key)
}

func (c *PositionCache) GetByCoordAndComponents(x int, y int, cmpTypes []string) (EntityList, bool) {
	key := strconv.Itoa(x) + "," + strconv.Itoa(y)
	entityList, ok := c.Get(key)
	if ok {
		return entityList.GetEntities(cmpTypes), true
	}
	return EntityList{}, false
}

func (c *PositionCache) GetOneByCoordAndComponents(x int, y int, cmpTypes []string) (*Entity, bool) {
	key := strconv.Itoa(x) + "," + strconv.Itoa(y)
	entityList, ok := c.Get(key)
	if ok {
		found := entityList.GetEntities(cmpTypes)
		if len(found) == 1 {
			return found[0], true
		}
	}
	return &Entity{}, false
}

func (c *PositionCache) Get(key string) (EntityList, bool) {
	set, ok := c.Entities[key]
	entities := make([]*Entity, 0, len(set))
	if ok {
		for entity := range set {
			entities = append(entities, entity)
		}

	}
	return entities, ok
}

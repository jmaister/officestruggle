package ecs

import (
	"fmt"
	"strconv"
	"strings"

	"jordiburgos.com/officestruggle/constants"
)

type Entity struct {
	Id         int
	Components map[string]Component
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
		Components: make(map[string]Component),
		Engine:     engine,
	}
	engine.currentId = engine.currentId + 1
	engine.Entities = append(engine.Entities, newEntity)
	return newEntity
}

func (engine *Engine) DestroyEntity(entity *Entity) {
	// Remove all components to trigger possible actions
	for k := range entity.Components {
		entity.RemoveComponent(k)
	}
	engine.Entities = engine.Entities.RemoveEntity(entity)

	entity.Components = nil
	entity.Engine = nil
	entity = nil
}

// Used to load a game state
func (engine *Engine) SetEntityList(entityList EntityList) {
	engine.Entities = entityList
	engine.currentId = entityList[len(entityList)-1].Id + 10

	// Set the engine for every entity
	for _, entity := range entityList {
		entity.Engine = engine
	}

	// Recreate the position cache
	for k := range engine.PosCache.Entities {
		delete(engine.PosCache.Entities, k)
	}

	found := engine.Entities.GetEntities([]string{constants.Position})
	for _, f := range found {
		cmp := f.GetComponent(constants.Position)
		// Triggers position cache update
		f.ReplaceComponent(cmp)
	}
}

/**
 * Entity
 */

func (entity Entity) String() string {
	var str = "Entity " + strconv.Itoa(entity.Id) + "["
	for _, c := range entity.Components {
		str += c.ComponentType() + ","
	}
	str += "]"
	return str
}

func (entity *Entity) AddComponent(component Component) {
	entity.Components[component.ComponentType()] = component

	// Call event if possible
	cmp, ok := component.(OnAddComponent)
	if ok {
		cmp.OnAdd(entity.Engine, entity)
	}
}

func (entity *Entity) RemoveComponent(componentType string) {
	component, ok := entity.Components[componentType]
	if ok {
		delete(entity.Components, componentType)

		// Call event if possible
		cmp, ok := component.(OnRemoveComponent)
		if ok {
			cmp.OnRemove(entity.Engine, entity)
		}
	}
}

func (entity *Entity) ReplaceComponent(newComponent Component) {
	entity.RemoveComponent(newComponent.ComponentType())
	entity.AddComponent(newComponent)
}

func (entity Entity) HasComponent(componentType string) bool {
	if _, ok := entity.Components[componentType]; ok {
		return true
	} else {
		return false
	}
}

func (entity Entity) HasComponents(componentTypes []string) bool {
	// Check to see if the entity has the given components
	containsAll := true
	for i := 0; i < len(componentTypes); i++ {
		if _, ok := entity.Components[componentTypes[i]]; !ok {
			containsAll = false
			break
		}
	}
	return containsAll
}

func (entity Entity) GetComponent(componentType string) Component {
	if cmp, ok := entity.Components[componentType]; ok {
		return cmp
	} else {
		return nil
	}
}

/**
 * EntityList
 */

// TODO: move to GameState to allow queries with values, ie Z=1
func (entityList EntityList) GetEntities(types []string) EntityList {
	found := EntityList{}
	for _, entity := range entityList {
		if entity.HasComponents(types) {
			found = append(found, entity)
		}
	}
	return found
}

// Only one Entity expected, nil if not
func (entityList EntityList) GetEntity(types []string) *Entity {
	found := entityList.GetEntities(types)
	if len(found) == 1 {
		return found[0]
	}
	fmt.Println("entities found", found)
	fmt.Println("Warning, more than one entity found for types: " + strings.Join(types, ","))
	return nil
}

func (entityList *EntityList) RemoveEntity(entity *Entity) EntityList {
	old := *entityList
	for i, e := range old {
		if e.Id == entity.Id {
			return append(old[:i], old[i+1:]...)
		}
	}
	return old
}

func (entityList EntityList) Concat(other EntityList) EntityList {
	return append(entityList, other...)
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

func (c *PositionCache) GetByCoord(x int, y int, z int) (EntityList, bool) {
	key := createKey(x, y, z)
	return c.Get(key)
}

func (c *PositionCache) GetByCoordAndComponents(x int, y int, z int, cmpTypes []string) (EntityList, bool) {
	key := createKey(x, y, z)
	entityList, ok := c.Get(key)
	if ok {
		return entityList.GetEntities(cmpTypes), true
	}
	return EntityList{}, false
}

func (c *PositionCache) GetOneByCoordAndComponents(x int, y int, z int, cmpTypes []string) (*Entity, bool) {
	key := createKey(x, y, z)
	entityList, ok := c.Get(key)
	if ok {
		found := entityList.GetEntities(cmpTypes)
		if len(found) == 1 {
			return found[0], true
		}
	}
	return &Entity{}, false
}

func createKey(x int, y int, z int) string {
	return strconv.Itoa(x) + "," + strconv.Itoa(y) + "," + strconv.Itoa(z)
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

func (c *PositionCache) GetZ(z int) EntityList {
	suffix := "," + strconv.Itoa(z)
	entities := []*Entity{}
	for k, v := range c.Entities {
		if strings.HasSuffix(k, suffix) {
			for entity, present := range v {
				if present {
					entities = append(entities, entity)
				}
			}
		}
	}
	return entities
}

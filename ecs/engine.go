package ecs

type Engine struct {
	currentId  int
	components []*Component
	entities   []*Entity
}

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

type Entity struct {
	id         int
	components []*Component
}

func (entity *Entity) AddComponent(component *Component) {
	entity.components = append(entity.components, component)
}

type Component struct {
	Type string
}

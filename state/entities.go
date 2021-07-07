package state

import "jordiburgos.com/officestruggle/ecs"

func NewPlayer(entity *ecs.Entity) *ecs.Entity {
	entity.AddComponent(Player, PlayerComponent{})
	entity.AddComponent(Description, DescriptionComponent{Name: "You"})
	entity.AddComponent(Apparence, ApparenceComponent{Color: "#0000FF", Char: '@'})
	entity.AddComponent(Layer400, Layer400Component{})
	return entity
}

func ApplyPosition(entity *ecs.Entity, x int, y int) *ecs.Entity {
	entity.AddComponent(Position, PositionComponent{X: x, Y: y})
	return entity
}

func NewGlobin(entity *ecs.Entity) *ecs.Entity {
	entity.AddComponent(AI, AIComponent{})
	entity.AddComponent(Description, DescriptionComponent{Name: "Goblin"})
	entity.AddComponent(IsBlocking, IsBlockingComponent{})
	entity.AddComponent(Apparence, ApparenceComponent{Color: "#00FC00", Char: 'g'})
	entity.AddComponent(Layer400, Layer400Component{})
	return entity
}

package state

import "jordiburgos.com/officestruggle/ecs"

func NewPlayer(entity *ecs.Entity) *ecs.Entity {
	entity.AddComponent(Player, PlayerComponent{})
	entity.AddComponent(Description, DescriptionComponent{Name: "You"})
	entity.AddComponent(Apparence, ApparenceComponent{Color: "#ffffff", Bg: "#FF0000", Char: '@'})
	entity.AddComponent(Layer400, Layer400Component{})
	entity.AddComponent(Stats, StatsComponent{
		statsValues: &statsValues{
			Health:     10,
			MaxHealth:  10,
			Defense:    3,
			MaxDefense: 10,
			Power:      2,
			MaxPower:   10,
			Fov:        10,
		},
	})
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
	entity.AddComponent(Stats, StatsComponent{
		statsValues: &statsValues{
			Health:     4,
			MaxHealth:  10,
			Defense:    1,
			MaxDefense: 10,
			Power:      4,
			MaxPower:   10,
			Fov:        6,
		},
	})
	return entity
}

func NewHealthPotion(entity *ecs.Entity) *ecs.Entity {
	entity.AddComponent(Consumable, ConsumableComponent{})
	entity.AddComponent(Description, DescriptionComponent{Name: "Health Potion"})
	entity.AddComponent(Apparence, ApparenceComponent{Color: "#FF0000", Char: 'o'})
	entity.AddComponent(Layer300, Layer400Component{})
	entity.AddComponent(Stats, StatsComponent{
		statsValues: &statsValues{
			Health:     5,
			MaxHealth:  1,
			Defense:    0,
			MaxDefense: 0,
			Power:      0,
			MaxPower:   0,
			Fov:        0,
		},
	})
	return entity
}

package state

import "jordiburgos.com/officestruggle/ecs"

func NewPlayer(entity *ecs.Entity) *ecs.Entity {
	entity.AddComponent(Player, PlayerComponent{})
	entity.AddComponent(Description, DescriptionComponent{Name: "Player"})
	entity.AddComponent(Apparence, ApparenceComponent{Color: "#ffffff", Bg: "#FF0000", Char: '@'})
	entity.AddComponent(Layer500, Layer400Component{})
	entity.AddComponent(Stats, StatsComponent{
		StatsValues: &StatsValues{
			Health:     10,
			MaxHealth:  10,
			Defense:    3,
			MaxDefense: 10,
			Power:      4,
			MaxPower:   10,
			Fov:        10,
		},
	})
	entity.AddComponent(Inventory, InventoryComponent{
		Items:    ecs.EntityList{},
		MaxItems: 10,
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
		StatsValues: &StatsValues{
			Health:     4,
			MaxHealth:  10,
			Defense:    1,
			MaxDefense: 10,
			Power:      4,
			MaxPower:   10,
			Fov:        6,
		},
	})
	entity.AddComponent(Consumable, ConsumableComponent{
		StatsValues: &StatsValues{
			Health:     -3,
			MaxHealth:  0,
			Defense:    0,
			MaxDefense: 0,
			Power:      1,
			MaxPower:   1,
			Fov:        0,
		},
	})
	return entity
}

func NewHealthPotion(entity *ecs.Entity) *ecs.Entity {
	entity.AddComponent(IsPickup, IsPickupComponent{})
	entity.AddComponent(Consumable, ConsumableComponent{
		StatsValues: &StatsValues{
			Health:     5,
			MaxHealth:  1,
			Defense:    0,
			MaxDefense: 0,
			Power:      0,
			MaxPower:   0,
			Fov:        0,
		},
	})
	entity.AddComponent(Description, DescriptionComponent{Name: "Health Potion"})
	entity.AddComponent(Apparence, ApparenceComponent{Color: "#FF0000", Char: 'o'})
	entity.AddComponent(Layer300, Layer400Component{})
	return entity
}

func NewSword(entity *ecs.Entity) *ecs.Entity {
	entity.AddComponent(IsPickup, IsPickupComponent{})
	entity.AddComponent(Equipable, EquipableComponent{
		StatsValues: &StatsValues{
			Health:     0,
			MaxHealth:  1,
			Defense:    0,
			MaxDefense: 0,
			Power:      5,
			MaxPower:   5,
			Fov:        1,
		},
	})
	entity.AddComponent(Description, DescriptionComponent{Name: "Sword"})
	entity.AddComponent(Apparence, ApparenceComponent{Color: "#1EFFFF", Char: '/'})
	entity.AddComponent(Layer300, Layer400Component{})
	return entity
}

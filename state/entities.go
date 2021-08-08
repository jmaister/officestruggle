package state

import (
	"jordiburgos.com/officestruggle/ecs"
)

func NewPlayer(entity *ecs.Entity) *ecs.Entity {
	entity.AddComponent(PlayerComponent{})
	entity.AddComponent(DescriptionComponent{Name: "Player"})
	entity.AddComponent(ApparenceComponent{Color: "#ffffff", Char: '@'})
	entity.AddComponent(Layer400Component{})
	entity.AddComponent(StatsComponent{
		StatsValues: &StatsValues{},
	})
	entity.AddComponent(InventoryComponent{
		Items:    ecs.EntityList{},
		MaxItems: 10,
	})
	entity.AddComponent(EquipmentComponent{
		Base: StatsValues{
			Health:     10,
			MaxHealth:  10,
			Defense:    3,
			MaxDefense: 10,
			Power:      4,
			MaxPower:   10,
			Fov:        10,
			MaxFov:     20,
		},
		Items: map[EquipPosition]*ecs.Entity{},
	})
	return entity
}

func ApplyPosition(entity *ecs.Entity, x int, y int, z int) *ecs.Entity {
	entity.AddComponent(PositionComponent{X: x, Y: y, Z: z})
	return entity
}

func NewGlobin(entity *ecs.Entity) *ecs.Entity {
	entity.AddComponent(AIComponent{})
	entity.AddComponent(DescriptionComponent{Name: "Goblin"})
	entity.AddComponent(IsBlockingComponent{})
	entity.AddComponent(ApparenceComponent{Color: "#00FC00", Char: 'g'})
	entity.AddComponent(Layer400Component{})
	entity.AddComponent(StatsComponent{
		StatsValues: &StatsValues{
			Health:     4,
			MaxHealth:  10,
			Defense:    1,
			MaxDefense: 10,
			Power:      4,
			MaxPower:   10,
			Fov:        6,
			MaxFov:     6,
		},
	})
	entity.AddComponent(ConsumableComponent{
		StatsValues: &StatsValues{
			Health:     -3,
			MaxHealth:  0,
			Defense:    0,
			MaxDefense: 0,
			Power:      1,
			MaxPower:   1,
			Fov:        0,
			MaxFov:     0,
		},
	})
	return entity
}

func NewHealthPotion(entity *ecs.Entity) *ecs.Entity {
	entity.AddComponent(IsPickupComponent{})
	entity.AddComponent(ConsumableComponent{
		StatsValues: &StatsValues{
			Health:     5,
			MaxHealth:  1,
			Defense:    0,
			MaxDefense: 0,
			Power:      0,
			MaxPower:   0,
			Fov:        0,
			MaxFov:     0,
		},
	})
	entity.AddComponent(DescriptionComponent{Name: "Health Potion"})
	entity.AddComponent(ApparenceComponent{Color: "#FF0000", Char: 'o'})
	entity.AddComponent(Layer300Component{})
	return entity
}

func NewSword(entity *ecs.Entity) *ecs.Entity {
	entity.AddComponent(IsPickupComponent{})
	entity.AddComponent(EquipableComponent{
		EquipSlot: EquipWeapon,
		StatsValues: &StatsValues{
			Health:     0,
			MaxHealth:  1,
			Defense:    0,
			MaxDefense: 0,
			Power:      5,
			MaxPower:   5,
			Fov:        1,
			MaxFov:     1,
		},
	})
	entity.AddComponent(DescriptionComponent{Name: "Sword"})
	entity.AddComponent(ApparenceComponent{Color: "#1EFFFF", Char: '/'})
	entity.AddComponent(Layer300Component{})
	return entity
}

func applyVisitableEntity(entity *ecs.Entity) *ecs.Entity {
	entity.AddComponent(Layer100Component{})
	entity.AddComponent(VisitableComponent{Explored: false, Visible: false})
	return entity
}

func NewWall(entity *ecs.Entity, x int, y int, z int) *ecs.Entity {
	ApplyPosition(entity, x, y, z)
	applyVisitableEntity(entity)
	entity.AddComponent(DescriptionComponent{Name: "Wall"})
	entity.AddComponent(IsBlockingComponent{})
	entity.AddComponent(ApparenceComponent{Color: "#1a1aff", Char: '#'})
	return entity
}

func NewFloor(entity *ecs.Entity, x int, y int, z int) *ecs.Entity {
	ApplyPosition(entity, x, y, z)
	applyVisitableEntity(entity)
	entity.AddComponent(DescriptionComponent{Name: "Floor"})
	entity.AddComponent(IsFloorComponent{})
	entity.AddComponent(ApparenceComponent{Color: "#e3e3e3", Char: '.'})
	return entity
}

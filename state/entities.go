package state

import (
	"image/color"

	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/palette"
)

func NewPlayer(entity *ecs.Entity) *ecs.Entity {
	entity.AddComponent(PlayerComponent{})
	entity.AddComponent(DescriptionComponent{Name: "Player"})
	entity.AddComponent(ApparenceComponent{Color: color.White, Char: '@'})
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
			Health:    10,
			MaxHealth: 10,
			Defense:   3,
			Power:     4,
			Fov:       10,
		},
		Items: map[constants.EquipSlot]*ecs.Entity{},
	})
	entity.AddComponent(LevelingComponent{
		CurrentLevel:  1,
		CurrentXP:     0,
		LevelUpBase:   0,
		LevelUpFactor: 150,
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
	entity.AddComponent(ApparenceComponent{Color: palette.PColor(palette.Green, 0.6), Char: 'g'})
	entity.AddComponent(Layer400Component{})
	entity.AddComponent(StatsComponent{
		StatsValues: &StatsValues{
			Health:    4,
			MaxHealth: 10,
			Defense:   1,
			Power:     4,
			Fov:       6,
		},
	})
	entity.AddComponent(ConsumableComponent{
		StatsValues: &StatsValues{
			Health:    -3,
			MaxHealth: 0,
			Defense:   0,
			Power:     1,
			Fov:       0,
		},
	})
	entity.AddComponent(XPGiverComponent{
		XPBase:     10,
		XPPerLevel: 20,
		Level:      1,
	})
	return entity
}

func NewDragon(entity *ecs.Entity) *ecs.Entity {
	entity.AddComponent(WinGameComponent{})
	entity.AddComponent(AIComponent{})
	entity.AddComponent(DescriptionComponent{Name: "Dragon - Final Boss"})
	entity.AddComponent(IsBlockingComponent{})
	entity.AddComponent(ApparenceComponent{Color: palette.PColor(palette.Green, 0.3), Char: 'D'})
	entity.AddComponent(Layer400Component{})
	entity.AddComponent(StatsComponent{
		StatsValues: &StatsValues{
			Health:    200,
			MaxHealth: 200,
			Defense:   1,
			Power:     8,
			Fov:       10,
		},
	})
	entity.AddComponent(ConsumableComponent{
		StatsValues: &StatsValues{
			Health:    -3,
			MaxHealth: 0,
			Defense:   0,
			Power:     1,
			Fov:       0,
		},
	})
	entity.AddComponent(XPGiverComponent{
		XPBase:     10,
		XPPerLevel: 20,
		Level:      1,
	})
	return entity
}

func NewHealthPotion(entity *ecs.Entity) *ecs.Entity {
	entity.AddComponent(IsPickupComponent{})
	entity.AddComponent(ConsumableComponent{
		StatsValues: &StatsValues{
			Health:    5,
			MaxHealth: 1,
			Defense:   0,
			Power:     0,
			Fov:       0,
		},
	})
	entity.AddComponent(DescriptionComponent{Name: "Health Potion"})
	entity.AddComponent(ApparenceComponent{Color: palette.PColor(palette.Red, 0.5), Char: 'o'})
	entity.AddComponent(Layer300Component{})
	entity.AddComponent(XPGiverComponent{
		XPBase:     5,
		XPPerLevel: 0,
		Level:      1,
	})
	return entity
}

func NewSword(entity *ecs.Entity) *ecs.Entity {
	entity.AddComponent(IsPickupComponent{})
	entity.AddComponent(EquipableComponent{
		EquipSlot: constants.EquipWeapon,
		StatsValues: &StatsValues{
			Health:    0,
			MaxHealth: 1,
			Defense:   0,
			Power:     5,
			Fov:       1,
		},
	})
	entity.AddComponent(DescriptionComponent{Name: "Sword"})
	entity.AddComponent(ApparenceComponent{Color: palette.PColor(palette.Cyan, 0.5), Char: '/'})
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
	entity.AddComponent(ApparenceComponent{Color: palette.PColor(palette.Blue, 0.5), Char: '#'})
	return entity
}

func NewFloor(entity *ecs.Entity, x int, y int, z int) *ecs.Entity {
	ApplyPosition(entity, x, y, z)
	applyVisitableEntity(entity)
	entity.AddComponent(DescriptionComponent{Name: "Floor"})
	entity.AddComponent(IsFloorComponent{})
	entity.AddComponent(ApparenceComponent{Color: palette.PColor(palette.Gray, 0.7), Char: '.'})
	return entity
}

func NewUpstairs(entity *ecs.Entity, x int, y int, z int, targetX int, targetY int, targetZ int) *ecs.Entity {
	ApplyPosition(entity, x, y, z)
	applyVisitableEntity(entity)
	entity.AddComponent(DescriptionComponent{Name: "Stairs going up"})
	entity.AddComponent(StairsComponent{
		GoingUp: true,
		TargetX: targetX,
		TargetY: targetY,
		TargetZ: targetZ,
	})
	entity.AddComponent(ApparenceComponent{Color: palette.PColor(palette.Sepia, 0.2), Char: '<'})
	return entity
}

func NewDownstairs(entity *ecs.Entity, x int, y int, z int, targetX int, targetY int, targetZ int) *ecs.Entity {
	ApplyPosition(entity, x, y, z)
	applyVisitableEntity(entity)
	entity.AddComponent(DescriptionComponent{Name: "Stairs going down"})
	entity.AddComponent(StairsComponent{
		GoingUp: false,
		TargetX: targetX,
		TargetY: targetY,
		TargetZ: targetZ,
	})
	entity.AddComponent(ApparenceComponent{Color: palette.PColor(palette.Sepia, 0.2), Char: '>'})
	return entity
}

func NewMoneyAmount(entity *ecs.Entity, amount int) *ecs.Entity {
	entity.AddComponent(IsPickupComponent{})
	entity.AddComponent(DescriptionComponent{Name: "Money amount"})
	entity.AddComponent(ApparenceComponent{Color: palette.PColor(palette.Amber, 0.6), Char: '$'})
	entity.AddComponent(Layer300Component{})
	entity.AddComponent(MoneyComponent{
		Coins: amount,
	})
	return entity
}

package state

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"jordiburgos.com/officestruggle/assets"
	"jordiburgos.com/officestruggle/ecs"
)

func NewPlayer(entity *ecs.Entity) *ecs.Entity {
	entity.AddComponent(Player, PlayerComponent{})
	entity.AddComponent(Description, DescriptionComponent{Name: "Player"})
	entity.AddComponent(Apparence, ApparenceComponent{Color: "#ffffff", Char: '@'})
	entity.AddComponent(Layer400, Layer400Component{})
	entity.AddComponent(Stats, StatsComponent{
		StatsValues: &StatsValues{},
	})
	entity.AddComponent(Inventory, InventoryComponent{
		Items:    ecs.EntityList{},
		MaxItems: 10,
	})
	entity.AddComponent(Equipment, EquipmentComponent{
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
			MaxFov:     6,
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
			MaxFov:     0,
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
			MaxFov:     0,
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
		Position: EquipWeapon,
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
	entity.AddComponent(Description, DescriptionComponent{Name: "Sword"})
	entity.AddComponent(Apparence, ApparenceComponent{Color: "#1EFFFF", Char: '/'})
	entity.AddComponent(Layer300, Layer400Component{})
	return entity
}

type damageAnimation struct {
	X                 int
	Y                 int
	Damage            string
	AnimationStart    time.Time
	AnimationDuration time.Duration
}

func (a damageAnimation) StartTime() time.Time {
	return a.AnimationStart
}
func (a damageAnimation) Duration() time.Duration {
	return a.AnimationDuration
}
func (a damageAnimation) Update(percent float64, screen *ebiten.Image) {
	x := a.X * 20
	y := a.Y*20 - int(3*20*(1-percent))

	fnt := assets.MplusFont(20)
	text.Draw(screen, a.Damage, fnt, x, y, color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 127,
	})
}

func NewDamageAnimation(entity *ecs.Entity, x int, y int, damage string) *ecs.Entity {
	entity.AddComponent(Layer500, Layer500Component{})
	entity.AddComponent(Animated, AnimatedComponent{
		Animation: damageAnimation{
			X:                 x,
			Y:                 y,
			Damage:            damage,
			AnimationStart:    time.Now(),
			AnimationDuration: 1 * time.Second,
		},
	})
	return entity
}

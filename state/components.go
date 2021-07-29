package state

import (
	"strconv"
	"strings"

	"jordiburgos.com/officestruggle/ecs"
)

const (
	Player      = "player"
	Apparence   = "apparence"
	Position    = "position"
	Move        = "move"
	IsBlocking  = "isBlocking"
	IsFloor     = "isFloor"
	Layer100    = "layer100" // Walls and floors
	Layer300    = "layer300" // Objects
	Layer400    = "layer400" // Player and enemies
	Layer500    = "layer500" // Animations
	Visitable   = "visitable"
	Description = "description"
	AI          = "ai"
	Stats       = "stats"
	Consumable  = "consumable"
	IsPickup    = "isPickup"
	Dead        = "dead"
	Inventory   = "inventory"
	Equipable   = "equipable"
	Equipment   = "equipment"
	Animated    = "animated"
)

type PlayerComponent struct {
}

func (a PlayerComponent) ComponentType() string {
	return Player
}

type ApparenceComponent struct {
	Color string
	Bg    string
	Char  rune
}

func (a ApparenceComponent) ComponentType() string {
	return Apparence
}

type PositionComponent struct {
	X int
	Y int
}

func (a PositionComponent) ComponentType() string {
	return Position
}

func GetPosition(entity *ecs.Entity) PositionComponent {
	return entity.GetComponent(Position).(PositionComponent)
}

func (a PositionComponent) GetKey() string {
	return strconv.Itoa(a.X) + "," + strconv.Itoa(a.Y)
}

func (a PositionComponent) String() string {
	return a.GetKey()
}

func (a PositionComponent) OnAdd(engine *ecs.Engine, entity *ecs.Entity) {
	engine.PosCache.Add(a.GetKey(), entity)
}

func (a PositionComponent) OnRemove(engine *ecs.Engine, entity *ecs.Entity) {
	engine.PosCache.Delete(a.GetKey(), entity)
}

type MoveComponent struct {
	X int
	Y int
}

func (a MoveComponent) ComponentType() string {
	return Move
}

type IsBlockingComponent struct{}

func (a IsBlockingComponent) ComponentType() string {
	return IsBlocking
}

type IsFloorComponent struct{}

func (a IsFloorComponent) ComponentType() string {
	return IsFloor
}

type Layer100Component struct{}

func (a Layer100Component) ComponentType() string {
	return Layer100
}

type Layer300Component struct{}

func (a Layer300Component) ComponentType() string {
	return Layer300
}

type Layer400Component struct{}

func (a Layer400Component) ComponentType() string {
	return Layer400
}

type Layer500Component struct{}

func (a Layer500Component) ComponentType() string {
	return Layer500
}

type VisitableComponent struct {
	// You already know this tile
	Explored bool
	// It is reachable from your sight, depends on current position of the player
	Visible bool
}

func (a VisitableComponent) ComponentType() string {
	return Visitable
}

type DescriptionComponent struct {
	Name string
}

func (a DescriptionComponent) ComponentType() string {
	return Description
}

func GetDescription(entity *ecs.Entity) string {
	cmp, ok := entity.GetComponent(Description).(DescriptionComponent)
	if ok {
		return cmp.Name
	}
	return "UnDeFiNeD!"
}

func GetLongDescription(entity *ecs.Entity) string {
	cmp, ok := entity.GetComponent(Description).(DescriptionComponent)
	if ok {
		str := cmp.Name

		if entity.HasComponent(Dead) {
			str = str + " corpse"
		}

		if entity.HasComponent(Stats) {
			stats := entity.GetComponent(Stats).(StatsComponent)
			str = str + " (" + stats.String() + ")"
		} else if entity.HasComponent(Consumable) {
			cons := entity.GetComponent(Consumable).(ConsumableComponent)
			str = str + " (" + cons.String() + ")"
		} else if entity.HasComponent(Equipable) {
			eq := entity.GetComponent(Equipable).(EquipableComponent)
			str = str + " (" + eq.String() + ")"
		}
		return str
	}
	return "UnDeFiNeD!"
}

type AIComponent struct {
}

func (a AIComponent) ComponentType() string {
	return AI
}

type StatsValues struct {
	Health     int
	MaxHealth  int
	Defense    int
	MaxDefense int
	Power      int
	MaxPower   int
	Fov        int
	MaxFov     int
}

func (a StatsValues) String() string {
	s := ""
	s += statDiff("Health", a.Health, a.MaxHealth)
	s += statDiff("Def", a.Defense, a.MaxDefense)
	s += statDiff("Pow", a.Power, a.MaxPower)
	s += statDiff("FOV", a.Fov, a.MaxFov)
	return strings.Trim(s, " ")
}

func (plStats StatsValues) Merge(other StatsValues) StatsValues {
	// First increases the max values

	plStats.MaxHealth += other.MaxHealth
	plStats.Health = increase(plStats.Health, plStats.MaxHealth, other.Health)
	plStats.MaxDefense += other.MaxDefense
	plStats.Defense = increase(plStats.Defense, plStats.MaxDefense, other.Defense)
	plStats.MaxPower += other.MaxPower
	plStats.Power = increase(plStats.Power, plStats.MaxPower, other.Power)
	plStats.MaxFov += other.MaxFov
	plStats.Fov = increase(plStats.Fov, plStats.MaxFov, other.Fov)

	return plStats
}

func increase(current int, max int, incr int) int {
	current = current + incr
	if current > max {
		return max
	}
	return current
}

func addSign(i int) string {
	if i >= 0 {
		return "+" + strconv.Itoa(i)
	}
	return strconv.Itoa(i)
}

func statDiff(name string, value int, max int) string {
	if value != 0 || max != 0 {
		return name + " " + addSign(value) + "/" + addSign(max) + " "
	}
	return ""
}

type StatsComponent struct {
	*StatsValues
}

func (a StatsComponent) ComponentType() string {
	return Stats
}

type ConsumableComponent struct {
	*StatsValues
}

func (a ConsumableComponent) ComponentType() string {
	return Consumable
}

type IsPickupComponent struct {
}

func (a IsPickupComponent) ComponentType() string {
	return IsPickup
}

type DeadComponent struct {
}

func (a DeadComponent) ComponentType() string {
	return Dead
}

type InventoryComponent struct {
	Items    ecs.EntityList
	MaxItems int
}

func (a *InventoryComponent) AddItem(entity *ecs.Entity) bool {
	if len(a.Items) >= a.MaxItems {
		return false
	}
	a.Items = append(a.Items, entity)
	return true
}

func (a *InventoryComponent) RemoveItem(entity *ecs.Entity) bool {
	for i, item := range a.Items {
		if item == entity {
			a.Items = append(a.Items[:i], a.Items[i+1:]...)
			return true
		}
	}
	return false
}

func (a InventoryComponent) ComponentType() string {
	return Inventory
}

type EquipableComponent struct {
	*StatsValues

	EquipSlot EquipPosition
}

func (a EquipableComponent) ComponentType() string {
	return Equipable
}

type EquipPosition string

const (
	EquipHead   EquipPosition = "head"
	EquipShield EquipPosition = "shield"
	EquipWeapon EquipPosition = "weapon"
	EquipBoots  EquipPosition = "boot"
	EquipArmor  EquipPosition = "armor"
)

var EquipmentPositions = []EquipPosition{
	EquipHead,
	EquipShield,
	EquipWeapon,
	EquipBoots,
	EquipArmor,
}

type EquipmentComponent struct {
	Base  StatsValues
	Items map[EquipPosition]*ecs.Entity
}

func (e *EquipmentComponent) UpdateStats(player *ecs.Entity) {
	newState := e.Base
	for _, item := range e.Items {
		itemStats := item.GetComponent(Equipable).(EquipableComponent)
		newState = newState.Merge(*itemStats.StatsValues)
	}
	player.ReplaceComponent(Stats, StatsComponent{
		StatsValues: &newState,
	})
}

func (a EquipmentComponent) OnAdd(engine *ecs.Engine, entity *ecs.Entity) {
	// Update player.StatsComponent
	a.UpdateStats(entity)
}

func (a EquipmentComponent) OnRemove(engine *ecs.Engine, entity *ecs.Entity) {
	// Update player.StatsComponent
	a.UpdateStats(entity)
}

func (a EquipmentComponent) ComponentType() string {
	return Equipment
}

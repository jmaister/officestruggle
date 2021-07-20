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
	Layer100    = "layer100"
	Layer300    = "layer300"
	Layer400    = "layer400"
	Layer500    = "layer500"
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

func toStr(i int) string {
	return strconv.Itoa(i)
}

func st(name string, value int, max int) string {
	return name + ": " + toStr(value) + "/" + toStr(max)
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

func (a *InventoryComponent) PickUp(entity *ecs.Entity) bool {
	if len(a.Items) >= a.MaxItems {
		return false
	}
	a.Items = append(a.Items, entity)
	return true
}

func (a *InventoryComponent) Drop(entity *ecs.Entity) bool {
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
}

func (a EquipableComponent) ComponentType() string {
	return Equipable
}

type EquipPosition string

var (
	EquipHead   EquipPosition = "hd"
	EquipShield EquipPosition = "sh"
	EquipWeapon EquipPosition = "wp"
	EquipBoots  EquipPosition = "bt"
)

type EquipmentComponent struct {
	Items map[EquipPosition]*ecs.Entity
}

func (a EquipmentComponent) ComponentType() string {
	return Equipment
}

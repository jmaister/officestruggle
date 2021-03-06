package state

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/interfaces"
	"jordiburgos.com/officestruggle/palette"
)

type PlayerComponent struct {
}

func (a PlayerComponent) ComponentType() string {
	return constants.Player
}

type ApparenceComponent struct {
	Color color.Color
	Bg    color.Color
	Char  rune
}

func (a ApparenceComponent) ComponentType() string {
	return constants.Apparence
}

type PositionComponent struct {
	X int
	Y int
	Z int
}

func (a PositionComponent) ComponentType() string {
	return constants.Position
}

func GetPosition(entity *ecs.Entity) PositionComponent {
	return entity.GetComponent(constants.Position).(PositionComponent)
}

func (a PositionComponent) GetKey() string {
	return strconv.Itoa(a.X) + "," + strconv.Itoa(a.Y) + "," + strconv.Itoa(a.Z)
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
	X        int
	Y        int
	Z        int
	Absolute bool // default false: incremental, true: absolute
}

func (a MoveComponent) ComponentType() string {
	return constants.Move
}

type IsBlockingComponent struct{}

func (a IsBlockingComponent) ComponentType() string {
	return constants.IsBlocking
}

type IsFloorComponent struct{}

func (a IsFloorComponent) ComponentType() string {
	return constants.IsFloor
}

type Layer100Component struct{}

func (a Layer100Component) ComponentType() string {
	return constants.Layer100
}

type Layer300Component struct{}

func (a Layer300Component) ComponentType() string {
	return constants.Layer300
}

type Layer400Component struct{}

func (a Layer400Component) ComponentType() string {
	return constants.Layer400
}

type Layer500Component struct{}

func (a Layer500Component) ComponentType() string {
	return constants.Layer500
}

type VisitableComponent struct {
	// You already know this tile
	Explored bool
	// It is reachable from your sight, depends on current position of the player
	Visible bool
}

func (a VisitableComponent) ComponentType() string {
	return constants.Visitable
}

type DescriptionComponent struct {
	Name string
}

func (a DescriptionComponent) ComponentType() string {
	return constants.Description
}

func GetDescription(entity *ecs.Entity) string {
	cmp, ok := entity.GetComponent(constants.Description).(DescriptionComponent)
	if ok {
		return cmp.Name
	}
	return "UnDeFiNeD!"
}

func GetLongDescription(entity *ecs.Entity) string {
	cmp, ok := entity.GetComponent(constants.Description).(DescriptionComponent)
	if ok {
		str := cmp.Name

		if entity.HasComponent(constants.Stats) {
			stats := entity.GetComponent(constants.Stats).(StatsComponent)
			str = str + " (" + stats.String() + ")"
		} else if entity.HasComponent(constants.Consumable) {
			cons := entity.GetComponent(constants.Consumable).(ConsumableComponent)
			str = str + " (" + cons.String() + ")"
		} else if entity.HasComponent(constants.Equipable) {
			eq := entity.GetComponent(constants.Equipable).(EquipableComponent)
			str = str + " (" + eq.String() + ")"
		}
		return str
	}
	return "UnDeFiNeD!"
}

type AIComponent struct {
}

func (a AIComponent) ComponentType() string {
	return constants.AI
}

type StatsValues struct {
	Health    int
	MaxHealth int
	Defense   int
	Power     int
	Fov       int
}

func (a StatsValues) String() string {
	s := ""
	s += statStrMax("Health", a.Health, a.MaxHealth)
	s += statStr("Pow", a.Power)
	s += statStr("Def", a.Defense)
	s += statStr("FOV", a.Fov)
	return strings.Trim(s, " ")
}

func (plStats StatsValues) Merge(other StatsValues) StatsValues {
	// First increases the max values

	plStats.MaxHealth += other.MaxHealth
	plStats.Health = increase(plStats.Health, plStats.MaxHealth, other.Health)
	plStats.Defense += other.Defense
	plStats.Power += other.Power
	plStats.Fov += other.Fov

	return plStats
}

func increase(current int, max int, incr int) int {
	current = current + incr
	if current > max {
		return max
	}
	return current
}

func (plStats StatsValues) ApplyMultiplier(other StatsValues) StatsValues {
	plStats.Health = applyMultiplier(plStats.Health, other.Health)
	plStats.MaxHealth = applyMultiplier(plStats.MaxHealth, other.MaxHealth)
	plStats.Defense = applyMultiplier(plStats.Defense, other.Defense)
	plStats.Power = applyMultiplier(plStats.Power, other.Power)
	plStats.Fov = applyMultiplier(plStats.Fov, other.Fov)

	return plStats
}

func applyMultiplier(value int, pct int) int {
	return int(float64(value) * (float64(pct) / 100.0))
}

func addSign(i int) string {
	if i >= 0 {
		return "+" + strconv.Itoa(i)
	}
	return strconv.Itoa(i)
}

func statStrMax(name string, value int, max int) string {
	if value != 0 || max != 0 {
		return name + " " + addSign(value) + "/" + addSign(max) + " "
	}
	return ""
}

func statStr(name string, value int) string {
	if value != 0 {
		return fmt.Sprintf("%s %d ", name, value)
	}
	return ""
}

type StatsComponent struct {
	*StatsValues
}

func (a StatsComponent) ComponentType() string {
	return constants.Stats
}

type ConsumableComponent struct {
	*StatsValues
}

func (a ConsumableComponent) ComponentType() string {
	return constants.Consumable
}

type IsPickupComponent struct {
}

func (a IsPickupComponent) ComponentType() string {
	return constants.IsPickup
}

type DeadComponent struct {
}

func (a DeadComponent) ComponentType() string {
	return constants.Dead
}

type InventoryComponent struct {
	Items    ecs.EntityList
	MaxItems int
	Coins    int
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
	return constants.Inventory
}

type EquipableComponent struct {
	*StatsValues

	EquipSlot constants.EquipSlot
	Level     int
	Tier      int
	SetName   string
}

func (a EquipableComponent) ComponentType() string {
	return constants.Equipable
}

type EquipmentComponent struct {
	Base  StatsValues
	Items map[constants.EquipSlot]*ecs.Entity
}

func (e *EquipmentComponent) UpdateStats(player *ecs.Entity) {
	newState := e.Base
	for _, item := range e.Items {
		itemStats := item.GetComponent(constants.Equipable).(EquipableComponent)
		newState = newState.Merge(*itemStats.StatsValues)
	}
	player.ReplaceComponent(StatsComponent{
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
	return constants.Equipment
}

type AnimatedComponent struct {
	Animation interfaces.Animation
}

func (a AnimatedComponent) ComponentType() string {
	return constants.Animated
}

type ConsumeEffectComponent struct {
	Targeting       gamestate.TargetingType
	TargetTypes     []string
	TargetCount     int
	EffectAnimation interfaces.Animation
	EffectFunction  gamestate.EffectFunction
}

func (a ConsumeEffectComponent) ComponentType() string {
	return constants.ConsumeEffect
}

type ParalizeComponent struct {
	TurnsLeft int
}

func (a ParalizeComponent) ComponentType() string {
	return constants.Paralize
}

// Implement gamestate.EffectInfoColor
func (a ParalizeComponent) EffectInfo() string {
	return fmt.Sprintf("P-%d", a.TurnsLeft)
}
func (a ParalizeComponent) EffectInfoColor() color.Color {
	return palette.PColor(palette.Orange, 0.7)
}

type StairsComponent struct {
	GoingUp bool
	TargetX int
	TargetY int
	TargetZ int
}

func (a StairsComponent) ComponentType() string {
	return constants.Stairs
}

type LevelingComponent struct {
	CurrentLevel  int // Current entity level
	CurrentXP     int // Current entity experience points
	LevelUpBase   int // The base number we decide for leveling up
	LevelUpFactor int // The number to multiply against the entity's current level
}

func (a LevelingComponent) GetNextLevelXP() int {
	return a.LevelUpBase + a.CurrentLevel*a.LevelUpFactor
}

func (a LevelingComponent) RequiresLevelUp() bool {
	return a.CurrentLevel >= a.GetNextLevelXP()
}

func (a LevelingComponent) ComponentType() string {
	return constants.Leveling
}

type XPGiverComponent struct {
	// XP given: XPBase + (XPPerLevel * Level)
	XPBase     int // XP given when entity dies/consumed
	XPPerLevel int // XP per level
	Level      int // Level of the XP giver
}

func (a XPGiverComponent) Calculate() int {
	return a.XPBase + (a.XPPerLevel * a.Level)
}

func (a XPGiverComponent) ComponentType() string {
	return constants.XPGiver
}

type LootDropComponent struct {
	Entities ecs.EntityList
	Coins    int
}

func (a LootDropComponent) ComponentType() string {
	return constants.LootDrop
}

/*
 * Coins are in Copper units.
 * 100 Copper = 1 Silver
 * 100 Silver = 1 Gold
 */
type MoneyComponent struct {
	Coins int
}

func (a MoneyComponent) ComponentType() string {
	return constants.Money
}

// Player wins the game when the entity is killed or consumed
type WinGameComponent struct {
	Coins int
}

func (a WinGameComponent) ComponentType() string {
	return constants.WinGame
}

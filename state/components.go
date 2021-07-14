package state

import (
	"strconv"

	"jordiburgos.com/officestruggle/ecs"
)

const (
	Player        = "player"
	Apparence     = "apparence"
	Position      = "position"
	Move          = "move"
	IsBlocking    = "isBlocking"
	IsFloor       = "isFloor"
	Layer100      = "layer100"
	Layer300      = "layer300"
	Layer400      = "layer400"
	Visitable     = "visitable"
	Description   = "description"
	AI            = "ai"
	Stats         = "stats"
	StatsModifier = "statsModifier"
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
	return Layer300
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
	return "something"
}

type AIComponent struct {
}

func (a AIComponent) ComponentType() string {
	return AI
}

type statsValues struct {
	Health     int
	MaxHealth  int
	Defense    int
	MaxDefense int
	Power      int
	MaxPower   int
	Fov        int
}

type StatsComponent struct {
	*statsValues
}

func toStr(i int) string {
	return strconv.Itoa(i)
}

func st(name string, value int, max int) string {
	return name + ": " + toStr(value) + "/" + toStr(max)
}

func (s StatsComponent) String() string {
	return st("HP", s.Health, s.MaxHealth) + " " +
		st("Def", s.Defense, s.MaxDefense) + " " +
		st("Pow", s.Power, s.MaxPower)
}

func (a StatsComponent) ComponentType() string {
	return Stats
}

/*
type StatsModifierComponent struct {
	*statsValues
}

func (a StatsModifierComponent) ComponentType() string {
	return StatsModifier
}
*/

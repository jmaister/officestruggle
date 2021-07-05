package state

const (
	Player     = "player"
	Apparence  = "apparence"
	Position   = "position"
	Move       = "move"
	IsBlocking = "isBlocking"
	Layer100   = "layer100"
	Layer300   = "layer300"
	Layer400   = "layer400"
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

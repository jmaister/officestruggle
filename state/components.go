package state

const (
	Apparence = "apparence"
	Position  = "position"
)

type PlayerComponent struct {
}

type ApparenceComponent struct {
	Color string
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

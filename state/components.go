package state

type Apparence struct {
	Color string
	Char  rune
}

func (a *Apparence) ComponentType() string {
	return "apparence"
}

type Position struct {
	X int
	Y int
}

func (a *Position) ComponentType() string {
	return "position"
}

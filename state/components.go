package state

import "jordiburgos.com/officestruggle/ecs"

type Apparence struct {
	*ecs.Component

	Color string
	Char  rune
}

func NewApparence(color string, char rune) *Apparence {
	return &Apparence{
		Component: &ecs.Component{
			Type: "apparence",
		},
		Color: color,
		Char:  char,
	}
}

type Position struct {
	*ecs.Component
	X int
	Y int
}

package state

import "jordiburgos.com/officestruggle/grid"

type GameState struct {
	Fov  *FieldOfVision
	Grid *grid.Grid
}

func NewGameState(grid *grid.Grid) *GameState {
	fov := FieldOfVision{}
	fov.SetTorchRadius(6)

	return &GameState{
		Fov:  &fov,
		Grid: grid,
	}
}

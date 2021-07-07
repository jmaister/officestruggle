package state

import (
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/grid"
)

type GameState struct {
	Fov          *FieldOfVision
	Grid         *grid.Grid
	Player       *ecs.Entity
	IsPlayerTurn bool
}

func NewGameState(grid *grid.Grid, player *ecs.Entity) *GameState {
	fov := FieldOfVision{}
	fov.SetTorchRadius(6)

	return &GameState{
		Fov:          &fov,
		Grid:         grid,
		Player:       player,
		IsPlayerTurn: true,
	}
}

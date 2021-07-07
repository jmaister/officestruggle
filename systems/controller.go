package systems

import (
	tl "github.com/JoelOtter/termloop"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/state"
)

type Controller struct {
	*tl.Entity
	Engine    *ecs.Engine
	Grid      *grid.Grid
	GameState *state.GameState
}

func (ctl *Controller) Tick(event tl.Event) {
	if event.Type == tl.EventKey && ctl.GameState.IsPlayerTurn {
		var move state.MoveComponent
		switch event.Key { // If so, switch on the pressed key.
		case tl.KeyArrowRight:
			move = state.MoveComponent{
				X: 1, Y: 0,
			}
		case tl.KeyArrowLeft:
			move = state.MoveComponent{
				X: -1, Y: 0,
			}
		case tl.KeyArrowUp:
			move = state.MoveComponent{
				X: 0, Y: -1,
			}
		case tl.KeyArrowDown:
			move = state.MoveComponent{
				X: 0, Y: 1,
			}
		}

		player := ctl.GameState.Player
		player.AddComponent(state.Move, move)

		ctl.GameState.IsPlayerTurn = false
	}

	// TODO: show description when clicking on a tile

	// This is what defines a turn step
	// systems.Render not needed, done in Draw(...) func
	if !ctl.GameState.IsPlayerTurn {
		AI(ctl.Engine, ctl.GameState)
		ctl.GameState.IsPlayerTurn = true
	}
	Movement(ctl.Engine, ctl.Grid)

}

func (ctl *Controller) Draw(screen *tl.Screen) {
	Render(ctl.Engine, ctl.GameState, screen)
}

package systems

import (
	tl "github.com/JoelOtter/termloop"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/game"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/state"
)

type Controller struct {
	*tl.Entity
	Engine    *ecs.Engine
	Grid      *grid.Grid
	GameState *game.GameState
}

func NewController(gs *game.GameState) *Controller {
	return &Controller{
		Engine:    gs.Engine,
		GameState: gs,
		Grid:      gs.Grid,
	}
}

func (ctl *Controller) Tick(event tl.Event) {

	if event.Type == tl.EventKey && ctl.GameState.IsPlayerTurn {
		var move state.MoveComponent
		switch event.Key {
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
		Movement(ctl.GameState, ctl.Engine, ctl.Grid)

		ctl.GameState.IsPlayerTurn = false
	}

	// TODO: show description when clicking on a tile

	// This is what defines a turn step
	// systems.Render not needed, done in Draw(...) func
	if !ctl.GameState.IsPlayerTurn {
		AI(ctl.Engine, ctl.GameState)
		Movement(ctl.GameState, ctl.Engine, ctl.Grid)

		ctl.GameState.IsPlayerTurn = true
	}
}

func (ctl *Controller) Draw(screen *tl.Screen) {
	Render(ctl.Engine, ctl.GameState, screen)
}

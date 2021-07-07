package systems

import (
	tl "github.com/JoelOtter/termloop"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/state"
)

type Controller struct {
	*tl.Entity
	Engine *ecs.Engine
	Grid   *grid.Grid
}

func (ctl *Controller) Tick(event tl.Event) {
	if event.Type == tl.EventKey {

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

		player := ctl.Engine.Entities.GetEntities([]string{state.Player})[0]
		player.AddComponent(state.Move, move)
	}

	// This is what defines a turn step
	// systems.Render not needed, done in Draw(...) func
	Movement(ctl.Engine, ctl.Grid)

}

func (ctl *Controller) Draw(screen *tl.Screen) {
	Render(ctl.Engine, screen)
}

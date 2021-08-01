package systems

import (
	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/state"
)

func ComputeFov(engine *ecs.Engine, gs *gamestate.GameState) {
	stats, _ := gs.Player.GetComponent(constants.Stats).(state.StatsComponent)
	position, _ := gs.Player.GetComponent(constants.Position).(state.PositionComponent)

	// TODO: check if radius changed to recalcualte
	gs.Fov.Compute(gs, position.X, position.Y, stats.Fov)
}

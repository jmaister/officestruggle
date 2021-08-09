package systems

import (
	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/state"
)

func UseStairs(gs *gamestate.GameState, stairs *ecs.Entity) {
	cmp, ok := stairs.GetComponent(constants.Stairs).(state.StairsComponent)
	if ok {
		gs.Player.AddComponent(state.MoveComponent{
			Absolute: true,
			X:        cmp.TargetX,
			Y:        cmp.TargetY,
			Z:        cmp.TargetZ,
		})
	} else {
		gs.Log(constants.Warn, "Not valid stairs.")
	}
}

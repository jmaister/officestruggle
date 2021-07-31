package systems

import (
	"math"

	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/state"
)

func HandleMovementKey(engine *ecs.Engine, gs *gamestate.GameState, dx int, dy int) {
	gs.Player.AddComponent(state.MoveComponent{X: dx, Y: dy})
}

func Movement(engine *ecs.Engine, gs *gamestate.GameState, g *grid.Grid) {
	movable := []string{constants.Move}

	for _, entity := range engine.Entities.GetEntities(movable) {
		move := entity.GetComponent(constants.Move).(state.MoveComponent)
		entity.RemoveComponent(constants.Move)

		position, _ := entity.GetComponent(constants.Position).(state.PositionComponent)

		mx := position.X + move.X
		my := position.Y + move.Y

		// Check map boundaries
		m := g.Map
		mx = int(math.Min(float64(m.Width+m.X-1), math.Max(float64(m.X), float64(mx))))
		my = int(math.Min(float64(m.Height+m.Y-1), math.Max(float64(m.Y), float64(my))))

		// Check for blockers
		newPosition := state.PositionComponent{
			X: mx,
			Y: my,
		}
		entitiesOnPosition, _ := engine.PosCache.Get(newPosition.GetKey())
		blockersOnPosition := entitiesOnPosition.GetEntities([]string{constants.IsBlocking})
		isBlocked := len(blockersOnPosition) > 0

		if isBlocked {
			Attack(engine, gs, entity, blockersOnPosition)
		} else {
			entity.ReplaceComponent(newPosition)
		}
	}

}

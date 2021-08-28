package systems

import (
	"fmt"
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
	entities := engine.Entities.GetEntities([]string{constants.Move})
	for _, entity := range entities {
		move := entity.GetComponent(constants.Move).(state.MoveComponent)
		entity.RemoveComponent(constants.Move)

		position, _ := entity.GetComponent(constants.Position).(state.PositionComponent)

		var newPosition state.PositionComponent
		if move.Absolute {
			newPosition = state.PositionComponent{
				X: move.X,
				Y: move.Y,
				Z: move.Z,
			}
		} else {
			mx := position.X + move.X
			my := position.Y + move.Y
			mz := position.Z + move.Z

			// Check map boundaries
			m := g.Map
			mx = int(math.Min(float64(m.Width+m.X-1), math.Max(float64(m.X), float64(mx))))
			my = int(math.Min(float64(m.Height+m.Y-1), math.Max(float64(m.Y), float64(my))))
			mz = int(math.Min(float64(gs.Grid.Levels), math.Max(0, float64(mz))))

			// Check for blockers
			newPosition = state.PositionComponent{
				X: mx,
				Y: my,
				Z: mz,
			}
		}

		entitiesOnPosition, _ := engine.PosCache.Get(newPosition.GetKey())
		blockersOnPosition := entitiesOnPosition.GetEntities([]string{constants.IsBlocking})
		isBlocked := len(blockersOnPosition) > 0

		if isBlocked {
			Attack(engine, gs, entity, blockersOnPosition)
		} else {
			entity.ReplaceComponent(newPosition)
		}

		if entity == gs.Player {
			// Move camera
			mapWidth := gs.Grid.Map.Width
			mapHeight := gs.Grid.Map.Height
			gs.Camera.MoveCamera(newPosition.X, newPosition.Y, mapWidth, mapHeight)

			// Player changes level
			if newPosition.Z != gs.CurrentZ {
				gs.Log(constants.Good, fmt.Sprintf("Moving to floor %d.", newPosition.Z+1))
				gs.CurrentZ = newPosition.Z

				ComputeFov(engine, gs)
			}

		}

	}

}

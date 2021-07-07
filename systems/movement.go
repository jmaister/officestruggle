package systems

import (
	"fmt"
	"math"

	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/state"
)

func Movement(engine *ecs.Engine, g *grid.Grid) {
	movable := []string{state.Move}

	for _, entity := range engine.Entities.GetEntities(movable) {
		move := entity.RemoveComponent(state.Move).(state.MoveComponent)
		position, _ := entity.GetComponent(state.Position).(state.PositionComponent)

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
		blockersOnPosition := entitiesOnPosition.GetEntities([]string{state.IsBlocking})
		isBlocked := len(blockersOnPosition) > 0
		for _, blocker := range blockersOnPosition {
			attackerStr := state.GetDescription(entity)
			blockerStr := state.GetDescription(blocker)
			fmt.Println(attackerStr + " kicked a " + blockerStr)
		}

		if !isBlocked {
			entity.AddComponent(state.Position, newPosition)

			// Update cache
			engine.PosCache.Delete(position.GetKey(), entity)
			engine.PosCache.Add(newPosition.GetKey(), entity)
		}
	}

}

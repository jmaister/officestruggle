package systems

import (
	"math"
	"strconv"

	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/state"
)

func Movement(gs *gamestate.GameState, engine *ecs.Engine, g *grid.Grid) {
	movable := []string{state.Move}

	for _, entity := range engine.Entities.GetEntities(movable) {
		gs.L.Println("Movement Entity " + entity.String())
		move := entity.RemoveComponent(state.Move).(state.MoveComponent)
		position, _ := entity.GetComponent(state.Position).(state.PositionComponent)

		gs.L.Println("from position " + position.String())
		mx := position.X + move.X
		my := position.Y + move.Y

		gs.L.Println("to position " + strconv.Itoa(mx) + " " + strconv.Itoa(my))

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

		if isBlocked {
			Attack(gs, entity, blockersOnPosition)
		} else {
			entity.AddComponent(state.Position, newPosition)

			// Update cache
			engine.PosCache.Delete(position.GetKey(), entity)
			engine.PosCache.Add(newPosition.GetKey(), entity)
		}
	}

}

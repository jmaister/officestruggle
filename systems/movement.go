package systems

import (
	"math"

	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/state"
)

/*
func getPosition(entity *ecs.Entity) (state.PositionComponent, bool) {
	position, ok := entity.GetComponent(state.Position).(state.PositionComponent)
	return position, ok
}
*/

func Movement(engine *ecs.Engine, g *grid.Grid) {
	movable := []string{state.Move}
	// blockers := []string{state.IsBlocking, state.Position}

	for _, entity := range engine.Entities.GetEntities(movable) {
		move := entity.RemoveComponent(state.Move).(state.MoveComponent)
		position, _ := entity.GetComponent(state.Position).(state.PositionComponent)

		mx := position.X + move.X
		my := position.Y + move.Y

		// Check map boundaries
		m := g.Map
		mx = int(math.Min(float64(m.Width+m.X-1), math.Max(21, float64(mx))))
		my = int(math.Min(float64(m.Height+m.Y-1), math.Max(3, float64(my))))

		// Check for blockers
		newPosition := state.PositionComponent{
			X: mx,
			Y: my,
		}
		entitiesOnPosition, found := engine.PosCache.Get(newPosition.GetKey())
		isBlocked := found
		for _, entity := range entitiesOnPosition {
			isBlocked = isBlocked && entity.HasComponent(state.IsBlocking)
		}

		if !isBlocked {
			entity.AddComponent(state.Position, newPosition)

			engine.PosCache.Delete(position.GetKey(), entity)
			engine.PosCache.Add(position.GetKey(), entity)
		}
	}

}

package systems

import (
	"math"

	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/state"
)

func Movement(engine *ecs.Engine, g *grid.Grid) {
	movable := []string{state.Move}
	for _, entity := range engine.GetEntities(movable) {
		move := entity.RemoveComponent(state.Move).(state.MoveComponent)
		position, _ := entity.GetComponent(state.Position).(state.PositionComponent)

		mx := position.X + move.X
		my := position.Y + move.Y

		// Check map boundaries
		m := g.Map
		mx = int(math.Min(float64(m.Width+m.X-1), math.Max(21, float64(mx))))
		my = int(math.Min(float64(m.Height+m.Y-1), math.Max(3, float64(my))))

		entity.AddComponent(state.Position, state.PositionComponent{
			X: mx,
			Y: my,
		})
	}

}

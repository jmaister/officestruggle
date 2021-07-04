package systems

import (
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/state"
)

func Movement(engine *ecs.Engine) {
	movable := []string{state.Move}
	for _, entity := range engine.GetEntities(movable) {
		move := entity.RemoveComponent(state.Move).(state.MoveComponent)
		position, _ := entity.GetComponent(state.Position).(state.PositionComponent)

		entity.AddComponent(state.Position, state.PositionComponent{
			X: position.X + move.X,
			Y: position.Y + move.Y,
		})
	}

}

package systems

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"jordiburgos.com/officestruggle/animations"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/state"
)

func AnimationSystem(engine *ecs.Engine, gameState *gamestate.GameState, screen *ebiten.Image) {
	entities := engine.Entities.GetEntities([]string{state.Animated})

	for _, entity := range entities {
		animatedCmp, ok := entity.GetComponent(state.Animated).(animations.AnimatedComponent)
		if ok {
			animation := animatedCmp.Animation

			now := time.Now()
			start := animation.StartTime()
			end := start.Add(animation.Duration())

			remaining := end.Sub(now).Milliseconds()
			if remaining >= 0 {
				percent := float64(remaining) / float64(animation.Duration().Milliseconds())
				animation.Update(percent, gameState, screen)
			} else {
				entity.RemoveComponent(state.Animated)
				engine.DestroyEntity(entity)
			}

		}
	}

}

package systems

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/state"
)

func Animation(engine *ecs.Engine, gameState *gamestate.GameState, screen *ebiten.Image) {
	entities := engine.Entities.GetEntities([]string{state.Animated})

	if len(entities) > 0 {
		fmt.Println("runnign animation system", len(entities))
	}

	for _, entity := range entities {
		animatedCmp, ok := entity.GetComponent(state.Animated).(state.AnimatedComponent)
		if ok {
			animation := animatedCmp.Animation

			now := time.Now()
			start := animation.StartTime()
			end := start.Add(animation.Duration())

			remaining := end.Sub(now).Milliseconds()
			if remaining >= 0 {
				percent := float64(remaining) / float64(animation.Duration().Milliseconds())
				animation.Update(percent, screen)
			} else {
				entity.RemoveComponent(state.Animated)
				engine.DestroyEntity(entity)
			}
		}
	}

}

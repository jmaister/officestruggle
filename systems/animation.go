package systems

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"jordiburgos.com/officestruggle/animations"
	"jordiburgos.com/officestruggle/assets"
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

type DamageAnimation struct {
	X                 int
	Y                 int
	Damage            string
	AnimationStart    time.Time
	AnimationDuration time.Duration
}

func (a DamageAnimation) StartTime() time.Time {
	return a.AnimationStart
}
func (a DamageAnimation) Duration() time.Duration {
	return a.AnimationDuration
}
func (a DamageAnimation) Update(percent float64, gs *gamestate.GameState, screen *ebiten.Image) {
	x, y := ToPixel(gs, a.X, a.Y)

	y = y - int(float64(3*gs.TileHeight)*(1-percent))

	fnt := assets.MplusFont(20)
	text.Draw(screen, a.Damage, fnt, x, y, color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	})
}

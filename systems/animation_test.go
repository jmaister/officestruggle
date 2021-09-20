package systems_test

import (
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/game"
	"jordiburgos.com/officestruggle/interfaces"
	"jordiburgos.com/officestruggle/state"
	"jordiburgos.com/officestruggle/systems"
)

func TestAnimation(t *testing.T) {

	engine := ecs.NewEngine()
	gs := game.NewGameState(engine)

	e := engine.NewEntity()

	e.AddComponent(state.AnimatedComponent{
		Animation: systems.LevelUpAnimation{
			AnimationInfo: interfaces.AnimationInfo{
				StartTime: time.Now(),
				Duration:  750 * time.Millisecond,
				Source: interfaces.Point{
					X: 30,
					Y: 30,
				},
			},
		}})

	screen := ebiten.NewImage(gs.ScreenWidth, gs.ScreenHeight)

	for i := 0; i < 1000; i++ {
		systems.AnimationSystem(engine, gs, screen)
		systems.Render(engine, gs, screen)
	}
}

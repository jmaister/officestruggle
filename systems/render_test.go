package systems_test

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/game"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/systems"
)

func TestRender(t *testing.T) {
	engine := ecs.NewEngine()
	gs := game.NewGameState(engine)

	gs.ScreenState = gamestate.GameScreen
	systems.ComputeFov(engine, gs)

	screen := ebiten.NewImage(gs.ScreenWidth, gs.ScreenHeight)

	systems.Render(engine, gs, screen)
}

func BenchmarkRender(b *testing.B) {
	engine := ecs.NewEngine()
	gs := game.NewGameState(engine)
	gs.ScreenState = gamestate.GameScreen

	gs.ScreenState = gamestate.GameScreen
	systems.ComputeFov(engine, gs)

	screen := ebiten.NewImage(gs.ScreenWidth, gs.ScreenHeight)

	for i := 0; i < b.N; i++ {
		systems.Render(engine, gs, screen)

	}
}

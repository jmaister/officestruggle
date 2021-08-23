package systems_test

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/game"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/systems"
)

func TestRender(t *testing.T) {
	g := game.NewGame()
	engine := g.Engine
	gs := g.GameState

	gs.ScreenState = gamestate.GameScreen
	systems.ComputeFov(engine, gs)
	gs.Log(constants.Good, "Testing the game.")

	screen := ebiten.NewImage(gs.ScreenWidth, gs.ScreenHeight)

	systems.Render(engine, gs, screen)
}

func BenchmarkRender(b *testing.B) {
	engine := ecs.NewEngine()
	gs := game.NewGameState(engine)

	gs.ScreenState = gamestate.GameScreen
	systems.ComputeFov(engine, gs)
	gs.Log(constants.Good, "Testing the game.")

	screen := ebiten.NewImage(gs.ScreenWidth, gs.ScreenHeight)

	for i := 0; i < b.N; i++ {
		systems.Render(engine, gs, screen)
	}
}

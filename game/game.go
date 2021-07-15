package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/state"
	"jordiburgos.com/officestruggle/systems"
)

type Game struct {
	Engine    *ecs.Engine
	GameState *gamestate.GameState
}

func NewGame() *Game {
	engine := ecs.NewEngine()

	return &Game{
		// ECS engine
		Engine:    engine,
		GameState: gamestate.NewGameState(engine),
	}
}

func (g *Game) Update() error {

	player := g.GameState.Player
	position := state.GetPosition(player)
	stats, _ := player.GetComponent(state.Stats).(state.StatsComponent)
	g.GameState.Fov.Compute(g.GameState, position.X, position.Y, stats.Fov)

	// Update the logical state
	keys := inpututil.PressedKeys()
	hasPressedKeys := len(keys) > 0 && inpututil.IsKeyJustPressed(keys[0])

	// TODO: https://github.com/hajimehoshi/ebiten/issues/648

	if g.GameState.IsPlayerTurn && hasPressedKeys {
		fmt.Println(keys)

		actionedKey := false
		dx := 0
		dy := 0
		switch keys[0] {
		case ebiten.KeyArrowRight:
			dx = 1
			actionedKey = true
		case ebiten.KeyArrowLeft:
			dx = -1
			actionedKey = true
		case ebiten.KeyArrowUp:
			dy = -1
			actionedKey = true
		case ebiten.KeyArrowDown:
			dy = 1
			actionedKey = true
		}

		if actionedKey {
			player.AddComponent(state.Move, state.MoveComponent{X: dx, Y: dy})
			systems.Movement(g.GameState, g.Engine, g.GameState.Grid)

			g.GameState.IsPlayerTurn = false
		}

	}

	if !g.GameState.IsPlayerTurn {
		systems.AI(g.Engine, g.GameState)
		systems.Movement(g.GameState, g.Engine, g.GameState.Grid)

		g.GameState.IsPlayerTurn = true
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Render the screen
	systems.Render(g.Engine, g.GameState, screen)
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return g.GameState.ScreenWidth, g.GameState.ScreenHeight
}

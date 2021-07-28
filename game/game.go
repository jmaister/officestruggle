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
	hasPressedKeys := len(keys) > 0 && repeatingKeyPressed(keys[0])

	if g.GameState.ScreenState == gamestate.GameScreen {
		if g.GameState.IsPlayerTurn && hasPressedKeys {
			fmt.Println(keys)

			movementKey := false
			inventoryKey := false

			dx := 0
			dy := 0
			switch keys[0] {
			case ebiten.KeyArrowRight:
				dx = 1
				movementKey = true
			case ebiten.KeyArrowLeft:
				dx = -1
				movementKey = true
			case ebiten.KeyArrowUp:
				dy = -1
				movementKey = true
			case ebiten.KeyArrowDown:
				dy = 1
				movementKey = true
			case ebiten.KeyG:
				inventoryKey = true
			case ebiten.KeyI:
				g.GameState.ScreenState = gamestate.InventoryScreen
			case ebiten.KeyZ:
				g.GameState.ScreenState = gamestate.TargetingScreen
			}

			if movementKey {
				player.AddComponent(state.Move, state.MoveComponent{X: dx, Y: dy})
				systems.Movement(g.Engine, g.GameState, g.GameState.Grid)

				g.GameState.IsPlayerTurn = false
			} else if inventoryKey {
				systems.InventoryPickUp(g.GameState)

				g.GameState.IsPlayerTurn = false
			}

		}

		if !g.GameState.IsPlayerTurn {
			systems.AI(g.Engine, g.GameState)
			systems.Movement(g.Engine, g.GameState, g.GameState.Grid)

			g.GameState.IsPlayerTurn = true
		}
	} else if g.GameState.ScreenState == gamestate.WelcomeScreen {
		if hasPressedKeys {
			if keys[0] == ebiten.KeyEnter {
				g.GameState.ScreenState = gamestate.GameScreen
			} else if keys[0] == ebiten.KeyT {
				g.GameState.ScreenState = gamestate.TestScreen
			}
		}
	} else if g.GameState.ScreenState == gamestate.InventoryScreen {
		if hasPressedKeys {
			fmt.Println(keys)
			// TODO: move key handling to inventory.go
			if keys[0] == ebiten.KeyI || keys[0] == ebiten.KeyEscape {
				// Hide inventory screen
				g.GameState.ScreenState = gamestate.GameScreen
			} else if keys[0] == ebiten.KeyUp {
				// Selected item up
				systems.InventoryKeyUp(g.GameState)
			} else if keys[0] == ebiten.KeyDown {
				// Selected item down
				systems.InventoryKeyDown(g.GameState)
			} else if keys[0] == ebiten.KeyLeft {
				// Selected item left
				systems.InventoryKeyLeft(g.GameState)
			} else if keys[0] == ebiten.KeyRight {
				// Selected item right
				systems.InventoryKeyRight(g.GameState)
			} else if keys[0] == ebiten.KeyC {
				// Consume
				systems.InventoryConsume(g.GameState)
			} else if keys[0] == ebiten.KeyD {
				// Drop
				systems.InventoryDrop(g.GameState)
			} else if keys[0] == ebiten.KeyE {
				// Equip
				systems.InventoryEquip(g.GameState)
			} else if keys[0] == ebiten.KeyU {
				// Unequip
				systems.InventoryUnequip(g.GameState)
			}
		}
	} else if g.GameState.ScreenState == gamestate.TargetingScreen {
		if hasPressedKeys {
			if keys[0] == ebiten.KeyZ || keys[0] == ebiten.KeyEscape {
				g.GameState.ScreenState = gamestate.GameScreen
			}
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mouseX, mouseY := ebiten.CursorPosition()
			systems.TargetingMouseClick(g.Engine, g.GameState, mouseX, mouseY)
		}
	} else if g.GameState.ScreenState == gamestate.TestScreen {
		if hasPressedKeys {
			fmt.Println(keys)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.GameState.ScreenState == gamestate.GameScreen {
		// Render the screen
		systems.Render(g.Engine, g.GameState, screen)
	} else if g.GameState.ScreenState == gamestate.TargetingScreen {
		// Render the screen
		systems.Render(g.Engine, g.GameState, screen)
		systems.RenderTargetingScreen(g.Engine, g.GameState, screen)
	} else if g.GameState.ScreenState == gamestate.WelcomeScreen {
		systems.RenderWelcomesScreen(g.Engine, g.GameState, screen)
	} else if g.GameState.ScreenState == gamestate.InventoryScreen {
		systems.RenderInventoryScreen(g.Engine, g.GameState, screen)
	} else if g.GameState.ScreenState == gamestate.TestScreen {
		systems.RenderTestScreen(g.Engine, g.GameState, screen)
	}

	// Update active animations
	systems.AnimationSystem(g.Engine, g.GameState, screen)
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return g.GameState.ScreenWidth, g.GameState.ScreenHeight
}

// repeatingKeyPressed return true when key is pressed considering the repeat state.
// https://github.com/hajimehoshi/ebiten/issues/648
func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 5
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

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
			case ebiten.KeyEscape:
				// Just for debug
				panic("ESCAPE !!!")
			case ebiten.KeyG:
				inventoryKey = true
			case ebiten.KeyI:
				g.GameState.ScreenState = gamestate.InventoryScreen
			}

			if movementKey {
				player.AddComponent(state.Move, state.MoveComponent{X: dx, Y: dy})
				systems.Movement(g.GameState, g.Engine, g.GameState.Grid)

				g.GameState.IsPlayerTurn = false
			} else if inventoryKey {
				systems.InventoryPickUp(g.GameState)

				g.GameState.IsPlayerTurn = false
			}

		}

		if !g.GameState.IsPlayerTurn {
			systems.AI(g.Engine, g.GameState)
			systems.Movement(g.GameState, g.Engine, g.GameState.Grid)

			g.GameState.IsPlayerTurn = true
		}
	} else if g.GameState.ScreenState == gamestate.WelcomeScreen {
		if hasPressedKeys && keys[0] == ebiten.KeyEnter {
			g.GameState.ScreenState = gamestate.GameScreen
		}
	} else if g.GameState.ScreenState == gamestate.InventoryScreen {
		if hasPressedKeys {
			inventory, _ := player.GetComponent(state.Inventory).(state.InventoryComponent)
			if keys[0] == ebiten.KeyI {
				// Hide inventory screen
				g.GameState.ScreenState = gamestate.GameScreen
			} else if keys[0] == ebiten.KeyUp {
				// Selected item up
				g.GameState.InventoryScreen.Selected = g.GameState.InventoryScreen.Selected - 1
			} else if keys[0] == ebiten.KeyDown {
				// Selected item down
				g.GameState.InventoryScreen.Selected = g.GameState.InventoryScreen.Selected + 1
			} else if keys[0] == ebiten.KeyC {
				// Consume
				sel := g.GameState.InventoryScreen.Selected
				if sel >= 0 && sel < len(inventory.Items) {
					consumable := inventory.Items[g.GameState.InventoryScreen.Selected]
					systems.InventoryConsume(g.GameState, consumable)
				}
			} else if keys[0] == ebiten.KeyD {
				// Drop
				sel := g.GameState.InventoryScreen.Selected
				if sel >= 0 && sel < len(inventory.Items) {
					consumable := inventory.Items[g.GameState.InventoryScreen.Selected]
					systems.InventoryDrop(g.GameState, consumable)
				}
			}

			// Update inventory as list could have changed
			inventory, _ = player.GetComponent(state.Inventory).(state.InventoryComponent)
			if g.GameState.InventoryScreen.Selected < 0 {
				g.GameState.InventoryScreen.Selected = 0
			}
			if len(inventory.Items) > 0 && g.GameState.InventoryScreen.Selected >= len(inventory.Items) {
				g.GameState.InventoryScreen.Selected = len(inventory.Items) - 1
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.GameState.ScreenState == gamestate.GameScreen {
		// Render the screen
		systems.Render(g.Engine, g.GameState, screen)
	} else if g.GameState.ScreenState == gamestate.WelcomeScreen {
		systems.RenderWelcomesScreen(g.Engine, g.GameState, screen)
	} else if g.GameState.ScreenState == gamestate.InventoryScreen {
		systems.RenderInventoryScreen(g.Engine, g.GameState, screen)
	}
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

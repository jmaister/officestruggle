package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"jordiburgos.com/officestruggle/assets"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
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
		GameState: NewGameState(engine),
	}
}

func (g *Game) Update() error {

	// Update the logical state
	keys := inpututil.PressedKeys()
	hasPressedKeys := len(keys) > 0 && repeatingKeyPressed(keys[0])

	if g.GameState.ScreenState == gamestate.GameScreen {
		systems.ComputeFov(g.Engine, g.GameState)
		if g.GameState.IsPlayerTurn && hasPressedKeys {
			fmt.Println(keys)

			movementKey := false
			inventoryPickupKey := false

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
				inventoryPickupKey = true
			case ebiten.KeyI:
				g.GameState.ScreenState = gamestate.InventoryScreen
			case ebiten.KeyZ:
				g.GameState.ScreenState = gamestate.TargetingScreen
			case ebiten.KeyS:
				// Save game
				if inpututil.IsKeyJustPressed(ebiten.KeyS) {
					saveFileName, err := systems.CreateTimestamSavegame()
					if err != nil {
						fmt.Printf("Error saving the game: %s\n", err)
					} else {
						systems.SaveGame(g.Engine, g.GameState, saveFileName)
					}
				}
			/* case ebiten.KeyL:
			// Load game
			if inpututil.IsKeyJustPressed(ebiten.KeyL) {
				g.GameState.ScreenState = gamestate.LoadingScreen
				go systems.LoadGame(g.Engine, g.GameState)
			}
			*/
			case ebiten.KeyEnter:
				g.GameState.ActionScreenState.Actions.IsFocused = true
				g.GameState.ActionScreenState.Actions.Selected = 0
				g.GameState.ActionScreenState.Items = nil
				g.GameState.ScreenState = gamestate.ActionDialog
			}

			if movementKey {
				systems.HandleMovementKey(g.Engine, g.GameState, dx, dy)
				systems.Movement(g.Engine, g.GameState, g.GameState.Grid)

				g.GameState.IsPlayerTurn = false
			} else if inventoryPickupKey {
				systems.InventoryPickUpItemsOnFloor(g.GameState)

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
				systems.InventoryConsume(g.Engine, g.GameState)
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
	} else if g.GameState.ScreenState == gamestate.GameoverScreen {
		if hasPressedKeys {
			if keys[0] == ebiten.KeyEnter {
				engine := ecs.NewEngine()
				g.Engine = engine
				g.GameState = NewGameState(engine)
			}
		}
	} else if g.GameState.ScreenState == gamestate.ActionDialog {
		if hasPressedKeys {
			if keys[0] == ebiten.KeyUp {
				// Selected item up
				systems.ActionDialogKeyUp(g.GameState)
			} else if keys[0] == ebiten.KeyDown {
				// Selected item down
				systems.ActionDialogKeyDown(g.GameState)
			} else if keys[0] == ebiten.KeyEnter {
				// Selected item down
				systems.ActionDialogActivate(g.Engine, g.GameState)
				systems.Movement(g.Engine, g.GameState, g.GameState.Grid)

				g.GameState.IsPlayerTurn = false
			} else if keys[0] == ebiten.KeyEscape || keys[0] == ebiten.KeyQ {
				// Selected item down
				systems.ActionDialogExit(g.GameState)
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.GameState.ScreenState == gamestate.GameScreen {
		// Render the screen
		systems.Render(g.Engine, g.GameState, screen)
		systems.EffectInfoSystem(g.Engine, g.GameState, screen)
		// Update active animations
		systems.AnimationSystem(g.Engine, g.GameState, screen)
	} else if g.GameState.ScreenState == gamestate.ActionDialog {
		// Render the screen
		systems.Render(g.Engine, g.GameState, screen)
		if g.GameState.ScreenState == gamestate.ActionDialog {
			systems.DrawActionDialog(g.Engine, g.GameState, screen)
		}
	} else if g.GameState.ScreenState == gamestate.TargetingScreen {
		// Render the screen
		systems.Render(g.Engine, g.GameState, screen)
		systems.EffectInfoSystem(g.Engine, g.GameState, screen)
		systems.RenderTargetingScreen(g.Engine, g.GameState, screen)
		// Update active animations
		systems.AnimationSystem(g.Engine, g.GameState, screen)
	} else if g.GameState.ScreenState == gamestate.WelcomeScreen {
		systems.RenderWelcomesScreen(g.Engine, g.GameState, screen)
	} else if g.GameState.ScreenState == gamestate.InventoryScreen {
		systems.RenderInventoryScreen(g.Engine, g.GameState, screen)
	} else if g.GameState.ScreenState == gamestate.TestScreen {
		systems.RenderTestScreen(g.Engine, g.GameState, screen)
	} else if g.GameState.ScreenState == gamestate.LoadingScreen {
		systems.DrawText(screen, g.GameState, 10, 10, assets.LoadFontCached(40), "Loading ...", color.White, color.Black)
	} else if g.GameState.ScreenState == gamestate.GameoverScreen {
		systems.DrawText(screen, g.GameState, 10, 10, assets.LoadFontCached(40), "Game Over!", color.White, color.Black)
		systems.DrawText(screen, g.GameState, 10, 20, assets.LoadFontCached(20), "You have been defeated.", color.White, color.Black)
		systems.DrawText(screen, g.GameState, 10, 30, assets.LoadFontCached(20), "Press [ENTER] to continue.", color.White, color.Black)
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

package systems

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"jordiburgos.com/officestruggle/assets"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/state"
)

func RenderInventoryScreen(engine *ecs.Engine, gameState *gamestate.GameState, screen *ebiten.Image) {

	text.Draw(screen, "In ventory", fnt40, 30, 35, color.White)
	text.Draw(screen, "C - Consume, D - Drop", fnt20, 30, 40, color.White)

	drawFullInventory(screen, gameState)
}

func drawFullInventory(screen *ebiten.Image, gs *gamestate.GameState) {
	fontSize := 14
	font := assets.MplusFont(float64(fontSize))

	position := gs.Grid.Inventory
	inventory, _ := gs.Player.GetComponent(state.Inventory).(state.InventoryComponent)

	cl := ParseHexColorFast("#FFFFFF")

	y := position.Y
	status := "Inventory " + strconv.Itoa(len(inventory.Items)) + "/" + strconv.Itoa(inventory.MaxItems) + " - Selected: " + strconv.Itoa(gs.InventoryScreen.Selected)
	text.Draw(screen, status, font, (position.X)*fontSize, y*fontSize, cl)

	if len(inventory.Items) > 0 {
		for i, entity := range inventory.Items {
			selected := (i == gs.InventoryScreen.Selected)

			str := strconv.Itoa(i) + "-" + state.GetLongDescription(entity)
			px := (position.X) * fontSize
			py := (y + 1) * fontSize
			if selected {
				DrawTextRect(screen, str, px, py, font, color.White)
				text.Draw(screen, str, font, px, py, color.Black)
			} else {
				text.Draw(screen, str, font, px, py, color.White)
			}
			y++
		}
	} else {
		text.Draw(screen, "- No items in the inventory -", font, (position.X)*fontSize, (y+1)*fontSize, cl)
	}
}

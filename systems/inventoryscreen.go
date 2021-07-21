package systems

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"jordiburgos.com/officestruggle/assets"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/state"
)

func RenderInventoryScreen(engine *ecs.Engine, gameState *gamestate.GameState, screen *ebiten.Image) {

	text.Draw(screen, "Inventory", fnt40, 30, 35, color.White)
	text.Draw(screen, "C - Consume, D - Drop, E - Equip, U - Unequip", fnt20, 300, 40, color.White)

	// Inventory
	inventory, _ := gameState.Player.GetComponent(state.Inventory).(state.InventoryComponent)
	inventoryTitle := fmt.Sprintf("Inventory %2d/%2d", len(inventory.Items), inventory.MaxItems)
	invStrItems := []string{}
	for _, item := range inventory.Items {
		invStrItems = append(invStrItems, state.GetLongDescription(item))
	}
	drawList(screen, gameState, &gameState.InventoryScreenState.InventoryState, invStrItems, gameState.Grid.Inventory, inventoryTitle)

	// Equipment
	equipment, _ := gameState.Player.GetComponent(state.Equipment).(state.EquipmentComponent)
	equipmentTitle := "Equipment"
	equipStrItems := []string{}
	for _, position := range state.EquipmentPositions {
		item, ok := equipment.Items[position]
		if ok {
			equipStrItems = append(equipStrItems, fmt.Sprintf("%6s: %s", position, state.GetLongDescription(item)))
		} else {
			equipStrItems = append(equipStrItems, fmt.Sprintf("%6s: - empty -", position))
		}
	}
	drawList(screen, gameState, &gameState.InventoryScreenState.EquipmentState, equipStrItems, gameState.Grid.Equipment, equipmentTitle)

}

func drawList(screen *ebiten.Image, gs *gamestate.GameState, listState *gamestate.ListState, items []string, position grid.Rect, title string) {
	fontSize := 12
	font := assets.MplusFont(float64(fontSize))

	y := position.Y
	text.Draw(screen, title, font, (position.X)*fontSize, y*fontSize, color.White)

	if len(items) > 0 {
		for i, itemStr := range items {
			selected := (i == listState.Selected)

			str := fmt.Sprintf("%2d - %s", i+1, itemStr)
			px := (position.X) * fontSize
			py := (y + 1) * fontSize
			if selected && listState.IsFocused {
				DrawTextRect(screen, str, px, py, font, color.White)
				text.Draw(screen, str, font, px, py, color.Black)
			} else {
				text.Draw(screen, str, font, px, py, color.White)
			}
			y++
		}
	} else {
		text.Draw(screen, "- No items -", font, (position.X)*fontSize, (y+1)*fontSize, color.White)
	}
}

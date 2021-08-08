package systems

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/state"
)

func RenderInventoryScreen(engine *ecs.Engine, gameState *gamestate.GameState, screen *ebiten.Image) {

	text.Draw(screen, "Inventory", fnt40, 30, 35, color.White)
	text.Draw(screen, "C - Consume, D - Drop, E - Equip, U - Unequip", fnt20, 300, 40, color.White)

	// Inventory
	inventory, _ := gameState.Player.GetComponent(constants.Inventory).(state.InventoryComponent)
	inventoryTitle := fmt.Sprintf("Inventory %2d/%2d", len(inventory.Items), inventory.MaxItems)
	invStrItems := []string{}
	for _, item := range inventory.Items {
		invStrItems = append(invStrItems, state.GetLongDescription(item))
	}
	DrawSelectionList(screen, &gameState.InventoryScreenState.InventoryState, invStrItems, gameState.Grid.Inventory, inventoryTitle)

	// Equipment
	equipment, _ := gameState.Player.GetComponent(constants.Equipment).(state.EquipmentComponent)
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
	DrawSelectionList(screen, &gameState.InventoryScreenState.EquipmentState, equipStrItems, gameState.Grid.Equipment, equipmentTitle)

}

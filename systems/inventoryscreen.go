package systems

import (
	"fmt"
	"image/color"
	"strconv"

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
	inventoryTitle := "Inventory " + strconv.Itoa(len(inventory.Items)) + "/" + strconv.Itoa(inventory.MaxItems)
	drawList(screen, gameState, &ListState{
		Selected:  0,
		IsFocused: true,
	}, inventory.Items, gameState.Grid.Inventory, inventoryTitle)

	// Equipment
	equipment, _ := gameState.Player.GetComponent(state.Equipment).(state.EquipmentComponent)
	equipmentTitle := "Equipment"
	items := ecs.EntityList{}
	for _, item := range equipment.Items {
		items = append(items, item)
	}
	drawList(screen, gameState, &ListState{
		Selected:  0,
		IsFocused: false,
	}, items, gameState.Grid.Equipment, equipmentTitle)
}

/*
func drawFullInventory(screen *ebiten.Image, gs *gamestate.GameState) {
	fontSize := 14
	font := assets.MplusFont(float64(fontSize))

	position := gs.Grid.Inventory
	inventory, _ := gs.Player.GetComponent(state.Inventory).(state.InventoryComponent)

	cl := ParseHexColorFast("#FFFFFF")

	y := position.Y
	status := "Inventory " + strconv.Itoa(len(inventory.Items)) + "/" + strconv.Itoa(inventory.MaxItems) + " - Selected: " + strconv.Itoa(gs.InventoryScreenState.Selected)
	text.Draw(screen, status, font, (position.X)*fontSize, y*fontSize, cl)

	if len(inventory.Items) > 0 {
		for i, entity := range inventory.Items {
			selected := (i == gs.InventoryScreenState.Selected)

			str := strconv.Itoa(i) + "-" + state.GetLongDescription(entity)
			px := (position.X) * fontSize
			py := (y + 1) * fontSize
			if selected && gs.InventoryScreenState.Focus == gamestate.InventoryFocus {
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

func drawEquipment(screen *ebiten.Image, gs *gamestate.GameState) {
	fontSize := 14
	font := assets.MplusFont(float64(fontSize))

	position := gs.Grid.Equipment
	equipment, _ := gs.Player.GetComponent(state.Equipment).(state.EquipmentComponent)

	y := position.Y
	status := "Equipment"
	text.Draw(screen, status, font, (position.X)*fontSize, y*fontSize, color.White)

	if len(equipment.Items) > 0 {
		i := 0
		for _, entity := range equipment.Items {
			selected := (i == gs.InventoryScreenState.Selected)

			str := strconv.Itoa(i) + "-" + state.GetLongDescription(entity)
			px := (position.X) * fontSize
			py := (y + 1) * fontSize
			if selected && gs.InventoryScreenState.Focus == gamestate.EquipmentFocus {
				DrawTextRect(screen, str, px, py, font, color.White)
				text.Draw(screen, str, font, px, py, color.Black)
			} else {
				text.Draw(screen, str, font, px, py, color.White)
			}
			y++
			i++
		}
	} else {
		text.Draw(screen, "- No items in the inventory -", font, (position.X)*fontSize, (y+1)*fontSize, color.White)
	}

}
*/

type ListState struct {
	Selected  int
	IsFocused bool
}

func drawList(screen *ebiten.Image, gs *gamestate.GameState, listState *ListState, items ecs.EntityList, position grid.Rect, title string) {
	fontSize := 14
	font := assets.MplusFont(float64(fontSize))

	y := position.Y
	text.Draw(screen, title, font, (position.X)*fontSize, y*fontSize, color.White)

	if len(items) > 0 {
		for i, entity := range items {
			selected := (i == listState.Selected)

			str := fmt.Sprintf("%2d - %s", i+1, state.GetLongDescription(entity))
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

func InventoryKeyUp(gs *gamestate.GameState) {
	UpdateInventorySelection(gs, -1)
}

func InventoryKeyDown(gs *gamestate.GameState) {
	UpdateInventorySelection(gs, 1)
}

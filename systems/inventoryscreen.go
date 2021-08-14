package systems

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"jordiburgos.com/officestruggle/assets"
	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/palette"
	"jordiburgos.com/officestruggle/state"
)

func RenderInventoryScreen(engine *ecs.Engine, gameState *gamestate.GameState, screen *ebiten.Image) {

	text.Draw(screen, "Inventory", fnt40, 30, 35, color.White)
	text.Draw(screen, "(C) Consume    (D) Drop    (E) Equip     (U) Unequip", fnt20, 300, 40, color.White)

	// Inventory
	inventory, _ := gameState.Player.GetComponent(constants.Inventory).(state.InventoryComponent)
	inventoryTitle := fmt.Sprintf("Inventory %2d/%2d", len(inventory.Items), inventory.MaxItems)
	invStrItems := []string{}
	for _, item := range inventory.Items {
		invStrItems = append(invStrItems, state.GetDescription(item))
	}
	DrawSelectionList(screen, gameState, &gameState.InventoryScreenState.InventoryState, invStrItems, gameState.Grid.Inventory, inventoryTitle)

	// Equipment
	equipment, _ := gameState.Player.GetComponent(constants.Equipment).(state.EquipmentComponent)
	equipmentTitle := "Equipment"
	equipStrItems := []string{}
	for _, position := range constants.EquipmentSlots {
		item, ok := equipment.Items[position]
		if ok {
			equipStrItems = append(equipStrItems, fmt.Sprintf("%6s: %s", position, state.GetDescription(item)))
		} else {
			equipStrItems = append(equipStrItems, fmt.Sprintf("%6s: - empty -", position))
		}
	}
	DrawSelectionList(screen, gameState, &gameState.InventoryScreenState.EquipmentState, equipStrItems, gameState.Grid.Equipment, equipmentTitle)

	if gameState.InventoryScreenState.InventoryState.Selected < len(inventory.Items) {
		candidate := inventory.Items[gameState.InventoryScreenState.InventoryState.Selected]
		if candidate != nil && candidate.HasComponent(constants.Equipable) {
			drawEquipDiff(screen, gameState, candidate, equipment)
		}
	}
}

var fnt = assets.MplusFont(20)

func drawEquipDiff(screen *ebiten.Image, gs *gamestate.GameState, candidate *ecs.Entity, equipment state.EquipmentComponent) {
	y := 25
	x := 1
	d := 30

	DrawText(screen, gs, x, y, fnt, "To be equipped", color.White, color.Black)
	DrawText(screen, gs, x+d, y, fnt, "Equipped", color.White, color.Black)
	y = y + 2

	candidateEquipable, _ := candidate.GetComponent(constants.Equipable).(state.EquipableComponent)
	candidateStats := candidateEquipable.StatsValues

	equipped := equipment.Items[candidateEquipable.EquipSlot]
	equippedEquipable := state.EquipableComponent{}
	equippedStats := &state.StatsValues{}
	if equipped != nil {
		equippedEquipable = equipped.GetComponent(constants.Equipable).(state.EquipableComponent)
		equippedStats = equippedEquipable.StatsValues
	}

	var levelColor color.Color = color.White
	leveling := gs.Player.GetComponent(constants.Leveling).(state.LevelingComponent)
	if candidateEquipable.Level > leveling.CurrentLevel {
		levelColor = palette.PColor(palette.Red, 0.6)
	} else if candidateEquipable.Level > equippedEquipable.Level {
		levelColor = palette.PColor(palette.Green, 0.6)
	}

	setColor := palette.PColor(palette.Turquoise, 0.5)

	DrawText(screen, gs, x, y, fnt, state.GetDescription(candidate), color.White, color.Black)
	if candidateEquipable.SetName != "" {
		DrawText(screen, gs, x, y+1, fnt, candidateEquipable.SetName, setColor, color.Black)
	}
	DrawText(screen, gs, x, y+2, fnt, fmt.Sprintf("Level %d", candidateEquipable.Level), levelColor, color.Black)
	drawStatImprovement(screen, gs, x, y+3, "Health", candidateStats.Health, candidateStats.MaxHealth, equippedStats.Health, equippedStats.MaxHealth)
	drawStatImprovement(screen, gs, x, y+4, "Def   ", candidateStats.Defense, candidateStats.MaxDefense, equippedStats.Defense, equippedStats.MaxDefense)
	drawStatImprovement(screen, gs, x, y+5, "Power ", candidateStats.Power, candidateStats.MaxPower, equippedStats.Power, equippedStats.MaxPower)
	drawStatImprovement(screen, gs, x, y+6, "FOV   ", candidateStats.Fov, candidateStats.MaxFov, equippedStats.Fov, equippedStats.MaxFov)

	if equipped != nil {
		DrawText(screen, gs, x+d, y, fnt, state.GetDescription(equipped), color.White, color.Black)
		if equippedEquipable.SetName != "" {
			DrawText(screen, gs, x+d, y+1, fnt, equippedEquipable.SetName, setColor, color.Black)
		}
		DrawText(screen, gs, x+d, y+2, fnt, fmt.Sprintf("Level %d", equippedEquipable.Level), color.White, color.Black)
		drawStatImprovement(screen, gs, x+d, y+3, "Health", equippedStats.Health, equippedStats.MaxHealth, candidateStats.Health, candidateStats.MaxHealth)
		drawStatImprovement(screen, gs, x+d, y+4, "Def   ", equippedStats.Defense, equippedStats.MaxDefense, candidateStats.Defense, candidateStats.MaxDefense)
		drawStatImprovement(screen, gs, x+d, y+5, "Power ", equippedStats.Power, equippedStats.MaxPower, candidateStats.Power, candidateStats.MaxPower)
		drawStatImprovement(screen, gs, x+d, y+6, "FOV   ", equippedStats.Fov, equippedStats.MaxFov, candidateStats.Fov, candidateStats.MaxFov)

	}
}

func drawStatImprovement(screen *ebiten.Image, gs *gamestate.GameState, x int, y int, name string, itemValue int, itemValueMax int, valueCompared int, valueComparedMax int) {
	var cl color.Color = color.White
	if (itemValue > valueCompared || itemValueMax > valueComparedMax) && valueCompared != -1 {
		cl = palette.PColor(palette.Green, 0.6)
	} else if itemValue < valueCompared || itemValueMax < valueComparedMax {
		cl = palette.PColor(palette.Orange, 0.6)
	}
	DrawText(screen, gs, x, y, fnt, fmt.Sprintf("%s +%d/+%d", name, itemValue, itemValueMax), cl, color.Black)
}

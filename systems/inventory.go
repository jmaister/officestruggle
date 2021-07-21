package systems

import (
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/state"
)

func InventoryPickUp(gs *gamestate.GameState) {
	player := gs.Player
	position := player.GetComponent(state.Position).(state.PositionComponent)

	pickables, ok := gs.Engine.PosCache.GetByCoordAndComponents(position.X, position.Y, []string{state.IsPickup})
	if ok && len(pickables) > 0 {
		inventory := player.GetComponent(state.Inventory).(state.InventoryComponent)
		for _, pickable := range pickables {
			pickUpOk := inventory.AddItem(pickable)
			if pickUpOk {
				gs.Log(gamestate.Info, state.GetDescription(player)+" picks up "+state.GetLongDescription(pickable))

				pickable.RemoveComponent(state.IsPickup)
				pickable.RemoveComponent(state.Position)
			} else {
				gs.Log(gamestate.Bad, state.GetDescription(player)+" can't pickup, inventory is full.")
			}
		}
		player.ReplaceComponent(state.Inventory, inventory)
	} else {
		gs.Log(gamestate.Warn, "No items to pickup found at this location.")
	}
}

func getCurrentInventoryItem(gs *gamestate.GameState) (*ecs.Entity, bool) {
	player := gs.Player
	inventory, _ := player.GetComponent(state.Inventory).(state.InventoryComponent)

	sel := gs.InventoryScreenState.InventoryState.Selected
	if sel >= 0 && sel < len(inventory.Items) {
		item := inventory.Items[sel]
		return item, true
	}
	return nil, false
}

func getCurrentEquipmentItem(gs *gamestate.GameState) (*ecs.Entity, bool) {
	player := gs.Player
	equipment, _ := player.GetComponent(state.Equipment).(state.EquipmentComponent)

	sel := gs.InventoryScreenState.EquipmentState.Selected
	if sel >= 0 && sel < len(state.EquipmentPositions) {
		pos := state.EquipmentPositions[sel]
		item, ok := equipment.Items[pos]
		if ok {
			return item, true
		}
	}
	return nil, false
}

func InventoryConsume(gs *gamestate.GameState) {
	player := gs.Player
	inventory, _ := player.GetComponent(state.Inventory).(state.InventoryComponent)

	consumable, ok := getCurrentInventoryItem(gs)
	isConsumable := consumable.HasComponent(state.Consumable)
	if ok && isConsumable {
		// TODO: move to it's own System
		// Consume by player
		conStats := consumable.GetComponent(state.Consumable).(state.ConsumableComponent)
		plStats := player.GetComponent(state.Stats).(state.StatsComponent)

		newStats := plStats.Merge(*conStats.StatsValues)
		player.ReplaceComponent(state.Stats, state.StatsComponent{
			StatsValues: &newStats,
		})

		gs.Log(gamestate.Info, "Consumed "+state.GetLongDescription(consumable))

		// Remove from inventory
		inventory.RemoveItem(consumable)
		player.ReplaceComponent(state.Inventory, inventory)

		// Destroy entity
		engine := consumable.Engine
		engine.DestroyEntity(consumable)

		updateInventorySelection(gs, 0)
	} else if !isConsumable {
		gs.Log(gamestate.Warn, state.GetDescription(consumable)+" can't be consumed.")
	} else {
		gs.Log(gamestate.Warn, "No items to consume.")
	}
}

func InventoryDrop(gs *gamestate.GameState) {
	player := gs.Player
	inventory, _ := player.GetComponent(state.Inventory).(state.InventoryComponent)
	position := player.GetComponent(state.Position).(state.PositionComponent)

	invetoryItem, ok := getCurrentInventoryItem(gs)
	if ok {
		// TODO: move to it's own System
		gs.Log(gamestate.Info, "You dropped "+state.GetLongDescription(invetoryItem))

		// Remove from inventory
		inventory.RemoveItem(invetoryItem)
		player.ReplaceComponent(state.Inventory, inventory)

		// Set new position
		invetoryItem.AddComponent(state.Position, state.PositionComponent{
			X: position.X,
			Y: position.Y,
		})
		invetoryItem.AddComponent(state.IsPickup, state.IsPickupComponent{})

		updateInventorySelection(gs, 0)
	} else {
		gs.Log(gamestate.Warn, "No items to drop.")
	}
}

func InventoryEquip(gs *gamestate.GameState) {
	player := gs.Player
	inventory, _ := player.GetComponent(state.Inventory).(state.InventoryComponent)

	item, ok := getCurrentInventoryItem(gs)
	equipable, isEquipable := item.GetComponent(state.Equipable).(state.EquipableComponent)
	if ok && isEquipable {
		// TODO: move to it's own System
		gs.Log(gamestate.Info, "You dropped "+state.GetLongDescription(item))

		// Remove from inventory
		inventory.RemoveItem(item)
		player.ReplaceComponent(state.Inventory, inventory)

		// Add to equip
		equipment := player.GetComponent(state.Equipment).(state.EquipmentComponent)
		equipment.Items[equipable.Position] = item

		// TODO: move to it's own System
		equipment.UpdateStats(player)

		updateInventorySelection(gs, 0)
	} else if !isEquipable {
		gs.Log(gamestate.Warn, state.GetDescription(item)+" can't be equiped.")
	} else {
		gs.Log(gamestate.Warn, "No items to drop.")
	}
}

func InventoryUnequip(gs *gamestate.GameState) {
	player := gs.Player

	item, ok := getCurrentEquipmentItem(gs)
	if ok {
		// TODO: move to it's own System
		gs.Log(gamestate.Info, "You unequipped "+state.GetLongDescription(item))

		// Remove from equip
		equipable, _ := item.GetComponent(state.Equipable).(state.EquipableComponent)
		equipment := player.GetComponent(state.Equipment).(state.EquipmentComponent)
		delete(equipment.Items, equipable.Position)

		// Add to inventory
		// TODO: check if there is space in the inventory
		inventory, _ := player.GetComponent(state.Inventory).(state.InventoryComponent)
		inventory.AddItem(item)
		player.ReplaceComponent(state.Inventory, inventory)

		// TODO: move to it's own System
		equipment.UpdateStats(player)

		updateInventorySelection(gs, 0)
	} else {
		gs.Log(gamestate.Warn, "No items to drop.")
	}
}

func updateInventorySelection(gs *gamestate.GameState, change int) {

	selected := gs.InventoryScreenState.InventoryState.Selected + change

	inventory, _ := gs.Player.GetComponent(state.Inventory).(state.InventoryComponent)
	max := len(inventory.Items)
	if selected < 0 {
		selected = 0
	} else if max > 0 && selected >= max {
		selected = max - 1
	}

	gs.InventoryScreenState.InventoryState.Selected = selected
}

func updateEquipmentSelection(gs *gamestate.GameState, change int) {

	selected := gs.InventoryScreenState.EquipmentState.Selected + change
	max := len(state.EquipmentPositions)
	if selected < 0 {
		selected = 0
	} else if max > 0 && selected >= max {
		selected = max - 1
	}

	gs.InventoryScreenState.EquipmentState.Selected = selected
}

func InventoryKeyUp(gs *gamestate.GameState) {
	if gs.InventoryScreenState.InventoryState.IsFocused {
		updateInventorySelection(gs, -1)
	} else if gs.InventoryScreenState.EquipmentState.IsFocused {
		updateEquipmentSelection(gs, -1)
	}
}

func InventoryKeyDown(gs *gamestate.GameState) {
	if gs.InventoryScreenState.InventoryState.IsFocused {
		updateInventorySelection(gs, 1)
	} else if gs.InventoryScreenState.EquipmentState.IsFocused {
		updateEquipmentSelection(gs, 1)
	}
}

func InventoryKeyLeft(gs *gamestate.GameState) {
	gs.InventoryScreenState.InventoryState.IsFocused = true
	gs.InventoryScreenState.EquipmentState.IsFocused = false
}

func InventoryKeyRight(gs *gamestate.GameState) {
	gs.InventoryScreenState.EquipmentState.IsFocused = true
	gs.InventoryScreenState.InventoryState.IsFocused = false
}

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

func getCurrentSelectedItem(gs *gamestate.GameState) (*ecs.Entity, bool) {
	player := gs.Player
	inventory, _ := player.GetComponent(state.Inventory).(state.InventoryComponent)

	sel := gs.InventoryScreenState.Selected
	if sel >= 0 && sel < len(inventory.Items) {
		item := inventory.Items[gs.InventoryScreenState.Selected]
		return item, true
	}
	return nil, false
}

func InventoryConsume(gs *gamestate.GameState) {
	player := gs.Player
	inventory, _ := player.GetComponent(state.Inventory).(state.InventoryComponent)

	consumable, ok := getCurrentSelectedItem(gs)
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

		UpdateInventorySelection(gs, 0)
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

	invetoryItem, ok := getCurrentSelectedItem(gs)
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

		UpdateInventorySelection(gs, 0)
	} else {
		gs.Log(gamestate.Warn, "No items to drop.")
	}
}

func InventoryEquip(gs *gamestate.GameState) {
	player := gs.Player
	inventory, _ := player.GetComponent(state.Inventory).(state.InventoryComponent)

	item, ok := getCurrentSelectedItem(gs)
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

		UpdateInventorySelection(gs, 0)
	} else if !isEquipable {
		gs.Log(gamestate.Warn, state.GetDescription(item)+" can't be equiped.")
	} else {
		gs.Log(gamestate.Warn, "No items to drop.")
	}
}

func UpdateInventorySelection(gs *gamestate.GameState, change int) {

	selected := gs.InventoryScreenState.Selected + change

	inventory, _ := gs.Player.GetComponent(state.Inventory).(state.InventoryComponent)
	if selected < 0 {
		selected = 0
	} else if len(inventory.Items) > 0 && selected >= len(inventory.Items) {
		selected = len(inventory.Items) - 1
	}

	gs.InventoryScreenState.Selected = selected
}

func InventoryKeyLeft(gs *gamestate.GameState) {
	gs.InventoryScreenState.Selected = 0
	gs.InventoryScreenState.Focus = gamestate.InventoryFocus
}

func InventoryKeyRight(gs *gamestate.GameState) {
	gs.InventoryScreenState.Selected = 0
	gs.InventoryScreenState.Focus = gamestate.EquipmentFocus
}

package systems

import (
	"fmt"

	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/state"
)

func InventoryPickUpItemsOnFloor(gs *gamestate.GameState) {
	player := gs.Player
	position := player.GetComponent(constants.Position).(state.PositionComponent)

	pickables, ok := gs.Engine.PosCache.GetByCoordAndComponents(position.X, position.Y, position.Z, []string{constants.IsPickup})
	if ok && len(pickables) > 0 {
		for _, pickable := range pickables {
			PickupEntity(gs, pickable)
		}
	} else {
		gs.Log(constants.Warn, "No items to pickup found at this location.")
	}
}

func PickupEntity(gs *gamestate.GameState, pickable *ecs.Entity) {
	player := gs.Player
	inventory := player.GetComponent(constants.Inventory).(state.InventoryComponent)

	moneyComponent, isMoney := pickable.GetComponent(constants.Money).(state.MoneyComponent)
	if isMoney {
		inventory.Coins += moneyComponent.Coins
		gold, silver, copper := Coins(moneyComponent.Coins)
		moneyStr := ""
		if gold > 0 {
			moneyStr += fmt.Sprintf("%d Gold ", gold)
		}
		if silver > 0 {
			moneyStr += fmt.Sprintf("%d Silver ", silver)
		}
		if copper > 0 {
			moneyStr += fmt.Sprintf("%d Copper", copper)
		}
		gs.Log(constants.Good, fmt.Sprintf("You found %s coins.", moneyStr))
		ecs.NewEngine().DestroyEntity(pickable)
	} else if pickable.HasComponent(constants.IsPickup) {
		pickUpOk := inventory.AddItem(pickable)
		if pickUpOk {
			gs.Log(constants.Info, state.GetDescription(player)+" picks up "+state.GetLongDescription(pickable))

			pickable.RemoveComponent(constants.IsPickup)
			pickable.RemoveComponent(constants.Position)
		} else {
			gs.Log(constants.Bad, state.GetDescription(player)+" can't pickup, inventory is full.")
		}
	} else {
		gs.Log(constants.Bad, fmt.Sprintf("%s can't be picked up.", state.GetDescription(pickable)))
	}
	player.ReplaceComponent(inventory)
}

func getCurrentInventoryItem(gs *gamestate.GameState) (*ecs.Entity, bool) {
	player := gs.Player
	inventory, _ := player.GetComponent(constants.Inventory).(state.InventoryComponent)

	sel := gs.InventoryScreenState.InventoryState.Selected
	if sel >= 0 && sel < len(inventory.Items) {
		item := inventory.Items[sel]
		return item, true
	}
	return nil, false
}

func getCurrentEquipmentItem(gs *gamestate.GameState) (*ecs.Entity, bool) {
	player := gs.Player
	equipment, _ := player.GetComponent(constants.Equipment).(state.EquipmentComponent)

	sel := gs.InventoryScreenState.EquipmentState.Selected
	if sel >= 0 && sel < len(constants.EquipmentSlots) {
		pos := constants.EquipmentSlots[sel]
		item, ok := equipment.Items[pos]
		if ok {
			return item, true
		}
	}
	return nil, false
}

func InventoryConsume(engine *ecs.Engine, gs *gamestate.GameState) {

	consumable, ok := getCurrentInventoryItem(gs)
	if ok {
		ConsumeConsumableComponent(engine, gs, consumable)

		updateInventorySelection(gs, 0)
	} else {
		gs.Log(constants.Warn, "No items to consume.")
	}
}

func InventoryDrop(gs *gamestate.GameState) {
	player := gs.Player
	inventory, _ := player.GetComponent(constants.Inventory).(state.InventoryComponent)
	position := player.GetComponent(constants.Position).(state.PositionComponent)

	invetoryItem, ok := getCurrentInventoryItem(gs)
	if ok {
		// TODO: move to it's own System
		gs.Log(constants.Info, "You dropped "+state.GetLongDescription(invetoryItem))

		// Remove from inventory
		inventory.RemoveItem(invetoryItem)
		player.ReplaceComponent(inventory)

		// Set new position
		invetoryItem.AddComponent(state.PositionComponent{
			X: position.X,
			Y: position.Y,
			Z: position.Z,
		})
		invetoryItem.AddComponent(state.IsPickupComponent{})

		updateInventorySelection(gs, 0)
	} else {
		gs.Log(constants.Warn, "No items to drop.")
	}
}

func InventoryEquip(gs *gamestate.GameState) {
	item, ok := getCurrentInventoryItem(gs)
	if ok {
		EquipEntity(gs, item)
	} else {
		gs.Log(constants.Warn, "No items to equip.")
	}
}

func EquipEntity(gs *gamestate.GameState, item *ecs.Entity) {
	player := gs.Player
	inventory, _ := player.GetComponent(constants.Inventory).(state.InventoryComponent)
	leveling, _ := player.GetComponent(constants.Leveling).(state.LevelingComponent)

	equipable, isEquipable := item.GetComponent(constants.Equipable).(state.EquipableComponent)

	if isEquipable {
		isLevelCorrect := leveling.CurrentLevel >= equipable.Level
		if isLevelCorrect {
			// Remove from inventory
			removed := inventory.RemoveItem(item)
			if removed {
				player.ReplaceComponent(inventory)
			}

			// Remove position if came from the floor
			item.RemoveComponent(constants.Position)

			// Add to equip
			equipment := player.GetComponent(constants.Equipment).(state.EquipmentComponent)
			current, ok := equipment.Items[equipable.EquipSlot]
			if ok {
				inventory.AddItem(current)
				player.ReplaceComponent(inventory)
			}
			equipment.Items[equipable.EquipSlot] = item
			player.ReplaceComponent(equipment)

			gs.Log(constants.Info, "You equipped "+state.GetLongDescription(item))

		} else {
			gs.Log(constants.Warn, fmt.Sprintf("%s can't be equiped. You must be at least level %d.", state.GetDescription(item), equipable.Level))
		}

	} else {
		gs.Log(constants.Warn, state.GetDescription(item)+" can't be equiped.")
	}
}

func InventoryUnequip(gs *gamestate.GameState) {
	player := gs.Player

	item, ok := getCurrentEquipmentItem(gs)
	if ok {
		// TODO: move to it's own System
		gs.Log(constants.Info, "You unequipped "+state.GetLongDescription(item))

		// Remove from equipment
		equipable := item.GetComponent(constants.Equipable).(state.EquipableComponent)
		equipment := player.GetComponent(constants.Equipment).(state.EquipmentComponent)
		delete(equipment.Items, equipable.EquipSlot)
		player.ReplaceComponent(equipment)

		// Add to inventory
		// TODO: check if there is space in the inventory
		inventory, _ := player.GetComponent(constants.Inventory).(state.InventoryComponent)
		inventory.AddItem(item)
		player.ReplaceComponent(inventory)

		updateInventorySelection(gs, 0)
	} else {
		gs.Log(constants.Warn, "No items to unequip.")
	}
}

func updateInventorySelection(gs *gamestate.GameState, change int) {

	selected := gs.InventoryScreenState.InventoryState.Selected + change

	inventory, _ := gs.Player.GetComponent(constants.Inventory).(state.InventoryComponent)
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
	max := len(constants.EquipmentSlots)
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

package systems

import (
	"fmt"

	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/state"
)

func EquipEntity(gs *gamestate.GameState, item *ecs.Entity) {
	player := gs.Player
	inventory, _ := player.GetComponent(constants.Inventory).(state.InventoryComponent)
	leveling, _ := player.GetComponent(constants.Leveling).(state.LevelingComponent)

	equipable, isEquipable := item.GetComponent(constants.Equipable).(state.EquipableComponent)

	if isEquipable {
		isLevelCorrect := leveling.CurrentLevel >= equipable.Level
		if isItemEquipped(gs, item) {
			gs.Log(constants.Warn, fmt.Sprintf("Item %s is already equipped.", state.GetLongDescription(item)))

		} else if isLevelCorrect {
			// Remove from inventory
			removed := inventory.RemoveItem(item)
			if removed {
				player.ReplaceComponent(inventory)
			}

			// Remove position if came from the floor
			item.RemoveComponent(constants.Position)

			// Remove current item, and put it into the inventory
			equipment := player.GetComponent(constants.Equipment).(state.EquipmentComponent)
			current, ok := equipment.Items[equipable.EquipSlot]
			inventoryOk := true
			if ok {
				inventoryOk = inventory.AddItem(current)
				player.ReplaceComponent(inventory)
			}
			if inventoryOk {
				// Add to equip
				equipment.Items[equipable.EquipSlot] = item
				player.ReplaceComponent(equipment)

				gs.Log(constants.Info, "You equipped "+state.GetLongDescription(item))
			} else {
				gs.Log(constants.Warn, fmt.Sprintf("Inventory is full. Dropping equipped %s to the floor.", state.GetLongDescription(item)))
				playerPosition := player.GetComponent(constants.Position).(state.PositionComponent)
				DropEntities(gs.Engine, gs, playerPosition, ecs.EntityList{item})
			}

		} else {
			gs.Log(constants.Warn, fmt.Sprintf("%s can't be equiped. You must be at least level %d.", state.GetDescription(item), equipable.Level))
		}

	} else {
		gs.Log(constants.Warn, state.GetDescription(item)+" can't be equiped.")
	}
}

func isItemEquipped(gs *gamestate.GameState, equipable *ecs.Entity) bool {
	player := gs.Player
	equipment := player.GetComponent(constants.Equipment).(state.EquipmentComponent)
	for _, item := range equipment.Items {
		if item != nil && item.Id == equipable.Id {
			return true
		}
	}
	return false
}

func UnequipItem(gs *gamestate.GameState, item *ecs.Entity) {
	player := gs.Player
	equipment := player.GetComponent(constants.Equipment).(state.EquipmentComponent)

	// Remove from equipment
	equipable := item.GetComponent(constants.Equipable).(state.EquipableComponent)
	delete(equipment.Items, equipable.EquipSlot)
	player.ReplaceComponent(equipment)

	// Add to inventory
	inventory, _ := player.GetComponent(constants.Inventory).(state.InventoryComponent)
	inventoryOk := inventory.AddItem(item)
	if inventoryOk {
		player.ReplaceComponent(inventory)
		gs.Log(constants.Info, "You unequipped "+state.GetLongDescription(item))
	} else {
		gs.Log(constants.Warn, fmt.Sprintf("Inventory is full. Dropping equipped %s to the floor.", state.GetLongDescription(item)))
		playerPosition := player.GetComponent(constants.Position).(state.PositionComponent)
		DropEntities(gs.Engine, gs, playerPosition, ecs.EntityList{item})
	}
}

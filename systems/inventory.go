package systems

import (
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
			pickUpOk := inventory.PickUp(pickable)
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
		gs.Log(gamestate.Warn, "No pickables found at this location")
	}
}

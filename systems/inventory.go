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

func InventoryConsume(gs *gamestate.GameState, consumable *ecs.Entity) {
	player := gs.Player
	inventory, _ := player.GetComponent(state.Inventory).(state.InventoryComponent)

	// Consume by player
	conStats := consumable.GetComponent(state.Consumable).(state.ConsumableComponent)
	plStats := player.GetComponent(state.Stats).(state.StatsComponent)
	// First increates the max
	plStats.MaxHealth += conStats.MaxHealth
	plStats.Health = increase(plStats.Health, plStats.MaxHealth, conStats.Health)
	plStats.MaxDefense += conStats.MaxDefense
	plStats.Defense = increase(plStats.Defense, plStats.MaxDefense, conStats.Defense)
	plStats.MaxPower += conStats.MaxPower
	plStats.Power = increase(plStats.Power, plStats.MaxPower, conStats.Power)
	plStats.Fov += conStats.Fov

	gs.Log(gamestate.Info, "Consumed "+state.GetLongDescription(consumable))

	// Remove from inventory
	inventory.Drop(consumable)
	player.ReplaceComponent(state.Inventory, inventory)

	// Destroy entity
	engine := consumable.Engine
	engine.DestroyEntity(consumable)
}

func increase(current int, max int, incr int) int {
	current = current + incr
	if current > max {
		return max
	}
	return current
}

func InventoryDrop(gs *gamestate.GameState, consumable *ecs.Entity) {
	player := gs.Player
	inventory, _ := player.GetComponent(state.Inventory).(state.InventoryComponent)
	position := player.GetComponent(state.Position).(state.PositionComponent)

	gs.Log(gamestate.Info, "You dropped "+state.GetLongDescription(consumable))

	// Remove from inventory
	inventory.Drop(consumable)
	player.ReplaceComponent(state.Inventory, inventory)

	// Set new position
	consumable.AddComponent(state.Position, state.PositionComponent{
		X: position.X,
		Y: position.Y,
	})
}

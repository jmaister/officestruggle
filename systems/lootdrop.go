package systems

import (
	"math/rand"

	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/state"
)

func LootDropSystem(engine *ecs.Engine, gs *gamestate.GameState, lootDropItem *ecs.Entity) {
	if lootDropItem.HasComponent(constants.LootDrop) {
		lootDrop := lootDropItem.GetComponent(constants.LootDrop).(state.LootDropComponent)
		position := lootDropItem.GetComponent(constants.Position).(state.PositionComponent)

		// Remove corpse position to place it in an empty floor
		lootDropItem.RemoveComponent(constants.Position)

		// Money
		money := state.NewMoneyAmount(gs.Engine.NewEntity(), lootDrop.Coins)

		// Items + Money + corpse
		itemsToPlace := append(lootDrop.Entities, money, lootDropItem)

		DropEntities(engine, gs, position, itemsToPlace)

		lootDropItem.RemoveComponent(constants.LootDrop)
	}

}

// Spawn entities around a position randomly. If it can't find a suitable position around, it uses the same.
func DropEntities(engine *ecs.Engine, gs *gamestate.GameState, position state.PositionComponent, itemsToPlace ecs.EntityList) {
	// Spawn the items and money around the corpse, not in the same position
	for _, item := range itemsToPlace {
		// Find a free position around the corpse
		positioned := false
		for radius := 1; radius < 5 && !positioned; radius++ {
			candidates := grid.GetCircle(grid.Tile{
				X: position.X,
				Y: position.Y,
			}, radius)
			rand.Shuffle(len(candidates), func(i int, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })

			for _, candidate := range candidates {
				elems, found := engine.PosCache.GetByCoord(candidate.X, candidate.Y, gs.CurrentZ)
				if found && (len(elems) == 1 && elems[0].HasComponent(constants.IsFloor)) {
					item.AddComponent(state.PositionComponent{X: candidate.X, Y: candidate.Y, Z: gs.CurrentZ})
					positioned = true
					break
				}
			}
		}
		// If we have not found a valid position, place it in the same position as the corpse
		if !positioned {
			item.AddComponent(position)
		}
	}

}

package systems

import (
	"strconv"

	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/game"
	"jordiburgos.com/officestruggle/state"
)

func Attack(gs *game.GameState, attacker *ecs.Entity, blockers ecs.EntityList) {

	// Try if attacker has Stats
	if attacker.HasComponent(state.Stats) {
		aStats := attacker.GetComponent(state.Stats).(state.StatsComponent)
		gs.L.Println("pl stats", aStats)

		// Check every blocker if it has Stats
		for _, blocker := range blockers {
			if blocker.HasComponent(state.Stats) {
				bStats := blocker.GetComponent(state.Stats).(state.StatsComponent)

				// Damage calculation and attack
				damage := aStats.Power - bStats.Defense
				if damage >= 0 {
					gs.L.Println(state.GetDescription(attacker) + " hit " + state.GetDescription(blocker) + " with " + strconv.Itoa(damage) + " points.")
					newHealth := bStats.Health - damage
					if newHealth <= 0 {
						Kill(blocker)
					} else {
						bStats.Health = newHealth
						blocker.ReplaceComponent(state.Stats, bStats)
					}
				} else {
					gs.L.Println(state.GetDescription(blocker) + " blocked attack from " + state.GetDescription(attacker) + ".")
				}
			}
		}
	}
}

func Kill(entity *ecs.Entity) {
	entity.RemoveComponent(state.AI)
	entity.RemoveComponent(state.IsBlocking)
	entity.RemoveComponent(state.Layer400)
	entity.AddComponent(state.Layer300, state.Layer300Component{})

	apparence, ok := entity.GetComponent(state.Apparence).(state.ApparenceComponent)
	if ok {
		apparence.Char = '%'
		entity.ReplaceComponent(state.Apparence, apparence)
	}
}

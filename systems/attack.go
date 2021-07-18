package systems

import (
	"strconv"

	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/state"
)

func Attack(gs *gamestate.GameState, attacker *ecs.Entity, blockers ecs.EntityList) {

	// Try if attacker has Stats
	if attacker.HasComponent(state.Stats) {
		aStats := attacker.GetComponent(state.Stats).(state.StatsComponent)

		// Check every blocker if it has Stats
		for _, blocker := range blockers {
			if blocker.HasComponent(state.Stats) {
				bStats := blocker.GetComponent(state.Stats).(state.StatsComponent)

				// Damage calculation and attack
				damage := aStats.Power - bStats.Defense
				if damage >= 0 {
					gs.Log(state.GetDescription(attacker) + " attacks " + state.GetDescription(blocker) + " with " + strconv.Itoa(damage) + " damage points.")
					newHealth := bStats.Health - damage
					if newHealth <= 0 {
						Kill(gs, blocker)
					} else {
						bStats.Health = newHealth
						blocker.ReplaceComponent(state.Stats, bStats)
					}
				} else {
					gs.Log(state.GetDescription(blocker) + " blocked attack from " + state.GetDescription(attacker) + ".")
				}
			} else if blocker.HasComponent(state.Visitable) {
				gs.Log(state.GetDescription(attacker) + " hits a " + state.GetDescription(blocker))
			}
		}
	}
}

func Kill(gs *gamestate.GameState, entity *ecs.Entity) {
	gs.Log(state.GetDescription(entity) + " is dead.")
	entity.RemoveComponent(state.AI)
	entity.RemoveComponent(state.IsBlocking)
	entity.RemoveComponent(state.Layer400)
	entity.AddComponent(state.Layer300, state.Layer300Component{})
	entity.AddComponent(state.Dead, state.DeadComponent{})

	apparence, ok := entity.GetComponent(state.Apparence).(state.ApparenceComponent)
	if ok {
		apparence.Char = '%'
		entity.ReplaceComponent(state.Apparence, apparence)
	}
}

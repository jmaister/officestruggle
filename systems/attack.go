package systems

import (
	"strconv"
	"time"

	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/interfaces"
	"jordiburgos.com/officestruggle/state"
)

func Attack(engine *ecs.Engine, gs *gamestate.GameState, attacker *ecs.Entity, blockers ecs.EntityList) {

	// Try if attacker has Stats
	if attacker.HasComponent(constants.Stats) {
		aStats := attacker.GetComponent(constants.Stats).(state.StatsComponent)

		// Check every blocker if it has Stats
		for _, blocker := range blockers {
			if blocker.HasComponent(constants.Stats) {
				bStats := blocker.GetComponent(constants.Stats).(state.StatsComponent)

				// Damage calculation and attack
				damage := aStats.Power - bStats.Defense
				if damage >= 0 {
					gs.Log(constants.Danger, state.GetDescription(attacker)+" attacks "+state.GetDescription(blocker)+" with "+strconv.Itoa(damage)+" damage points.")
					newHealth := bStats.Health - damage

					CreateDamageAnimation(engine, attacker, blocker, strconv.Itoa(damage))

					if newHealth <= 0 {
						Kill(gs, blocker)
					} else {
						bStats.Health = newHealth
						blocker.ReplaceComponent(bStats)
					}
				} else {
					gs.Log(constants.Danger, state.GetDescription(blocker)+" blocked attack from "+state.GetDescription(attacker)+".")
				}
			} else if blocker.HasComponent(constants.Visitable) {
				gs.Log(constants.Warn, state.GetDescription(attacker)+" hits a "+state.GetDescription(blocker))
			}
		}
	}
}

func AttackWithItem(engine *ecs.Engine, gs *gamestate.GameState, attacker *ecs.Entity, blocker *ecs.Entity, entityUsed *ecs.Entity, damage int) {

	// Check if blocker has Stats
	if blocker.HasComponent(constants.Stats) {
		bStats := blocker.GetComponent(constants.Stats).(state.StatsComponent)

		// Damage calculation and attack
		damage := damage - bStats.Defense
		if damage >= 0 {
			gs.Log(constants.Danger, state.GetDescription(attacker)+" attacks "+state.GetDescription(blocker)+" using "+state.GetDescription(entityUsed)+" with "+strconv.Itoa(damage)+" damage points.")
			newHealth := bStats.Health - damage

			if newHealth <= 0 {
				Kill(gs, blocker)
			} else {
				bStats.Health = newHealth
				blocker.ReplaceComponent(bStats)
			}
		} else {
			gs.Log(constants.Danger, state.GetDescription(blocker)+" blocked attack from "+state.GetDescription(attacker)+".")
		}
	} else if blocker.HasComponent(constants.Visitable) {
		gs.Log(constants.Warn, state.GetDescription(attacker)+" hits a "+state.GetDescription(blocker))
	}
}

func Kill(gs *gamestate.GameState, entity *ecs.Entity) {
	gs.Log(constants.Good, state.GetDescription(entity)+" is dead.")
	entity.RemoveComponent(constants.AI)
	entity.RemoveComponent(constants.IsBlocking)
	entity.RemoveComponent(constants.Layer400)
	entity.RemoveComponent(constants.Stats)
	entity.AddComponent(state.Layer300Component{})
	entity.AddComponent(state.DeadComponent{})
	entity.AddComponent(state.IsPickupComponent{})

	apparence, ok := entity.GetComponent(constants.Apparence).(state.ApparenceComponent)
	if ok {
		apparence.Char = '%'
		entity.ReplaceComponent(apparence)
	}
	// TODO: if player is killed, change screen state
}

func CreateDamageAnimation(engine *ecs.Engine, source *ecs.Entity, target *ecs.Entity, str string) {
	animationEntity := engine.NewEntity()
	animationEntity.AddComponent(state.Layer500Component{})

	aPos := source.GetComponent(constants.Position).(state.PositionComponent)
	bPos := target.GetComponent(constants.Position).(state.PositionComponent)

	dir := grid.UP
	if aPos.X == bPos.X {
		if aPos.Y > bPos.Y {
			// b
			// @
			dir = grid.UP_RIGHT
		} else {
			// @
			// b
			dir = grid.DOWN_RIGHT
		}
	} else if aPos.Y == bPos.Y {
		if aPos.X > bPos.X {
			// b@
			dir = grid.UP_LEFT
		} else {
			// @b
			dir = grid.UP_RIGHT
		}
	}

	animationEntity.AddComponent(state.AnimatedComponent{
		Animation: DamageAnimation{
			AnimationInfo: interfaces.AnimationInfo{
				StartTime: time.Now(),
				Duration:  750 * time.Millisecond,
				Source:    source,
				Target:    target,
			},
			Direction: dir,
			Damage:    str,
		},
	})

}

package systems

import (
	"strconv"
	"time"

	"jordiburgos.com/officestruggle/animations"
	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/grid"
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
					gs.Log(gamestate.Danger, state.GetDescription(attacker)+" attacks "+state.GetDescription(blocker)+" with "+strconv.Itoa(damage)+" damage points.")
					newHealth := bStats.Health - damage

					aPos := attacker.GetComponent(constants.Position).(state.PositionComponent)
					bPos := blocker.GetComponent(constants.Position).(state.PositionComponent)
					createDamageAnimation(engine, aPos, bPos, strconv.Itoa(damage))

					if newHealth <= 0 {
						Kill(gs, blocker)
					} else {
						bStats.Health = newHealth
						blocker.ReplaceComponent(bStats)
					}
				} else {
					gs.Log(gamestate.Danger, state.GetDescription(blocker)+" blocked attack from "+state.GetDescription(attacker)+".")
				}
			} else if blocker.HasComponent(constants.Visitable) {
				gs.Log(gamestate.Warn, state.GetDescription(attacker)+" hits a "+state.GetDescription(blocker))
			}
		}
	}
}

func Kill(gs *gamestate.GameState, entity *ecs.Entity) {
	gs.Log(gamestate.Good, state.GetDescription(entity)+" is dead.")
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

func createDamageAnimation(engine *ecs.Engine, aPos state.PositionComponent, bPos state.PositionComponent, str string) {
	animationEntity := engine.NewEntity()
	animationEntity.AddComponent(state.Layer500Component{})

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

	animationEntity.AddComponent(animations.AnimatedComponent{
		Animation: animations.DamageAnimation{
			X:                 bPos.X,
			Y:                 bPos.Y,
			Direction:         dir,
			Damage:            str,
			AnimationStart:    time.Now(),
			AnimationDuration: 750 * time.Millisecond,
		},
	})

}

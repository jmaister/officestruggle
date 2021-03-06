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
			// Attacks can be done by the player or to the player, avoid fights between mobs
			if blocker.HasComponent(constants.Stats) && (attacker.HasComponent(constants.Player) || blocker.HasComponent(constants.Player)) {
				bStats := blocker.GetComponent(constants.Stats).(state.StatsComponent)

				// Damage calculation and attack
				damage := aStats.Power - bStats.Defense
				if damage >= 0 {
					t := constants.Danger
					if attacker == gs.Player {
						t = constants.Good
					}
					gs.Log(t, state.GetDescription(attacker)+" attacks "+state.GetDescription(blocker)+" with "+strconv.Itoa(damage)+" damage points.")
					newHealth := bStats.Health - damage

					aPos := state.GetPosition(attacker)
					bPos := state.GetPosition(blocker)
					CreateDamageAnimation(engine, interfaces.Point{
						X: aPos.X,
						Y: aPos.Y,
					}, interfaces.Point{
						X: bPos.X,
						Y: bPos.Y,
					}, strconv.Itoa(damage))

					if newHealth <= 0 {
						Kill(engine, gs, attacker, blocker)
					} else {
						bStats.Health = newHealth
						blocker.ReplaceComponent(bStats)
					}
				} else {
					t := constants.Danger
					if blocker == gs.Player {
						t = constants.Good
					}
					gs.Log(t, state.GetDescription(blocker)+" blocked attack from "+state.GetDescription(attacker)+".")
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
			t := constants.Danger
			if attacker == gs.Player {
				t = constants.Good
			}
			gs.Log(t, state.GetDescription(attacker)+" attacks "+state.GetDescription(blocker)+" using "+state.GetDescription(entityUsed)+" with "+strconv.Itoa(damage)+" damage points.")
			newHealth := bStats.Health - damage

			if newHealth <= 0 {
				Kill(engine, gs, attacker, blocker)
			} else {
				bStats.Health = newHealth
				blocker.ReplaceComponent(bStats)
			}
		} else {
			t := constants.Danger
			if blocker == gs.Player {
				t = constants.Good
			}
			gs.Log(t, state.GetDescription(blocker)+" blocked attack from "+state.GetDescription(attacker)+".")
		}
	} else if blocker.HasComponent(constants.Visitable) {
		gs.Log(constants.Warn, state.GetDescription(attacker)+" hits a "+state.GetDescription(blocker))
	}
}

func Kill(engine *ecs.Engine, gs *gamestate.GameState, attacker *ecs.Entity, entity *ecs.Entity) {
	gs.Log(constants.Good, state.GetDescription(entity)+" is dead.")

	entity.AddComponent(state.DeadComponent{})
	entity.RemoveComponent(constants.AI)
	entity.RemoveComponent(constants.Stats)
	entity.RemoveComponent(constants.IsBlocking)
	entity.RemoveComponent(constants.Layer400)
	entity.AddComponent(state.Layer300Component{})
	entity.AddComponent(state.IsPickupComponent{})

	description := entity.GetComponent(constants.Description).(state.DescriptionComponent)
	description.Name += " corpse"
	entity.ReplaceComponent(description)

	apparence := entity.GetComponent(constants.Apparence).(state.ApparenceComponent)
	apparence.Char = '%'
	entity.ReplaceComponent(apparence)

	if entity == gs.Player {
		gs.ScreenState = gamestate.GameoverScreen
	} else if attacker == gs.Player {
		GiveXP(gs, gs.Player, entity)
	}

	// Win the game killed entity has a "WinGame" component
	if entity.HasComponent(constants.WinGame) {
		gs.ScreenState = gamestate.GameWinScreen
	}

	// Default XP for eating a corpse
	entity.ReplaceComponent(state.XPGiverComponent{
		XPBase: 10,
	})

	// LootDrop
	LootDropSystem(engine, gs, entity)
}

func CreateDamageAnimation(engine *ecs.Engine, source interfaces.Point, target interfaces.Point, str string) {
	animationEntity := engine.NewEntity()
	animationEntity.AddComponent(state.Layer500Component{})

	aPos := source
	bPos := target

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
		Animation: LeavingStringAnimation{
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

package systems

import (
	"math/rand"
	"time"

	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/state"
)

func ConsumeConsumableComponent(gs *gamestate.GameState, consumable *ecs.Entity) bool {

	isConsumable := consumable.HasComponent(constants.Consumable)
	isConsumeEffect := consumable.HasComponent(constants.ConsumeEffect)
	if isConsumable {

		// Consume by player
		player := gs.Player
		conStats := consumable.GetComponent(constants.Consumable).(state.ConsumableComponent)
		plStats := player.GetComponent(constants.Stats).(state.StatsComponent)
		apparence := player.GetComponent(constants.Apparence).(state.ApparenceComponent)

		newStats := plStats.Merge(*conStats.StatsValues)
		player.ReplaceComponent(state.StatsComponent{
			StatsValues: &newStats,
		})
		player.AddComponent(state.AnimatedComponent{
			Animation: HealthPotionAnimation{
				AnimationStart:    time.Now(),
				AnimationDuration: 1 * time.Second,
				StartingApparence: apparence,
			},
		})

		gs.Log(constants.Info, "Consumed "+state.GetLongDescription(consumable))
		return true
	} else if isConsumeEffect {
		consumeEffect := consumable.GetComponent(constants.ConsumeEffect).(state.ConsumeEffectComponent)

		damagePerEnemy := consumeEffect.Damage
		if consumeEffect.DamageType == gamestate.DamageSharedType {
			damagePerEnemy = damagePerEnemy / consumeEffect.TargetCount
		}

		enemiesInFov := getEnemiesInFov(gs, consumeEffect)

		switch consumeEffect.Targeting {
		case gamestate.RandomAcquisitionType:
			// Find enemies on FOV
			// Select n randomly
			if len(enemiesInFov) > 0 {
				// Attack enemies
				for i := 0; i < consumeEffect.TargetCount && len(enemiesInFov) > 0; i++ {
					target := enemiesInFov[rand.Intn(len(enemiesInFov))]
					AttackWithItem(gs.Engine, gs, gs.Player, target, consumable, damagePerEnemy)

					animation := consumeEffect.Animation
					if animation.NeedsInit() {
						animation = animation.Init(gs.Player, target)
					}
					e := gs.Engine.NewEntity()
					e.AddComponent(state.AnimatedComponent{
						Animation: animation,
					})

					// Recalculate as enemies could be killed
					enemiesInFov = getEnemiesInFov(gs, consumeEffect)
				}
			} else {
				gs.Log(constants.Warn, state.GetLongDescription(consumable)+" is used but no targets found.")
			}
			return true
		default:
			// TODO: add other target types
			panic("Add other target types")
		}
	} else {
		gs.Log(constants.Warn, state.GetLongDescription(consumable)+" is not consumable.")
		return false
	}

}

func getEnemiesInFov(gs *gamestate.GameState, consumeEffect state.ConsumeEffectComponent) ecs.EntityList {
	enemiesInFov := ecs.EntityList{}
	for _, enemy := range gs.Engine.Entities.GetEntities(consumeEffect.TargetTypes) {
		position := enemy.GetComponent(constants.Position).(state.PositionComponent)
		if gs.Fov.IsVisible(position.X, position.Y) {
			enemiesInFov = append(enemiesInFov, enemy)
		}
	}
	return enemiesInFov
}

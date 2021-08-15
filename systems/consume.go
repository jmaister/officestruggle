package systems

import (
	"math/rand"
	"time"

	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/interfaces"
	"jordiburgos.com/officestruggle/state"
)

func ConsumeConsumableComponent(engine *ecs.Engine, gs *gamestate.GameState, consumable *ecs.Entity) bool {

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
				AnimationInfo: interfaces.AnimationInfo{
					StartTime: time.Now(),
					Duration:  1 * time.Second,
					Source:    &ecs.Entity{},
					Target:    &ecs.Entity{},
				},
				StartingApparence: apparence,
			},
		})

		gs.Log(constants.Info, "Consumed "+state.GetLongDescription(consumable))
		removeAndDestroy(gs.Engine, gs, consumable)

		if newStats.Health <= 0 {
			Kill(engine, gs, consumable, player)
		}

		return true
	} else if isConsumeEffect {
		consumeEffect := consumable.GetComponent(constants.ConsumeEffect).(state.ConsumeEffectComponent)

		enemiesInFov := getEnemiesInFov(gs, consumeEffect)

		switch consumeEffect.Targeting {
		case gamestate.RandomAcquisitionType:
			// Find enemies on FOV
			// Select n randomly
			if len(enemiesInFov) > 0 {
				// Apply EffectFunction on each target
				for i := 0; i < consumeEffect.TargetCount && len(enemiesInFov) > 0; i++ {
					target := enemiesInFov[rand.Intn(len(enemiesInFov))]

					consumeEffect.EffectFunction(gs.Engine, gs, consumable, gs.Player, target)

					animation := consumeEffect.EffectAnimation
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
			removeAndDestroy(gs.Engine, gs, consumable)
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

	enemies := gs.Engine.Entities.GetEntities(consumeEffect.TargetTypes)
	enemies = FilterZ(enemies, gs.CurrentZ)

	for _, enemy := range enemies {
		position := enemy.GetComponent(constants.Position).(state.PositionComponent)
		if gs.Fov.IsVisible(position.X, position.Y) {
			enemiesInFov = append(enemiesInFov, enemy)
		}
	}
	return enemiesInFov
}

func removeAndDestroy(engine *ecs.Engine, gs *gamestate.GameState, consumable *ecs.Entity) {
	player := gs.Player
	inventory, _ := player.GetComponent(constants.Inventory).(state.InventoryComponent)

	if consumable.HasComponent(constants.XPGiver) {
		GiveXP(gs, gs.Player, consumable)
	}

	// Remove from inventory
	inventory.RemoveItem(consumable)
	gs.Player.ReplaceComponent(inventory)

	// Destroy entity
	engine.DestroyEntity(consumable)
}

package systems

import (
	"fmt"
	"time"

	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/interfaces"
	"jordiburgos.com/officestruggle/palette"
	"jordiburgos.com/officestruggle/state"
)

func NewLightningScroll(entity *ecs.Entity) *ecs.Entity {
	entity.AddComponent(state.IsPickupComponent{})
	entity.AddComponent(state.DescriptionComponent{Name: "Lightning Scroll"})
	entity.AddComponent(state.ApparenceComponent{Color: "#DAA520", Char: '♪'})
	entity.AddComponent(state.Layer300Component{})
	entity.AddComponent(state.ConsumeEffectComponent{
		Targeting:   gamestate.RandomAcquisitionType,
		TargetTypes: []string{constants.AI},
		TargetCount: 3,
		EffectAnimation: FallingCharAnimation{
			AnimationInfo: interfaces.AnimationInfo{
				StartTime: time.Time{},
				Duration:  750 * time.Millisecond,
				Source:    &ecs.Entity{},
				Target:    &ecs.Entity{},
			},
			Direction: grid.DOWN,
			Char:      "♪",
			Color:     palette.PColor(palette.Blue, 0.5),
			Text:      "10",
		},
		EffectFunction: func(engine *ecs.Engine, gs *gamestate.GameState, item *ecs.Entity, source *ecs.Entity, target *ecs.Entity) {
			AttackWithItem(gs.Engine, gs, gs.Player, target, item, 10)
		},
	})
	return entity
}

func NewParalizeScroll(entity *ecs.Entity) *ecs.Entity {
	entity.AddComponent(state.IsPickupComponent{})
	entity.AddComponent(state.DescriptionComponent{Name: "Paralize Scroll"})
	entity.AddComponent(state.ApparenceComponent{Color: "#FFB0B0", Char: '♪'})
	entity.AddComponent(state.Layer300Component{})
	entity.AddComponent(state.ConsumeEffectComponent{
		// Targeting:   gamestate.ManualAcquisitionType,
		Targeting:   gamestate.RandomAcquisitionType,
		TargetTypes: []string{constants.AI},
		TargetCount: 3,
		EffectAnimation: FallingCharAnimation{
			AnimationInfo: interfaces.AnimationInfo{
				StartTime: time.Time{},
				Duration:  750 * time.Millisecond,
				Source:    &ecs.Entity{},
				Target:    &ecs.Entity{},
			},
			Direction: grid.DOWN,
			Char:      "♪",
			Color:     palette.PColor(palette.Orange, 0.7),
			Text:      "Paralized",
		},
		EffectFunction: func(engine *ecs.Engine, gs *gamestate.GameState, item *ecs.Entity, source *ecs.Entity, target *ecs.Entity) {
			turnsLeft := 5
			if target.HasComponent(constants.Paralize) {
				current, _ := target.GetComponent(constants.Paralize).(state.ParalizeComponent)
				turnsLeft = turnsLeft + current.TurnsLeft
			}

			// TODO: add visual tags to ecs.Entity to show paralized, freezed, dizzy, blind, ...
			target.ReplaceComponent(state.ParalizeComponent{
				TurnsLeft: turnsLeft,
			})
			gs.Log(constants.Info, fmt.Sprintf("%s got paralized by %d turns.", state.GetDescription(target), turnsLeft))
		},
	})
	return entity
}

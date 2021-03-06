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

func FilterZ(entitylist ecs.EntityList, z int) ecs.EntityList {
	filteredList := ecs.EntityList{}

	for _, e := range entitylist {
		pos, ok := e.GetComponent(constants.Position).(state.PositionComponent)
		if ok && pos.Z == z {
			filteredList = append(filteredList, e)
		}
	}

	return filteredList
}

func FilterNot(entitylist ecs.EntityList, componentType string) ecs.EntityList {
	filteredList := ecs.EntityList{}

	for _, e := range entitylist {
		if !e.HasComponent(componentType) {
			filteredList = append(filteredList, e)
		}
	}

	return filteredList
}

func FilterFunc(entitylist ecs.EntityList, fn func(*ecs.Entity) bool) ecs.EntityList {
	filteredList := ecs.EntityList{}

	for _, e := range entitylist {
		if fn(e) {
			filteredList = append(filteredList, e)
		}
	}

	return filteredList
}

func NewLightningScroll(entity *ecs.Entity) *ecs.Entity {
	entity.AddComponent(state.IsPickupComponent{})
	entity.AddComponent(state.DescriptionComponent{Name: "Lightning Scroll"})
	entity.AddComponent(state.ApparenceComponent{Color: palette.PColor(palette.Orange, 0.6), Char: '~'})
	entity.AddComponent(state.Layer300Component{})
	entity.AddComponent(state.ConsumeEffectComponent{
		Targeting:   gamestate.RandomAcquisitionType,
		TargetTypes: []string{constants.AI},
		TargetCount: 3,
		EffectAnimation: FallingCharAnimation{
			AnimationInfo: interfaces.AnimationInfo{
				StartTime: time.Time{},
				Duration:  750 * time.Millisecond,
			},
			Direction: grid.DOWN,
			Char:      "~",
			Color:     palette.PColor(palette.Blue, 0.5),
			Text:      "10",
		},
		EffectFunction: func(engine *ecs.Engine, gs *gamestate.GameState, item *ecs.Entity, source *ecs.Entity, target *ecs.Entity) {
			AttackWithItem(gs.Engine, gs, gs.Player, target, item, 10)
		},
	})
	entity.AddComponent(state.XPGiverComponent{
		XPBase:     5,
		XPPerLevel: 10,
		Level:      1,
	})

	return entity
}

func NewParalizeScroll(entity *ecs.Entity) *ecs.Entity {
	entity.AddComponent(state.IsPickupComponent{})
	entity.AddComponent(state.DescriptionComponent{Name: "Paralize Scroll"})
	entity.AddComponent(state.ApparenceComponent{Color: palette.PColor(palette.Pink, 0.6), Char: '~'})
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
			},
			Direction: grid.DOWN,
			Char:      "~",
			Color:     palette.PColor(palette.Orange, 0.7),
			Text:      "Paralized",
		},
		EffectFunction: func(engine *ecs.Engine, gs *gamestate.GameState, item *ecs.Entity, source *ecs.Entity, target *ecs.Entity) {
			turnsLeft := 5
			if target.HasComponent(constants.Paralize) {
				current, _ := target.GetComponent(constants.Paralize).(state.ParalizeComponent)
				turnsLeft = turnsLeft + current.TurnsLeft
			}

			target.ReplaceComponent(state.ParalizeComponent{
				TurnsLeft: turnsLeft,
			})
			gs.Log(constants.Info, fmt.Sprintf("%s got paralized for %d turns.", state.GetDescription(target), turnsLeft))
		},
	})
	entity.AddComponent(state.XPGiverComponent{
		XPBase:     5,
		XPPerLevel: 10,
		Level:      1,
	})
	return entity
}

package systems

import (
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
	entity.AddComponent(state.DescriptionComponent{Name: "Lightning scroll"})
	entity.AddComponent(state.ApparenceComponent{Color: "#DAA520", Char: '♪'})
	entity.AddComponent(state.Layer300Component{})
	entity.AddComponent(state.ConsumeEffectComponent{
		Targeting:   gamestate.RandomAcquisitionType,
		TargetTypes: []string{constants.AI},
		TargetCount: 3,
		Damage:      10,
		DamageType:  gamestate.DamageEachType,
		Animation: FallingCharAnimation{
			AnimationInfo: interfaces.AnimationInfo{
				StartTime: time.Time{},
				Duration:  1 * time.Second,
				Source:    &ecs.Entity{},
				Target:    &ecs.Entity{},
			},
			Direction: grid.DOWN,
			Char:      "♪",
			Color:     palette.PColor(palette.Blue, 0.5),
		},
	})
	return entity
}

package actions

import (
	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
)

type TargetingType string

const RandomAcquisitionType = "random"
const SelectedAcquisitionType = "selected"
const AreaAcquisitionType = "area"

type RequiresTargetComponent struct {
	Targeting   TargetingType
	TargetTypes []string
	OnSelect    func(engine *ecs.Engine, gs *gamestate.GameState, attacker *ecs.Entity, targets ecs.EntityList)
}

func (a RequiresTargetComponent) ComponentType() string {
	return constants.RequiresTarget
}

func LightningScrollRandomTarget(engine *ecs.Engine, gs *gamestate.GameState, attacker *ecs.Entity, targets ecs.EntityList) {
	//rndEnemy := targets[rand.Intn(len(targets))]

	// TODO: create Attack system with more options: send damage
	// TODO: create targeting system, calculate targets, log if no targets available, destroy used item? add Used() method
	// systems.Attack(engine, gs, attacker, ecs.EntityList{rndEnemy})
}

package systems

import (
	"math/rand"

	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
)

func TargetingSystem() {

}

func LightningScrollRandomTargetfunc(engine *ecs.Engine, gs *gamestate.GameState, item *ecs.Entity, itemUser *ecs.Entity, targets ecs.EntityList) {
	rndEnemy := targets[rand.Intn(len(targets))]

	// TODO: create Attack system with more options: send damage
	// TODO: create targeting system, calculate targets, log if no targets available, destroy used item? add Used() method
	// systems.Attack(engine, gs, attacker, ecs.EntityList{rndEnemy})
	Attack(engine, gs, itemUser, ecs.EntityList{rndEnemy})
}
